/*
 * Copyright (c) 2019, WSO2 Inc. (http://www.wso2.com). All Rights Reserved.
 *
 * This software is the property of WSO2 Inc. and its suppliers, if any.
 * Dissemination of any information or reproduction of any material contained
 * herein is strictly forbidden, unless permitted by WSO2 in accordance with
 * the WSO2 Commercial License available at http://wso2.com/licenses. For specific
 * language governing the permissions and limitations under this license,
 * please see the license as well as any agreement you’ve entered into with
 * WSO2 governing the purchase of this software and any associated services.
 */

package login

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/wso2/choreo/components/cli/internal/pkg/cmd/common"
	"golang.org/x/oauth2"
)

func NewLoginCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "login",
		Short:   "Login to " + common.ProductName,
		Example: common.GetAbsoluteCommandName("login"),
		Args:    cobra.NoArgs,
		Run:     runLogin,
	}
}

func runLogin(cmd *cobra.Command, args []string) {
	authCode := getAuthCode()
	accessToken := getAccessToken(authCode)
	persistAccessToken(accessToken)
}

func persistAccessToken(accessToken string) {
	panic("persistAccessToken Method not implemented")
}

func getAccessToken(authCode string) string {
	log.Println("Received: " + authCode)
	panic("getAccessToken Method not implemented")
}

func getAuthCode() string {
	codeServicePort := common.GetFirstOpenPort(callBackDefaultPort)

	authCodeChannel := startAuthCodeReceivingService(codeServicePort, callbackUrlContext)

	openBrowserForAuthentication(callbackUrlContext, codeServicePort)

	authCode := <-authCodeChannel

	return authCode
}

func startAuthCodeReceivingService(port int, urlContext string) chan string {
	authCodeChannel := make(chan string)
	go listenForAuthCode(common.GetLocalBindAddress(port), urlContext, authCodeChannel)
	return authCodeChannel
}

func openBrowserForAuthentication(context string, port int) {
	callBackUrlTemplate := "http://localhost:%d" + context
	redirectUrl := fmt.Sprintf(callBackUrlTemplate, port)
	hubAuthUrl := getHubUrl(redirectUrl)
	log.Println("Callback URL: " + hubAuthUrl)
}

func getHubUrl(redirectUrl string) string {
	conf := &oauth2.Config{
		ClientID:    clientId,
		RedirectURL: redirectUrl,
		Endpoint: oauth2.Endpoint{
			AuthURL: authUrl,
		},
	}

	return conf.AuthCodeURL("state")
}

func listenForAuthCode(addrString string, callbackUrlContext string, authCodeChannel chan<- string) {
	mux := http.NewServeMux()
	server := http.Server{Addr: addrString, Handler: mux}
	mux.HandleFunc(callbackUrlContext, func(writer http.ResponseWriter, request *http.Request) {
		if err := request.ParseForm(); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			common.ExitWithError("Error parsing received query parameters", err)
		}

		code := request.Form.Get("code")

		if code == "" {
			writer.WriteHeader(http.StatusBadRequest)
			common.ExitWithErrorMessage("Blank auth code received from IDP")
		}

		writer.WriteHeader(http.StatusOK)
		authCodeChannel <- code
	})

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		common.ExitWithError("Error while initializing auth code accepting service", err)
	}
}

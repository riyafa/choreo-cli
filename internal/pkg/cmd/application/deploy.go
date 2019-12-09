/*
 * Copyright (c) 2019, WSO2 Inc. (http://www.wso2.com). All Rights Reserved.
 *
 * This software is the property of WSO2 Inc. and its suppliers, if any.
 * Dissemination of any information or reproduction of any material contained
 * herein in any form is strictly forbidden, unless permitted by WSO2 expressly.
 * You may not alter or remove any copyright or other notice from copies of this content.
 */

package application

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wso2/choreo-cli/internal/pkg/client"
	"github.com/wso2/choreo-cli/internal/pkg/cmd/common"
	"github.com/wso2/choreo-cli/internal/pkg/cmd/runtime"
)

func NewDeployCommand(cliContext runtime.CliContext) *cobra.Command {
	const cmdDeploy = "deploy"
	cmd := &cobra.Command{
		Use:   cmdDeploy + " GITHUB_REPO_URL",
		Short: "Deploy an application",
		Example: fmt.Sprint(common.GetAbsoluteCommandName(cmdApplication, cmdDeploy),
			" https://github.com/wso2/choreo-ballerina-hello"),
		Args: cobra.ExactArgs(1),
		Run:  runDeployAppCommand(cliContext),
	}
	cmd.Flags().StringP("name", "n", "", "the name to be used for the created application")
	return cmd
}

func runDeployAppCommand(cliContext runtime.CliContext) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if !client.IsUserLoggedIn(cliContext) {
			common.ExitWithErrorMessage(cliContext.Out(), "Please login first")
		}

		appName, err := cmd.Flags().GetString("name")
		if err !=nil {
			common.ExitWithError(cliContext.Out(), "Error while reading the application name flag value", err)
		}

		if appName != "" {
			deploymentDetails, err := cliContext.Client().CreateAndDeployAppWithName(appName, args[0])
			printDeployResponse(err, cliContext, deploymentDetails)
		} else {
			deploymentDetails, err := cliContext.Client().CreateAndDeployApp(args[0])
			printDeployResponse(err, cliContext, deploymentDetails)
		}

	}
}

func printDeployResponse(err error, cliContext runtime.CliContext, deploymentDetails runtime.DeploymentDetails) {
	if err != nil {
		common.ExitWithError(cliContext.Out(), "Error occurred while deploying the application", err)
	} else {
		common.PrintInfo(cliContext.Out(), "A new application is created for deployment with Id: "+
			deploymentDetails.ApplicationId+"\nOnce deployed, the app can be accessed from "+
			deploymentDetails.DeploymentUrl)
	}
}

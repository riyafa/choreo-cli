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

package main

import (
	"github.com/spf13/cobra"
	"github.com/wso2/choreo/components/cli/internal/pkg/cmd"
	cmdCommon "github.com/wso2/choreo/components/cli/internal/pkg/cmd/common"
	"github.com/wso2/choreo/components/cli/internal/pkg/cmd/login"
)

func main() {

	command := cobra.Command{
		Use:   "choreo <command>",
		Short: "Manage integration applications with Choreo platform",
	}

	command.AddCommand(cmd.NewVersionCommand())
	command.AddCommand(login.NewLoginCommand())

	if err := command.Execute(); err != nil {
		cmdCommon.ExitWithError("Error executing choreo command", err)
	}
}

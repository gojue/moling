/*
 *
 *  Copyright 2025 CFC4N <cfc4n.cs@gmail.com>. All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 *  Repository: https://github.com/gojue/moling
 *
 */

package cmd

import (
	"github.com/gojue/moling/client"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Provides automated access to MoLing MCP Server for local MCP clients, Cline, Roo Code, and Claude, etc.",
	Long: `Automatically checks the MCP clients installed on the current computer, displays them, and automatically adds the MoLing MCP Server configuration to enable one-click activation, reducing the hassle of manual configuration.
Currently supports the following clients: Cline, Roo Code, Claude
    moling client -l --list   List the current installed MCP clients
    moling client -i --install Add MoLing MCP Server configuration to the currently installed MCP clients on this computer
`,
	RunE: ClientCommandFunc,
}

var (
	list    bool
	install bool
)

// ClientCommandFunc executes the "config" command.
func ClientCommandFunc(command *cobra.Command, args []string) error {
	logger := initLogger(mlConfig.BasePath)
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	multi := zerolog.MultiLevelWriter(consoleWriter, logger)
	logger = zerolog.New(multi).With().Timestamp().Logger()
	mlConfig.SetLogger(logger)
	logger.Info().Msg("Start to show MCP Clients")
	mcpConfig := client.NewMCPServerConfig(CliDescription, CliName, MCPServerName)
	exePath, err := os.Executable()
	if err == nil {
		logger.Info().Str("exePath", exePath).Msg("executable path, will use this path to find the config file")
		mcpConfig.Command = exePath
	}
	cm := client.NewManager(logger, mcpConfig)
	if install {
		logger.Info().Msg("Start to add MCP Server configuration into MCP Clients.")
		cm.SetupConfig()
		logger.Info().Msg("Add MCP Server configuration into MCP Clients successfully.")
		return nil
	}
	logger.Info().Msg("Start to list MCP Clients")
	cm.ListClient()
	logger.Info().Msg("List MCP Clients successfully.")
	return nil
}

func init() {
	clientCmd.PersistentFlags().BoolVar(&list, "list", false, "List the current installed MCP clients")
	clientCmd.PersistentFlags().BoolVarP(&install, "install", "i", false, "Add MoLing MCP Server configuration to the currently installed MCP clients on this computer. default is all")
	rootCmd.AddCommand(clientCmd)
}

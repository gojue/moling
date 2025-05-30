// Copyright 2025 CFC4N <cfc4n.cs@gmail.com>. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Repository: https://github.com/gojue/moling

// Package services Description: This file contains the implementation of the CommandServer interface for macOS and  Linux.
package command

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog"

	"github.com/gojue/moling/pkg/comm"
	"github.com/gojue/moling/pkg/config"
	"github.com/gojue/moling/pkg/services/abstract"
	"github.com/gojue/moling/pkg/utils"
)

var (
	// ErrCommandNotFound is returned when the command is not found.
	ErrCommandNotFound = fmt.Errorf("command not found")
	// ErrCommandNotAllowed is returned when the command is not allowed.
	ErrCommandNotAllowed = fmt.Errorf("command not allowed")
)

const (
	CommandServerName comm.MoLingServerType = "Command"
)

// CommandServer implements the Service interface and provides methods to execute named commands.
type CommandServer struct {
	abstract.MLService
	config    *CommandConfig
	osName    string
	osVersion string
}

// NewCommandServer creates a new CommandServer with the given allowed commands.
func NewCommandServer(ctx context.Context) (abstract.Service, error) {
	var err error
	cc := NewCommandConfig()
	gConf, ok := ctx.Value(comm.MoLingConfigKey).(*config.MoLingConfig)
	if !ok {
		return nil, fmt.Errorf("CommandServer: invalid config type")
	}

	lger, ok := ctx.Value(comm.MoLingLoggerKey).(zerolog.Logger)
	if !ok {
		return nil, fmt.Errorf("CommandServer: invalid logger type")
	}

	loggerNameHook := zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, msg string) {
		e.Str("Service", string(CommandServerName))
	})

	cs := &CommandServer{
		MLService: abstract.NewMLService(ctx, lger.Hook(loggerNameHook), gConf),
		config:    cc,
	}

	err = cs.InitResources()
	if err != nil {
		return nil, err
	}

	return cs, nil
}

func (cs *CommandServer) Init() error {
	var err error
	pe := abstract.PromptEntry{
		PromptVar: mcp.Prompt{
			Name:        "command_prompt",
			Description: "get command prompt",
			//Arguments:   make([]mcp.PromptArgument, 0),
		},
		HandlerFunc: cs.handlePrompt,
	}
	cs.AddPrompt(pe)
	cs.AddTool(mcp.NewTool(
		"execute_command",
		mcp.WithDescription("Execute a named command.Only support command execution on macOS and will strictly follow safety guidelines, ensuring that commands are safe and secure"),
		mcp.WithString("command",
			mcp.Description("The command to execute"),
			mcp.Required(),
		),
	), cs.handleExecuteCommand)
	return err
}

func (cs *CommandServer) handlePrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "",
		Messages: []mcp.PromptMessage{
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf(cs.config.prompt, cs.MlConfig().SystemInfo),
				},
			},
		},
	}, nil
}

// handleExecuteCommand handles the execution of a named command.
func (cs *CommandServer) handleExecuteCommand(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	command, ok := args["command"].(string)
	if !ok {
		return mcp.NewToolResultError(fmt.Errorf("command must be a string").Error()), nil
	}

	// Check if the command is allowed
	if !cs.isAllowedCommand(command) {
		cs.Logger.Err(ErrCommandNotAllowed).Str("command", command).Msgf("If you want to allow this command, add it to %s", filepath.Join(cs.MlConfig().BasePath, "config", cs.MlConfig().ConfigFile))
		return mcp.NewToolResultError(fmt.Sprintf("Error: Command '%s' is not allowed", command)), nil
	}

	// Execute the command
	output, err := ExecCommand(command)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error executing command: %v", err)), nil
	}

	return mcp.NewToolResultText(output), nil
}

// isAllowedCommand checks if the command is allowed based on the configuration.
func (cs *CommandServer) isAllowedCommand(command string) bool {
	// 检查命令是否在允许的列表中
	for _, allowed := range cs.config.allowedCommands {
		if strings.HasPrefix(command, allowed) {
			return true
		}
	}

	// 如果命令包含管道符，进一步检查每个子命令
	if strings.Contains(command, "|") {
		parts := strings.Split(command, "|")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if !cs.isAllowedCommand(part) {
				return false
			}
		}
		return true
	}

	if strings.Contains(command, "&") {
		parts := strings.Split(command, "&")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if !cs.isAllowedCommand(part) {
				return false
			}
		}
		return true
	}

	return false
}

// Config returns the configuration of the service as a string.
func (cs *CommandServer) Config() string {
	cs.config.AllowedCommand = strings.Join(cs.config.allowedCommands, ",")
	cfg, err := json.Marshal(cs.config)
	if err != nil {
		cs.Logger.Err(err).Msg("failed to marshal config")
		return "{}"
	}
	cs.Logger.Debug().Str("config", string(cfg)).Msg("CommandServer config")
	return string(cfg)
}

func (cs *CommandServer) Name() comm.MoLingServerType {
	return CommandServerName
}

func (cs *CommandServer) Close() error {
	// Cancel the context to stop the browser
	cs.Logger.Debug().Msg("CommandServer closed")
	return nil
}

// LoadConfig loads the configuration from a JSON object.
func (cs *CommandServer) LoadConfig(jsonData map[string]interface{}) error {
	err := utils.MergeJSONToStruct(cs.config, jsonData)
	if err != nil {
		return err
	}
	// split the AllowedCommand string into a slice
	cs.config.allowedCommands = strings.Split(cs.config.AllowedCommand, ",")
	return cs.config.Check()
}

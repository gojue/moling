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
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog"
	"path/filepath"
	"strings"
)

var (
	// ErrCommandNotFound is returned when the command is not found.
	ErrCommandNotFound = fmt.Errorf("command not found")
	// ErrCommandNotAllowed is returned when the command is not allowed.
	ErrCommandNotAllowed = fmt.Errorf("command not allowed")
)

// CommandServer implements the Service interface and provides methods to execute named commands.
type CommandServer struct {
	MLService
	config       *CommandConfig
	globalConfig *MoLingConfig
	osName       string
	osVersion    string
}

// NewCommandServer creates a new CommandServer with the given allowed commands.
func NewCommandServer(ctx context.Context, args []string) (Service, error) {
	var err error
	cc := NewCommandConfig(args)
	gConf, ok := ctx.Value(MoLingConfigKey).(*MoLingConfig)
	if !ok {
		return nil, fmt.Errorf("CommandServer: invalid config type")
	}

	lger, ok := ctx.Value(MoLingLoggerKey).(zerolog.Logger)
	if !ok {
		return nil, fmt.Errorf("CommandServer: invalid logger type")
	}

	loggerNameHook := zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, msg string) {
		e.Str("Service", "CommandServer")
	})

	cs := &CommandServer{
		MLService: MLService{
			ctx:    ctx,
			logger: lger.Hook(loggerNameHook),
		},
		config:       cc,
		globalConfig: gConf,
	}

	err = cs.init()
	if err != nil {
		return nil, err
	}

	pe := PromptEntry{
		prompt: mcp.Prompt{
			Name:        "command_prompt",
			Description: fmt.Sprintf("You are a command-line tool assistant, using macOS 15.3.3 system commands to help users troubleshoot network issues, system performance, file searching, and statistics, among other things."),
			//Arguments:   make([]mcp.PromptArgument, 0),
		},
		phf: cs.handlePrompt,
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

	return cs, nil
}

func (cs *CommandServer) handlePrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: fmt.Sprintf(""),
		Messages: []mcp.PromptMessage{
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("You are a command-line tool assistant, using %s system commands to help users troubleshoot network issues, system performance, among other things.", cs.globalConfig.SystemInfo),
				},
			},
		},
	}, nil
}

// handleExecuteCommand handles the execution of a named command.
func (cs *CommandServer) handleExecuteCommand(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	command, ok := request.Params.Arguments["command"].(string)
	if !ok {
		return nil, fmt.Errorf("command must be a string")
	}

	// Check if the command is allowed
	if !cs.isAllowedCommand(command) {
		cs.logger.Err(ErrCommandNotAllowed).Str("command", command).Msgf("If you want to allow this command, add it to %s", filepath.Join(cs.globalConfig.BasePath, "config", cs.globalConfig.ConfigFile))
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error: Command '%s' is not allowed", command),
				},
			},
			IsError: true,
		}, nil
	}

	// Execute the command
	output, err := ExecCommand(command)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error executing command: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: output,
			},
		},
	}, nil
}

// isAllowedCommand checks if the command is allowed based on the configuration.
func (cs *CommandServer) isAllowedCommand(command string) bool {

	// 检查命令是否在允许的列表中
	for _, allowed := range cs.config.AllowedCommands {
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

	return false
}

// Config returns the configuration of the service as a string.
func (cs *CommandServer) Config() string {
	cfg, err := json.Marshal(cs.config)
	if err != nil {
		cs.logger.Err(err).Msg("failed to marshal config")
		return "{}"
	}
	cs.logger.Debug().Str("config", string(cfg)).Msg("CommandServer config")
	return string(cfg)
}

func (cs *CommandServer) Name() string {
	return "CommandServer"
}

func (bs *CommandServer) Close() error {
	// Cancel the context to stop the browser
	return nil
}

func init() {
	RegisterServ(NewCommandServer)
}

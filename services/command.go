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

const (
	CommandServerName = "CommandServer"
)

// CommandServer implements the Service interface and provides methods to execute named commands.
type CommandServer struct {
	MLService
	config    *CommandConfig
	osName    string
	osVersion string
}

// NewCommandServer creates a new CommandServer with the given allowed commands.
func NewCommandServer(ctx context.Context) (Service, error) {
	var err error
	cc := NewCommandConfig()
	gConf, ok := ctx.Value(MoLingConfigKey).(*MoLingConfig)
	if !ok {
		return nil, fmt.Errorf("CommandServer: invalid config type")
	}

	lger, ok := ctx.Value(MoLingLoggerKey).(zerolog.Logger)
	if !ok {
		return nil, fmt.Errorf("CommandServer: invalid logger type")
	}

	loggerNameHook := zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, msg string) {
		e.Str("Service", CommandServerName)
	})

	cs := &CommandServer{
		MLService: NewMLService(ctx, lger.Hook(loggerNameHook), gConf),
		config:    cc,
	}

	err = cs.init()
	if err != nil {
		return nil, err
	}

	return cs, nil
}

func (cs *CommandServer) Init() error {
	var err error
	pe := PromptEntry{
		prompt: mcp.Prompt{
			Name:        "command_prompt",
			Description: fmt.Sprintf("You are a command-line tool assistant, using %s system commands to help users troubleshoot network issues, system performance, file searching, and statistics, among other things.", cs.MlConfig().SystemInfo),
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
	return err
}

func (cs *CommandServer) handlePrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: fmt.Sprintf(""),
		Messages: []mcp.PromptMessage{
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf(""),
				},
			},
		},
	}, nil
}

// handleExecuteCommand handles the execution of a named command.
func (cs *CommandServer) handleExecuteCommand(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	command, ok := request.Params.Arguments["command"].(string)
	if !ok {
		return cs.CallToolResultErr(fmt.Errorf("command must be a string").Error()), nil
	}

	// Check if the command is allowed
	if !cs.isAllowedCommand(command) {
		cs.logger.Err(ErrCommandNotAllowed).Str("command", command).Msgf("If you want to allow this command, add it to %s", filepath.Join(cs.MlConfig().BasePath, "config", cs.MlConfig().ConfigFile))
		return cs.CallToolResultErr(fmt.Sprintf("Error: Command '%s' is not allowed", command)), nil
	}

	// Execute the command
	output, err := ExecCommand(command)
	if err != nil {
		return cs.CallToolResultErr(fmt.Sprintf("Error executing command: %v", err)), nil
	}

	return cs.CallToolResult(output), nil
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
		cs.logger.Err(err).Msg("failed to marshal config")
		return "{}"
	}
	cs.logger.Debug().Str("config", string(cfg)).Msg("CommandServer config")
	return string(cfg)
}

func (cs *CommandServer) Name() string {
	return CommandServerName
}

func (cs *CommandServer) Close() error {
	// Cancel the context to stop the browser
	cs.logger.Debug().Msg("CommandServer closed")
	return nil
}

// LoadConfig loads the configuration from a JSON object.
func (cs *CommandServer) LoadConfig(jsonData map[string]interface{}) error {
	err := mergeJSONToStruct(cs.config, jsonData)
	if err != nil {
		return err
	}
	// split the AllowedCommand string into a slice
	cs.config.allowedCommands = strings.Split(cs.config.AllowedCommand, ",")
	return cs.config.Check()
}

func init() {
	RegisterServ(CommandServerName, NewCommandServer)
}

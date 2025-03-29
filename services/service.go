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

package services

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"
	"sync"
)

type contextKey string

// MoLingConfigKey is a context key for storing the version of MoLing
const (
	MoLingConfigKey contextKey = "moling_config"
	MoLingLoggerKey contextKey = "moling_logger"
)

// Service defines the interface for a service with various handlers and tools.
type Service interface {
	Ctx() context.Context
	// Resources returns a map of resources and their corresponding handler functions.
	Resources() map[mcp.Resource]server.ResourceHandlerFunc
	// ResourceTemplates returns a map of resource templates and their corresponding handler functions.
	ResourceTemplates() map[mcp.ResourceTemplate]server.ResourceTemplateHandlerFunc
	// Prompts returns a map of prompts and their corresponding handler functions.
	Prompts() []PromptEntry
	// Tools returns a slice of server tools.
	Tools() []server.ServerTool
	// NotificationHandlers returns a map of notification handlers.
	NotificationHandlers() map[string]server.NotificationHandlerFunc

	// Config returns the configuration of the service as a string.
	Config() string

	MlConfig() *MoLingConfig

	// Name returns the name of the service.
	Name() string

	// Close closes the service and releases any resources it holds.
	Close() error
}

type PromptEntry struct {
	prompt mcp.Prompt
	phf    server.PromptHandlerFunc
}

func (pe *PromptEntry) Prompt() mcp.Prompt {
	return pe.prompt
}

func (pe *PromptEntry) Handler() server.PromptHandlerFunc {
	return pe.phf

}

// MLService implements the Service interface and provides methods to manage resources, templates, prompts, tools, and notification handlers.
type MLService struct {
	ctx                  context.Context
	lock                 *sync.Mutex
	resources            map[mcp.Resource]server.ResourceHandlerFunc
	resourcesTemplates   map[mcp.ResourceTemplate]server.ResourceTemplateHandlerFunc
	prompts              []PromptEntry
	tools                []server.ServerTool
	notificationHandlers map[string]server.NotificationHandlerFunc
	logger               zerolog.Logger // The logger for the service
	mlConfig             *MoLingConfig  // The configuration for the service
}

// init initializes the MLService with empty maps and a mutex.
func (mls *MLService) init() error {
	mls.lock = &sync.Mutex{}
	mls.resources = make(map[mcp.Resource]server.ResourceHandlerFunc)
	mls.resourcesTemplates = make(map[mcp.ResourceTemplate]server.ResourceTemplateHandlerFunc)
	mls.prompts = make([]PromptEntry, 0)
	mls.notificationHandlers = make(map[string]server.NotificationHandlerFunc)
	mls.tools = []server.ServerTool{}
	return nil
}

// Ctx returns the context of the MLService.
func (mls *MLService) Ctx() context.Context {
	return mls.ctx
}

// AddResource adds a resource and its handler function to the service.
func (mls *MLService) AddResource(rs mcp.Resource, hr server.ResourceHandlerFunc) {
	mls.lock.Lock()
	defer mls.lock.Unlock()
	mls.resources[rs] = hr
}

// AddResourceTemplate adds a resource template and its handler function to the service.
func (mls *MLService) AddResourceTemplate(rt mcp.ResourceTemplate, hr server.ResourceTemplateHandlerFunc) {
	mls.lock.Lock()
	defer mls.lock.Unlock()
	mls.resourcesTemplates[rt] = hr
}

// AddPrompt adds a prompt and its handler function to the service.
func (mls *MLService) AddPrompt(pe PromptEntry) {
	mls.lock.Lock()
	defer mls.lock.Unlock()
	mls.prompts = append(mls.prompts, pe)
}

// AddTool adds a tool and its handler function to the service.
func (mls *MLService) AddTool(tool mcp.Tool, handler server.ToolHandlerFunc) {
	mls.lock.Lock()
	defer mls.lock.Unlock()
	mls.tools = append(mls.tools, server.ServerTool{Tool: tool, Handler: handler})
}

// AddNotificationHandler adds a notification handler to the service.
func (mls *MLService) AddNotificationHandler(name string, handler server.NotificationHandlerFunc) {
	mls.lock.Lock()
	defer mls.lock.Unlock()
	mls.notificationHandlers[name] = handler
}

// Resources returns the map of resources and their handler functions.
func (mls *MLService) Resources() map[mcp.Resource]server.ResourceHandlerFunc {
	mls.lock.Lock()
	defer mls.lock.Unlock()
	return mls.resources
}

// ResourceTemplates returns the map of resource templates and their handler functions.
func (mls *MLService) ResourceTemplates() map[mcp.ResourceTemplate]server.ResourceTemplateHandlerFunc {
	mls.lock.Lock()
	defer mls.lock.Unlock()
	return mls.resourcesTemplates
}

// Prompts returns the map of prompts and their handler functions.
func (mls *MLService) Prompts() []PromptEntry {
	mls.lock.Lock()
	defer mls.lock.Unlock()
	return mls.prompts
}

// Tools returns the slice of server tools.
func (mls *MLService) Tools() []server.ServerTool {
	mls.lock.Lock()
	defer mls.lock.Unlock()
	return mls.tools
}

// NotificationHandlers returns the map of notification handlers.
func (mls *MLService) NotificationHandlers() map[string]server.NotificationHandlerFunc {
	mls.lock.Lock()
	defer mls.lock.Unlock()
	return mls.notificationHandlers
}

// CallToolResult return a CallToolResult with the given message and success status.
func (mls *MLService) CallToolResult(msg string) *mcp.CallToolResult {
	return mls.callToolResult(false, msg)
}

// CallToolResultErr return a CallToolResult with the given message and error status.
func (mls *MLService) CallToolResultErr(msg string) *mcp.CallToolResult {
	return mls.callToolResult(true, msg)
}

// callToolResult returns a CallToolResult with the given message and error status. Not allowed to be called directly.
func (mls *MLService) callToolResult(isError bool, msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: msg,
			},
		},
		IsError: isError,
	}
}

func (mls *MLService) MlConfig() *MoLingConfig {
	return mls.mlConfig
}

// Config returns the configuration of the service as a string.
func (mls *MLService) Config() string {
	panic("not implemented yet") // TODO: Implement
}

func (mls *MLService) Name() string {
	panic("not implemented yet") // TODO: Implement
}

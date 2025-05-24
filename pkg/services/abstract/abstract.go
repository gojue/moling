// Copyright 2025 CFC4N <cfc4n.cs@gmail.com>. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Repository: https://github.com/gojue/moling

package abstract

import (
	"context"

	"github.com/gojue/moling/pkg/comm"
	"github.com/gojue/moling/pkg/config"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ServiceFactory func(ctx context.Context) (Service, error)

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
	// LoadConfig loads the configuration for the service from a map.
	LoadConfig(jsonData map[string]interface{}) error

	// Init initializes the service with the given context and configuration.
	Init() error

	MlConfig() *config.MoLingConfig

	// Name returns the name of the service.
	Name() comm.MoLingServerType

	// Close closes the service and releases any resources it holds.
	Close() error
}

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

package abstract

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestMLService_AddResource(t *testing.T) {
	service := &MLService{}
	err := service.InitResources()
	if err != nil {
		t.Fatalf("Failed to initialize MLService: %v", err)
	}
	resource := mcp.Resource{Name: "testResource"}
	handler := func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				Text:     "text",
				URI:      "uri",
				MIMEType: "text/plain",
			},
		}, nil
	}

	service.AddResource(resource, handler)

	if len(service.resources) != 1 {
		t.Errorf("Expected 1 resource, got %d", len(service.resources))
	}
	if service.resources[resource] == nil {
		t.Errorf("Handler for resource not found")
	}
}

func TestMLService_AddResourceTemplate(t *testing.T) {
	service := &MLService{}
	err := service.InitResources()
	if err != nil {
		t.Fatalf("Failed to initialize MLService: %v", err)
	}
	template := mcp.ResourceTemplate{Name: "testTemplate"}
	handler := func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				Text:     "text",
				URI:      "uri",
				MIMEType: "text/plain",
			},
		}, nil
	}

	service.AddResourceTemplate(template, handler)

	if len(service.resourcesTemplates) != 1 {
		t.Errorf("Expected 1 resource template, got %d", len(service.resourcesTemplates))
	}
	if service.resourcesTemplates[template] == nil {
		t.Errorf("Handler for resource template not found")
	}
}

func TestMLService_AddPrompt(t *testing.T) {
	service := &MLService{}
	err := service.InitResources()
	if err != nil {
		t.Fatalf("Failed to initialize MLService: %v", err)
	}
	prompt := "testPrompt"
	handler := func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		pms := make([]mcp.PromptMessage, 0)
		pms = append(pms, mcp.PromptMessage{
			Role: mcp.RoleUser,
			Content: mcp.TextContent{
				Type: "text",
				Text: "Prompt response",
			},
		})
		return &mcp.GetPromptResult{
			Description: "prompt description",
			Messages:    pms,
		}, nil
	}
	pe := PromptEntry{
		PromptVar:   mcp.Prompt{Name: "testPrompt"},
		HandlerFunc: handler,
	}
	service.AddPrompt(pe)

	if len(service.prompts) != 1 {
		t.Errorf("Expected 1 prompt, got %d", len(service.prompts))
	}
	for _, p := range service.prompts {
		if p.PromptVar.Name != prompt {
			t.Errorf("Expected prompt name %s, got %s", prompt, p.PromptVar.Name)
		}
	}
}

func TestMLService_AddTool(t *testing.T) {
	service := &MLService{}
	err := service.InitResources()
	if err != nil {
		t.Fatalf("Failed to initialize MLService: %v", err)
	}
	tool := mcp.Tool{Name: "testTool"}
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Prompt response",
				},
			},
		}, nil
	}

	service.AddTool(tool, handler)

	if len(service.tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(service.tools))
	}

	// After
	if service.tools[0].Tool.Name != tool.Name {
		t.Errorf("Tool not added correctly")
	}
}

func TestMLService_AddNotificationHandler(t *testing.T) {
	service := &MLService{}
	err := service.InitResources()
	if err != nil {
		t.Fatalf("Failed to initialize MLService: %v", err)
	}
	name := "testHandler"
	handler := func(ctx context.Context, n mcp.JSONRPCNotification) {
		t.Logf("Received notification: %s", n.Method)
	}

	service.AddNotificationHandler(name, handler)

	if len(service.notificationHandlers) != 1 {
		t.Errorf("Expected 1 notification handler, got %d", len(service.notificationHandlers))
	}
	if service.notificationHandlers[name] == nil {
		t.Errorf("Handler for notification not found")
	}
}

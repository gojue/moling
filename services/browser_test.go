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
	"github.com/rs/zerolog"
	"os"
	"path/filepath"
	"testing"
)

// initTestEnv initializes the test environment by creating a temporary log file and setting up the logger.
func initTestEnv() (zerolog.Logger, context.Context, error) {
	logFile := filepath.Join(os.TempDir(), "moling.log")
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	var logger zerolog.Logger
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return zerolog.Logger{}, nil, err
	}
	logger = zerolog.New(f).With().Timestamp().Logger()
	mlConfig := &MoLingConfig{
		ConfigFile: filepath.Join("config", "test_config.json"),
		BasePath:   os.TempDir(),
	}
	ctx := context.WithValue(context.Background(), MoLingConfigKey, mlConfig)
	ctx = context.WithValue(ctx, MoLingLoggerKey, logger)
	return logger, ctx, nil
}

func TestBrowserServer(t *testing.T) {
	//
	//cfg := &BrowserConfig{
	//	Headless:        true,
	//	Timeout:         30,
	//	UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
	//	DefaultLanguage: "en-US",
	//	URLTimeout:      10,
	//	CSSTimeout:      10,
	//}
	logger, ctx, err := initTestEnv()
	if err != nil {
		t.Fatalf("Failed to initialize test environment: %v", err)
	}
	logger.Info().Msg("TestBrowserServer")

	_, err = NewBrowserServer(ctx)
	if err != nil {
		t.Fatalf("Failed to create BrowserServer: %v", err)
	}

	/*
		t.Run("TestNavigate", func(t *testing.T) {
			request := mcp.CallToolRequest{
				Request: mcp.Request{
					Method: "tools/call",
				},
			}
			request.Params.Arguments = map[string]interface{}{
				"url": "https://www.baidu.com",
			}
			result, err := bs.handleNavigate(ctx, request)
			if err != nil {
				t.Fatalf("handleNavigate failed: %v", err)
			}
			if result.Content[0].(mcp.TextContent).Text != "Navigated to https://www.example.com" {
				t.Errorf("Unexpected result: %v", result.Content[0].(mcp.TextContent).Text)
			}
		})
	*/
	//
	//t.Run("TestScreenshot", func(t *testing.T) {
	//	request := mcp.CallToolRequest{
	//		Params: mcp.ToolParams{
	//			Arguments: map[string]interface{}{
	//				"name": "test_screenshot",
	//			},
	//		},
	//	}
	//	_, err := bs.handleScreenshot(ctx, request)
	//	if err != nil {
	//		t.Fatalf("handleScreenshot failed: %v", err)
	//	}
	//})
	//
	//t.Run("TestClick", func(t *testing.T) {
	//	request := mcp.CallToolRequest{
	//		Params: mcp.ToolParams{
	//			Arguments: map[string]interface{}{
	//				"selector": "body",
	//			},
	//		},
	//	}
	//	_, err := bs.handleClick(ctx, request)
	//	if err != nil {
	//		t.Fatalf("handleClick failed: %v", err)
	//	}
	//})
	//
	//t.Run("TestFill", func(t *testing.T) {
	//	request := mcp.CallToolRequest{
	//		Params: mcp.ToolParams{
	//			Arguments: map[string]interface{}{
	//				"selector": "input[name='q']",
	//				"value":    "test",
	//			},
	//		},
	//	}
	//	_, err := bs.handleFill(ctx, request)
	//	if err != nil {
	//		t.Fatalf("handleFill failed: %v", err)
	//	}
	//})
	//
	//t.Run("TestSelect", func(t *testing.T) {
	//	request := mcp.CallToolRequest{
	//		Params: mcp.ToolParams{
	//			Arguments: map[string]interface{}{
	//				"selector": "select[name='dropdown']",
	//				"value":    "option1",
	//			},
	//		},
	//	}
	//	_, err := bs.handleSelect(ctx, request)
	//	if err != nil {
	//		t.Fatalf("handleSelect failed: %v", err)
	//	}
	//})
	//
	//t.Run("TestHover", func(t *testing.T) {
	//	request := mcp.CallToolRequest{
	//		Params: mcp.ToolParams{
	//			Arguments: map[string]interface{}{
	//				"selector": "body",
	//			},
	//		},
	//	}
	//	_, err := bs.handleHover(ctx, request)
	//	if err != nil {
	//		t.Fatalf("handleHover failed: %v", err)
	//	}
	//})
	//
	//t.Run("TestEvaluate", func(t *testing.T) {
	//	request := mcp.CallToolRequest{
	//		Params: mcp.ToolParams{
	//			Arguments: map[string]interface{}{
	//				"script": "document.title",
	//			},
	//		},
	//	}
	//	result, err := bs.handleEvaluate(ctx, request)
	//	if err != nil {
	//		t.Fatalf("handleEvaluate failed: %v", err)
	//	}
	//	if result.Content[0].(mcp.TextContent).Text == "" {
	//		t.Errorf("Unexpected result: %v", result.Content[0].(mcp.TextContent).Text)
	//	}
	//})
}

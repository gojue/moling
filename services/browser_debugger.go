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

// Package services provides a set of services for the MoLing application.
package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"github.com/mark3labs/mcp-go/mcp"
)

// handleDebugEnable handles the enabling and disabling of debugging in the browser.
func (bs *BrowserServer) handleDebugEnable(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	enabled, ok := request.Params.Arguments["enabled"].(bool)
	if !ok {
		return mcp.NewToolResultError("enabled must be a boolean"), nil
	}

	var err error
	rctx, cancel := context.WithCancel(bs.ctx)
	defer cancel()

	if enabled {
		err = chromedp.Run(rctx, chromedp.ActionFunc(func(ctx context.Context) error {
			t := chromedp.FromContext(ctx).Target
			// 使用Execute方法执行AttachToTarget命令
			params := target.AttachToTarget(t.TargetID).WithFlatten(true)
			return t.Execute(ctx, "Target.attachToTarget", params, nil)
		}))
	} else {
		err = chromedp.Run(rctx, chromedp.ActionFunc(func(ctx context.Context) error {
			t := chromedp.FromContext(ctx).Target
			// 使用Execute方法执行DetachFromTarget命令
			params := target.DetachFromTarget().WithSessionID(t.SessionID)
			return t.Execute(ctx, "Target.detachFromTarget", params, nil)
		}))
	}

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to %s debugging: %v",
			map[bool]string{true: "enable", false: "disable"}[enabled], err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("Debugging %s",
		map[bool]string{true: "enabled", false: "disabled"}[enabled])), nil
}

// handleSetBreakpoint handles setting a breakpoint in the browser.
func (bs *BrowserServer) handleSetBreakpoint(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, ok := request.Params.Arguments["url"].(string)
	if !ok {
		return mcp.NewToolResultError("url must be a string"), nil
	}

	line, ok := request.Params.Arguments["line"].(float64)
	if !ok {
		return mcp.NewToolResultError("line must be a number"), nil
	}

	column, _ := request.Params.Arguments["column"].(float64)
	condition, _ := request.Params.Arguments["condition"].(string)

	var breakpointID string
	rctx, cancel := context.WithCancel(bs.ctx)
	defer cancel()
	err := chromedp.Run(rctx, chromedp.ActionFunc(func(ctx context.Context) error {
		t := chromedp.FromContext(ctx).Target
		params := map[string]interface{}{
			"url":       url,
			"line":      int(line),
			"column":    int(column),
			"condition": condition,
		}

		var result map[string]interface{}
		// 使用Execute方法执行Debugger.setBreakpoint命令
		if err := t.Execute(ctx, "Debugger.setBreakpoint", params, &result); err != nil {
			return err
		}

		breakpointID, ok = result["breakpointId"].(string)
		if !ok {
			breakpointID = ""
			return fmt.Errorf("failed to get breakpoint ID")
		}
		return nil
	}))

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to set breakpoint: %v", err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("Breakpoint set with ID: %s", breakpointID)), nil
}

// handleRemoveBreakpoint handles removing a breakpoint in the browser.
func (bs *BrowserServer) handleRemoveBreakpoint(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	breakpointID, ok := request.Params.Arguments["breakpointId"].(string)
	if !ok {
		return mcp.NewToolResultError("breakpointId must be a string"), nil
	}
	rctx, cancel := context.WithCancel(bs.ctx)
	defer cancel()
	err := chromedp.Run(rctx, chromedp.ActionFunc(func(ctx context.Context) error {
		t := chromedp.FromContext(ctx).Target
		// 使用Execute方法执行Debugger.removeBreakpoint命令
		return t.Execute(ctx, "Debugger.removeBreakpoint", map[string]interface{}{
			"breakpointId": breakpointID,
		}, nil)
	}))

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to remove breakpoint: %v", err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("Breakpoint %s removed", breakpointID)), nil
}

// handlePause handles pausing the JavaScript execution in the browser.
func (bs *BrowserServer) handlePause(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	rctx, cancel := context.WithCancel(bs.ctx)
	defer cancel()
	err := chromedp.Run(rctx, chromedp.ActionFunc(func(ctx context.Context) error {
		t := chromedp.FromContext(ctx).Target
		// 使用Execute方法执行Debugger.pause命令
		return t.Execute(ctx, "Debugger.pause", nil, nil)
	}))

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to pause execution: %v", err)), nil
	}
	return mcp.NewToolResultText("JavaScript execution paused"), nil
}

// handleResume handles resuming the JavaScript execution in the browser.
func (bs *BrowserServer) handleResume(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	rctx, cancel := context.WithCancel(bs.ctx)
	defer cancel()
	err := chromedp.Run(rctx, chromedp.ActionFunc(func(ctx context.Context) error {
		t := chromedp.FromContext(ctx).Target
		// 使用Execute方法执行Debugger.resume命令
		return t.Execute(ctx, "Debugger.resume", nil, nil)
	}))

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to resume execution: %v", err)), nil
	}
	return mcp.NewToolResultText("JavaScript execution resumed"), nil
}

// handleStepOver handles stepping over the next line of JavaScript code in the browser.
func (bs *BrowserServer) handleGetCallstack(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var callstack interface{}
	rctx, cancel := context.WithCancel(bs.ctx)
	defer cancel()
	err := chromedp.Run(rctx, chromedp.ActionFunc(func(ctx context.Context) error {
		t := chromedp.FromContext(ctx).Target
		// 使用Execute方法执行Debugger.getStackTrace命令
		return t.Execute(ctx, "Debugger.getStackTrace", nil, &callstack)
	}))

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get call stack: %v", err)), nil
	}

	callstackJSON, err := json.Marshal(callstack)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal call stack: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Current call stack: %s", string(callstackJSON))), nil
}

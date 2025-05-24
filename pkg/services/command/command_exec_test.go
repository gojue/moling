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

package command

import (
	"context"
	"errors"
	"os/exec"
	"reflect"
	"testing"
	"time"

	"github.com/gojue/moling/pkg/comm"
)

// MockCommandServer is a mock implementation of CommandServer for testing purposes.
type MockCommandServer struct {
	CommandServer
}

// TestExecuteCommand tests the ExecCommand function.
func TestExecuteCommand(t *testing.T) {
	execCmd := "echo 'Hello, World!'"
	// Test a simple command
	output, err := ExecCommand(execCmd)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expectedOutput := "Hello, World!\n"
	if output != expectedOutput {
		t.Errorf("Expected output %q, got %q", expectedOutput, output)
	}
	t.Logf("Command output: %s", output)
	// Test a command with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	execCmd = "curl ifconfig.me|grep Time"
	output, err = ExecCommand(execCmd)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	t.Logf("Command output: %s", output)
	cmd := exec.CommandContext(ctx, "sleep", "1")
	err = cmd.Run()
	if err == nil {
		t.Fatalf("Expected timeout error, got nil")
	}
	if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Errorf("Expected context deadline exceeded error, got %v", ctx.Err())
	}
}

func TestAllowCmd(t *testing.T) {
	// Test with a command that is allowed
	_, ctx, err := comm.InitTestEnv()
	if err != nil {
		t.Fatalf("Failed to initialize test environment: %v", err)
	}

	cs, err := NewCommandServer(ctx)
	if err != nil {
		t.Fatalf("Failed to create CommandServer: %v", err)
	}

	cc := StructToMap(NewCommandConfig())
	t.Logf("CommandConfig: %v", cc)
	err = cs.LoadConfig(cc)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	cmd := "cd /var/logs/notfound && git log --since=\"today\" --pretty=format:\"%h - %an, %ar : %s\""
	cs1 := cs.(*CommandServer)
	if !cs1.isAllowedCommand(cmd) {
		t.Errorf("Command '%s' is not allowed", cmd)
	}
	t.Log("Command is allowed:", cmd)
}

// 将 struct 转换为 map
func StructToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		// 跳过未导出的字段
		if field.PkgPath != "" {
			continue
		}
		// 获取字段的 json tag（如果存在）
		key := field.Name
		if tag := field.Tag.Get("json"); tag != "" {
			key = tag
		}
		result[key] = value.Interface()
	}
	return result
}

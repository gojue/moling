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
	"errors"
	"os/exec"
	"testing"
	"time"
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

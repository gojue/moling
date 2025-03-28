//go:build windows

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

// Package services Description: This file contains the implementation of the CommandServer interface for Windows.
package services

import (
	"os/exec"
)

// ExecCommand executes a command and returns its output.
func ExecCommand(command string) (string, error) {
	var cmd *exec.Cmd
	cmd = exec.Command("cmd", "/C", command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

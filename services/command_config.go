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
	"fmt"
)

// CommandConfig represents the configuration for allowed commands.
type CommandConfig struct {
	AllowedCommand  string `json:"allowed_command"` // AllowedCommand is a list of allowed command. split by comma. e.g. ls,cat,echo
	allowedCommands []string
}

var (
	allowedCmdDefault = []string{
		"ls", "cat", "echo", "pwd", "head", "tail", "grep", "find", "stat", "df",
		"du", "free", "top", "ps", "uptime", "who", "w", "last", "uname", "hostname",
		"ifconfig", "netstat", "ping", "traceroute", "route", "ip", "ss", "lsof", "vmstat",
		"iostat", "mpstat", "sar", "uptime", "cut", "sort", "uniq", "wc", "awk", "sed",
		"diff", "cmp", "comm", "file", "basename", "dirname", "chmod", "chown", "curl",
		"nslookup", "dig", "host", "ssh", "scp", "sftp", "ftp", "wget", "tar", "gzip",
		"scutil", "networksetup",
	}
)

// NewCommandConfig creates a new CommandConfig with the given allowed commands.
func NewCommandConfig() *CommandConfig {
	return &CommandConfig{
		allowedCommands: allowedCmdDefault,
	}
}

// Check validates the allowed commands in the CommandConfig.
func (cc *CommandConfig) Check() error {
	var cnt int
	cnt = len(cc.allowedCommands)

	// Check if any command is empty
	for _, cmd := range cc.allowedCommands {
		if cmd == "" {
			cnt -= 1
		}
	}

	if cnt <= 0 {
		return fmt.Errorf("no allowed commands specified")
	}

	return nil
}

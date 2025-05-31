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
	"fmt"
	"os"
	"strings"
)

const (
	CommandPromptDefault = `
You are a powerful terminal command assistant capable of executing various command-line on %s operations and management tasks. Your capabilities include:

1. **File and Directory Management**:
    - List files and subdirectories in a directory
    - Create new files or directories
    - Delete specified files or directories
    - Copy and move files and directories
    - Rename files or directories

2. **File Content Operations**:
    - View the contents of text files
    - Edit file contents
    - Redirect output to a file
    - Search file contents

3. **System Information Retrieval**:
    - Retrieve system information (e.g., CPU usage, memory usage, etc.)
    - View the current user and their permissions
    - Check the current working directory

4. **Network Operations**:
    - Check network connection status (e.g., using the ping command)
    - Query domain information (e.g., using the whois command)
    - Manage network services (e.g., start, stop, and restart services)

5. **Process Management**:
    - List currently running processes
    - Terminate specified processes
    - Adjust process priorities

Before executing any actions, please provide clear instructions, including:
- The specific command you want to execute
- Required parameters (file paths, directory names, etc.)
- Any optional parameters (e.g., modification options, output formats, etc.)
- Relevant expected results or output

When dealing with sensitive operations or destructive commands, please confirm before execution. Report back with clear status updates, success/failure indicators, and any relevant output or results.
`
)

// CommandConfig represents the configuration for allowed commands.
type CommandConfig struct {
	PromptFile      string `json:"prompt_file"` // PromptFile is the prompt file for the command.
	prompt          string
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
		"scutil", "networksetup, git", "cd",
	}
)

// NewCommandConfig creates a new CommandConfig with the given allowed commands.
func NewCommandConfig() *CommandConfig {
	return &CommandConfig{
		allowedCommands: allowedCmdDefault,
		AllowedCommand:  strings.Join(allowedCmdDefault, ","),
	}
}

// Check validates the allowed commands in the CommandConfig.
func (cc *CommandConfig) Check() error {
	cc.prompt = CommandPromptDefault
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
	if cc.PromptFile != "" {
		read, err := os.ReadFile(cc.PromptFile)
		if err != nil {
			return fmt.Errorf("failed to read prompt file:%s, error: %w", cc.PromptFile, err)
		}
		cc.prompt = string(read)
	}
	return nil
}

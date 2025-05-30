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

package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// FileSystemPromptDefault is the default prompt for the file system.
	FileSystemPromptDefault = `
You are a powerful local filesystem management assistant capable of performing various file operations and management tasks. Your capabilities include:

1. **File Browsing**: Navigate to specified directories to load lists of files and folders.

2. **File Operations**:
   - Create new files or folders
   - Delete specified files or folders
   - Copy and move files and folders
   - Rename files or folders

3. **File Content Operations**:
   - Read the contents of text files and return them
   - Write text to specified files
   - Append content to existing files

4. **File Information Retrieval**:
   - Retrieve properties of files or folders (e.g., size, creation date, modification date)
   - Check if files or folders exist

5. **Search Functionality**:
   - Search for files in specified directories, supporting wildcard matching
   - Filter search results by file type or modification date

For all actions, please provide clear instructions, including:
- The specific action you want to perform
- Required parameters (directory paths, filenames, content, etc.)
- Any optional parameters (e.g., new filenames, search patterns, etc.)
- Relevant expected outcomes

You should confirm actions before execution when dealing with sensitive operations or destructive commands. Report back with clear status updates, success/failure indicators, and any relevant output or results.
`
)

var (
	allowedDirsDefault = os.TempDir()
)

// FileSystemConfig represents the configuration for the file system.
type FileSystemConfig struct {
	PromptFile  string `json:"prompt_file"` // PromptFile is the prompt file for the file system.
	prompt      string
	AllowedDir  string `json:"allowed_dir"` // AllowedDirs is a list of allowed directories. split by comma. e.g. /tmp,/var/tmp
	allowedDirs []string
	CachePath   string `json:"cache_path"` // CachePath is the root path for the file system.
}

// NewFileSystemConfig creates a new FileSystemConfig with the given allowed directories.
func NewFileSystemConfig(path string) *FileSystemConfig {
	paths := strings.Split(path, ",")
	path = ""
	if strings.TrimSpace(path) == "" {
		path = allowedDirsDefault
		paths = []string{allowedDirsDefault}
	}

	return &FileSystemConfig{
		AllowedDir:  path,
		CachePath:   path,
		allowedDirs: paths,
	}
}

// Check validates the allowed directories in the FileSystemConfig.
func (fc *FileSystemConfig) Check() error {
	fc.prompt = FileSystemPromptDefault
	normalized := make([]string, 0, len(fc.allowedDirs))
	for _, dir := range fc.allowedDirs {
		abs, err := filepath.Abs(strings.TrimSpace(dir))
		if err != nil {
			return fmt.Errorf("failed to resolve path %s: %w", dir, err)
		}
		info, err := os.Stat(abs)
		if err != nil {
			return fmt.Errorf("failed to access directory %s: %w", abs, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("path is not a directory: %s", abs)
		}

		normalized = append(normalized, filepath.Clean(abs)+string(filepath.Separator))
	}
	fc.allowedDirs = normalized

	if fc.PromptFile != "" {
		read, err := os.ReadFile(fc.PromptFile)
		if err != nil {
			return fmt.Errorf("failed to read prompt file:%s, error: %w", fc.PromptFile, err)
		}
		fc.prompt = string(read)
	}

	return nil
}

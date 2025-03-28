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
	"os"
	"path/filepath"
)

var (
	// TODO : read from config file
	allowedDirsDefault = []string{
		os.TempDir(),
	}
)

// FileSystemConfig represents the configuration for the file system.
type FileSystemConfig struct {
	AllowedDirs []string `json:"allowed_dirs"` // AllowedDirs is a list of allowed directories.
	CachePath   string   `json:"cache_path"`   // CachePath is the root path for the file system.
}

// NewFileSystemConfig creates a new FileSystemConfig with the given allowed directories.
func NewFileSystemConfig(path []string) *FileSystemConfig {
	if len(path) == 0 {
		path = allowedDirsDefault
	}

	return &FileSystemConfig{
		AllowedDirs: path,
		CachePath:   path[0],
	}
}

// Check validates the allowed directories in the FileSystemConfig.
func (fc *FileSystemConfig) Check() error {
	normalized := make([]string, 0, len(fc.AllowedDirs))
	for _, dir := range fc.AllowedDirs {
		abs, err := filepath.Abs(dir)
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
	fc.AllowedDirs = normalized
	return nil
}

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
// Source: https://github.com/mark3labs/mcp-filesystem-server

// Package services provides the implementation of the FileSystemServer, which allows access to files and directories on the local file system.
package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// MaxInlineSize Maximum size for inline content (5MB)
	MaxInlineSize = 1024 * 1024 * 5
	// MaxBase64Size Maximum size for base64 encoding (1MB)
	MaxBase64Size = 1024 * 1024 * 1
)
const (
	FilesystemServerName MoLingServerType = "FileSystem"
)

type FileInfo struct {
	Size        int64     `json:"size"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
	Accessed    time.Time `json:"accessed"`
	IsDirectory bool      `json:"isDirectory"`
	IsFile      bool      `json:"isFile"`
	Permissions string    `json:"permissions"`
}

type FilesystemServer struct {
	MLService
	config *FileSystemConfig
}

func NewFilesystemServer(ctx context.Context) (Service, error) {
	// Validate the config
	var err error
	globalConf := ctx.Value(MoLingConfigKey).(*MoLingConfig)
	userDataDir := filepath.Join(globalConf.BasePath, "data")

	fc := NewFileSystemConfig(userDataDir)

	lger, ok := ctx.Value(MoLingLoggerKey).(zerolog.Logger)
	if !ok {
		return nil, fmt.Errorf("FilesystemServer: invalid logger type")
	}

	loggerNameHook := zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, msg string) {
		e.Str("Service", string(FilesystemServerName))
	})

	fs := &FilesystemServer{
		MLService: NewMLService(ctx, lger.Hook(loggerNameHook), globalConf),
		config:    fc,
	}

	err = fs.init()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize filesystem server: %v", err)
	}

	return fs, nil
}

func (fs *FilesystemServer) Init() error {
	// Register resource handlers
	fs.AddResource(mcp.NewResource("file://", "File System",
		mcp.WithResourceDescription("Access to files and directories on the local file system"),
	), fs.handleReadResource)

	pe := PromptEntry{
		prompt: mcp.Prompt{
			Name:        "filesystem_prompt",
			Description: fmt.Sprintf("Get the relevant functions and prompts of the FileSystem MCP Server."),
		},
		phf: fs.handlePrompt,
	}
	fs.AddPrompt(pe)

	// Register tool handlers
	fs.AddTool(mcp.NewTool("read_file",
		mcp.WithDescription("Read the complete contents of a file from the file system."),
		mcp.WithString("path",
			mcp.Description("Relative path to the file to read"),
			mcp.Required(),
		),
	), fs.handleReadFile)

	fs.AddTool(mcp.NewTool(
		"write_file",
		mcp.WithDescription("Create a new file or overwrite an existing file with new content."),
		mcp.WithString("path",
			mcp.Description("Relative Path where to write the file"),
			mcp.Required(),
		),
		mcp.WithString("content",
			mcp.Description("Content to write to the file"),
			mcp.Required(),
		),
	), fs.handleWriteFile)

	fs.AddTool(mcp.NewTool(
		"list_directory",
		mcp.WithDescription("Get a detailed listing of all files and directories in a specified path."),
		mcp.WithString("path",
			mcp.Description("Relative Path of the directory to list"),
			mcp.Required(),
		),
	), fs.handleListDirectory)

	fs.AddTool(mcp.NewTool(
		"create_directory",
		mcp.WithDescription("Create a new directory or ensure a directory exists."),
		mcp.WithString("path",
			mcp.Description("Relative Path of the directory to create"),
			mcp.Required(),
		),
	), fs.handleCreateDirectory)

	fs.AddTool(mcp.NewTool(
		"move_file",
		mcp.WithDescription("Move or rename files and directories."),
		mcp.WithString("source",
			mcp.Description("Relative Source path of the file or directory"),
			mcp.Required(),
		),
		mcp.WithString("destination",
			mcp.Description("Relative Destination path"),
			mcp.Required(),
		),
	), fs.handleMoveFile)

	fs.AddTool(mcp.NewTool(
		"search_files",
		mcp.WithDescription("Recursively search for files and directories matching a pattern."),
		mcp.WithString("path",
			mcp.Description("Relative Starting path for the search"),
			mcp.Required(),
		),
		mcp.WithString("pattern",
			mcp.Description("Relative Search pattern to match against file names"),
			mcp.Required(),
		),
	), fs.handleSearchFiles)

	fs.AddTool(mcp.NewTool(
		"get_file_info",
		mcp.WithDescription("Retrieve detailed metadata about a file or directory."),
		mcp.WithString("path",
			mcp.Description("Relative Path to the file or directory"),
			mcp.Required(),
		),
	), fs.handleGetFileInfo)

	fs.AddTool(mcp.NewTool(
		"list_allowed_directories",
		mcp.WithDescription("Returns the list of directories that this server is allowed to access."),
	), fs.handleListAllowedDirectories)
	return nil
}

// handlePrompt handles the prompt request for the FilesystemServer
func (fs *FilesystemServer) handlePrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: fmt.Sprintf(""),
		Messages: []mcp.PromptMessage{
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf(`
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
`),
				},
			},
		},
	}, nil
}

// isPathInAllowedDirs checks if a path is within any of the allowed directories
func (fs *FilesystemServer) isPathInAllowedDirs(path string) bool {
	// Ensure path is absolute and clean
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// Add trailing separator to ensure we're checking a directory or a file within a directory
	// and not a prefix match (e.g., /tmp/foo should not match /tmp/foobar)
	if !strings.HasSuffix(absPath, string(filepath.Separator)) {
		// If it'fss a file, we need to check its directory
		if info, err := os.Stat(absPath); err == nil && !info.IsDir() {
			absPath = filepath.Dir(absPath) + string(filepath.Separator)
		} else {
			absPath = absPath + string(filepath.Separator)
		}
	}

	// Check if the path is within any of the allowed directories
	for _, dir := range fs.config.allowedDirs {
		if strings.HasPrefix(absPath, dir) {
			return true
		}
	}
	return false
}

func (fs *FilesystemServer) validatePath(requestedPath string) (string, error) {
	// Always convert to absolute path first
	var hasPrefix bool
	var firstDir string
	for _, dir := range fs.config.allowedDirs {
		if firstDir == "" {
			firstDir = dir
		}
		if strings.HasPrefix(requestedPath, dir) {
			hasPrefix = true
			break
		}
	}
	if !hasPrefix {
		requestedPath = filepath.Join(firstDir, requestedPath)
	}
	abs, err := filepath.Abs(requestedPath)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}

	// Check if path is within allowed directories
	if !fs.isPathInAllowedDirs(abs) {
		return "", fmt.Errorf("access denied - path outside allowed directories: %s", abs)
	}

	// Handle symlinks
	realPath, err := filepath.EvalSymlinks(abs)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
		// For new files, check parent directory
		parent := filepath.Dir(abs)
		realParent, err := filepath.EvalSymlinks(parent)
		if err != nil {
			return "", fmt.Errorf("parent directory does not exist: %s", parent)
		}

		if !fs.isPathInAllowedDirs(realParent) {
			return "", fmt.Errorf(
				"access denied - parent directory outside allowed directories",
			)
		}
		return abs, nil
	}

	// Check if the real path (after resolving symlinks) is still within allowed directories
	if !fs.isPathInAllowedDirs(realPath) {
		return "", fmt.Errorf(
			"access denied - symlink target outside allowed directories",
		)
	}

	return realPath, nil
}

func (fs *FilesystemServer) getFileStats(path string) (FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return FileInfo{}, err
	}

	return FileInfo{
		Size:        info.Size(),
		Created:     info.ModTime(), // Note: ModTime used as birth time isn't always available
		Modified:    info.ModTime(),
		Accessed:    info.ModTime(), // Note: Access time isn't always available
		IsDirectory: info.IsDir(),
		IsFile:      !info.IsDir(),
		Permissions: fmt.Sprintf("%o", info.Mode().Perm()),
	}, nil
}

func (fs *FilesystemServer) searchFiles(rootPath, pattern string) ([]string, error) {
	var results []string
	pattern = strings.ToLower(pattern)

	err := filepath.Walk(
		rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip errors and continue
			}

			// Try to validate path
			if _, err := fs.validatePath(path); err != nil {
				return nil // Skip invalid paths
			}

			if strings.Contains(strings.ToLower(info.Name()), pattern) {
				results = append(results, path)
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Resource handler
func (fs *FilesystemServer) handleReadResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	uri := request.Params.URI
	fs.logger.Debug().Str("uri", uri).Msg("handleReadResource")

	// Check if it'fss a file:// URI
	if !strings.HasPrefix(uri, "file://") {
		return nil, fmt.Errorf("unsupported URI scheme: %s", uri)
	}

	// Extract the path from the URI
	path := strings.TrimPrefix(uri, "file://")

	// Validate the path
	validPath, err := fs.validatePath(path)
	if err != nil {
		return nil, err
	}

	// Get file info
	fileInfo, err := os.Stat(validPath)
	if err != nil {
		return nil, err
	}

	// If it'fss a directory, return a listing
	if fileInfo.IsDir() {
		entries, err := os.ReadDir(validPath)
		if err != nil {
			return nil, err
		}

		var result strings.Builder
		result.WriteString(fmt.Sprintf("Directory listing for: %s\n\n", validPath))

		for _, entry := range entries {
			entryPath := filepath.Join(validPath, entry.Name())
			entryURI := pathToResourceURI(entryPath)

			if entry.IsDir() {
				result.WriteString(fmt.Sprintf("[DIR]  %s (%s)\n", entry.Name(), entryURI))
			} else {
				info, err := entry.Info()
				if err == nil {
					result.WriteString(fmt.Sprintf("[FILE] %s (%s) - %d bytes\n",
						entry.Name(), entryURI, info.Size()))
				} else {
					result.WriteString(fmt.Sprintf("[FILE] %s (%s)\n", entry.Name(), entryURI))
				}
			}
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      uri,
				MIMEType: "text/plain",
				Text:     result.String(),
			},
		}, nil
	}

	// It'fss a file, determine how to handle it
	mimeType := detectMimeType(validPath)

	// Check file size
	if fileInfo.Size() > MaxInlineSize {
		// File is too large to inline, return a reference instead
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      uri,
				MIMEType: "text/plain",
				Text:     fmt.Sprintf("File is too large to display inline (%d bytes). Use the read_file tool to access specific portions.", fileInfo.Size()),
			},
		}, nil
	}

	// Read the file content
	content, err := os.ReadFile(validPath)
	if err != nil {
		return nil, err
	}

	// Handle based on content type
	if isTextFile(mimeType) {
		// It'fss a text file, return as text
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      uri,
				MIMEType: mimeType,
				Text:     string(content),
			},
		}, nil
	} else {
		// It'fss a binary file
		if fileInfo.Size() <= MaxBase64Size {
			// Small enough for base64 encoding
			return []mcp.ResourceContents{
				mcp.BlobResourceContents{
					URI:      uri,
					MIMEType: mimeType,
					Blob:     base64.StdEncoding.EncodeToString(content),
				},
			}, nil
		} else {
			// Too large for base64, return a reference
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      uri,
					MIMEType: "text/plain",
					Text:     fmt.Sprintf("Binary file (%s, %d bytes). Use the read_file tool to access specific portions.", mimeType, fileInfo.Size()),
				},
			}, nil
		}
	}
}

// Tool handlers

func (fs *FilesystemServer) handleReadFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return mcp.NewToolResultError("Path must be a string"), nil
	}

	// 判断 前缀是不是已经包含了
	//path = filepath.Join(fss.config.CachePath, path)
	validPath, err := fs.validatePath(path)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("validate Path Error: %v", err)), nil
	}

	// Check if it'fss a directory
	info, err := os.Stat(validPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("check directory error: %v", err)), nil
	}

	if info.IsDir() {
		// For directories, return a resource reference instead
		resourceURI := pathToResourceURI(validPath)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("This is a directory. Use the resource URI to browse its contents: %s", resourceURI),
				},
				mcp.EmbeddedResource{
					Type: "resource",
					Resource: mcp.TextResourceContents{
						URI:      resourceURI,
						MIMEType: "text/plain",
						Text:     fmt.Sprintf("Directory: %s", validPath),
					},
				},
			},
		}, nil
	}

	// Determine MIME type
	mimeType := detectMimeType(validPath)

	// Check file size
	if info.Size() > MaxInlineSize {
		// File is too large to inline, return a resource reference
		resourceURI := pathToResourceURI(validPath)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("File is too large to display inline (%d bytes). Access it via resource URI: %s", info.Size(), resourceURI),
				},
				mcp.EmbeddedResource{
					Type: "resource",
					Resource: mcp.TextResourceContents{
						URI:      resourceURI,
						MIMEType: "text/plain",
						Text:     fmt.Sprintf("Large file: %s (%s, %d bytes)", validPath, mimeType, info.Size()),
					},
				},
			},
		}, nil
	}

	// Read file content
	content, err := os.ReadFile(validPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error reading file: %v", err)), nil
	}

	// Handle based on content type
	if isTextFile(mimeType) {
		// It'fss a text file, return as text
		return mcp.NewToolResultText(string(content)), nil
	} else if isImageFile(mimeType) {
		// It'fss an image file, return as image content
		if info.Size() <= MaxBase64Size {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("Image file: %s (%s, %d bytes)", validPath, mimeType, info.Size()),
					},
					mcp.ImageContent{
						Type:     "image",
						Data:     base64.StdEncoding.EncodeToString(content),
						MIMEType: mimeType,
					},
				},
			}, nil
		} else {
			// Too large for base64, return a reference
			resourceURI := pathToResourceURI(validPath)
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("Image file is too large to display inline (%d bytes). Access it via resource URI: %s", info.Size(), resourceURI),
					},
					mcp.EmbeddedResource{
						Type: "resource",
						Resource: mcp.TextResourceContents{
							URI:      resourceURI,
							MIMEType: "text/plain",
							Text:     fmt.Sprintf("Large image: %s (%s, %d bytes)", validPath, mimeType, info.Size()),
						},
					},
				},
			}, nil
		}
	} else {
		// It'fss another type of binary file
		resourceURI := pathToResourceURI(validPath)

		if info.Size() <= MaxBase64Size {
			// Small enough for base64 encoding
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("Binary file: %s (%s, %d bytes)", validPath, mimeType, info.Size()),
					},
					mcp.EmbeddedResource{
						Type: "resource",
						Resource: mcp.BlobResourceContents{
							URI:      resourceURI,
							MIMEType: mimeType,
							Blob:     base64.StdEncoding.EncodeToString(content),
						},
					},
				},
			}, nil
		} else {
			// Too large for base64, return a reference
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("Binary file: %s (%s, %d bytes). Access it via resource URI: %s", validPath, mimeType, info.Size(), resourceURI),
					},
					mcp.EmbeddedResource{
						Type: "resource",
						Resource: mcp.TextResourceContents{
							URI:      resourceURI,
							MIMEType: "text/plain",
							Text:     fmt.Sprintf("Binary file: %s (%s, %d bytes)", validPath, mimeType, info.Size()),
						},
					},
				},
			}, nil
		}
	}
}

func (fs *FilesystemServer) handleWriteFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return mcp.NewToolResultError("Path must be a string"), nil
	}
	content, ok := request.Params.Arguments["content"].(string)
	if !ok {
		return mcp.NewToolResultError("Content must be a string"), nil
	}

	//path = filepath.Join(fss.config.CachePath, path)

	validPath, err := fs.validatePath(path)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// Check if it'fss a directory
	if info, err := os.Stat(validPath); err == nil && info.IsDir() {
		return mcp.NewToolResultError(fmt.Sprintf("Error: Cannot write to a directory:%s", validPath)), nil
	}

	// Create parent directories if they don't exist
	parentDir := filepath.Dir(validPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error creating parent directories: %v", err)), nil
	}

	if err := os.WriteFile(validPath, []byte(content), 0644); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error writing file: %v", err)), nil
	}

	// Get file info for the response
	info, err := os.Stat(validPath)
	if err != nil {
		// File was written but we couldn't get info
		return mcp.NewToolResultText(fmt.Sprintf("Successfully wrote to %s", path)), nil
	}

	resourceURI := pathToResourceURI(validPath)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Successfully wrote %d bytes to %s", info.Size(), path),
			},
			mcp.EmbeddedResource{
				Type: "resource",
				Resource: mcp.TextResourceContents{
					URI:      resourceURI,
					MIMEType: "text/plain",
					Text:     fmt.Sprintf("File: %s (%d bytes)", validPath, info.Size()),
				},
			},
		},
	}, nil
}

func (fs *FilesystemServer) handleListDirectory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return mcp.NewToolResultError("Path must be a string"), nil
	}

	validPath, err := fs.validatePath(path)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("validate path error: %v, path:%s", err, validPath)), nil
	}

	// Check if it'fss a directory
	info, err := os.Stat(validPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Check directory %s Error: %v", validPath, err)), nil
	}

	if !info.IsDir() {
		return mcp.NewToolResultError(fmt.Sprintf("Error: Path is not a directory:%s", validPath)), nil
	}

	entries, err := os.ReadDir(validPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error reading directory: %v", err)), nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Directory listing for: %s\n\n", validPath))

	for _, entry := range entries {
		entryPath := filepath.Join(validPath, entry.Name())
		resourceURI := pathToResourceURI(entryPath)

		if entry.IsDir() {
			result.WriteString(fmt.Sprintf("[DIR]  %s (%s)\n", entry.Name(), resourceURI))
		} else {
			info, err := entry.Info()
			if err == nil {
				result.WriteString(fmt.Sprintf("[FILE] %s (%s) - %d bytes\n",
					entry.Name(), resourceURI, info.Size()))
			} else {
				result.WriteString(fmt.Sprintf("[FILE] %s (%s)\n", entry.Name(), resourceURI))
			}
		}
	}

	// Return both text content and embedded resource
	resourceURI := pathToResourceURI(validPath)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: result.String(),
			},
			mcp.EmbeddedResource{
				Type: "resource",
				Resource: mcp.TextResourceContents{
					URI:      resourceURI,
					MIMEType: "text/plain",
					Text:     fmt.Sprintf("Directory: %s", validPath),
				},
			},
		},
	}, nil
}

func (fs *FilesystemServer) handleCreateDirectory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return mcp.NewToolResultError("path must be a string"), nil
	}

	validPath, err := fs.validatePath(path)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
	}

	// Check if path already exists
	if info, err := os.Stat(validPath); err == nil {
		if info.IsDir() {
			resourceURI := pathToResourceURI(validPath)
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("Directory already exists: %s", path),
					},
					mcp.EmbeddedResource{
						Type: "resource",
						Resource: mcp.TextResourceContents{
							URI:      resourceURI,
							MIMEType: "text/plain",
							Text:     fmt.Sprintf("Directory: %s", validPath),
						},
					},
				},
			}, nil
		}
		return mcp.NewToolResultError(fmt.Sprintf("Error: Path exists but is not a directory: %s", path)), nil
	}

	if err := os.MkdirAll(validPath, 0755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error creating directory: %v", err)), nil
	}

	resourceURI := pathToResourceURI(validPath)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Successfully created directory %s", path),
			},
			mcp.EmbeddedResource{
				Type: "resource",
				Resource: mcp.TextResourceContents{
					URI:      resourceURI,
					MIMEType: "text/plain",
					Text:     fmt.Sprintf("Directory: %s", validPath),
				},
			},
		},
	}, nil
}

func (fs *FilesystemServer) handleMoveFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	source, ok := request.Params.Arguments["source"].(string)
	if !ok {
		return mcp.NewToolResultError("source must be a string"), nil
	}
	destination, ok := request.Params.Arguments["destination"].(string)
	if !ok {
		return mcp.NewToolResultError("destination must be a string"), nil
	}

	validSource, err := fs.validatePath(source)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error with source path: %v", err)), nil
	}

	// Check if source exists
	if _, err := os.Stat(validSource); os.IsNotExist(err) {
		return mcp.NewToolResultError(fmt.Sprintf("Error: Source does not exist: %s", source)), nil
	}

	validDest, err := fs.validatePath(destination)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error with destination path: %v", err)), nil
	}

	// Create parent directory for destination if it doesn't exist
	destDir := filepath.Dir(validDest)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error creating destination directory: %v", err)), nil
	}

	if err := os.Rename(validSource, validDest); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error moving file: %v", err)), nil
	}

	resourceURI := pathToResourceURI(validDest)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf(
					"Successfully moved %s to %s",
					source,
					destination,
				),
			},
			mcp.EmbeddedResource{
				Type: "resource",
				Resource: mcp.TextResourceContents{
					URI:      resourceURI,
					MIMEType: "text/plain",
					Text:     fmt.Sprintf("Moved file: %s", validDest),
				},
			},
		},
	}, nil
}

func (fs *FilesystemServer) handleSearchFiles(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return mcp.NewToolResultError("path must be a string"), nil
	}
	pattern, ok := request.Params.Arguments["pattern"].(string)
	if !ok {
		return mcp.NewToolResultError("pattern must be a string"), nil
	}

	validPath, err := fs.validatePath(path)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
	}

	// Check if it'fss a directory
	info, err := os.Stat(validPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
	}

	if !info.IsDir() {
		return mcp.NewToolResultError("Error: Search path must be a directory"), nil
	}

	results, err := fs.searchFiles(validPath, pattern)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error searching files: %v", err)), nil
	}

	if len(results) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No files found matching pattern '%s' in %s", pattern, path)), nil
	}

	// Format results with resource URIs
	var formattedResults strings.Builder
	formattedResults.WriteString(fmt.Sprintf("Found %d results:\n\n", len(results)))

	for _, result := range results {
		resourceURI := pathToResourceURI(result)
		info, err := os.Stat(result)
		if err == nil {
			if info.IsDir() {
				formattedResults.WriteString(fmt.Sprintf("[DIR]  %s (%s)\n", result, resourceURI))
			} else {
				formattedResults.WriteString(fmt.Sprintf("[FILE] %s (%s) - %d bytes\n",
					result, resourceURI, info.Size()))
			}
		} else {
			formattedResults.WriteString(fmt.Sprintf("%s (%s)\n", result, resourceURI))
		}
	}

	return mcp.NewToolResultText(formattedResults.String()), nil
}

func (fs *FilesystemServer) handleGetFileInfo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return mcp.NewToolResultError(fmt.Errorf("path %v must be a string", request.Params.Arguments["path"]).Error()), nil
	}

	validPath, err := fs.validatePath(path)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	info, err := fs.getFileStats(validPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error getting file info: %v", err)), nil
	}

	// Get MIME type for files
	mimeType := "directory"
	if info.IsFile {
		mimeType = detectMimeType(validPath)
	}

	resourceURI := pathToResourceURI(validPath)

	// Determine file type text
	var fileTypeText string
	if info.IsDirectory {
		fileTypeText = "Directory"
	} else {
		fileTypeText = "File"
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf(
					"File information for: %s\n\nSize: %d bytes\nCreated: %s\nModified: %s\nAccessed: %s\nIsDirectory: %v\nIsFile: %v\nPermissions: %s\nMIME Type: %s\nResource URI: %s",
					validPath,
					info.Size,
					info.Created.Format(time.RFC3339),
					info.Modified.Format(time.RFC3339),
					info.Accessed.Format(time.RFC3339),
					info.IsDirectory,
					info.IsFile,
					info.Permissions,
					mimeType,
					resourceURI,
				),
			},
			mcp.EmbeddedResource{
				Type: "resource",
				Resource: mcp.TextResourceContents{
					URI:      resourceURI,
					MIMEType: "text/plain",
					Text: fmt.Sprintf("%s: %s (%s, %d bytes)",
						fileTypeText,
						validPath,
						mimeType,
						info.Size),
				},
			},
		},
	}, nil
}

func (fs *FilesystemServer) handleListAllowedDirectories(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Remove the trailing separator for display purposes
	displayDirs := make([]string, len(fs.config.allowedDirs))
	for i, dir := range fs.config.allowedDirs {
		displayDirs[i] = strings.TrimSuffix(dir, string(filepath.Separator))
	}

	var result strings.Builder
	result.WriteString("Allowed directories:")

	for _, dir := range displayDirs {
		resourceURI := pathToResourceURI(dir)
		result.WriteString(fmt.Sprintf("%s (%s)\n", dir, resourceURI))
	}

	return mcp.NewToolResultText(result.String()), nil
}

// Config returns the configuration of the service as a string.
func (fs *FilesystemServer) Config() string {
	fs.config.AllowedDir = strings.Join(fs.config.allowedDirs, ",")
	cfg, err := json.Marshal(fs.config)
	if err != nil {
		fs.logger.Err(err).Msg("failed to marshal config")
		return "{}"
	}
	return string(cfg)
}

func (fs *FilesystemServer) Name() MoLingServerType {
	return FilesystemServerName
}

func (fs *FilesystemServer) Close() error {
	// Cancel the context to stop the browser
	fs.logger.Debug().Msg("closing FilesystemServer")
	return nil
}

// LoadConfig loads the configuration from a JSON object.
func (fs *FilesystemServer) LoadConfig(jsonData map[string]interface{}) error {
	err := mergeJSONToStruct(fs.config, jsonData)
	if err != nil {
		return err
	}
	fs.config.allowedDirs = strings.Split(fs.config.AllowedDir, ",")
	return fs.config.Check()
}

func init() {
	RegisterServ(FilesystemServerName, NewFilesystemServer)
}

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
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// detectMimeType tries to determine the MIME type of a file
func detectMimeType(path string) string {
	// First try by extension
	ext := filepath.Ext(path)
	if ext != "" {
		mimeType := mime.TypeByExtension(ext)
		if mimeType != "" {
			return mimeType
		}
	}

	// If that fails, try to read a bit of the file
	file, err := os.Open(path)
	if err != nil {
		return "application/octet-stream" // Default
	}
	defer file.Close()

	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		return "application/octet-stream" // Default
	}

	// Use http.DetectContentType
	return http.DetectContentType(buffer[:n])
}

// isTextFile determines if a file is likely a text file based on MIME type
func isTextFile(mimeType string) bool {
	return strings.HasPrefix(mimeType, "text/") ||
		mimeType == "application/json" ||
		mimeType == "application/xml" ||
		mimeType == "application/javascript" ||
		mimeType == "application/x-javascript" ||
		strings.Contains(mimeType, "+xml") ||
		strings.Contains(mimeType, "+json")
}

// isImageFile determines if a file is an image based on MIME type
func isImageFile(mimeType string) bool {
	return strings.HasPrefix(mimeType, "image/")
}

// pathToResourceURI converts a file path to a resource URI
func pathToResourceURI(path string) string {
	return "file://" + path
}

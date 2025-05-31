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

package utils

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// CreateDirectory checks if a directory exists, and creates it if it doesn't
func CreateDirectory(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0o755)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

// StringInSlice checks if a string is in a slice of strings
func StringInSlice(s string, modules []string) bool {
	for _, module := range modules {
		if module == s {
			return true
		}
	}
	return false
}

// MergeJSONToStruct 将JSON中的字段合并到结构体中
func MergeJSONToStruct(target any, jsonMap map[string]any) error {
	// 获取目标结构体的反射值
	val := reflect.ValueOf(target).Elem()
	typ := val.Type()

	// 遍历JSON map中的每个字段
	for jsonKey, jsonValue := range jsonMap {
		// 遍历结构体的每个字段
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			// 检查JSON字段名是否与结构体的JSON tag匹配
			if field.Tag.Get("json") == jsonKey {
				// 获取结构体字段的反射值
				fieldVal := val.Field(i)
				// 检查字段是否可设置
				if fieldVal.CanSet() {
					// 将JSON值转换为结构体字段的类型
					jsonVal := reflect.ValueOf(jsonValue)
					if jsonVal.Type().ConvertibleTo(fieldVal.Type()) {
						fieldVal.Set(jsonVal.Convert(fieldVal.Type()))
					} else {
						return fmt.Errorf("type mismatch for field %s, value:%v", jsonKey, jsonValue)
					}
				}
			}
		}
	}
	return nil
}

// DetectMimeType tries to determine the MIME type of a file
func DetectMimeType(path string) string {
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

// IsTextFile determines if a file is likely a text file based on MIME type
func IsTextFile(mimeType string) bool {
	return strings.HasPrefix(mimeType, "text/") ||
		mimeType == "application/json" ||
		mimeType == "application/xml" ||
		mimeType == "application/javascript" ||
		mimeType == "application/x-javascript" ||
		strings.Contains(mimeType, "+xml") ||
		strings.Contains(mimeType, "+json")
}

// IsImageFile determines if a file is an image based on MIME type
func IsImageFile(mimeType string) bool {
	return strings.HasPrefix(mimeType, "image/")
}

// PathToResourceURI converts a file path to a resource URI
func PathToResourceURI(path string) string {
	return "file://" + path
}

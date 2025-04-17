/*
 * Copyright 2025 CFC4N <cfc4n.cs@gmail.com>. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Repository: https://github.com/gojue/moling
 */

package utils

import "os"

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

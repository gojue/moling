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

import (
	"fmt"
	"os"
)

var pidFile *os.File

// CreatePIDFile creates and locks a PID file to prevent multiple instances.
func CreatePIDFile(pidFilePath string) error {
	// Open or create the PID file
	file, err := os.OpenFile(pidFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open PID file: %w", err)
	}

	// Try to lock the file using platform-specific code
	locked, err := lockFile(file)
	if err != nil {
		_ = file.Close()
		return fmt.Errorf("failed to lock PID file: %w", err)
	}
	if !locked {
		_ = file.Close()
		return fmt.Errorf("another instance is already running: %s", pidFilePath)
	}

	// Write the current PID to the file
	err = file.Truncate(0)
	if err != nil {
		_ = unlockFile(file)
		_ = file.Close()
		return fmt.Errorf("failed to truncate PID file: %w", err)
	}
	_, err = file.WriteString(fmt.Sprintf("%d\n", os.Getpid()))
	if err != nil {
		_ = unlockFile(file)
		_ = file.Close()
		return fmt.Errorf("failed to write PID to file: %w", err)
	}

	// Keep the file open to maintain the lock
	pidFile = file
	return nil
}

// RemovePIDFile releases the lock and removes the PID file.
func RemovePIDFile(pidFilePath string) error {
	if pidFile != nil {
		err := unlockFile(pidFile)
		if err != nil {
			return fmt.Errorf("failed to unlock PID file: %w", err)
		}
		_ = pidFile.Close()
		pidFile = nil
		err = os.Remove(pidFilePath)
		if err != nil {
			return fmt.Errorf("failed to remove PID file: %w", err)
		}
	}
	return nil
}

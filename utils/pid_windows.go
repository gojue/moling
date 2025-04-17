//go:build windows

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
	"errors"
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32     = syscall.NewLazyDLL("kernel32.dll")
	lockFileEx   = kernel32.NewProc("LockFileEx")
	unlockFileEx = kernel32.NewProc("UnlockFileEx")
)

const (
	LockfileExclusiveLock   = 2
	LockfileFailImmediately = 1
)
const ErrorLockViolation = syscall.Errno(33) // 0x21

// lockFile locks the given file using Windows API.
func lockFile(file *os.File) (bool, error) {
	handle := syscall.Handle(file.Fd())
	var overlapped syscall.Overlapped

	flags := LockfileExclusiveLock | LockfileFailImmediately
	r, _, err := lockFileEx.Call(
		uintptr(handle),
		uintptr(flags),
		0,
		1,
		0,
		uintptr(unsafe.Pointer(&overlapped)),
	)

	if r == 0 {
		if !errors.Is(err, ErrorLockViolation) {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

// unlockFile unlocks the given file using Windows API.
func unlockFile(file *os.File) error {
	handle := syscall.Handle(file.Fd())
	var overlapped syscall.Overlapped

	r, _, err := unlockFileEx.Call(
		uintptr(handle),
		0,
		1,
		0,
		uintptr(unsafe.Pointer(&overlapped)),
	)

	if r == 0 {
		return err
	}

	return nil
}

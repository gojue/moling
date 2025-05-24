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
	"os"
	"sync"
)

// RotateWriter 是一个简单的日志轮转写入器
type RotateWriter struct {
	filePaths    [2]string // 两个日志文件的路径
	currentIndex int       // 当前写入的文件索引（0或1）
	maxSize      int64     // 文件大小阈值（字节）
	count        uint16    // 当前文件的写入次数
	mu           sync.Mutex
	file         *os.File // 当前打开的文件句柄
}

// NewRotateWriter 创建一个新的 RotateWriter 实例
func NewRotateWriter(filePath string, maxSize int64) (*RotateWriter, error) {
	rw := &RotateWriter{
		filePaths:    [2]string{filePath + ".1", filePath + ".2"},
		currentIndex: 0,
		maxSize:      maxSize,
		count:        0,
	}

	// 初始化文件，确保两个文件都存在
	for _, path := range rw.filePaths {
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	// 打开当前文件
	file, err := os.OpenFile(rw.filePaths[rw.currentIndex], os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	rw.file = file

	return rw, nil
}

// Write 实现 io.Writer 接口
func (rw *RotateWriter) Write(p []byte) (n int, err error) {
	err = rw.CheckFileSize()
	if err != nil {
		return 0, err
	}
	// 写入日志
	return rw.file.Write(p)
}

// CheckFileSize 检查当前文件大小
func (rw *RotateWriter) CheckFileSize() error {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	if rw.count >= 10 {
		// 检查当前文件大小
		fileInfo, err := rw.file.Stat()
		if err != nil {
			return err
		}

		// 如果文件大小超过阈值，切换到另一个文件并清空
		if fileInfo.Size() >= rw.maxSize {
			// 关闭当前文件
			rw.file.Close()

			// 切换到另一个文件
			rw.currentIndex = (rw.currentIndex + 1) % 2

			// 清空目标文件
			file, err := os.OpenFile(rw.filePaths[rw.currentIndex], os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
			if err != nil {
				return err
			}
			rw.file = file
		}
		rw.count = 0
	}

	rw.count++
	return nil
}

// Close 关闭 RotateWriter
func (rw *RotateWriter) Close() error {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	return rw.file.Close()
}

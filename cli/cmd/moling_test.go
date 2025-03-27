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

package cmd

import (
	"context"
	"github.com/gojue/moling/services"
	"path/filepath"
	"testing"
)

func TestNewMLServer(t *testing.T) {
	err := CreateDirectory(mlConfig.BasePath)
	if err != nil {
		t.Errorf("Failed to create base directory: %v", err)
	}
	for _, dirName := range []string{"logs", "config", "browser"} {
		err = CreateDirectory(filepath.Join(mlConfig.BasePath, dirName))
		if err != nil {
			t.Errorf("Failed to create directory %s: %v", dirName, err)
		}
	}
	loger := initLogger(mlConfig.BasePath)
	mlConfig.SetLogger(loger)

	ctx := context.WithValue(context.Background(), services.MoLingConfigKey, mlConfig)
	ctx = context.WithValue(ctx, services.MoLingLoggerKey, loger)
	// Create a new server with the filesystem service
	fs, err := services.NewFilesystemServer(ctx, []string{"/tmp/"})
	if err != nil {
		t.Errorf("Failed to create filesystem server: %v", err)
	}
	srvs := []services.Service{
		fs,
	}
	srv, err := NewMoLingServer(ctx, srvs)
	if err != nil {
		t.Errorf("Failed to create server: %v", err)
	}
	err = srv.Serve()
	if err != nil {
		t.Errorf("Failed to start server: %v", err)
	}
	t.Logf("Server started successfully: %v", srv)
}

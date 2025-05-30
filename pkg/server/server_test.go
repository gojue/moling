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

package server

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gojue/moling/pkg/comm"
	"github.com/gojue/moling/pkg/config"
	"github.com/gojue/moling/pkg/services/abstract"
	"github.com/gojue/moling/pkg/services/filesystem"
	"github.com/gojue/moling/pkg/utils"
)

func TestNewMLServer(t *testing.T) {
	// Create a new MoLingConfig
	mlConfig := config.MoLingConfig{
		BasePath: filepath.Join(os.TempDir(), "moling_test"),
	}
	mlDirectories := []string{
		"logs",    // log file
		"config",  // config file
		"browser", // browser cache
		"data",    // data
		"cache",
	}
	err := utils.CreateDirectory(mlConfig.BasePath)
	if err != nil {
		t.Errorf("Failed to create base directory: %w", err)
	}
	for _, dirName := range mlDirectories {
		err = utils.CreateDirectory(filepath.Join(mlConfig.BasePath, dirName))
		if err != nil {
			t.Errorf("Failed to create directory %s: %w", dirName, err)
		}
	}
	logger, ctx, err := comm.InitTestEnv()
	if err != nil {
		t.Fatalf("Failed to initialize test environment: %w", err)
	}
	logger.Info().Msg("TestBrowserServer")
	mlConfig.SetLogger(logger)

	// Create a new server with the filesystem service
	fs, err := filesystem.NewFilesystemServer(ctx)
	if err != nil {
		t.Errorf("Failed to create filesystem server: %w", err)
	}
	err = fs.Init()
	if err != nil {
		t.Errorf("Failed to initialize filesystem server: %w", err)
	}
	srvs := []abstract.Service{
		fs,
	}
	srv, err := NewMoLingServer(ctx, srvs, mlConfig)
	if err != nil {
		t.Errorf("Failed to create server: %w", err)
	}
	err = srv.Serve()
	if err != nil {
		t.Errorf("Failed to start server: %w", err)
	}
	t.Logf("Server started successfully: %v", srv)
}

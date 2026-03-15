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
	"net/http"
	"net/http/httptest"
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
		t.Errorf("Failed to create base directory: %s", err.Error())
	}
	for _, dirName := range mlDirectories {
		err = utils.CreateDirectory(filepath.Join(mlConfig.BasePath, dirName))
		if err != nil {
			t.Errorf("Failed to create directory %s: %s", dirName, err.Error())
		}
	}
	logger, ctx, err := comm.InitTestEnv()
	if err != nil {
		t.Fatalf("Failed to initialize test environment: %s", err.Error())
	}
	logger.Info().Msg("TestBrowserServer")
	mlConfig.SetLogger(logger)

	// Create a new server with the filesystem service
	fs, err := filesystem.NewFilesystemServer(ctx)
	if err != nil {
		t.Errorf("Failed to create filesystem server: %s", err.Error())
	}
	err = fs.Init()
	if err != nil {
		t.Errorf("Failed to initialize filesystem server: %s", err.Error())
	}
	srvs := []abstract.Service{
		fs,
	}
	srv, err := NewMoLingServer(ctx, srvs, mlConfig)
	if err != nil {
		t.Errorf("Failed to create server: %s", err.Error())
	}
	err = srv.Serve()
	if err != nil {
		t.Errorf("Failed to start server: %s", err.Error())
	}
	t.Logf("Server started successfully: %v", srv)
}

// TestRequireJSONContentType verifies that the middleware blocks POST requests
// with non-application/json Content-Types (which browsers treat as "simple
// requests" and therefore never trigger a CORS preflight).
func TestRequireJSONContentType(t *testing.T) {
	// A simple downstream handler that always returns 200.
	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := requireJSONContentType(ok)

	tests := []struct {
		name           string
		method         string
		contentType    string
		expectedStatus int
	}{
		// Non-POST requests must pass through regardless of Content-Type.
		{"GET no CT", http.MethodGet, "", http.StatusOK},
		{"GET text/plain", http.MethodGet, "text/plain", http.StatusOK},
		// POST with application/json (with and without charset param) must pass.
		{"POST application/json", http.MethodPost, "application/json", http.StatusOK},
		{"POST application/json; charset=utf-8", http.MethodPost, "application/json; charset=utf-8", http.StatusOK},
		// POST with "simple" Content-Types that bypass CORS preflight must be rejected.
		{"POST text/plain", http.MethodPost, "text/plain", http.StatusUnsupportedMediaType},
		{"POST application/x-www-form-urlencoded", http.MethodPost, "application/x-www-form-urlencoded", http.StatusUnsupportedMediaType},
		{"POST multipart/form-data", http.MethodPost, "multipart/form-data", http.StatusUnsupportedMediaType},
		// POST with empty or missing Content-Type must also be rejected.
		{"POST no CT", http.MethodPost, "", http.StatusUnsupportedMediaType},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/message", nil)
			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rr.Code)
			}
		})
	}
}

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

package comm

import (
	"context"
	"os"
	"path/filepath"

	"github.com/gojue/moling/pkg/config"
	"github.com/rs/zerolog"
)

type MoLingServerType string

type contextKey string

// MoLingConfigKey is a context key for storing the version of MoLing
const (
	MoLingConfigKey contextKey = "moling_config"
	MoLingLoggerKey contextKey = "moling_logger"
)

// InitTestEnv initializes the test environment by creating a temporary log file and setting up the logger.
func InitTestEnv() (zerolog.Logger, context.Context, error) {
	logFile := filepath.Join(os.TempDir(), "moling.log")
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	var logger zerolog.Logger
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return zerolog.Logger{}, nil, err
	}
	logger = zerolog.New(f).With().Timestamp().Logger()
	mlConfig := &config.MoLingConfig{
		ConfigFile: filepath.Join("config", "test_config.json"),
		BasePath:   os.TempDir(),
	}
	ctx := context.WithValue(context.Background(), MoLingConfigKey, mlConfig)
	ctx = context.WithValue(ctx, MoLingLoggerKey, logger)
	return logger, ctx, nil
}

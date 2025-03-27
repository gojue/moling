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

import "github.com/rs/zerolog"

// Config is an interface that defines a method for checking configuration validity.
type Config interface {
	// Check validates the configuration and returns an error if the configuration is invalid.
	Check() error
}

// MoLingConfig is a struct that holds the configuration for the MoLing server.
type MoLingConfig struct {
	ConfigFile string `json:"config_file"` // The path to the configuration file.
	BasePath   string `json:"base_path"`   // The base path for the server, used for storing files. automatically created if not exists.
	//AllowDir   []string `json:"allow_dir"`   // The directories that are allowed to be accessed by the server.
	Version    string `json:"version"`     // The version of the MoLing server.
	ListenAddr string `json:"listen_addr"` // The address to listen on for SSE mode.
	Debug      bool   `json:"debug"`       // Debug mode, if true, the server will run in debug mode.
	logger     zerolog.Logger
}

func (cfg *MoLingConfig) Check() error {
	panic("not implemented yet") // TODO: Implement Check
}

func (cfg *MoLingConfig) Logger() zerolog.Logger {
	return cfg.logger
}

func (cfg *MoLingConfig) SetLogger(logger zerolog.Logger) {
	cfg.logger = logger
}

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

// Package services provides a set of services for the MoLing application.
package services

import (
	"fmt"
	"os"
	"path/filepath"
)

type BrowserConfig struct {
	Headless        bool   `json:"headless"`
	Timeout         int    `json:"timeout"`
	Proxy           string `json:"proxy"`
	UserAgent       string `json:"user_agent"`
	DefaultLanguage string `json:"default_language"`
	URLTimeout      int    `json:"url_timeout"`
	CSSTimeout      int    `json:"css_timeout"`
	DataPath        string `json:"data_path"` // DataPath is the path to the data directory.
	//logger          *zerolog.Logger
}

func (cfg *BrowserConfig) Check() error {
	if cfg.Timeout <= 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}
	if cfg.URLTimeout <= 0 {
		return fmt.Errorf("URL timeout must be greater than 0")
	}
	if cfg.CSSTimeout <= 0 {
		return fmt.Errorf("CSS timeout must be greater than 0")
	}
	return nil
}

// NewBrowserConfig creates a new BrowserConfig with default values.
func NewBrowserConfig() *BrowserConfig {
	return &BrowserConfig{
		Headless:        false,
		Timeout:         30,
		URLTimeout:      10,
		CSSTimeout:      10,
		UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3 (GoJue MoLing)",
		DefaultLanguage: "en-US",
		DataPath:        filepath.Join(os.TempDir(), ".moling", "data"),
	}
}

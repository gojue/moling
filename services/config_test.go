/*
 *
 *  Copyright 2025 CFC4N <cfc4n.cs@gmail.com>. All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 *  Repository: https://github.com/gojue/moling
 *
 */

package services

import (
	"encoding/json"
	"os"
	"testing"
)

// TestConfigLoad tests the loading of the configuration from a JSON file.
func TestConfigLoad(t *testing.T) {
	configFile := "config_test.json"
	cfg := &MoLingConfig{}
	cfg.ConfigFile = "config.json"
	cfg.BasePath = "/tmp/moling"
	cfg.Version = "1.0.0"
	cfg.ListenAddr = ":8080"
	cfg.Debug = true
	cfg.Username = "user1"
	cfg.HomeDir = "/Users/user1"
	cfg.SystemInfo = "Darwin 15.3.3"

	jsonData, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	mlConfig, ok := jsonMap["MoLingConfig"].(map[string]interface{})
	if !ok {
		t.Fatalf("failed to parse MoLingConfig from JSON")
	}
	if err := mergeJSONToStruct(cfg, mlConfig); err != nil {
		t.Fatalf("failed to merge JSON to struct: %v", err)
	}
	t.Logf("Config loaded, MoLing Config.BasePath: %s", cfg.BasePath)
	if cfg.BasePath != "/newpath/.moling" {
		t.Fatalf("expected BasePath to be '/newpath/.moling', got '%s'", cfg.BasePath)
	}
}

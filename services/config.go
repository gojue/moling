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

import (
	"fmt"
	"github.com/rs/zerolog"
	"reflect"
)

// Config is an interface that defines a method for checking configuration validity.
type Config interface {
	// Check validates the configuration and returns an error if the configuration is invalid.
	Check() error
}

// MoLingConfig is a struct that holds the configuration for the MoLing server.
type MoLingConfig struct {
	ConfigFile string `json:"config_file"` // The path to the configuration file.
	BasePath   string `json:"base_path"`   // The base path for the server, used for storing files. automatically created if not exists. eg: /Users/user1/.moling
	//AllowDir   []string `json:"allow_dir"`   // The directories that are allowed to be accessed by the server.
	Version    string `json:"version"`     // The version of the MoLing server.
	ListenAddr string `json:"listen_addr"` // The address to listen on for SSE mode.
	Debug      bool   `json:"debug"`       // Debug mode, if true, the server will run in debug mode.
	Username   string // The username of the user running the server.
	HomeDir    string // The home directory of the user running the server. macOS: /Users/user1, Linux: /home/user1
	SystemInfo string // The system information of the user running the server. macOS: Darwin 15.3.3, Linux: Ubuntu 20.04.1 LTS
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

// mergeJSONToStruct 将JSON中的字段合并到结构体中

func mergeJSONToStruct(target interface{}, jsonMap map[string]interface{}) error {
	// 获取目标结构体的反射值
	val := reflect.ValueOf(target).Elem()
	typ := val.Type()

	// 遍历JSON map中的每个字段
	for jsonKey, jsonValue := range jsonMap {
		// 遍历结构体的每个字段
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			// 检查JSON字段名是否与结构体的JSON tag匹配
			if field.Tag.Get("json") == jsonKey {
				// 获取结构体字段的反射值
				fieldVal := val.Field(i)
				// 检查字段是否可设置
				if fieldVal.CanSet() {
					// 将JSON值转换为结构体字段的类型
					jsonVal := reflect.ValueOf(jsonValue)
					if jsonVal.Type().ConvertibleTo(fieldVal.Type()) {
						fieldVal.Set(jsonVal.Convert(fieldVal.Type()))
					} else {
						return fmt.Errorf("type mismatch for field %s, value:%v", jsonKey, jsonValue)
					}
				}
			}
		}
	}
	return nil
}

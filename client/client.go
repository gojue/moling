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

package client

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog"
	"os"
)

var (
	// ClineConfigPath is the path to the Cline config file.
	clientLists = make(map[string]string, 3)
)

const MCPServersKey = "mcpServers"

// MCPServerConfig represents the configuration for the MCP Client.
type MCPServerConfig struct {
	Description string   `json:"description"`       // Description of the MCP Server
	IsActive    bool     `json:"isActive"`          // Is the MCP Server active
	Command     string   `json:"command,omitempty"` // Command to start the MCP Server, STDIO mode only
	Args        []string `json:"args,omitempty"`    // Arguments to pass to the command, STDIO mode only
	BaseUrl     string   `json:"baseUrl,omitempty"` // Base URL of the MCP Server, SSE mode only
	TimeOut     uint16   `json:"timeout,omitempty"` // Timeout for the MCP Server, default is 300 seconds
	ServerName  string
}

// NewMCPServerConfig creates a new MCPServerConfig instance.
func NewMCPServerConfig(description string, command string, srvName string) MCPServerConfig {
	return MCPServerConfig{
		Description: description,
		IsActive:    true,
		Command:     command,
		Args:        []string{}, // not used
		BaseUrl:     "",
		ServerName:  srvName,
		TimeOut:     300,
	}
}

// Manager manages the configuration of different clients.
type Manager struct {
	logger    zerolog.Logger
	clients   map[string]string
	mcpConfig MCPServerConfig
}

// NewManager creates a new ClientManager instance.
func NewManager(lger zerolog.Logger, mcpConfig MCPServerConfig) (cm *Manager) {
	cm = &Manager{
		clients:   make(map[string]string, 3),
		logger:    lger,
		mcpConfig: mcpConfig,
	}
	cm.clients = clientLists
	return cm
}

// ListClient lists all the clients and checks if they exist.
func (c *Manager) ListClient() {
	for name, path := range c.clients {
		c.logger.Debug().Msgf("Client %s: %s", name, path)
		if !c.checkExist(path) {
			// path not exists
			c.logger.Info().Str("Client Name", name).Bool("exist", false).Msg("Client is not exist")
		} else {
			c.logger.Info().Str("Client Name", name).Bool("exist", true).Msg("Client is exist")
		}
	}
	return
}

// SetupConfig sets up the configuration for the clients.
func (c *Manager) SetupConfig() {
	for name, path := range c.clients {
		c.logger.Debug().Msgf("Client %s: %s", name, path)
		if !c.checkExist(path) {
			continue
		}
		// read config file
		file, err := os.ReadFile(path)
		if err != nil {
			c.logger.Error().Str("Client Name", name).Msgf("Failed to open config file %s: %s", path, err)
			continue
		}
		c.logger.Debug().Str("Client Name", name).Str("config", string(file)).Send()
		b, err := c.appendConfig(c.mcpConfig.ServerName, file)
		if err != nil {
			c.logger.Error().Str("Client Name", name).Msgf("Failed to append config file %s: %s", path, err)
			continue
		}
		c.logger.Debug().Str("Client Name", name).Str("newConfig", string(b)).Send()
		// write config file
		err = os.WriteFile(path, b, 0644)
		if err != nil {
			c.logger.Error().Str("Client Name", name).Msgf("Failed to write config file %s: %s", path, err)
			continue
		}
		c.logger.Info().Str("Client Name", name).Msgf("Successfully added config to %s", path)
	}
	return
}

// appendConfig appends the mlMCPConfig to the client config.
func (c *Manager) appendConfig(name string, payload []byte) ([]byte, error) {
	var err error
	var jsonMap map[string]interface{}
	var jsonBytes []byte
	err = json.Unmarshal(payload, &jsonMap)
	if err != nil {
		return nil, err
	}
	jsonMcpServer, ok := jsonMap[MCPServersKey].(map[string]interface{})
	if !ok {
		return nil, errors.New("MCPServersKey not found in JSON")
	}
	jsonMcpServer[name] = c.mcpConfig
	jsonMap[MCPServersKey] = jsonMcpServer
	jsonBytes, err = json.MarshalIndent(jsonMap, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

// checkExist checks if the file at the given path exists.
func (c *Manager) checkExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			c.logger.Debug().Msgf("Client config file %s does not exist", path)
			return false
		}
		c.logger.Info().Msgf("check file failed, error:%v", err)
		return false
	}
	return true
}

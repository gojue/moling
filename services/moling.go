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
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"
	"log"
	"os"
	"strings"
	"time"
)

type MoLingServerType string // MoLingServerType is the type of the server

type MoLingServer struct {
	ctx        context.Context
	server     *server.MCPServer
	services   []Service
	logger     zerolog.Logger
	mlConfig   MoLingConfig
	listenAddr string // SSE mode listen address, if empty, use STDIO mode.
}

func NewMoLingServer(ctx context.Context, srvs []Service, mlConfig MoLingConfig) (*MoLingServer, error) {
	mcpServer := server.NewMCPServer(
		mlConfig.ServerName,
		mlConfig.Version,
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithPromptCapabilities(true),
	)
	// Set the context for the server
	ms := &MoLingServer{
		ctx:        ctx,
		server:     mcpServer,
		services:   srvs,
		listenAddr: mlConfig.ListenAddr,
		logger:     ctx.Value(MoLingLoggerKey).(zerolog.Logger),
		mlConfig:   mlConfig,
	}
	err := ms.init()
	return ms, err
}

func (m *MoLingServer) init() error {
	var err error
	for _, srv := range m.services {
		m.logger.Debug().Str("serviceName", string(srv.Name())).Msg("Loading service")
		err = m.loadService(srv)
		if err != nil {
			m.logger.Info().Err(err).Str("serviceName", string(srv.Name())).Msg("Failed to load service")
		}
	}
	return err
}

func (m *MoLingServer) loadService(srv Service) error {

	// Add resources
	for r, rhf := range srv.Resources() {
		m.server.AddResource(r, rhf)
	}

	// Add Resource Templates
	for rt, rthf := range srv.ResourceTemplates() {
		m.server.AddResourceTemplate(rt, rthf)
	}

	// Add Tools
	m.server.AddTools(srv.Tools()...)

	// Add Notification Handlers
	for n, nhf := range srv.NotificationHandlers() {
		m.server.AddNotificationHandler(n, nhf)
	}

	// Add Prompts
	for _, pe := range srv.Prompts() {
		// Add Prompt
		m.server.AddPrompt(pe.Prompt(), pe.Handler())
	}
	return nil
}

func (s *MoLingServer) Serve() error {
	mLogger := log.New(s.logger, s.mlConfig.ServerName, 0)
	if s.listenAddr != "" {
		ltnAddr := fmt.Sprintf("http://%s", strings.TrimPrefix(s.listenAddr, "http://"))
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		multi := zerolog.MultiLevelWriter(consoleWriter, s.logger)
		s.logger = zerolog.New(multi).With().Timestamp().Logger()
		s.logger.Info().Str("listenAddr", s.listenAddr).Str("BaseURL", ltnAddr).Msg("Starting SSE server")
		s.logger.Warn().Msgf("The SSE server URL must be: %s. Please do not make mistakes, even if it is another IP or domain name on the same computer, it cannot be mixed.", ltnAddr)
		return server.NewSSEServer(s.server, server.WithBaseURL(ltnAddr)).Start(s.listenAddr)
	}
	s.logger.Info().Msg("Starting STDIO server")
	return server.ServeStdio(s.server, server.WithErrorLogger(mLogger))
}

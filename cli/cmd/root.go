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

package cmd

import (
	"context"
	"fmt"
	"github.com/gojue/moling/cli/cobrautl"
	"github.com/gojue/moling/services"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"path/filepath"
)

const (
	CliName            = "moling"
	CliNameZh          = "魔灵"
	MCPServerName      = "Gojue MoLing Server"
	CliDescription     = "A local office automation magic assistant, Moling can help you perform file viewing, reading, and other related tasks, as well as assist with command-line operations."
	CliDescriptionZh   = "本地办公自动化魔法助手，可以帮你实现文件查看、阅读等操作，也可以帮你执行命令行操作。"
	CliHomepage        = "https://gojue.cc/moling"
	CliAuthor          = "CFC4N <cfc4ncs@gmail.com>"
	CliGithubRepo      = "https://github.com/gojue/moling"
	CliDescriptionLong = `
MoLing is a computer-based MCP Server that implements system interaction through operating system APIs, enabling file system operations such as reading, writing, merging, statistics, and aggregation, as well as the ability to execute system commands. It is a dependency-free local office automation assistant.

Requiring no installation of any dependencies, Moling can be run directly and is compatible with multiple operating systems, including Windows, Linux, and macOS. This eliminates the hassle of dealing with environment conflicts involving Node.js, Python, and other development environments.

Usage:
  moling config
  moling -l 127.0.0:8080
  moling -h
`
	CliDescriptionLongZh = `Moling（魔灵） 是一个computer-use的MCP Server，基于操作系统API实现了系统交互，可以实现文件系统的读写、合并、统计、聚合等操作，也可以执行系统命令操作。是一个无需任何依赖的本地办公自动化助手。
没有任何安装依赖，直接运行，兼容Windows、Linux、macOS等操作系统。再也不用苦恼NodeJS、Python等环境冲突等问题。

Usage:
  moling config
  moling -l 127.0.0:8080
  moling -h
`
)

const (
	MLConfigName = "config.json" // config file name of MoLing Server
)

var (
	GitVersion = "linux-arm64-202503222008-v0.0.1"
	mlConfig   = &services.MoLingConfig{
		Version:    GitVersion,
		ConfigFile: filepath.Join("config", MLConfigName),
		BasePath:   filepath.Join(os.TempDir(), ".moling"), // will set in mlsCommandPreFunc
	}

	// mlDirectories is a list of directories to be created in the base path
	mlDirectories = []string{
		"logs",
		"config",
		"browser",
		"data",
		"cache",
	}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:        CliName,
	Short:      CliDescription,
	SuggestFor: []string{"molin", "moli", "mling"},

	Long: CliDescriptionLong,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE:              mlsCommandFunc,
	PersistentPreRunE: mlsCommandPreFunc,
}

func usageFunc(c *cobra.Command) error {
	return cobrautl.UsageFunc(c, GitVersion)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.SetUsageFunc(usageFunc)
	rootCmd.SetHelpTemplate(`{{.UsageString}}`)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Version = GitVersion
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version:\t%s" .Version}}
`)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// set default config file path
	currentUser, err := user.Current()
	if err == nil {
		mlConfig.BasePath = filepath.Join(currentUser.HomeDir, ".moling")
	}

	cobra.EnablePrefixMatching = true
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVar(&mlConfig.BasePath, "base_path", mlConfig.BasePath, "MoLing Base Data Path, automatically set by the system, cannot be changed, display only.")
	rootCmd.PersistentFlags().BoolVarP(&mlConfig.Debug, "debug", "d", false, "Debug mode, default is false.")
	rootCmd.SilenceUsage = true
}

// initLogger init logger
func initLogger(mlDataPath string) zerolog.Logger {
	var logger zerolog.Logger
	var err error
	logFile := filepath.Join(mlDataPath, "logs", "moling.log")
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if mlConfig.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	var f *os.File
	f, err = os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		panic(fmt.Sprintf("failed to open log file %s: %v", logFile, err))
	}
	logger = zerolog.New(f).With().Timestamp().Logger()
	logger.Info().Str("ServerName", MCPServerName).Str("version", GitVersion).Msg("start")
	return logger
}

func mlsCommandFunc(command *cobra.Command, args []string) error {
	loger := initLogger(mlConfig.BasePath)
	mlConfig.SetLogger(loger)

	ctx := context.WithValue(context.Background(), services.MoLingConfigKey, mlConfig)
	ctx = context.WithValue(ctx, services.MoLingLoggerKey, loger)

	var srvs []services.Service
	var closers []func() error
	for _, nsv := range services.ServiceList() {
		srv, err := nsv(ctx, args)
		if err != nil {
			return err
		}
		srvs = append(srvs, srv)
		closers = append(closers, srv.Close)
	}

	// MCPServer
	srv, err := NewMoLingServer(ctx, srvs)
	if err != nil {
		return err
	}
	err = srv.Serve()
	if err != nil {
		return err
	}

	// close all services
	for i, closer := range closers {
		err = closer()
		if err != nil {
			loger.Error().Err(err).Msgf("failed to close service %d", i)
		} else {
			loger.Info().Msgf("service %d closed", i)
		}
	}
	return err
}

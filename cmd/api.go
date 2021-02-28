package cmd

import (
	"github.com/michelaquino/golang_api_skeleton/config"
	"github.com/spf13/cobra"

	"github.com/michelaquino/golang_api_skeleton/src/log"
	"github.com/michelaquino/golang_api_skeleton/src/server"
)

var (
	logger = log.GetLogger()
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Starts API server",
	Long:  `Starts API server.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.Init()
		server.Start()
	},
}

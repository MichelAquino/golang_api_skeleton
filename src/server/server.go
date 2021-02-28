package server

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/michelaquino/golang_api_skeleton/src/log"
	"github.com/spf13/viper"
)

var (
	logger = log.GetLogger()
)

// Start HTTP server
func Start() {
	echoInstance := echo.New()
	ctx := context.Background()

	// Configure New Relic
	configureNewRelic(ctx, echoInstance)

	// Middlewares
	configureMiddlewares(echoInstance)

	// Configure routes
	configureRoutes(echoInstance)

	port := viper.GetInt("api.host.port")
	logger.Info(ctx, "start api", fmt.Sprintf("Started at %d", port), nil)

	echoInstance.Logger.Fatal(echoInstance.Start(fmt.Sprintf(":%d", port)))
}

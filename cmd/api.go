package cmd

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/michelaquino/golang_api_skeleton/config"
	"github.com/michelaquino/golang_api_skeleton/src/handlers"
	"github.com/michelaquino/golang_api_skeleton/src/repository"
	newrelic "github.com/newrelic/go-agent"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/michelaquino/golang_api_skeleton/src/log"
	apiMiddleware "github.com/michelaquino/golang_api_skeleton/src/middleware"
)

var (
	logger = log.GetLogger()
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Starts API service",
	Long:  `Starts API service.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.Init()
		start()
	},
}

func start() {
	echoInstance := echo.New()
	ctx := context.Background()

	// Configure New Relic
	configureNewRelic(ctx, echoInstance)

	// Middlewares
	echoInstance.Pre(middleware.RemoveTrailingSlash())
	echoInstance.Use(apiMiddleware.AssignRequestID)
	echoInstance.Use(apiMiddleware.RequestLogger)
	echoInstance.Use(middleware.Recover())
	echoInstance.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		ErrorHandler: func(err error, e echo.Context) error {
			return echo.NewHTTPError(http.StatusGatewayTimeout)
		},
		Timeout: viper.GetDuration("api.handler.timeout"),
	}))

	// Configure routes
	configureAllRoutes(echoInstance)

	port := viper.GetInt("api.host.port")
	logger.Info(ctx, "start api", fmt.Sprintf("Started at %d", port), nil)
	route := fmt.Sprintf(":%d", port)
	echoInstance.Logger.Fatal(echoInstance.Start(route))
}

func configureAllRoutes(echoInstance *echo.Echo) {
	// Metrics by Prometheus
	configureMetrics(echoInstance)

	// Healthcheck
	configureHealthcheckRoute(echoInstance)

	// User routes
	configureUserRoutes(echoInstance)
}

func configureHealthcheckRoute(echoInstance *echo.Echo) {
	echoInstance.GET("/healthcheck", handlers.Healthcheck)
}

func configureUserRoutes(echoInstance *echo.Echo) {
	userRepository := new(repository.UserMongoRepository)
	userHandler := handlers.NewUserHandler(userRepository)

	userGroup := echoInstance.Group("/user")
	userGroup.POST("", userHandler.CreateUser)
}

// configureMetrics is the method that configures Prometheus metrics' route.
func configureMetrics(echoInstance *echo.Echo) {
	echoInstance.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}

// configureNewRelic is the method that enables New Relic.
func configureNewRelic(ctx context.Context, echoInstance *echo.Echo) {
	newRelicEnable := viper.GetBool("new_relic.is.enabled")
	if !newRelicEnable {
		return
	}

	newRelicApp, err := createNewRelicApp(ctx)
	if err != nil {
		logger.Error(ctx, "enabling New Relic", err.Error(), nil)
	}

	echoInstance.Use(apiMiddleware.NewRelicMiddleware(newRelicApp))
	logger.Info(ctx, "enabling New Relic", "success", nil)
}

// createNewRelicApp is the method that creates New Relic configuration.
func createNewRelicApp(ctx context.Context) (newrelic.Application, error) {
	licenseKeyEnvVar := viper.GetString("new_relic.licence.key")

	config := newrelic.NewConfig("My Awesome API", licenseKeyEnvVar)
	proxyURL, err := url.Parse(viper.GetString("new_relic.proxy.url"))
	if err != nil {
		logger.Error(ctx, "parse proxy url from env var", err.Error(), nil)
	}

	config.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	newRelicApp, err := newrelic.NewApplication(config)
	if err != nil {
		logger.Error(ctx, "create New Relic APP ", err.Error(), nil)
		return nil, err
	}

	return newRelicApp, nil
}

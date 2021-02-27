package cmd

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo"
	"github.com/michelaquino/golang_api_skeleton/src/handlers"
	"github.com/michelaquino/golang_api_skeleton/src/repository"
	newrelic "github.com/newrelic/go-agent"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	apiContext "github.com/michelaquino/golang_api_skeleton/src/context"
	apiMiddleware "github.com/michelaquino/golang_api_skeleton/src/middleware"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Starts API service",
	Long:  `Starts API service.`,
	Run: func(cmd *cobra.Command, args []string) {
		// config.Init()
		start()
	},
}

func start() {
	logger := apiContext.GetLogger()
	echoInstance := echo.New()

	// Configure New Relic
	configureNewRelic(echoInstance)

	// Middlewares
	echoInstance.Use(apiMiddleware.RequestLogDataMiddleware())

	// Configure routes
	configureAllRoutes(echoInstance)

	port := viper.GetInt("api.host.port")
	logger.Info("Main", "main", "", "", "", fmt.Sprintf("Started at %d", port), "")
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
func configureNewRelic(echoInstance *echo.Echo) {
	logger := apiContext.GetLogger()

	newRelicEnable := viper.GetBool("new_relic.is.enabled")
	if newRelicEnable {
		if newRelicApp, err := createNewRelicApp(); err == nil {
			logger.Info("Main", "configureNewRelic", "", "", "Enabling New Relic", "Success", "New Relic ENABLED")
			echoInstance.Use(apiMiddleware.NewRelicMiddleware(newRelicApp))
		} else {
			logger.Error("Main", "configureNewRelic", "", "", "Enabling New Relic", "Error", err.Error())
		}

		return
	}

	logger.Info("Main", "configureNewRelic", "", "", "Enabling New Relic", "Success", "New Relic DISABLED")
}

// createNewRelicApp is the method that creates New Relic configuration.
func createNewRelicApp() (newrelic.Application, error) {
	log := apiContext.GetLogger()
	licenseKeyEnvVar := viper.GetString("new_relic.licence.key")

	config := newrelic.NewConfig("My Awesome API", licenseKeyEnvVar)
	proxyURL, err := url.Parse(viper.GetString("new_relic.proxy.url"))
	if err != nil {
		log.Error("Main", "createNewRelicApp", "", "", "Parse proxy url from env var", "Error", err.Error())
	}

	config.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	newRelicApp, err := newrelic.NewApplication(config)
	if err != nil {
		log.Error("Main", "createNewRelicApp", "", "", "Create New Relic APP ", "Error", err.Error())
		return nil, err
	}

	return newRelicApp, nil
}

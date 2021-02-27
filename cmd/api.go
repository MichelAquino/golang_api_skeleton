package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/labstack/echo"
	"github.com/michelaquino/golang_api_skeleton/src/handlers"
	"github.com/michelaquino/golang_api_skeleton/src/repository"
	newrelic "github.com/newrelic/go-agent"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"

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
	apiConfig := apiContext.GetAPIConfig()
	echoInstance := echo.New()

	// Configure New Relic
	configureNewRelic(echoInstance)

	// Middlewares
	echoInstance.Use(apiMiddleware.RequestLogDataMiddleware())

	// Configure routes
	configureAllRoutes(echoInstance)

	logger.Info("Main", "main", "", "", "", fmt.Sprintf("Started at %d", apiConfig.HostPort), "")
	route := fmt.Sprintf(":%d", apiConfig.HostPort)
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

	newRelicEnvVar := os.Getenv("ENABLE_NEW_RELIC")
	newRelicEnable, err := strconv.ParseBool(newRelicEnvVar)

	if err == nil && newRelicEnable {
		if newRelicApp, err := createNewRelicApp(); err == nil {
			logger.Info("Main", "configureNewRelic", "", "", "Enabling New Relic", "Success", "New Relic ENABLED")
			echoInstance.Use(apiMiddleware.NewRelicMiddleware(newRelicApp))
		} else {
			logger.Error("Main", "configureNewRelic", "", "", "Enabling New Relic", "Error", err.Error())
		}

		return
	}

	if err != nil {
		logger.Error("Main", "configureNewRelic", "", "", "Enabling New Relic", "Error", err.Error())
	}

	logger.Info("Main", "configureNewRelic", "", "", "Enabling New Relic", "Success", "New Relic DISABLED")
}

// createNewRelicApp is the method that creates New Relic configuration.
func createNewRelicApp() (newrelic.Application, error) {
	log := apiContext.GetLogger()
	licenseKeyEnvVar := os.Getenv("NEW_RELIC_LICENSE_KEY")

	config := newrelic.NewConfig("My Awesome API", licenseKeyEnvVar)
	proxyURL, err := url.Parse(os.Getenv("NEW_RELIC_PROXY_URL"))
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

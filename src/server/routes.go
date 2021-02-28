package server

import (
	"github.com/labstack/echo/v4"
	"github.com/michelaquino/golang_api_skeleton/src/handlers"
	"github.com/michelaquino/golang_api_skeleton/src/repository"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func configureRoutes(echoInstance *echo.Echo) {
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

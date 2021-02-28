package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	apiMiddleware "github.com/michelaquino/golang_api_skeleton/src/middleware"
	"github.com/spf13/viper"
)

func configureMiddlewares(echoInstance *echo.Echo) {
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
}

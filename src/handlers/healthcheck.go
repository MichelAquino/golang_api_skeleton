package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/michelaquino/golang_api_skeleton/src/log"
)

var (
	logger = log.GetLogger()
)

// Healthcheck is a method that responds only WORKING.
func Healthcheck(echoContext echo.Context) error {
	logger.Info(echoContext.Request().Context(), "healthcheck", "success", nil)
	return echoContext.String(http.StatusOK, "WORKING")
}

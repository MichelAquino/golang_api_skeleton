package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/michelaquino/golang_api_skeleton/src/log"
)

var (
	logger = log.GetLogger()
)

// RequestLogger is a middleware to log the request data
func RequestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := c.Request()
		response := c.Response()

		start := time.Now()
		if err := next(c); err != nil {
			logger.Error(c.Request().Context(), "request logger middleware", err.Error(), nil)
			c.Error(err)
		}
		stop := time.Now()

		extraFields := map[string]string{
			"method":       request.Method,
			"request_uri":  request.RequestURI,
			"status":       strconv.Itoa(response.Status),
			"request_time": stop.Sub(start).String(),
		}

		logger.Info(c.Request().Context(), "request logger middleware", "success", extraFields)
		return nil
	}
}

package middleware

import (
	"github.com/labstack/echo"
	"github.com/michelaquino/golang_api_skeleton/src/models"
	uuid "github.com/satori/go.uuid"
)

// RequestIDKey is the key to set request ID on context
const RequestIDKey = "requestLogDataContextKey"

// RequestLogDataMiddleware is a middleware to send request info to New Relic
func RequestLogDataMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echoContext echo.Context) error {
			requestID := echoContext.Request().Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = uuid.NewV4().String()
			}

			requestLogData := models.RequestLogData{
				ID:       requestID,
				OriginIP: echoContext.Request().RemoteAddr,
			}

			echoContext.Set(RequestIDKey, requestLogData)
			return next(echoContext)
		}
	}
}

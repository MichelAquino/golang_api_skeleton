package middleware

import (
	"bytes"
	"io/ioutil"

	"github.com/labstack/echo/v4"
	newrelic "github.com/newrelic/go-agent"
)

// NewRelicMiddleware is a middleware to send request info to New Relic.
func NewRelicMiddleware(newrelicApp newrelic.Application) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echoContext echo.Context) error {
			responseWriter := echoContext.Response().Writer

			// Copy struct request to remove body.
			requestCopy := *echoContext.Request()

			// Set body empty to send to New Relic
			requestCopy.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))

			transaction := newrelicApp.StartTransaction(requestCopy.URL.Path, responseWriter, &requestCopy)
			defer transaction.End()

			return next(echoContext)
		}
	}
}

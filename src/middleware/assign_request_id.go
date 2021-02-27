package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/michelaquino/golang_api_skeleton/src/context"
	uuid "github.com/satori/go.uuid"
)

// AssignRequestID is a middleware to set a request_id to response (as a header)
// and to request (in the underlying context). If the value is not found in
// X-Request-ID header from request the identifier will be generated. The
// request_id will be used by http clients in external requests
func AssignRequestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()

		reqID := req.Header.Get(echo.HeaderXRequestID)
		if len(reqID) == 0 {
			reqID = uuid.NewV4().String()
		}

		ctx := context.SetRequestID(req.Context(), reqID)
		reqWithContext := req.WithContext(ctx)
		c.SetRequest(reqWithContext)

		res.Header().Set(echo.HeaderXRequestID, reqID)

		return next(c)
	}
}

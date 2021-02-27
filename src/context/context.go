package context

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type httpRequestIDKey string

const requestIDKey httpRequestIDKey = "X-Request-ID"

// SetRequestID is a method that adds to context an X-Request-ID value
func SetRequestID(ctx context.Context, reqID string) context.Context {
	if reqID == "" {
		return ctx
	}

	if ctx.Value(requestIDKey) != nil {
		return ctx
	}

	return context.WithValue(ctx, requestIDKey, reqID)
}

// GetRequestID returns a requestID value based on context
func GetRequestID(ctx context.Context) string {
	if requestID := ctx.Value(requestIDKey); requestID != nil {
		return fmt.Sprintf("%s", requestID)
	}

	return uuid.New().String()
}

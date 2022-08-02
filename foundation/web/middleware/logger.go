package middleware

import (
	"context"
	"github.com/mchusovlianov/geodata/foundation/web"
	"go.uber.org/zap"
	"net/http"
)

// Logger writes some information about the request to the logs in the
// format: TraceID : (200) GET /foo -> IP ADDR (latency)
func Logger(log *zap.SugaredLogger) web.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			log.Infow("request", "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr)

			// Call the next handler.
			err := handler(ctx, w, r)

			// Return the error so it can be handled further up the chain.
			return err
		}

		return h
	}

	return m
}

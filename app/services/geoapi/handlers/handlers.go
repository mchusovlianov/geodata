// Package handlers manages the different versions of the API.
package handlers

import (
	"context"
	"github.com/jmoiron/sqlx"
	v1 "github.com/mchusovlianov/geodata/app/services/geoapi/handlers/v1"
	"github.com/mchusovlianov/geodata/foundation/web"
	"github.com/mchusovlianov/geodata/foundation/web/middleware"
	"go.uber.org/zap"
	"net/http"
)

// Options represent optional parameters.
type Options struct {
	corsOrigin string
}

// WithCORS provides configuration options for CORS.
func WithCORS(origin string) func(opts *Options) {
	return func(opts *Options) {
		opts.corsOrigin = origin
	}
}

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Log *zap.SugaredLogger
	DB  *sqlx.DB
}

// APIMux constructs a http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig, options ...func(opts *Options)) http.Handler {
	var opts Options
	for _, option := range options {
		option(&opts)
	}

	// Construct the web.App which holds all routes as well as common Middleware.
	var app *web.App

	// Do we need CORS?
	if opts.corsOrigin != "" {
		app = web.NewApp(
			middleware.Logger(cfg.Log),
			middleware.Errors(cfg.Log),
			middleware.Cors(opts.corsOrigin),
			middleware.Panics(),
		)

		// Accept CORS 'OPTIONS' preflight requests if config has been provided.
		// Don't forget to apply the CORS middleware to the routes that need it.
		// Example Config: `conf:"default:https://MY_DOMAIN.COM"`
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			return nil
		}
		app.Handle(http.MethodOptions, "", "/*", h, middleware.Cors(opts.corsOrigin))
	}

	if app == nil {
		app = web.NewApp(
			middleware.Logger(cfg.Log),
			middleware.Errors(cfg.Log),
			middleware.Panics(),
		)
	}

	// Load the v1 routes.
	v1.Routes(app, v1.Config{
		Log: cfg.Log,
		DB:  cfg.DB,
	})

	return app
}

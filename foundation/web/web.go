// Package web contains a small web framework extension.
package web

import (
	"context"
	"github.com/dimfeld/httptreemux/v5"
	"net/http"
)

// A Handler is a type that handles a http request within our own little mini
// framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	mux *httptreemux.ContextMux
	mw  []Middleware
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(mw ...Middleware) *App {
	mux := httptreemux.NewContextMux()
	return &App{
		mux: mux,
		mw:  mw,
	}
}

// ServeHTTP implements the http.Handler interface
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *App) Handle(method string, group string, path string, handler Handler, mw ...Middleware) {

	// First wrap handler specific middleware around this handler.
	handler = wrapMiddleware(mw, handler)

	// Add the application's general middleware to the handler chain.
	handler = wrapMiddleware(a.mw, handler)

	// The function to execute for each request.
	h := func(w http.ResponseWriter, r *http.Request) {

		// Pull the context from the request and
		// use it as a separate parameter.
		ctx := r.Context()

		// Call the wrapped handler functions.
		if err := handler(ctx, w, r); err != nil {
			return
		}
	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}
	a.mux.Handle(method, finalPath, h)
}

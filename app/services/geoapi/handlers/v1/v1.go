// Package v1 contains the full set of handler functions and routes
// supported by the v1 web api.
package v1

import (
	"github.com/mchusovlianov/geodata/app/services/geoapi/handlers/v1/locationgrp"
	"github.com/mchusovlianov/geodata/business/core/city"
	"github.com/mchusovlianov/geodata/business/core/country"
	"github.com/mchusovlianov/geodata/business/core/location"
	"github.com/mchusovlianov/geodata/foundation/web"
	"net/http"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log *zap.SugaredLogger
	DB  *sqlx.DB
}

// Routes binds all the version 1 routes.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	// Register user management and authentication endpoints.
	loch := locationgrp.Handlers{
		Location: location.NewCore(cfg.Log, cfg.DB, nil),
		Country:  country.NewCore(cfg.Log, cfg.DB, nil),
		City:     city.NewCore(cfg.Log, cfg.DB, nil),
	}
	app.Handle(http.MethodGet, version, "/location/:ip", loch.QueryByIP)
}

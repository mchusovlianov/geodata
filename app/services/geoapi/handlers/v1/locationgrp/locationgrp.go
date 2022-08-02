package locationgrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/mchusovlianov/geodata/business/core/city"
	"github.com/mchusovlianov/geodata/business/core/country"
	"github.com/mchusovlianov/geodata/business/core/location"
	"github.com/mchusovlianov/geodata/foundation/web"
	"net/http"
)

// Handlers manages the set of location endpoints.
type Handlers struct {
	Location location.Core
	Country  country.Core
	City     city.Core
}

// LocationResponse
type LocationResponse struct {
	Location location.Location `json:"location"`
	Country  country.Country   `json:"country"`
	City     city.City         `json:"city"`
}

// QueryByID returns a location by its IP.
func (h Handlers) QueryByIP(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ip := web.Param(r, "ip")

	loc, err := h.Location.QueryByIP(ctx, ip)
	if err != nil {
		switch {
		case errors.Is(err, location.ErrNotFound):
			return web.NewRequestError(err, http.StatusNotFound)
		case errors.Is(err, location.ErrValidation):
			return web.NewRequestError(err, http.StatusBadRequest)

		default:
			return fmt.Errorf("IP[%s]: %w", ip, err)
		}
	}

	cit, err := h.City.QueryByUUID(ctx, loc.CityUUID)
	if err != nil {
		switch {
		case errors.Is(err, city.ErrNotFound):
			return web.NewRequestError(err, http.StatusNotFound)
		case errors.Is(err, city.ErrValidation):
			return web.NewRequestError(err, http.StatusBadRequest)

		default:
			return fmt.Errorf("UUID[%s]: %w", loc.CityUUID, err)
		}
	}

	countr, err := h.Country.QueryByUUID(ctx, cit.CountryUUID)
	if err != nil {
		switch {
		case errors.Is(err, country.ErrNotFound):
			return web.NewRequestError(err, http.StatusNotFound)
		case errors.Is(err, country.ErrValidation):
			return web.NewRequestError(err, http.StatusBadRequest)

		default:
			return fmt.Errorf("UUID[%s]: %w", cit.CountryUUID, err)
		}
	}

	return web.Respond(ctx, w, LocationResponse{
		Location: loc,
		Country:  countr,
		City:     cit,
	}, http.StatusOK)
}

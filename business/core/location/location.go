// Package location provides a core business API. Right now these
// calls are just wrapping the data/store layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package location

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mchusovlianov/geodata/business/core/location/db"
	"github.com/mchusovlianov/geodata/foundation/database"
	"go.uber.org/zap"
	"time"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound   = errors.New("location not found")
	ErrValidation = errors.New("validation failed")
)

// Core manages the set of APIs for location access.
type Core struct {
	store db.Store
}

// NewCore constructs a core for location api access.
func NewCore(log *zap.SugaredLogger, dbConn *sqlx.DB, tx *sqlx.Tx) Core {
	return Core{
		store: db.NewStore(log, dbConn, tx),
	}
}

// Create adds a Location to the database. It returns the created Location with
// fields like ID and DateCreated populated.
func (c Core) Create(ctx context.Context, location NewLocation, now time.Time) (Location, error) {
	validate := validator.New()
	err := validate.Struct(location)
	if err != nil {
		return Location{}, ErrValidation
	}

	dbLocation := db.Location{
		UUID:         uuid.New().String(),
		CityUUID:     location.CityUUID,
		Longitude:    location.Longitude,
		Latitude:     location.Latitude,
		IP:           location.IP,
		MysteryValue: location.MysteryValue,
		DateCreated:  now,
		DateUpdated:  now,
	}

	if err := c.store.Create(ctx, dbLocation); err != nil {
		return Location{}, fmt.Errorf("create: %w", err)
	}

	return toLocation(dbLocation), nil
}

// QueryByUUID gets the specified location from the database by uuid.
func (c Core) QueryByUUID(ctx context.Context, locationUUID string) (Location, error) {
	dbLocation, err := c.store.QueryByUUID(ctx, locationUUID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Location{}, ErrNotFound
		}
		return Location{}, fmt.Errorf("query: %w", err)
	}

	return toLocation(dbLocation), nil
}

// QueryByIP gets the specified location from the database by ip-address.
func (c Core) QueryByIP(ctx context.Context, ip string) (Location, error) {
	validate := validator.New()
	err := validate.Struct(struct {
		IP string `validate:"required,ip"`
	}{IP: ip})
	if err != nil {
		return Location{}, ErrValidation
	}

	dbLocation, err := c.store.QueryByIP(ctx, ip)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Location{}, ErrNotFound
		}
		return Location{}, fmt.Errorf("query: %w", err)
	}

	return toLocation(dbLocation), nil
}

// QueryAll gets all locations from the database.
func (c Core) QueryAll(ctx context.Context) ([]Location, error) {
	dbLocations, err := c.store.QueryAll(ctx)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return []Location{}, ErrNotFound
		}
		return []Location{}, fmt.Errorf("query: %w", err)
	}

	return toLocationSlice(dbLocations), nil
}

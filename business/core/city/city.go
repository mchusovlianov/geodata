// Package location provides a core business API. Right now these
// calls are just wrapping the data/store layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package city

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mchusovlianov/geodata/business/core/city/db"
	"github.com/mchusovlianov/geodata/foundation/database"
	"go.uber.org/zap"
	"time"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound   = errors.New("city not found")
	ErrValidation = errors.New("validation failed")
	ErrDuplicate  = errors.New("city already exist")
)

// Core manages the set of APIs for city access.
type Core struct {
	store db.Store
}

// NewCore constructs a core for city api access.
func NewCore(log *zap.SugaredLogger, dbConn *sqlx.DB, tx *sqlx.Tx) Core {
	return Core{
		store: db.NewStore(log, dbConn, tx),
	}
}

// Create adds a City to the database. It returns the created City with
// fields like ID and DateCreated populated.
func (c Core) Create(ctx context.Context, city NewCity, now time.Time) (City, error) {
	validate := validator.New()
	err := validate.Struct(city)
	if err != nil {
		return City{}, ErrValidation
	}

	dbCity := db.City{
		UUID:        uuid.New().String(),
		CountryUUID: city.CountryUUID,
		Name:        city.Name,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := c.store.Create(ctx, dbCity); err != nil {
		if errors.Is(err, database.ErrDBDuplicatedEntry) {
			return City{}, ErrDuplicate
		}

		return City{}, fmt.Errorf("create: %w", err)
	}

	return toCity(dbCity), nil
}

// QueryByUUID gets the specified city from the database.
func (c Core) QueryByUUID(ctx context.Context, cityUUID string) (City, error) {
	dbCity, err := c.store.QueryByUUID(ctx, cityUUID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return City{}, ErrNotFound
		}
		return City{}, fmt.Errorf("query: %w", err)
	}

	return toCity(dbCity), nil
}

// QueryByCountryUUID gets the all cities with same country uuid from the database.
func (c Core) QueryByCountryUUID(ctx context.Context, countryUUID string) ([]City, error) {
	dbCities, err := c.store.QueryByCountryUUID(ctx, countryUUID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return []City{}, ErrNotFound
		}
		return []City{}, fmt.Errorf("query: %w", err)
	}

	return toCitySlice(dbCities), nil
}

// QueryAll gets all cities from the database.
func (c Core) QueryAll(ctx context.Context) ([]City, error) {
	dbCities, err := c.store.QueryAll(ctx)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return []City{}, ErrNotFound
		}
		return []City{}, fmt.Errorf("query: %w", err)
	}

	return toCitySlice(dbCities), nil
}

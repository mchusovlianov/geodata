// Package country provides a core business API. Right now these
// calls are just wrapping the data/store layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package country

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mchusovlianov/geodata/business/core/country/db"
	"github.com/mchusovlianov/geodata/foundation/database"
	"go.uber.org/zap"
	"time"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound   = errors.New("country not found")
	ErrValidation = errors.New("validation failed")
	ErrDuplicate  = errors.New("country already exist")
)

// Core manages the set of APIs for country access.
type Core struct {
	store db.Store
}

// NewCore constructs a core for country api access.
func NewCore(log *zap.SugaredLogger, dbConn *sqlx.DB, tx *sqlx.Tx) Core {
	return Core{
		store: db.NewStore(log, dbConn, tx),
	}
}

// Create adds a Country to the database. It returns the created Country with
// fields like ID and DateCreated populated.
func (c Core) Create(ctx context.Context, country NewCountry, now time.Time) (Country, error) {
	validate := validator.New()
	err := validate.Struct(country)
	if err != nil {
		return Country{}, ErrValidation
	}

	dbCountry := db.Country{
		UUID:        uuid.New().String(),
		Code:        country.Code,
		Name:        country.Name,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := c.store.Create(ctx, dbCountry); err != nil {
		if errors.Is(err, database.ErrDBDuplicatedEntry) {
			return Country{}, ErrDuplicate
		}

		return Country{}, fmt.Errorf("create: %w", err)
	}

	return toCountry(dbCountry), nil
}

// QueryByUUID gets the specified country from the database.
func (c Core) QueryByUUID(ctx context.Context, countryUUID string) (Country, error) {
	dbCountry, err := c.store.QueryByUUID(ctx, countryUUID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Country{}, ErrNotFound
		}
		return Country{}, fmt.Errorf("query: %w", err)
	}

	return toCountry(dbCountry), nil
}

// QueryByCode gets the specified country from the database by country code.
func (c Core) QueryByCode(ctx context.Context, code string) (Country, error) {
	dbCountry, err := c.store.QueryByCode(ctx, code)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Country{}, ErrNotFound
		}
		return Country{}, fmt.Errorf("query: %w", err)
	}

	return toCountry(dbCountry), nil
}

// QueryAll gets all countries from the database.
func (c Core) QueryAll(ctx context.Context) ([]Country, error) {
	dbCountries, err := c.store.QueryAll(ctx)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return []Country{}, ErrNotFound
		}
		return []Country{}, fmt.Errorf("query: %w", err)
	}

	return toCountrySlice(dbCountries), nil
}

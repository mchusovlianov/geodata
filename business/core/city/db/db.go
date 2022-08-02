package db

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mchusovlianov/geodata/foundation/database"
	"go.uber.org/zap"
)

type Store struct {
	log *zap.SugaredLogger
	db  *sqlx.DB
	tx  *sqlx.Tx
}

// NewStore constructs a data for api access.
func NewStore(log *zap.SugaredLogger, db *sqlx.DB, tx *sqlx.Tx) Store {
	return Store{
		log: log,
		db:  db,
		tx:  tx,
	}
}

// getConn returns a required execution context: transaction or database connection
func (s Store) getConn() sqlx.ExtContext {
	if s.tx != nil {
		return s.tx
	}

	return s.db
}

// Create adds a City to the database. It returns the created City with
// fields like ID and DateCreated populated.
func (s Store) Create(ctx context.Context, city City) error {
	const q = `
	INSERT INTO cities
		(uuid, country_uuid, name, date_created, date_updated)
	VALUES
		(:uuid, :country_uuid, :name, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.getConn(), q, city); err != nil {
		return fmt.Errorf("inserting city: %w", err)
	}

	return nil
}

// QueryByUUID gets the specified city from the database.
func (s Store) QueryByUUID(ctx context.Context, cityUUID string) (City, error) {
	data := struct {
		UUID string `db:"uuid"`
	}{
		UUID: cityUUID,
	}

	const q = `
	SELECT
		*
	FROM
		cities
	WHERE 
		uuid = :uuid`

	var city City
	if err := database.NamedQueryStruct(ctx, s.getConn(), q, data, &city); err != nil {
		return City{}, fmt.Errorf("selecting cityUUID[%q]: %w", cityUUID, err)
	}

	return city, nil
}

// QueryByCountryUUID gets the all cities with same country uuid from the database.
func (s Store) QueryByCountryUUID(ctx context.Context, countryUUID string) ([]City, error) {
	data := struct {
		CountryUUID string `db:"country_uuid"`
	}{
		CountryUUID: countryUUID,
	}

	const q = `
	SELECT
		*
	FROM
		cities
	WHERE 
		country_uuid = :country_uuid`

	var cities []City
	if err := database.NamedQuerySlice(ctx, s.getConn(), q, data, &cities); err != nil {
		return []City{}, fmt.Errorf("selecting countryUUID[%q]: %w", countryUUID, err)
	}

	return cities, nil
}

// QueryAll gets all cities from the database.
func (s Store) QueryAll(ctx context.Context) ([]City, error) {
	const q = `
	SELECT
		*
	FROM
		cities`

	var cities []City
	if err := database.NamedQuerySlice(ctx, s.getConn(), q, struct{}{}, &cities); err != nil {
		return []City{}, fmt.Errorf("selecting all cities: %w", err)
	}

	return cities, nil
}

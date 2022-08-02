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

// Create adds a Location to the database. It returns the created Location with
// fields like ID and DateCreated populated.
func (s Store) Create(ctx context.Context, location Location) error {
	const q = `
	INSERT INTO locations
		(uuid, city_uuid, mystery_value, ip, latitude, longitude, date_created, date_updated)
	VALUES
		(:uuid, :city_uuid, :mystery_value, :ip, :latitude, :longitude, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.getConn(), q, location); err != nil {
		return fmt.Errorf("inserting location: %w", err)
	}

	return nil
}

// QueryByUUID gets the specified location from the database by uuid.
func (s Store) QueryByUUID(ctx context.Context, locationUUID string) (Location, error) {
	data := struct {
		UUID string `db:"uuid"`
	}{
		UUID: locationUUID,
	}

	const q = `
	SELECT
		*
	FROM
		locations
	WHERE 
		uuid = :uuid`

	var location Location
	if err := database.NamedQueryStruct(ctx, s.getConn(), q, data, &location); err != nil {
		return Location{}, fmt.Errorf("selecting locationUUID[%q]: %w", locationUUID, err)
	}

	return location, nil
}

// QueryByIP gets the specified location from the database by ip-address.
func (s Store) QueryByIP(ctx context.Context, ip string) (Location, error) {
	data := struct {
		IP string `db:"ip"`
	}{
		IP: ip,
	}

	const q = `
	SELECT
		*
	FROM
		locations
	WHERE 
		ip = :ip`

	var location Location
	if err := database.NamedQueryStruct(ctx, s.getConn(), q, data, &location); err != nil {
		return Location{}, fmt.Errorf("selecting locationIP[%q]: %w", ip, err)
	}

	return location, nil
}

// QueryAll gets all countries from the database.
func (s Store) QueryAll(ctx context.Context) ([]Location, error) {
	const q = `
	SELECT
		*
	FROM
		locations`

	var locations []Location
	if err := database.NamedQuerySlice(ctx, s.getConn(), q, struct{}{}, &locations); err != nil {
		return []Location{}, fmt.Errorf("selecting all locations: %w", err)
	}

	return locations, nil
}

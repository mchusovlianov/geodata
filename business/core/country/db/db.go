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

// Create adds a Country to the database. It returns the created Country with
// fields like ID and DateCreated populated.
func (s Store) Create(ctx context.Context, country Country) error {
	const q = `
	INSERT INTO countries
		(uuid, code, name, date_created, date_updated)
	VALUES
		(:uuid, :code, :name, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.getConn(), q, country); err != nil {
		return fmt.Errorf("inserting country: %w", err)
	}

	return nil
}

// QueryByUUID gets the specified country from the database.
func (s Store) QueryByUUID(ctx context.Context, countryUUID string) (Country, error) {
	data := struct {
		UUID string `db:"uuid"`
	}{
		UUID: countryUUID,
	}

	const q = `
	SELECT
		*
	FROM
		countries
	WHERE 
		uuid = :uuid`

	var country Country
	if err := database.NamedQueryStruct(ctx, s.getConn(), q, data, &country); err != nil {
		return Country{}, fmt.Errorf("selecting countryUUID[%q]: %w", countryUUID, err)
	}

	return country, nil
}

// QueryByCode gets the specified country from the database by code.
func (s Store) QueryByCode(ctx context.Context, code string) (Country, error) {
	data := struct {
		Code string `db:"code"`
	}{
		Code: code,
	}

	const q = `
	SELECT
		*
	FROM
		countries
	WHERE 
		code = :code`

	var country Country
	if err := database.NamedQueryStruct(ctx, s.getConn(), q, data, &country); err != nil {
		return Country{}, fmt.Errorf("selecting countryCODE[%q]: %w", code, err)
	}

	return country, nil
}

// QueryAll gets all countries from the database.
func (s Store) QueryAll(ctx context.Context) ([]Country, error) {
	const q = `
	SELECT
		*
	FROM
		countries`

	var countries []Country
	if err := database.NamedQuerySlice(ctx, s.getConn(), q, struct{}{}, &countries); err != nil {
		return []Country{}, fmt.Errorf("selecting all countries: %w", err)
	}

	return countries, nil
}

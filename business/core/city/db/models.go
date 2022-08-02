// Package db is a package for keeping all db-related logic
package db

import (
	"time"
)

type City struct {
	UUID        string    `db:"uuid"`         // Unique identifier.
	CountryUUID string    `db:"country_uuid"` // Unique identifier of the linked country.
	Name        string    `db:"name"`         // Name of the city.
	DateCreated time.Time `db:"date_created"` // When the city was added.
	DateUpdated time.Time `db:"date_updated"` // When the city was last modified.
}

// Package db is a package for keeping all db-related logic
package db

import (
	"time"
)

type Location struct {
	UUID         string    `db:"uuid"`          // Unique identifier.
	CityUUID     string    `db:"city_uuid"`     // Unique identifier of the linked city.
	IP           string    `db:"ip"`            // IP of the location.
	Latitude     float64   `db:"latitude"`      // Latitude of the location.
	Longitude    float64   `db:"longitude"`     // Longitude of the location.
	MysteryValue int64     `db:"mystery_value"` // Mystery value of the location.
	DateCreated  time.Time `db:"date_created"`  // When the location was added.
	DateUpdated  time.Time `db:"date_updated"`  // When the location was last modified.
}

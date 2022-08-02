// Package db is a package for keeping all db-related logic
package db

import (
	"time"
)

type Country struct {
	UUID        string    `db:"uuid"`         // Unique identifier.
	Name        string    `db:"name"`         // Name of the country.
	Code        string    `db:"code"`         // Code of the country.
	DateCreated time.Time `db:"date_created"` // When the country was added.
	DateUpdated time.Time `db:"date_updated"` // When the country  was last modified.
}

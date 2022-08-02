package city

import (
	"github.com/mchusovlianov/geodata/business/core/city/db"
	"time"
	"unsafe"
)

// City
type City struct {
	UUID        string    `json:"uuid"`         // Unique identifier.
	CountryUUID string    `json:"country_uuid"` // Unique identifier of the linked country.
	Name        string    `json:"name"`         // Name of city.
	DateCreated time.Time `json:"date_created"` // When the city was added.
	DateUpdated time.Time `json:"date_updated"` // When the city was last modified.
}

// NewCity is what we require from clients when adding a City.
type NewCity struct {
	Name        string `json:"name" validate:"required"`
	CountryUUID string `json:"country_uuid" validate:"required"`
}

func toCity(dbCity db.City) City {
	ci := (*City)(unsafe.Pointer(&dbCity))
	return *ci
}

func toCitySlice(dbCities []db.City) []City {
	cis := make([]City, len(dbCities))
	for i, dbCit := range dbCities {
		cis[i] = toCity(dbCit)
	}
	return cis
}

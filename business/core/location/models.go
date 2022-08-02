package location

import (
	"github.com/mchusovlianov/geodata/business/core/location/db"
	"time"
	"unsafe"
)

// Location
type Location struct {
	UUID         string    `json:"uuid"`          // Unique identifier.
	CityUUID     string    `json:"city_uuid"`     // Unique identifier of the linked city.
	IP           string    `json:"ip"`            // IP of the location.
	Latitude     float64   `json:"latitude"`      // Latitude of the location.
	Longitude    float64   `json:"longitude"`     // Longitude of the location.
	MysteryValue int64     `json:"mystery_value"` // MysteryValue of the location.
	DateCreated  time.Time `json:"date_created"`  // When the location was added.
	DateUpdated  time.Time `json:"date_updated"`  // When the location was last modified.
}

// NewLocation is what we require from clients when adding a Location.
type NewLocation struct {
	IP           string  `json:"ip" validate:"required,ip"`
	Longitude    float64 `json:"longitude" validate:"required"`
	Latitude     float64 `json:"latitude" validate:"required"`
	MysteryValue int64   `json:"mystery_value" validate:"required"`
	CityUUID     string  `json:"city_uuid" validate:"required"`
}

func toLocation(dbLocation db.Location) Location {
	ci := (*Location)(unsafe.Pointer(&dbLocation))
	return *ci
}

func toLocationSlice(dbCities []db.Location) []Location {
	cis := make([]Location, len(dbCities))
	for i, dbCit := range dbCities {
		cis[i] = toLocation(dbCit)
	}
	return cis
}

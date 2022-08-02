package country

import (
	"github.com/mchusovlianov/geodata/business/core/country/db"
	"time"
	"unsafe"
)

// Country
type Country struct {
	UUID        string    `json:"uuid"`         // Unique identifier.
	Name        string    `json:"name"`         // Name of country.
	Code        string    `json:"code"`         // Code of country.
	DateCreated time.Time `json:"date_created"` // When the country was added.
	DateUpdated time.Time `json:"date_updated"` // When the country was last modified.
}

// NewCountry is what we require from clients when adding a Country.
type NewCountry struct {
	Name string `json:"name" validate:"required"`
	Code string `json:"code" validate:"required,alpha"`
}

func toCountry(dbCountry db.Country) Country {
	ci := (*Country)(unsafe.Pointer(&dbCountry))
	return *ci
}

func toCountrySlice(dbCities []db.Country) []Country {
	cis := make([]Country, len(dbCities))
	for i, dbCit := range dbCities {
		cis[i] = toCountry(dbCit)
	}
	return cis
}

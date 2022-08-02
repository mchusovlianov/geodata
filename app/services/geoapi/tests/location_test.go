package tests

import (
	"encoding/json"
	"fmt"
	"github.com/mchusovlianov/geodata/app/services/geoapi/handlers"
	"github.com/mchusovlianov/geodata/app/services/geoapi/handlers/v1/locationgrp"
	"github.com/mchusovlianov/geodata/business/data/tests"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// LocationTests holds methods for each location subtest. This type allows passing
// dependencies for tests while still providing a convenient syntax when
// subtests are registered.
type LocationTests struct {
	app http.Handler
}

// TestLocations is the entry point for testing location management functions.
func TestLocations(t *testing.T) {
	test := tests.NewIntegration(
		t,
		tests.DBContainer{
			Image:  "percona",
			Port:   "3306",
			Name:   "geodatatest",
			IsSeed: true,
			Args:   []string{"-e", "MYSQL_ROOT_PASSWORD=root"},
		},
	)
	t.Cleanup(test.Teardown)

	tests := LocationTests{
		app: handlers.APIMux(handlers.APIMuxConfig{
			Log: test.Log,
			DB:  test.DB,
		}),
	}

	t.Run("getLocation400", tests.getLocation400)
	t.Run("getLocation404", tests.getLocation404)
	t.Run("getLocation200", tests.getLocation200)
}

// getLocation400 validates a location request for a malformed ip.
func (lt *LocationTests) getLocation400(t *testing.T) {
	ip := "12345"

	r := httptest.NewRequest(http.MethodGet, "/v1/location/"+ip, nil)
	w := httptest.NewRecorder()

	lt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a location with a malformed ip.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new location %s.", testID, ip)
		{
			fmt.Println(w.Body.String())
			if w.Code != http.StatusBadRequest {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 400 for the response : %v", tests.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 400 for the response.", tests.Success, testID)
		}
	}
}

// getLocation404 valuuidates a location request for a location that does not exist with the endpoint.
func (lt *LocationTests) getLocation404(t *testing.T) {
	ip := "123.12.12.2"

	r := httptest.NewRequest(http.MethodGet, "/v1/location/"+ip, nil)
	w := httptest.NewRecorder()

	lt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a location with an unknown ip.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new location %s.", testID, ip)
		{
			if w.Code != http.StatusNotFound {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 404 for the response : %v", tests.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 404 for the response.", tests.Success, testID)
		}
	}
}

// getLocation200 validates a location request for an existing ip.
func (lt *LocationTests) getLocation200(t *testing.T) {
	ip := "32.123.12.2"
	r := httptest.NewRequest(http.MethodGet, "/v1/location/"+ip, nil)
	w := httptest.NewRecorder()

	lt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a location that exsits.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new location %s.", testID, ip)
		{
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 200 for the response : %v", tests.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 200 for the response.", tests.Success, testID)

			var got locationgrp.LocationResponse
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response : %v", tests.Failed, testID, err)
			}

			// Define what we wanted to receive.
			exp := got
			exp.Location.IP = ip
			exp.Country.UUID = "77eabf6e-30a8-44d0-8952-029d2ca06872"
			exp.City.UUID = "45b5fbd3-755f-4379-8f07-a58d4a30fa2f"

			if diff := cmp.Diff(got, exp); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get the expected result. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get the expected result.", tests.Success, testID)
		}
	}
}

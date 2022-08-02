package location_test

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/mchusovlianov/geodata/business/core/location"
	"github.com/mchusovlianov/geodata/business/data/tests"
	"github.com/mchusovlianov/geodata/foundation/docker"
	"testing"
	"time"
)

var c *docker.Container

func TestMain(m *testing.M) {
	m.Run()
}

func Test_Location(t *testing.T) {
	test := tests.NewIntegration(
		t,
		tests.DBContainer{
			Image: "percona",
			Port:  "3306",
			Name:  "testlocation",
			Args:  []string{"-e", "MYSQL_ROOT_PASSWORD=root"},
		},
	)
	t.Cleanup(test.Teardown)

	core := location.NewCore(test.Log, test.DB, nil)

	t.Log("Given the need to work with Location records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single Location.", testID)
		{
			ctx := context.Background()

			_, err := core.QueryByUUID(ctx, uuid.New().String())
			if err == nil || err != location.ErrNotFound {
				t.Fatalf("\t%s\tTest %d:\tShould not be able to retrieve not-exist location by ID: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould not be able to retrieve not-exist location by ID: %s.", tests.Success, testID, err)

			now := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)

			np := location.NewLocation{
				IP:           "200.106.141.15",
				Longitude:    7.206435933364332,
				Latitude:     -84.87503094689836,
				MysteryValue: 7823011346,
			}

			cty, err := core.Create(ctx, np, now)
			if err == nil {
				t.Fatalf("\t%s\tTest %d:\tShould not be able to create a location.", tests.Failed, testID)
			}
			t.Logf("\t%s\tTest %d:\tShould not be able to create a location: %s.", tests.Success, testID, err)

			np = location.NewLocation{
				IP:           "200.106.141.15",
				Longitude:    7.206435933364332,
				Latitude:     -84.87503094689836,
				MysteryValue: 7823011346,
				CityUUID:     "6c1a1d32-456f-4a20-91d0-cf962c3d6d67",
			}

			cty, err = core.Create(ctx, np, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create a location : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create a location.", tests.Success, testID)

			saved, err := core.QueryByUUID(ctx, cty.UUID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve location by ID: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve location by ID.", tests.Success, testID)

			if diff := cmp.Diff(cty, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same location. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same location.", tests.Success, testID)

			ip := "200.106.141.15"
			saved, err = core.QueryByIP(ctx, ip)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve location by IP: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve location by IP.", tests.Success, testID)

			if diff := cmp.Diff(cty, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same location. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same location.", tests.Success, testID)

			locations, err := core.QueryAll(ctx)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve all locations: %s.", tests.Failed, testID, err)
			}

			if len(locations) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould get back one location. Got: %v", tests.Failed, testID, len(locations))
			}

			if diff := cmp.Diff(saved, locations[0]); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same location. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve all locations.", tests.Success, testID)

		}
	}
}

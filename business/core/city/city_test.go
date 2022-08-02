package city_test

import (
	"context"
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/mchusovlianov/geodata/business/core/city"
	"github.com/mchusovlianov/geodata/business/data/tests"
	"github.com/mchusovlianov/geodata/foundation/docker"
	"testing"
	"time"
)

var c *docker.Container

func TestMain(m *testing.M) {
	m.Run()
}

func Test_City(t *testing.T) {
	test := tests.NewIntegration(
		t,
		tests.DBContainer{
			Image: "percona",
			Port:  "3306",
			Name:  "testcity",
			Args:  []string{"-e", "MYSQL_ROOT_PASSWORD=root"},
		},
	)
	t.Cleanup(test.Teardown)

	core := city.NewCore(test.Log, test.DB, nil)

	t.Log("Given the need to work with City records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single City.", testID)
		{
			ctx := context.Background()
			_, err := core.QueryByUUID(ctx, uuid.New().String())
			if err == nil || err != city.ErrNotFound {
				t.Fatalf("\t%s\tTest %d:\tShould not be able to retrieve not-exist city by ID: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould not be able to retrieve not-exist city by ID: %s.", tests.Success, testID, err)

			now := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)

			np := city.NewCity{
				Name: "Amsterdam",
			}

			cty, err := core.Create(ctx, np, now)
			if err == nil || !errors.Is(err, city.ErrValidation) {
				t.Fatalf("\t%s\tTest %d:\tShould not be able to create a city.", tests.Failed, testID)
			}
			t.Logf("\t%s\tTest %d:\tShould not be able to create a city: %s.", tests.Success, testID, err)

			np = city.NewCity{
				CountryUUID: "6c1a1d32-456f-4a20-91d0-cf962c3d6d67",
				Name:        "Amsterdam",
			}

			cty, err = core.Create(ctx, np, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create a city : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create a city.", tests.Success, testID)

			saved, err := core.QueryByUUID(ctx, cty.UUID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve city by ID: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve city by ID.", tests.Success, testID)

			if diff := cmp.Diff(cty, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same city. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same city.", tests.Success, testID)

			cities, err := core.QueryByCountryUUID(ctx, cty.CountryUUID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve cities by Country UUID: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve cities by Country UUID.", tests.Success, testID)

			if diff := cmp.Diff(cty, cities[0]); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same city. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same city.", tests.Success, testID)

			cty, err = core.Create(ctx, np, now)
			if err == nil || err != city.ErrDuplicate {
				t.Fatalf("\t%s\tTest %d:\tShould not be able to create a duplicate of the city : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould not be able to create a duplicate of the city.", tests.Success, testID)

			cities, err = core.QueryAll(ctx)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve all cities: %s.", tests.Failed, testID, err)
			}

			if len(cities) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould get back one city. Got: %v", tests.Failed, testID, len(cities))
			}

			if diff := cmp.Diff(saved, cities[0]); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same city. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve all cities.", tests.Success, testID)

		}
	}
}

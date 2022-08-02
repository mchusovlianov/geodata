package country_test

import (
	"context"
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/mchusovlianov/geodata/business/core/country"
	"github.com/mchusovlianov/geodata/business/data/tests"
	"github.com/mchusovlianov/geodata/foundation/docker"
	"testing"
	"time"
)

var c *docker.Container

func TestMain(m *testing.M) {
	m.Run()
}

func Test_Country(t *testing.T) {
	test := tests.NewIntegration(
		t,
		tests.DBContainer{
			Image: "percona",
			Port:  "3306",
			Name:  "testcountry",
			Args:  []string{"-e", "MYSQL_ROOT_PASSWORD=root"},
		},
	)
	t.Cleanup(test.Teardown)

	core := country.NewCore(test.Log, test.DB, nil)

	t.Log("Given the need to work with Country records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single Country.", testID)
		{
			ctx := context.Background()

			_, err := core.QueryByUUID(ctx, uuid.New().String())
			if err == nil || err != country.ErrNotFound {
				t.Fatalf("\t%s\tTest %d:\tShould not be able to retrieve not-exist country by ID: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould not be able to retrieve not-exist country by ID: %s.", tests.Success, testID, err)

			now := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)

			failedCountry := country.NewCountry{
				Code: "NE",
			}

			cty, err := core.Create(ctx, failedCountry, now)
			if err == nil || !errors.Is(err, country.ErrValidation) {
				t.Fatalf("\t%s\tTest %d:\tShould not be able to create a country.", tests.Failed, testID)
			}

			t.Logf("\t%s\tTest %d:\tShould not be able to create a country: %s. ", tests.Success, testID, err)

			np := country.NewCountry{
				Code: "NE",
				Name: "Nepal",
			}

			cty, err = core.Create(ctx, np, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create a country : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create a country.", tests.Success, testID)
			t.Logf("%+v\n", cty)
			saved, err := core.QueryByUUID(ctx, cty.UUID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve country by ID: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve country by ID.", tests.Success, testID)

			if diff := cmp.Diff(cty, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same country. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same country.", tests.Success, testID)

			saved, err = core.QueryByCode(ctx, np.Code)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve country by Code: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve country by Code.", tests.Success, testID)

			if diff := cmp.Diff(cty, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same country by Code. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same country by Code.", tests.Success, testID)

			cty, err = core.Create(ctx, np, now)
			if err == nil || err != country.ErrDuplicate {
				t.Fatalf("\t%s\tTest %d:\tShould not be able to create a duplicate of the country : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould not be able to create a duplicate of the country.", tests.Success, testID)

			countries, err := core.QueryAll(ctx)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve all countries: %s.", tests.Failed, testID, err)
			}

			if len(countries) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould get back one country. Got: %v", tests.Failed, testID, len(countries))
			}

			if diff := cmp.Diff(saved, countries[0]); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same country. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve all countries.", tests.Success, testID)

		}
	}
}

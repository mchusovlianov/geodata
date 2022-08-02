package importer_test

import (
	"bytes"
	"context"
	cityCore "github.com/mchusovlianov/geodata/business/core/city"
	countryCore "github.com/mchusovlianov/geodata/business/core/country"
	"github.com/mchusovlianov/geodata/business/core/importer"
	locationCore "github.com/mchusovlianov/geodata/business/core/location"
	"github.com/mchusovlianov/geodata/business/data/tests"
	"github.com/mchusovlianov/geodata/foundation/docker"
	"testing"
)

var c *docker.Container

func TestMain(m *testing.M) {
	m.Run()
}

func Test_Importer(t *testing.T) {
	test := tests.NewIntegration(
		t,
		tests.DBContainer{
			Image: "percona",
			Port:  "3306",
			Name:  "testimporter",
			Args:  []string{"-e", "MYSQL_ROOT_PASSWORD=root"},
		},
	)
	t.Cleanup(test.Teardown)

	core, err := importer.NewCore(test.Log, test.DB)
	if err != nil {
		t.Fatalf("Can't create an importer core %s", err)
	}

	t.Log("Given the need to work with the Import operation.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen import file in a wrong format.", testID)
		{
			b := []byte(`
country_code,country,city,latitude,longitude,mystery_value
LI,Guyana,Port Karson,-78.2274228596799,-163.26218895343357,1337885276`)

			bytesReader := bytes.NewReader(b)

			ctx := context.Background()
			_, err := core.Import(ctx, bytesReader, 2)
			if err == nil {
				t.Fatalf("\t%s\tTest %d:\tShould not be able to import file in the wrong format.", tests.Failed, testID)
			}

			t.Logf("\t%s\tTest %d:\tShould not be able to import file in the wrong format: %s.", tests.Success, testID, err)
		}

		testID += 1
		t.Logf("\tTest %d:\tWhen import file in a good format with only good records.", testID)
		{
			b := []byte(`ip_address,country_code,country,city,latitude,longitude,mystery_value
200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346
200.106.141.16,SI,Nepal,TestCity2,-84.87503094689832,7.206435933364332,7823011346
160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115`)

			bytesReader := bytes.NewReader(b)

			ctx := context.Background()
			stat, err := core.Import(ctx, bytesReader, 2)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to import a good file without error %s.", tests.Failed, testID, err)
			}

			if stat.GoodLines != 3 {
				t.Fatalf("\t%s\tTest %d:\tShould be able to import a good file without error. Wrong amount of inserted good lines: %v. Expected - %v.", tests.Failed, testID, stat.GoodLines, 3)
			}

			countryCoreInst := countryCore.NewCore(test.Log, test.DB, nil)
			countries, err := countryCoreInst.QueryAll(ctx)
			if len(countries) != 2 {
				t.Fatalf("\t%s\tTest %d:\tShould be created 2 countries during import, got - %v.", tests.Failed, testID, len(countries))
			}

			cityCoreInst := cityCore.NewCore(test.Log, test.DB, nil)
			cities, err := cityCoreInst.QueryAll(ctx)
			if len(cities) != 3 {
				t.Fatalf("\t%s\tTest %d:\tShould be created 3 cities during import, got - %v.", tests.Failed, testID, len(cities))
			}

			locationCoreInst := locationCore.NewCore(test.Log, test.DB, nil)
			locations, err := locationCoreInst.QueryAll(ctx)
			if len(locations) != 3 {
				t.Fatalf("\t%s\tTest %d:\tShould be created 3 locations during import, got - %v.", tests.Failed, testID, len(cities))
			}

			t.Logf("\t%s\tTest %d:\tShould be able to import a good file without error.", tests.Success, testID)
		}
	}
}

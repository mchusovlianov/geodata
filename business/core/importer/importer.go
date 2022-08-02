package importer

import (
	"context"
	"encoding/csv"
	"errors"
	"github.com/jmoiron/sqlx"
	cityCore "github.com/mchusovlianov/geodata/business/core/city"
	countryCore "github.com/mchusovlianov/geodata/business/core/country"
	locationCore "github.com/mchusovlianov/geodata/business/core/location"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"io"
	"strconv"
	"sync"
	"time"
)

var (
	ErrNotValidFormat = errors.New("file in wrong format")
)

const taskChanSize = 10

type Statistic struct {
	TotalLines  int
	FailedLines int
	GoodLines   int
	Duration    time.Duration
}

// Core manages the set of APIs for location access.
type Core struct {
	db     *sqlx.DB
	log    *zap.SugaredLogger
	good   atomic.Int32
	failed atomic.Int32
}

// NewCore constructs a core for location api access.
func NewCore(log *zap.SugaredLogger, dbConn *sqlx.DB) (Core, error) {
	return Core{
		log: log,
		db:  dbConn,
	}, nil
}

// worker - worker function to parse csv file and import to the database
func (c *Core) worker(wg *sync.WaitGroup, tasks chan []string) error {
	defer wg.Done()

	// =========================================================================
	// init cores
	countryCoreInst := countryCore.NewCore(c.log, c.db, nil)
	cityCoreInst := cityCore.NewCore(c.log, c.db, nil)
	locationCoreInst := locationCore.NewCore(c.log, c.db, nil)

	// =========================================================================
	// prepare caches for countries and cities to avoid unnecessary calls to db
	countryCache := make(map[string]string)
	cityCache := make(map[string]string)

	for task := range tasks {
		now := time.Now()
		ctx := context.Background()

		// =========================================================================
		// creating vars from csv line for next usage
		// this code looks repetitive but it allows to avoid
		// unclean indexes, like record[2]. It might lead to
		// hardly detected errors
		ipIdx, err := getRecordValue("ip_address")
		if err != nil {
			c.failed.Inc()
			continue
		}
		ip := task[ipIdx]

		codeIdx, err := getRecordValue("country_code")
		if err != nil {
			c.failed.Inc()
			continue
		}
		countryCode := task[codeIdx]

		countryIdx, err := getRecordValue("country")
		if err != nil {
			c.failed.Inc()
			continue
		}
		countryName := task[countryIdx]

		cityIdx, err := getRecordValue("city")
		if err != nil {
			c.failed.Inc()
			continue
		}
		cityName := task[cityIdx]

		latIdx, err := getRecordValue("latitude")
		if err != nil {
			c.failed.Inc()
			continue
		}
		latStr := task[latIdx]

		lonIdx, err := getRecordValue("longitude")
		if err != nil {
			c.failed.Inc()
			continue
		}
		lonStr := task[lonIdx]

		mysteryIdx, err := getRecordValue("mystery_value")
		if err != nil {
			c.failed.Inc()
			continue
		}
		mysteryValueStr := task[mysteryIdx]

		// =========================================================================
		// try ti get country uuid from cache
		countryUUID, ok := countryCache[countryCode]
		if !ok {
			// create new country
			createdCountry, err := countryCoreInst.Create(ctx, countryCore.NewCountry{
				Name: countryName,
				Code: countryCode,
			}, now)

			if errors.Is(err, countryCore.ErrValidation) {
				c.failed.Inc()
				continue
			}

			if errors.Is(err, countryCore.ErrDuplicate) {
				err = c.loadCaches(countryCode, countryCache, cityCache)
				if err != nil {
					c.failed.Inc()
					continue
				}

				countryUUID, ok = countryCache[countryCode]
				if !ok {
					c.failed.Inc()
					continue
				}
			}

			if err != nil {
				c.failed.Inc()
				continue
			}

			//update city cache if we successfully created  city
			if createdCountry.UUID != "" {
				// update country cache
				countryUUID = createdCountry.UUID
				countryCache[countryCode] = countryUUID
			}
		}

		// try ti get city uuid from cache
		cityUUID := cityCache[cityName+"#"+countryUUID]
		if cityUUID == "" {
			// create new city
			createdCity, err := cityCoreInst.Create(ctx, cityCore.NewCity{
				CountryUUID: countryUUID,
				Name:        cityName,
			}, now)

			if errors.Is(err, cityCore.ErrValidation) {
				c.failed.Inc()
				continue
			}

			if err != nil {
				c.failed.Inc()
				continue
			}

			//update city cache if we successfully created  city
			if createdCity.UUID != "" {
				cityUUID = createdCity.UUID
				cityCache[cityName+"#"+countryUUID] = cityUUID
			}
		}
		// prepare values for creating new location
		mysteryValue, err := strconv.ParseInt(mysteryValueStr, 10, 64)
		if err != nil {
			c.failed.Inc()
			continue
		}

		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			c.failed.Inc()
			continue
		}

		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			c.failed.Inc()
			continue
		}

		// create new location task
		_, err = locationCoreInst.Create(ctx, locationCore.NewLocation{
			IP:           ip,
			Longitude:    lon,
			Latitude:     lat,
			MysteryValue: mysteryValue,
			CityUUID:     cityUUID,
		}, now)

		if errors.Is(err, locationCore.ErrValidation) {
			c.failed.Inc()
			continue
		}

		if err != nil {
			c.failed.Inc()
			continue
		}

		c.good.Inc()
	}

	return nil
}

// Import - import csv file to the database
func (c *Core) Import(ctx context.Context, f io.Reader, workersCount int) (Statistic, error) {
	// prepare csv-reader
	csvReader := csv.NewReader(f)
	isFirstLine := true

	st := Statistic{}
	start := time.Now()

	// create task chan
	tasks := make([]chan []string, workersCount)
	for idx, _ := range tasks {
		tasks[idx] = make(chan []string, taskChanSize)
	}

	// worker usage array
	usage := make([]int, workersCount)
	// Map to assign a specific country to a specific worker
	countryToWorker := make(map[string]int)

	// run worker routines
	wg := sync.WaitGroup{}
	wg.Add(workersCount)

	c.log.Infow("start workers", "count", workersCount)
	for i := 0; i < workersCount; i++ {
		go c.worker(&wg, tasks[i])
	}

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return Statistic{}, err
		}

		if len(record) == 0 {
			continue
		}

		st.TotalLines += 1

		if isFirstLine {
			// check if field names are equal to which are they considered
			if !checkFileFormat(record) {
				return Statistic{}, ErrNotValidFormat
			}

			isFirstLine = false

			// pre-warm headersReverseMap
			// there is used map without mutex so there will be a problem
			// in case of simultaneously write by multiple go-routines
			// so run write to map operation before any of the go-routines started
			// this is the only write operation in headersReverseMap
			_, _ = getRecordValue("ip_address")
			continue
		}

		idx, _ := getRecordValue("country_code")
		str := record[idx]

		if taskChanIdx, ok := countryToWorker[str]; ok {
			tasks[taskChanIdx] <- record
		} else {
			var taskChanIdx int
			for idx, row := range usage {
				if usage[taskChanIdx] > row {
					taskChanIdx = idx
				}
			}

			usage[taskChanIdx] += 1
			countryToWorker[str] = taskChanIdx
			tasks[taskChanIdx] <- record
		}

		if st.TotalLines%10000 == 0 {
			c.log.Infow("processed", "lines", st.TotalLines)
		}
	}

	for _, task := range tasks {
		close(task)
	}
	wg.Wait()

	st.GoodLines = int(c.good.Load())
	st.FailedLines = int(c.failed.Load())

	st.Duration = time.Since(start)
	return st, nil
}

// loadCaches - load country and all related cities
func (c *Core) loadCaches(countryCode string, countryCache, cityCache map[string]string) error {
	countryCoreInst := countryCore.NewCore(c.log, c.db, nil)
	cityCoreInst := cityCore.NewCore(c.log, c.db, nil)

	ctx := context.Background()
	country, err := countryCoreInst.QueryByCode(ctx, countryCode)
	if err != nil {
		return err
	}

	cities, err := cityCoreInst.QueryByCountryUUID(ctx, country.UUID)
	if err != nil {
		return err
	}

	countryCache[countryCode] = country.UUID

	for _, city := range cities {
		cityCache[city.Name+"#"+country.UUID] = city.UUID
	}

	return nil
}

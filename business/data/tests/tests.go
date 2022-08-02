// Package tests contains supporting code for running tests.
package tests

import (
	"bytes"
	"context"
	"fmt"
	"github.com/mchusovlianov/geodata/business/data/dbschema"
	"github.com/mchusovlianov/geodata/foundation/database"
	"github.com/mchusovlianov/geodata/foundation/docker"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

// DBContainer provides configuration for a container to run.
type DBContainer struct {
	Image  string
	Port   string
	Name   string
	IsSeed bool
	Args   []string
}

func newLogger() (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]interface{}{
		"service": "test-logger",
	}

	log, err := config.Build()
	if err != nil {
		return nil, err
	}

	return log.Sugar(), nil
}

// NewUnit creates a test database inside a Docker container. It creates the
// required table structure but the database is otherwise empty. It returns
// the database to use as well as a function to call at the end of the test.
func NewUnit(t *testing.T, dbc DBContainer) (*zap.SugaredLogger, *sqlx.DB, func()) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	dbc.Args = append(dbc.Args, "-e", "MYSQL_DATABASE="+dbc.Name)
	c, err := docker.StartContainer(dbc.Image, dbc.Port, "mysqladmin ping --silent", dbc.Args...)
	if err != nil {
		t.Error(err)
	}

	var i int
	for ; i < 10; i++ {
		status, _ := docker.Status(c.ID)
		if status == "healthy" {
			break
		}

		time.Sleep(time.Second * time.Duration(i))
	}

	if i == 10 {
		docker.DumpContainerLogs(t, c.ID)
		docker.StopContainer(c.ID)
		t.Fatalf("can't wait until db is up")
	}

	db, err := database.Open(database.Config{
		User:     "root",
		Password: "root",
		Host:     c.Host,
		Name:     dbc.Name,
	})
	if err != nil {
		t.Fatalf("Opening database connection: %v", err)
	}

	t.Log("Waiting for database to be ready ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := dbschema.Migrate(ctx, db); err != nil {
		docker.DumpContainerLogs(t, c.ID)
		docker.StopContainer(c.ID)
		t.Fatalf("Migrating error: %s", err)
	}

	if dbc.IsSeed {
		if err := dbschema.Seed(ctx, db); err != nil {
			docker.DumpContainerLogs(t, c.ID)
			docker.StopContainer(c.ID)
			t.Fatalf("Seeding error: %s", err)
		}
	}

	log, err := newLogger()
	if err != nil {
		t.Fatalf("logger error: %s", err)
	}

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		db.Close()
		docker.StopContainer(c.ID)

		log.Sync()

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		os.Stdout = old
		fmt.Println("******************** LOGS ********************")
		fmt.Print(buf.String())
		fmt.Println("******************** LOGS ********************")
	}

	return log, db, teardown
}

// Test owns state for running and shutting down tests.
type Test struct {
	DB       *sqlx.DB
	Log      *zap.SugaredLogger
	Teardown func()

	t *testing.T
}

// NewIntegration creates a database, seeds it, constructs an authenticator.
func NewIntegration(t *testing.T, dbc DBContainer) *Test {
	log, db, teardown := NewUnit(t, dbc)

	test := Test{
		DB:       db,
		Log:      log,
		t:        t,
		Teardown: teardown,
	}

	return &test
}

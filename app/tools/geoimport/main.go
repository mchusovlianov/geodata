package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/ardanlabs/conf/v3"
	"github.com/mchusovlianov/geodata/business/core/importer"
	"github.com/mchusovlianov/geodata/business/data/dbschema"
	"github.com/mchusovlianov/geodata/foundation/database"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"runtime"
)

const serviceName = "data-importer"

func newLogger() (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]interface{}{
		"service": serviceName,
	}

	log, err := config.Build()
	if err != nil {
		return nil, err
	}

	return log.Sugar(), nil
}

func main() {
	log, err := newLogger()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	var filePath string
	flag.StringVar(&filePath, "filepath", "", "")
	flag.Parse()

	// Perform the startup and shutdown sequence.
	if err := run(log, filePath); err != nil {
		fmt.Println(err)
		log.Sync()
		os.Exit(1)
	}
}
func run(log *zap.SugaredLogger, filePath string) error {
	// =========================================================================
	// GOMAXPROCS

	// Want to see what maxprocs reports.
	opt := maxprocs.Logger(log.Infof)

	// Set the correct number of threads for the service
	// based on what is available either by the machine or quotas.
	if _, err := maxprocs.Set(opt); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}
	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// =========================================================================
	// Configuration
	cfg := struct {
		DB struct {
			User         string `conf:"default:root"`
			Password     string `conf:"default:mysql,mask"`
			Host         string `conf:"default:localhost"`
			Name         string `conf:"default:geodata"`
			MaxIdleConns int    `conf:"default:0"`
			MaxOpenConns int    `conf:"default:0"`
		}
		WorkersCount int `conf:"default:8""`
	}{}

	const prefix = "GEOIMPORT"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Infow("startup", "config", out)

	// =========================================================================
	// Database Support

	// Create connectivity to the database.
	log.Infow("startup", "status", "initializing database support", "host", cfg.DB.Host)

	db, err := database.Open(database.Config{
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		Host:         cfg.DB.Host,
		Name:         cfg.DB.Name,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer func() {
		log.Infow("shutdown", "status", "stopping database support", "host", cfg.DB.Host)
		db.Close()
	}()

	// =========================================================================
	// apply database schema
	ctx := context.Background()
	if err := dbschema.Migrate(ctx, db); err != nil {
		return fmt.Errorf("db schema migration error: %w", err)
	}

	// create new importer core
	importerCore, err := importer.NewCore(log, db)
	if err != nil {
		return fmt.Errorf("creating importer: %w", err)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("can't open file: %w", err)
	}
	defer f.Close()

	stat, err := importerCore.Import(ctx, f, cfg.WorkersCount)
	if err != nil {
		return fmt.Errorf("can't : %w", err)
	}

	out, err = conf.String(&stat)
	if err != nil {
		return fmt.Errorf("generating import result for output: %w", err)
	}

	log.Infow("result of import operation ", "statistic", out)

	return nil
}

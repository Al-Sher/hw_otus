package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/config"
	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/logger"
	sqlstorage "github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/storage/pgsql"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

var (
	configFile                string
	configFormat              string
	ErrUnsupportedStorageType = errors.New("unsupported storage type")
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
	flag.StringVar(&configFormat, "format", "", "Format configuration file")
}

func main() {
	flag.Parse()

	if configFormat == "" {
		configFormat = config.ParseFormatFile(configFile)
	}

	c, err := config.NewConfig(configFile, configFormat)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logg := logger.New(c.LoggerLevel(), c.LoggerPath())
	var db *sql.DB

	switch c.StorageType() {
	case sqlstorage.Type:
		if db, err = goose.OpenDBWithDriver("postgres", c.StorageDsn()); err != nil {
			logg.Fatal(err)
		}

		defer func() {
			err := db.Close()
			if err != nil {
				logg.Error(err)
			}
		}()
	default:
		logg.Fatal(ErrUnsupportedStorageType)
	}

	if err := goose.Up(db, c.MigrationPath()); err != nil {
		logg.Fatal(err)
	}
}

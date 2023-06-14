package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/app"
	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/config"
	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/server/http"
	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/storage/pgsql"
)

var (
	configFile   string
	configFormat string
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
	flag.StringVar(&configFormat, "format", "", "Format configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	if configFormat == "" {
		configFormat = config.ParseFormatFile(configFile)
	}

	c, err := config.NewConfig(configFile, configFormat)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx := context.Background()
	logg, err := logger.New(c.LoggerLevel(), c.LoggerPath())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var s storage.Storage

	switch c.StorageType() {
	case memorystorage.Type:
		s = memorystorage.New()
	case sqlstorage.Type:
		s = sqlstorage.New()
	}

	err = s.Connect(ctx, c.StorageDsn())
	if err != nil {
		logg.Error(err)
		os.Exit(1)
	}

	defer func() {
		err := s.Close(ctx)
		if err != nil {
			logg.Error(err)
		}
	}()

	calendar := app.New(s)

	server := internalhttp.NewServer(calendar, logg, c)
	grpcServer := grpc.NewServer(calendar)

	ctx, cancel := signal.NotifyContext(ctx,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if c.HTTPAddr() != "" {
			if err := server.Stop(ctx); err != nil {
				logg.Error("failed to stop http server: " + err.Error())
			}
		}

		if c.GRPCAddr() != "" {
			grpcServer.Stop()
		}
		logg.Info("calendar is shutdown...")
	}()

	logg.Info("calendar is running...")

	wg := sync.WaitGroup{}

	if c.HTTPAddr() != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := server.Start(ctx); err != nil {
				logg.Error("failed to start http server: " + err.Error())
				cancel()
				os.Exit(1)
			}
		}()
	}

	if c.GRPCAddr() != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := grpcServer.Start(ctx); err != nil {
				logg.Error("failed to start http server: " + err.Error())
				cancel()
				os.Exit(1)
			}
		}()
	}

	wg.Wait()
}

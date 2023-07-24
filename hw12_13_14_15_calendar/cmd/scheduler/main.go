package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/config"
	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/storage/pgsql"
	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/pkg/producer"
	_ "github.com/lib/pq"
	"github.com/rabbitmq/amqp091-go"
)

var (
	configFile   string
	configFormat string
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
	flag.StringVar(&configFormat, "format", "", "Format configuration file")
}

var publisher *producer.Producer

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

	logg, err := logger.New(c.LoggerLevel(), c.LoggerPath())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx := context.Background()
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

	conn, err := amqp091.Dial(c.RabbitAddr())
	if err != nil {
		logg.Error(err)
		os.Exit(1)
	}
	publisher = producer.New(c.RabbitExchange(), c.RabbitExchangeType(), conn)

	ctx, cancel := signal.NotifyContext(ctx,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		defer cancel()
		timer := time.NewTicker(time.Duration(c.SchedulerInterval()) * time.Second)
		timerForClear := time.NewTicker(time.Duration(c.ClearStorageInterval()) * time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				events, err := s.EventsForNotification(ctx)
				if err != nil {
					logg.Error(err)
					continue
				}
				ids, err := publishEvents(ctx, events, c.RabbitRoutingKey())
				if err != nil {
					logg.Error(err)
					continue
				}
				if err := s.ClearNotificationDates(ctx, ids); err != nil {
					logg.Error(err)
				}
			case <-timerForClear.C:
				if err := s.ClearOldEvents(ctx); err != nil {
					logg.Error(err)
				}
			}
		}
	}()

	logg.Info("scheduler is running...")
	<-ctx.Done()
	logg.Info("scheduler is shutdown...")
}

func publishEvents(ctx context.Context, events []storage.Event, routingKey string) ([]string, error) {
	eventsForClearNotification := make([]string, 0)
	for _, event := range events {
		if err := publish(ctx, event, routingKey); err != nil {
			return nil, err
		}
		eventsForClearNotification = append(eventsForClearNotification, event.ID)
	}

	return eventsForClearNotification, nil
}

func publish(ctx context.Context, event storage.Event, routingKey string) error {
	t, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = publisher.Publish(ctx, routingKey, t)
	return err
}

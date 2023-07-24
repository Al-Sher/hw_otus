package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/config"
	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/pkg/consumer"
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

	conn, err := amqp091.Dial(c.RabbitAddr())
	if err != nil {
		logg.Error(err)
		os.Exit(1)
	}
	listener, err := consumer.New(
		c.ConsumerName(),
		conn,
		c.RabbitExchange(),
		c.RabbitExchangeType(),
		c.QueueName(),
		c.RabbitRoutingKey(),
	)
	if err != nil {
		logg.Error(err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(ctx,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	ch, err := listener.Consume(ctx, c.QueueName(), notification)
	if err != nil {
		logg.Error(err)
		os.Exit(1) //nolint:gocritic
	}

	logg.Info("scheduler is running...")

	<-ch
	<-ctx.Done()
	logg.Info("scheduler is shutdown...")
}

func notification(msg []byte) error {
	fmt.Println(string(msg))

	return nil
}

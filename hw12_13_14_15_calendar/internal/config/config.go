package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"

	memorystorage "github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/storage/pgsql"
	"github.com/BurntSushi/toml"
)

var (
	ErrNotValidFormat     = errors.New("невалидный формат конфигурационного файла")
	ErrInvalidStorageType = errors.New("невалидный тип хранилища")
)

const (
	DefaultPathForLogger        = "/dev/stdout"
	DefaultReadTimeout          = 15
	DefaultMigrationPath        = "./migrations/"
	DefaultRabbitMQExchangeType = "direct"
	DefaultSchedulerInterval    = 10
	DefaultConsumerName         = "calendar-consumer"
	DefaultQueueName            = "calendar-queue"
	DefaultClearStorageInterval = 60 * 60
)

type configLogger struct {
	Level string `json:"level" toml:"level"`
	Path  string `json:"path" toml:"path"`
}

type server struct {
	HTTPAddr        string  `json:"httpAddr" toml:"httpAddr"`
	HTTPReadTimeout float64 `json:"httpReadTimeout" toml:"httpReadTimeout"`
	GRPCAddr        string  `json:"grpcAddr" toml:"grpcAddr"`
}

type storage struct {
	StorageType          string  `json:"storageType" toml:"storageType"`
	Dsn                  string  `json:"dsn" toml:"dsn"`
	MigrationPath        string  `json:"migrationPath" toml:"migrationPath"`
	ClearStorageInterval float64 `json:"clearStorageInterval" toml:"clearStorageInterval"`
}

type config struct {
	Logger   configLogger `json:"logger" toml:"logger"`
	Storage  storage      `json:"storage" toml:"storage"`
	Server   server       `json:"server" toml:"server"`
	RabbitMQ rabbitmq     `json:"rabbitMq" toml:"rabbitMq"`
}

type rabbitmq struct {
	Addr              string  `json:"addr" toml:"addr"`
	Exchange          string  `json:"exchange" toml:"exchange"`
	ExchangeType      string  `json:"exchangeType" toml:"exchangeType"`
	RabbitRoutingKey  string  `json:"rabbitRoutingKey" toml:"rabbitRoutingKey"`
	SchedulerInterval float64 `json:"schedulerInterval" toml:"schedulerInterval"`
	ConsumerName      string  `json:"consumerName" toml:"consumerName"`
	QueueName         string  `json:"queueName" toml:"queueName"`
}

type Config interface {
	LoggerLevel() string
	LoggerPath() string

	HTTPAddr() string
	HTTPReadTimeout() float64
	GRPCAddr() string

	StorageType() string
	StorageDsn() string
	MigrationPath() string
	ClearStorageInterval() float64

	RabbitAddr() string
	RabbitExchange() string
	RabbitExchangeType() string
	RabbitRoutingKey() string
	SchedulerInterval() float64
	ConsumerName() string
	QueueName() string
}

func (c *config) LoggerLevel() string {
	return c.Logger.Level
}

func (c *config) LoggerPath() string {
	return c.Logger.Path
}

func (c *config) HTTPAddr() string {
	return c.Server.HTTPAddr
}

func (c *config) HTTPReadTimeout() float64 {
	return c.Server.HTTPReadTimeout
}

func (c *config) GRPCAddr() string {
	return c.Server.GRPCAddr
}

func (c *config) StorageDsn() string {
	return c.Storage.Dsn
}

func (c *config) StorageType() string {
	return c.Storage.StorageType
}

func (c *config) MigrationPath() string {
	return c.Storage.MigrationPath
}

func (c *config) ClearStorageInterval() float64 {
	return c.Storage.ClearStorageInterval
}

func (c *config) RabbitAddr() string {
	return c.RabbitMQ.Addr
}

func (c *config) RabbitExchange() string {
	return c.RabbitMQ.Exchange
}

func (c *config) RabbitExchangeType() string {
	return c.RabbitMQ.ExchangeType
}

func (c *config) RabbitRoutingKey() string {
	return c.RabbitMQ.RabbitRoutingKey
}

func (c *config) SchedulerInterval() float64 {
	return c.RabbitMQ.SchedulerInterval
}

func (c *config) ConsumerName() string {
	return c.RabbitMQ.ConsumerName
}

func (c *config) QueueName() string {
	return c.RabbitMQ.QueueName
}

func NewConfig(path string, format string) (Config, error) {
	var err error
	cfg := &config{}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	value, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	switch format {
	case "json":
		cfg, err = newFromJSON(value)
	case "toml":
		cfg, err = newFromToml(value)
	default:
		return nil, ErrNotValidFormat
	}

	if err != nil {
		return nil, err
	}

	cfg.setDefaultValues()

	return cfg, cfg.validate()
}

func newFromJSON(value []byte) (*config, error) {
	c := &config{}
	if err := json.Unmarshal(value, &c); err != nil {
		return nil, err
	}

	return c, nil
}

func newFromToml(value []byte) (*config, error) {
	c := &config{}
	if err := toml.Unmarshal(value, &c); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *config) validate() error {
	switch c.StorageType() {
	case memorystorage.Type, sqlstorage.Type:
		break
	default:
		return ErrInvalidStorageType
	}

	return nil
}

func (c *config) setDefaultValues() {
	if c.Logger.Path == "" {
		c.Logger.Path = DefaultPathForLogger
	}

	if c.Server.HTTPReadTimeout == 0 {
		c.Server.HTTPReadTimeout = DefaultReadTimeout
	}

	if c.Storage.MigrationPath == "" {
		c.Storage.MigrationPath = DefaultMigrationPath
	}

	if c.RabbitMQ.ExchangeType == "" {
		c.RabbitMQ.ExchangeType = DefaultRabbitMQExchangeType
	}

	if c.RabbitMQ.SchedulerInterval == 0 {
		c.RabbitMQ.SchedulerInterval = DefaultSchedulerInterval
	}

	if c.RabbitMQ.ConsumerName == "" {
		c.RabbitMQ.ConsumerName = DefaultConsumerName
	}

	if c.RabbitMQ.QueueName == "" {
		c.RabbitMQ.QueueName = DefaultQueueName
	}

	if c.Storage.ClearStorageInterval == 0 {
		c.Storage.ClearStorageInterval = DefaultClearStorageInterval
	}
}

func ParseFormatFile(path string) string {
	configFormat := ""
	if l := strings.LastIndex(path, "."); l != -1 {
		configFormat = path[l+1:]
	}

	return configFormat
}

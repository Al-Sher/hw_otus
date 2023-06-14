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
	DefaultPathForLogger = "/dev/stdout"
	DefaultReadTimeout   = 15
	DefaultMigrationPath = "./migrations/"
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
	StorageType   string `json:"storageType" toml:"storageType"`
	Dsn           string `json:"dsn" toml:"dsn"`
	MigrationPath string `json:"migrationPath" toml:"migrationPath"`
}

type config struct {
	Logger  configLogger `json:"logger" toml:"logger"`
	Storage storage      `json:"storage" toml:"storage"`
	Server  server       `json:"server" toml:"server"`
}

type Config interface {
	afterCreate() error

	LoggerLevel() string
	LoggerPath() string

	HTTPAddr() string
	HTTPReadTimeout() float64
	GRPCAddr() string

	StorageType() string
	StorageDsn() string
	MigrationPath() string
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

func NewConfig(path string, format string) (Config, error) {
	var err error
	var cfg Config

	switch format {
	case "json":
		cfg, err = newFromJSON(path)
	case "toml":
		cfg, err = newFromToml(path)
	default:
		return nil, ErrNotValidFormat
	}

	if err != nil {
		return nil, err
	}

	return cfg, cfg.afterCreate()
}

func newFromJSON(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	value, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	c := &config{}
	if err := json.Unmarshal(value, &c); err != nil {
		return nil, err
	}

	return c, nil
}

func newFromToml(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	value, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	c := &config{}

	if err := toml.Unmarshal(value, &c); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *config) afterCreate() error {
	switch c.StorageType() {
	case memorystorage.Type, sqlstorage.Type:
		break
	default:
		return ErrInvalidStorageType
	}

	if c.Logger.Path == "" {
		c.Logger.Path = DefaultPathForLogger
	}

	if c.Server.HTTPReadTimeout == 0 {
		c.Server.HTTPReadTimeout = DefaultReadTimeout
	}

	if c.Storage.MigrationPath == "" {
		c.Storage.MigrationPath = DefaultMigrationPath
	}

	return nil
}

func ParseFormatFile(path string) string {
	configFormat := ""
	if l := strings.LastIndex(path, "."); l != -1 {
		configFormat = path[l+1:]
	}

	return configFormat
}

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
	LoggerLevel() string
	LoggerPath() string

	HTTPAddr() string
	HTTPReadTimeout() float64

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
}

func ParseFormatFile(path string) string {
	configFormat := ""
	if l := strings.LastIndex(path, "."); l != -1 {
		configFormat = path[l+1:]
	}

	return configFormat
}

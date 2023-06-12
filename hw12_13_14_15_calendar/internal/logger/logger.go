package logger

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(msg interface{})
	Error(msg interface{})
	Debug(msg interface{})
	Warning(msg interface{})
	Fatal(msg interface{})
}

type logger struct {
	logger *zap.Logger
}

func New(level string, path string) Logger {
	l, err := config(levelByString(level), path).Build()
	if err != nil {
		fmt.Printf("Ошибка инициализации логгера: %s", err)
		os.Exit(1)
	}

	return &logger{
		logger: l,
	}
}

func (l *logger) Info(msg interface{}) {
	l.logger.Info(toString(msg))
}

func (l *logger) Error(msg interface{}) {
	l.logger.Error(toString(msg))
}

func (l *logger) Debug(msg interface{}) {
	l.logger.Debug(toString(msg))
}

func (l *logger) Warning(msg interface{}) {
	l.logger.Warn(toString(msg))
}

func (l *logger) Fatal(msg interface{}) {
	l.logger.Fatal(toString(msg))
}

func levelByString(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "fatal":
		return zapcore.FatalLevel
	case "error":
		return zapcore.ErrorLevel
	case "warning":
		return zapcore.WarnLevel
	case "debug":
		return zapcore.DebugLevel
	default:
		return zapcore.InfoLevel
	}
}

func config(lvl zapcore.Level, path string) zap.Config {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(lvl)
	cfg.OutputPaths = []string{path}

	return cfg
}

func toString(msg interface{}) string {
	return fmt.Sprintf("%v", msg)
}

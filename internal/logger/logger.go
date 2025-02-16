package logger

import (
	"log/slog"
	"os"
)

type Config struct {
	Level     slog.Level `envconfig:"LEVEL" default:"info"`
	AddSource bool       `envconfig:"ADDSOURCE" default:"false"`
}

func Init(cfg *Config) *slog.Logger {
	logLevel := &slog.LevelVar{} // INFO log level by default
	logLevel.Set(cfg.Level)

	opts := &slog.HandlerOptions{
		AddSource: cfg.AddSource,
		Level:     logLevel,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)

	logger := slog.New(handler)

	slog.SetDefault(logger)

	return logger
}

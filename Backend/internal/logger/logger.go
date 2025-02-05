package logger

import (
	"log/slog"
	"os"
	"sync"
)

const (
	lvlLocal = "local"
	lvlDev   = "dev"
	lvlProd  = "prod"
)

type LogWrapper struct {
	logger *slog.Logger
	sync.Once
}

var log *LogWrapper

func MustInit(level string) {
	log = &LogWrapper{}
	log.Do(func() {
		switch level {
		case lvlLocal:
			log.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		case lvlDev:
			log.logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		case lvlProd:
			log.logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		}
	})
}

func Info(msg string, fields ...slog.Attr) {
	var args []any
	for i := range fields {
		args = append(args, slog.Attr{Key: fields[i].Key, Value: slog.StringValue(fields[i].Value.String())})
	}
	l := log.logger.With(args...)
	l.Info(msg)
}

func Debug(msg string, fields ...slog.Attr) {
	var args []any
	for i := range fields {
		args = append(args, slog.Attr{Key: fields[i].Key, Value: slog.StringValue(fields[i].Value.String())})
	}
	l := log.logger.With(args...)
	l.Debug(msg)
}

func Error(msg string, fields ...slog.Attr) {
	var args []any
	for i := range fields {
		args = append(args, slog.Attr{Key: fields[i].Key, Value: slog.StringValue(fields[i].Value.String())})
	}
	l := log.logger.With(args...)
	l.Error(msg)
}

func Warn(msg string, fields ...slog.Attr) {
	var args []any
	for i := range fields {
		args = append(args, slog.Attr{Key: fields[i].Key, Value: slog.StringValue(fields[i].Value.String())})
	}
	l := log.logger.With(args...)
	l.Warn(msg)
}

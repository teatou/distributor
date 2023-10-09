package mylogger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Error(...any)
	Errorf(string, ...any)
	Info(...any)
	Infof(string, ...any)
	Debug(...any)
	Debugf(string, ...any)
	Fatal(...any)
	Fatalf(string, ...any)
	With(...any) Logger
	Sync() error
}

func NewZapLogger(level string) (Logger, error) {
	atom, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, fmt.Errorf("parsing level: %w", err)
	}
	cfg := zap.Config{
		Encoding:         "json",
		Level:            atom,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("build mylogger: %w", err)
	}
	return &sLogger{
		zapLogger.Sugar(),
	}, nil
}

type sLogger struct {
	*zap.SugaredLogger
}

func (s *sLogger) With(args ...any) Logger {
	return &sLogger{
		s.SugaredLogger.With(args...),
	}
}

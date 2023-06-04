package logger

import (
	"fmt"

	"github.com/go-logr/logr"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

type Option func(*zap.Options) error

func WithLogLevel(level string) Option {
	return func(o *zap.Options) error {
		var lvl zapcore.Level
		if err := lvl.UnmarshalText([]byte(level)); err != nil {
			return err
		}
		o.Level = lvl
		return nil
	}
}

func WithTimeEncoder(encoder string) Option {
	return func(o *zap.Options) error {
		var enc zapcore.TimeEncoder
		if err := enc.UnmarshalText([]byte(encoder)); err != nil {
			return err
		}
		o.TimeEncoder = enc
		return nil
	}
}

func WithDevelopment(dev bool) Option {
	return func(o *zap.Options) error {
		o.Development = dev
		return nil
	}
}

func NewLogger(opts ...Option) (logr.Logger, error) {
	zapOpts := zap.Options{
		Level:       zapcore.InfoLevel,
		TimeEncoder: zapcore.EpochTimeEncoder,
		Development: false,
	}
	for _, setpOpt := range opts {
		if err := setpOpt(&zapOpts); err != nil {
			return logr.Logger{}, fmt.Errorf("error setting logger option: %v", err)
		}
	}
	return zap.New(zap.UseFlagOptions(&zapOpts)), nil
}

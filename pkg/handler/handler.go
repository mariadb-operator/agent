package handler

import (
	"github.com/go-logr/logr"
	"github.com/mariadb-operator/agent/pkg/filemanager"
	"github.com/mariadb-operator/agent/pkg/handler/bootstrap"
	"github.com/mariadb-operator/agent/pkg/handler/galerastate"
	"github.com/mariadb-operator/agent/pkg/handler/jsonencoder"
	"github.com/mariadb-operator/agent/pkg/handler/recovery"
)

type Handler struct {
	Bootstrap   *bootstrap.Bootstrap
	GaleraState *galerastate.GaleraState
	Recovery    *recovery.Recovery
}

type Options struct {
	bootstrap []bootstrap.Option
	recovery  []recovery.Option
}

type Option func(*Options)

func WithBootstrapOptions(opts ...bootstrap.Option) Option {
	return func(o *Options) {
		o.bootstrap = append(o.bootstrap, opts...)
	}
}

func WithRecoveryOptions(opts ...recovery.Option) Option {
	return func(o *Options) {
		o.recovery = append(o.recovery, opts...)
	}
}

func NewHandler(fileManager *filemanager.FileManager, logger *logr.Logger, handlerOpts ...Option) *Handler {
	opts := &Options{}
	for _, setOpts := range handlerOpts {
		setOpts(opts)
	}
	bootstrapLogger := logger.WithName("bootstrap")
	galeraStateLogger := logger.WithName("galerastate")
	recoveryLogger := logger.WithName("recovery")

	bootstrap := bootstrap.NewBootstrap(fileManager, &bootstrapLogger, opts.bootstrap...)
	galerastate := galerastate.NewGaleraState(fileManager, jsonencoder.NewJSONEncoder(&galeraStateLogger), &galeraStateLogger)
	recovery := recovery.NewRecover(fileManager, jsonencoder.NewJSONEncoder(&recoveryLogger), &recoveryLogger, opts.recovery...)

	return &Handler{
		Bootstrap:   bootstrap,
		GaleraState: galerastate,
		Recovery:    recovery,
	}
}

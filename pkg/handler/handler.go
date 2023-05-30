package handler

import (
	"github.com/go-logr/logr"
)

type Handler struct {
	GaleraState *GaleraState
	Bootstrap   *Bootstrap
	Recovery    *Recovery
	Mysld       *Mysld
}

func NewHandler(logger *logr.Logger) *Handler {
	galeraStateLogger := logger.WithName("galerastate")
	bootstrapLogger := logger.WithName("bootstrap")
	mysldLogger := logger.WithName("mysqld")
	recoveryLogger := logger.WithName("recovery")

	return &Handler{
		GaleraState: &GaleraState{
			logger: &galeraStateLogger,
		},
		Bootstrap: &Bootstrap{
			logger: &bootstrapLogger,
		},
		Mysld: &Mysld{
			logger: &mysldLogger,
		},
		Recovery: &Recovery{
			logger: &recoveryLogger,
		},
	}
}

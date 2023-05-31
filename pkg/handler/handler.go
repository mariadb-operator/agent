package handler

import (
	"github.com/go-logr/logr"
	"github.com/mariadb-operator/agent/pkg/filemanager"
)

type Handler struct {
	GaleraState *GaleraState
	Bootstrap   *Bootstrap
	Mysld       *Mysld
	Recovery    *Recovery
}

func NewHandler(fileManager *filemanager.FileManager, logger *logr.Logger) *Handler {
	galeraStateLogger := logger.WithName("galerastate")
	bootstrapLogger := logger.WithName("bootstrap")
	mysldLogger := logger.WithName("mysqld")
	recoveryLogger := logger.WithName("recovery")

	return &Handler{
		GaleraState: &GaleraState{
			fileManager: fileManager,
			logger:      &galeraStateLogger,
		},
		Bootstrap: &Bootstrap{
			fileManager: fileManager,
			logger:      &bootstrapLogger,
		},
		Mysld: &Mysld{
			logger: &mysldLogger,
		},
		Recovery: &Recovery{
			fileManager: fileManager,
			logger:      &recoveryLogger,
		},
	}
}

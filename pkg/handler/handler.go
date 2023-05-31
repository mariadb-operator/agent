package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/mariadb-operator/agent/pkg/filemanager"
)

type Handler struct {
	Bootstrap   *Bootstrap
	GaleraState *GaleraState
	Mysld       *Mysld
	Recovery    *Recovery
}

func NewHandler(fileManager *filemanager.FileManager, logger *logr.Logger) *Handler {
	galeraStateLogger := logger.WithName("galerastate")
	bootstrapLogger := logger.WithName("bootstrap")
	mysldLogger := logger.WithName("mysqld")
	recoveryLogger := logger.WithName("recovery")

	return &Handler{
		Bootstrap: &Bootstrap{
			fileManager: fileManager,
			logger:      &bootstrapLogger,
		},
		GaleraState: &GaleraState{
			fileManager: fileManager,
			jsonEncoder: &jsonEncoder{
				logger: &galeraStateLogger,
			},
			logger: &galeraStateLogger,
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

type jsonEncoder struct {
	logger *logr.Logger
}

func (j *jsonEncoder) encode(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		j.logger.Error(err, "error encoding json")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

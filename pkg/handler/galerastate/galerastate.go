package galerastate

import (
	"net/http"
	"os"

	"github.com/go-logr/logr"
	"github.com/mariadb-operator/agent/pkg/filemanager"
	"github.com/mariadb-operator/agent/pkg/galera"
	"github.com/mariadb-operator/agent/pkg/handler/jsonencoder"
)

type GaleraState struct {
	fileManager *filemanager.FileManager
	jsonEncoder *jsonencoder.JSONEncoder
	logger      *logr.Logger
}

func NewGaleraState(fileManager *filemanager.FileManager, jsonEncoder *jsonencoder.JSONEncoder, logger *logr.Logger) *GaleraState {
	return &GaleraState{
		fileManager: fileManager,
		jsonEncoder: jsonEncoder,
		logger:      logger,
	}
}

func (g *GaleraState) Get(w http.ResponseWriter, r *http.Request) {
	bytes, err := g.fileManager.ReadStateFile(galera.GaleraStateFileName)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		g.logger.Error(err, "error reading file")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var galeraState galera.GaleraState
	if err := galeraState.Unmarshal(bytes); err != nil {
		g.logger.Error(err, "error unmarshalling galera state")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	g.jsonEncoder.Encode(w, galeraState)
}

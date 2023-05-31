package handler

import (
	"net/http"
	"os"

	"github.com/go-logr/logr"
	"github.com/mariadb-operator/agent/pkg/filemanager"
	"github.com/mariadb-operator/agent/pkg/galera"
)

type GaleraState struct {
	fileManager *filemanager.FileManager
	jsonEncoder *jsonEncoder
	logger      *logr.Logger
}

func (h *GaleraState) Get(w http.ResponseWriter, r *http.Request) {
	bytes, err := h.fileManager.ReadStateFile(galera.GaleraStateFile)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		h.logger.Error(err, "error reading file")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var galeraState galera.GaleraState
	if err := galeraState.UnmarshalText(bytes); err != nil {
		h.logger.Error(err, "error unmarshalling galera state")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.jsonEncoder.encode(w, galeraState)
}

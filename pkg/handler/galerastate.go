package handler

import (
	"net/http"
	"os"

	"github.com/go-logr/logr"
	"github.com/mariadb-operator/agent/pkg/filemanager"
	"github.com/mariadb-operator/agent/pkg/galerastate"
)

var (
	galeraStateFile = "grastate.dat"
)

type GaleraState struct {
	fileManager *filemanager.FileManager
	jsonEncoder *jsonEncoder
	logger      *logr.Logger
}

func (h *GaleraState) Get(w http.ResponseWriter, r *http.Request) {
	bytes, err := h.fileManager.ReadStateFile(galeraStateFile)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		h.logger.Error(err, "error reading file")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var galeraState galerastate.GaleraState
	if err := galeraState.Unmarshal(bytes); err != nil {
		h.logger.Error(err, "error unmarshalling galera state")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.jsonEncoder.encode(w, galeraState)
}

func (h *GaleraState) Post(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

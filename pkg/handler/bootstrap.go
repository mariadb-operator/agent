package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-logr/logr"
	"github.com/mariadb-operator/agent/pkg/filemanager"
	"github.com/mariadb-operator/agent/pkg/galera"
)

type Bootstrap struct {
	fileManager *filemanager.FileManager
	logger      *logr.Logger
}

func (h *Bootstrap) Put(w http.ResponseWriter, r *http.Request) {
	if err := h.setSafeToBootstrap(); err != nil {
		h.logger.Error(err, "error setting safe to bootstrap")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Bootstrap) Delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Bootstrap) setSafeToBootstrap() error {
	bytes, err := h.fileManager.ReadStateFile(galera.GaleraStateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("error reading galera state: %v", err)
	}

	var galeraState galera.GaleraState
	if err := galeraState.UnmarshalText(bytes); err != nil {
		return fmt.Errorf("error unmarshaling galera state: %v", err)
	}

	galeraState.SafeToBootstrap = true
	bytes, err = galeraState.MarshalText()
	if err != nil {
		return fmt.Errorf("error marshallng galera state: %v", err)
	}

	if err := h.fileManager.WriteStateFile(galera.GaleraStateFile, bytes); err != nil {
		return fmt.Errorf("error writing galera state: %v", err)
	}
	return nil
}
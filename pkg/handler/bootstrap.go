package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-logr/logr"
	"github.com/mariadb-operator/agent/pkg/filemanager"
	"github.com/mariadb-operator/agent/pkg/galera"
	"github.com/mariadb-operator/agent/pkg/mariadbd"
)

type Bootstrap struct {
	fileManager *filemanager.FileManager
	logger      *logr.Logger
}

func (h *Bootstrap) Put(w http.ResponseWriter, r *http.Request) {
	var bootstrap galera.Bootstrap
	if err := json.NewDecoder(r.Body).Decode(&bootstrap); err != nil {
		h.logger.Error(err, "error decoding bootstrap")
		http.Error(w, "invalid body: a valid bootstrap object must be provided", http.StatusBadRequest)
		return
	}

	if err := h.setSafeToBootstrap(&bootstrap); err != nil {
		h.logger.Error(err, "error setting safe to bootstrap")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := h.fileManager.WriteConfigFile(galera.BootstrapFileName, []byte(galera.BootstrapFile)); err != nil {
		h.logger.Error(err, "error writing bootstrap file")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("reloading mariadbd process")
	if err := mariadbd.Reload(); err != nil {
		h.logger.Error(err, "error reloading mariadbd process")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.logger.Info("mariadbd process reloaded")

	w.WriteHeader(http.StatusOK)
}

func (h *Bootstrap) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.fileManager.DeleteConfigFile(galera.BootstrapFileName); err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		h.logger.Error(err, "error deleting bootstrap file")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Bootstrap) setSafeToBootstrap(bootstrap *galera.Bootstrap) error {
	bytes, err := h.fileManager.ReadStateFile(galera.GaleraStateFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("error reading galera state: %v", err)
	}

	var galeraState galera.GaleraState
	if err := galeraState.Unmarshal(bytes); err != nil {
		return fmt.Errorf("error unmarshaling galera state: %v", err)
	}

	galeraState.UUID = bootstrap.UUID
	galeraState.Seqno = bootstrap.Seqno
	galeraState.SafeToBootstrap = true
	bytes, err = galeraState.Marshal()
	if err != nil {
		return fmt.Errorf("error marshallng galera state: %v", err)
	}

	if err := h.fileManager.WriteStateFile(galera.GaleraStateFileName, bytes); err != nil {
		return fmt.Errorf("error writing galera state: %v", err)
	}
	return nil
}

package handler

import (
	"net/http"
	"os"

	"github.com/go-logr/logr"
	"github.com/mariadb-operator/agent/pkg/filemanager"
	"github.com/mariadb-operator/agent/pkg/galera"
	"github.com/mariadb-operator/agent/pkg/mariadbd"
)

type Recovery struct {
	fileManager *filemanager.FileManager
	logger      *logr.Logger
}

func (h *Recovery) Put(w http.ResponseWriter, r *http.Request) {
	if err := h.fileManager.WriteConfigFile(galera.RecoveryFileName, []byte(galera.RecoveryFile)); err != nil {
		h.logger.Error(err, "error writing recovery file")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("reloading mariadbd process")
	if err := mariadbd.Reload(); err != nil {
		h.logger.Error(err, "error reloading mariadbd process")
	} else {
		h.logger.Info("mariadbd process reloaded")
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Recovery) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.fileManager.DeleteConfigFile(galera.RecoveryFileName); err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		h.logger.Error(err, "error deleting recovery file")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

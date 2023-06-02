package handler

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/mariadb-operator/agent/pkg/filemanager"
	"github.com/mariadb-operator/agent/pkg/galera"
	"github.com/mariadb-operator/agent/pkg/mariadbd"
)

var (
	recoverRetries = 10
	recoverWait    = 3 * time.Second
)

type Recovery struct {
	fileManager *filemanager.FileManager
	jsonEncoder *jsonEncoder
	logger      *logr.Logger
}

func (h *Recovery) Put(w http.ResponseWriter, r *http.Request) {
	if err := h.fileManager.DeleteStateFile(galera.RecoveryLog); err != nil && !os.IsNotExist(err) {
		h.logger.Error(err, "error deleting existing recovery log")
	}

	if err := h.fileManager.WriteConfigFile(galera.RecoveryFileName, []byte(galera.RecoveryFile)); err != nil {
		h.logger.Error(err, "error writing recovery config")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("reloading mariadbd process")
	if err := mariadbd.Reload(); err != nil {
		h.logger.Error(err, "error reloading mariadbd process")
	} else {
		h.logger.Info("mariadbd process reloaded")
	}

	recover, err := h.recover()
	if err != nil {
		h.logger.Error(err, "error recovering galera")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := h.fileManager.DeleteConfigFile(galera.RecoveryFileName); err != nil {
		h.logger.Error(err, "error deleting recovery file")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.jsonEncoder.encode(w, recover)
}

func (h *Recovery) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.fileManager.DeleteConfigFile(galera.RecoveryFileName); err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Recovery) recover() (*galera.Recover, error) {
	for i := 0; i < recoverRetries; i++ {
		time.Sleep(recoverWait)

		bytes, err := h.fileManager.ReadStateFile(galera.RecoveryLog)
		if err != nil {
			h.logger.Error(err, "error recovering galera from recovery log", "retry", i, "max-retries", recoverRetries)
			continue
		}

		var recover galera.Recover
		err = recover.UnmarshalText(bytes)
		if err == nil {
			return &recover, nil
		}

		h.logger.Error(err, "error recovering galera from recovery log", "retry", i, "max-retries", recoverRetries)
	}
	return nil, fmt.Errorf("maximum retries (%d) reached attempting to recover galera from recovery log", recoverRetries)
}

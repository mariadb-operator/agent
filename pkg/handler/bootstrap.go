package handler

import (
	"net/http"

	"github.com/go-logr/logr"
	"github.com/mariadb-operator/agent/pkg/filemanager"
)

type Bootstrap struct {
	fileManager *filemanager.FileManager
	logger      *logr.Logger
}

func (h *Bootstrap) Put(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Bootstrap) Delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

package handler

import (
	"net/http"

	"github.com/go-logr/logr"
	"github.com/mariadb-operator/agent/pkg/filemanager"
)

type Recovery struct {
	fileManager *filemanager.FileManager
	logger      *logr.Logger
}

func (h *Recovery) Get(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Recovery) Put(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Recovery) Delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

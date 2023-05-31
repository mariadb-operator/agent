package handler

import (
	"net/http"

	"github.com/go-logr/logr"
	"github.com/mariadb-operator/agent/pkg/filemanager"
)

type GaleraState struct {
	fileManager *filemanager.FileManager
	logger      *logr.Logger
}

func (h *GaleraState) Get(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *GaleraState) Post(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

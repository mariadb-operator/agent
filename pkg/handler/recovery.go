package handler

import (
	"net/http"

	"github.com/go-logr/logr"
)

type Recovery struct {
	logger *logr.Logger
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

package handler

import (
	"net/http"

	"github.com/go-logr/logr"
)

type Bootstrap struct {
	logger *logr.Logger
}

func (h *Bootstrap) Put(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Bootstrap) Delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

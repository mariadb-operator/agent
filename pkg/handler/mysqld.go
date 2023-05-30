package handler

import (
	"net/http"

	"github.com/go-logr/logr"
)

type Mysld struct {
	logger *logr.Logger
}

func (h *Mysld) Post(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

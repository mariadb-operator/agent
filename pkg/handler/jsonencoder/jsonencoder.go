package jsonencoder

import (
	"encoding/json"
	"net/http"

	"github.com/go-logr/logr"
)

type JSONEncoder struct {
	logger *logr.Logger
}

func NewJSONEncoder(logger *logr.Logger) *JSONEncoder {
	return &JSONEncoder{
		logger: logger,
	}
}

func (j *JSONEncoder) Encode(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		j.logger.Error(err, "error encoding json")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

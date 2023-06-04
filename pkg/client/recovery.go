package client

import (
	"context"
	"net/http"

	"github.com/mariadb-operator/agent/pkg/galera"
)

type Recovery struct {
	*Client
}

func (r *Recovery) Start(ctx context.Context) (*galera.Bootstrap, error) {
	req, err := r.newRequestWithContext(ctx, http.MethodPut, "/api/recovery", nil)
	if err != nil {
		return nil, err
	}
	var bootstrap galera.Bootstrap
	if err := r.do(req, &bootstrap); err != nil {
		return nil, err
	}
	return &bootstrap, nil
}

func (r *Recovery) Stop(ctx context.Context) error {
	req, err := r.newRequestWithContext(ctx, http.MethodDelete, "/api/recovery", nil)
	if err != nil {
		return err
	}
	return r.do(req, nil)
}

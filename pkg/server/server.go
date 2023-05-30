package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-logr/logr"
)

type Server struct {
	httpServer *http.Server
	logger     *logr.Logger
}

func NewServer(addr string, handler http.Handler, logger *logr.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		logger: logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	serverContext, stopServer := context.WithCancel(ctx)
	errChan := make(chan error)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		defer stopServer()

		s.logger.Info("shutting down server")
		if err := s.httpServer.Shutdown(serverContext); err != nil {
			errChan <- fmt.Errorf("error shutting down server: %v", err)
		}
	}()

	s.logger.Info("server listening", "addr", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		errChan <- fmt.Errorf("error starting server: %v", err)
	}

	select {
	case <-serverContext.Done():
		return nil
	case err := <-errChan:
		return err
	}
}

package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-logr/logr"
)

type Option func(*Server)

func WithGracefulShutdown(timeout time.Duration) Option {
	return func(s *Server) {
		s.gracefulShutdown = timeout
	}
}

type Server struct {
	httpServer       *http.Server
	logger           *logr.Logger
	gracefulShutdown time.Duration
}

func NewServer(addr string, handler http.Handler, logger *logr.Logger, opts ...Option) *Server {
	srv := &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		logger:           logger,
		gracefulShutdown: 30 * time.Second,
	}
	for _, setOpt := range opts {
		setOpt(srv)
	}
	return srv
}

func (s *Server) Start(ctx context.Context) error {
	serverContext, stopServer := context.WithCancel(ctx)
	errChan := make(chan error)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		defer stopServer()

		shutdownCtx, cancel := context.WithTimeout(serverContext, s.gracefulShutdown)
		defer cancel()
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				s.logger.Info("graceful shutdown timed out... forcing exit")
				os.Exit(1)
			}
		}()

		s.logger.Info("shutting down server")
		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			errChan <- fmt.Errorf("error shutting down server: %v", err)
		}
	}()

	s.logger.Info("server listening", "addr", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		errChan <- fmt.Errorf("error starting server: %v", err)
	}

	select {
	case err := <-errChan:
		return err
	case <-serverContext.Done():
		return nil
	}
}

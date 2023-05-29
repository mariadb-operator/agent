package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
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

		log.Println("shutting down server")
		if err := s.httpServer.Shutdown(serverContext); err != nil {
			errChan <- fmt.Errorf("error shutting down server: %v", err)
		}
	}()

	log.Printf("server listening at %s", s.httpServer.Addr)
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

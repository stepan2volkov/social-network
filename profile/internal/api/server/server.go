package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	srv     *http.Server
	errChan chan error
}

func New(
	handler http.Handler,
	port int,
	writeTimeout, readTimeout time.Duration,
) *Server {
	addr := fmt.Sprintf(":%d", port)

	return &Server{
		srv: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
		errChan: make(chan error, 1),
	}
}

func (s *Server) Run() <-chan error {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				s.errChan <- err
			}
		}
	}()
	return s.errChan
}

func (s *Server) Stop(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		s.errChan <- err
	}

	close(s.errChan)
}

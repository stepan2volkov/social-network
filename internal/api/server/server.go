package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/stepan2volkov/social-network/internal/config"
)

type Server struct {
	srv     http.Server
	errChan chan error
}

// NewServer creates http.Server with settings from config.Config
func NewServer(conf config.Config, h http.Handler) *Server {
	s := &Server{
		errChan: make(chan error, 1),
	}
	s.srv = http.Server{
		Addr:              ":" + conf.AppPort,
		Handler:           h,
		ReadTimeout:       time.Duration(conf.ReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(conf.WriteTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(conf.ReadHeaderTimeout) * time.Second,
	}
	return s
}

func (s *Server) Start() <-chan error {
	go func() {
		log.Println("app started")
		if err := s.srv.ListenAndServe(); err != nil {
			s.errChan <- err
		}
		close(s.errChan)
	}()
	return s.errChan
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}

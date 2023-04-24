package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/stepan2volkov/social-network/profile/internal/api/router"
	"github.com/stepan2volkov/social-network/profile/internal/api/server"
	"github.com/stepan2volkov/social-network/profile/internal/app"
	"github.com/stepan2volkov/social-network/profile/internal/config"
)

const (
	shutdownTimeout = 30 * time.Second
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("Error to getting config: %s\n", err)
	}

	app := app.New()

	router := router.New(app)

	srv := server.New(router, cfg.Server.Port, cfg.Server.WriteTimeout, cfg.Server.ReadTimeout)
	log.Printf("Starting server on port http://0.0.0.0:%d\n", cfg.Server.Port)

	srvErrChan := srv.Run()

	select {
	case <-ctx.Done():
		srv.Stop(shutdownTimeout)
	case err := <-srvErrChan:
		log.Printf("Error to start server: %s\n", err)
	}

	for err := range srvErrChan {
		log.Printf("Starting server error: %s\n", err)
	}

	log.Println("Server stopped")
}

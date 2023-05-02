package main

import (
	"context"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"

	"github.com/stepan2volkov/social-network/profile/internal/app/userapp"
	"github.com/stepan2volkov/social-network/profile/internal/config"
	"github.com/stepan2volkov/social-network/profile/internal/domain"
	"github.com/stepan2volkov/social-network/profile/internal/repo/userrepo"
	"github.com/stepan2volkov/social-network/profile/internal/router"
	"github.com/stepan2volkov/social-network/profile/internal/server"
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

	userDB, err := userrepo.InitDB(cfg.DB.DSN)
	if err != nil {
		log.Fatalf("error to init user db connection: %s\n", err)
	}
	defer userDB.Close()

	// Migrate.
	if err = userrepo.Migrate(userDB); err != nil {
		log.Fatalf("error to migrate: %s", err)
	}

	// Initializing storage layout.
	userRepo := userrepo.New(userDB)

	// Initializing application layout.
	userApp := userapp.New(userRepo)

	// Initializing session management.
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(userDB)
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Secure = true
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	gob.Register(domain.UserIdentity{})

	// Initializing external API layout.
	router := router.New(
		sessionManager,
		userApp,
	)

	srv := server.New(router, cfg.Server.Port, cfg.Server.WriteTimeout, cfg.Server.ReadTimeout)
	log.Printf("Starting server on port http://0.0.0.0:%d\n", cfg.Server.Port)

	srvErrChan := srv.Run()

	select {
	case err := <-srvErrChan:
		log.Printf("Error to start server: %s\n", err)
	case <-ctx.Done():
	}

	srv.Stop(shutdownTimeout)

	for err := range srvErrChan {
		log.Printf("Starting server error: %s\n", err)
	}

	log.Println("Server stopped")
}

package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/stepan2volkov/social-network/internal/api/router"
	"github.com/stepan2volkov/social-network/internal/api/server"
	"github.com/stepan2volkov/social-network/internal/app/authapp"
	"github.com/stepan2volkov/social-network/internal/app/profileapp"
	"github.com/stepan2volkov/social-network/internal/config"
	"github.com/stepan2volkov/social-network/internal/store/mysqlstore/mysqlauthstore"
	"github.com/stepan2volkov/social-network/internal/store/mysqlstore/mysqlprofilestore"
	"github.com/stepan2volkov/social-network/internal/store/mysqlstore/mysqlstarter"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("error when getting config: %v", err)
	}

	db, err := mysqlstarter.GetDB(cfg.DbDsn)
	if err != nil {
		log.Fatalf("error when getting db connection: %v", err)
	}

	// Init storages
	authStore := mysqlauthstore.New(db)
	profileStore := mysqlprofilestore.New(db)

	// Init applications
	authApp := authapp.New(
		"social-network-auth",
		authStore,
		[]byte(cfg.SecretKey),
		cfg.SecretExpiredIn,
		[]byte(cfg.RefreshSecretKey),
		cfg.RefreshSecretExpiredIn)
	profileApp := profileapp.New(profileStore)

	// Version
	version := router.VersionInfo{
		BuildCommit: config.BuildCommit,
		BuildTime:   config.BuildTime,
	}

	// Init API
	rt := router.New(version, authApp, profileApp)
	srv := server.NewServer(cfg, rt)

	errChan := srv.Start()

	select {
	case err = <-errChan:
		log.Fatalln(err)
	case <-ctx.Done():
		srv.Stop()
		log.Println("App has been stopped")
	}
}

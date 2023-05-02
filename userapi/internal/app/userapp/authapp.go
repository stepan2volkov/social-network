package userapp

import (
	"github.com/stepan2volkov/social-network/profile/internal/repo/userrepo"
)

type App struct {
	repo *userrepo.Repository
}

func New(repo *userrepo.Repository) *App {
	return &App{
		repo: repo,
	}
}

package userapp

import (
	"context"

	"github.com/stepan2volkov/social-network/profile/internal/domain"
)

func (a *App) SearchProfiles(ctx context.Context, firstname, lastname string) ([]domain.Profile, error) {
	ret, err := a.repo.SearchProfiles(ctx, firstname, lastname)
	if err != nil {
		return nil, err
	}

	return ret, nil

}

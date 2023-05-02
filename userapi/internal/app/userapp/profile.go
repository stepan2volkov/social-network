package userapp

import (
	"context"

	"github.com/stepan2volkov/social-network/profile/internal/domain"
)

func (a *App) GetProfileByUsername(ctx context.Context, username string) (domain.Profile, error) {
	user, err := a.repo.GetByUsername(ctx, username)
	if err != nil {
		return domain.Profile{}, err
	}

	return user.Profile, nil
}

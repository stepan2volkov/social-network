package userapp

import (
	"context"

	"github.com/stepan2volkov/social-network/profile/internal/domain"
)

func (a *App) Login(ctx context.Context, creds domain.Credentials) (domain.UserIdentity, error) {
	user, err := a.repo.GetByUsername(ctx, creds.Username)
	if err != nil {
		return domain.UserIdentity{}, err
	}

	if !domain.CheckPassword(creds.Password, user.PasswordHash) {
		return domain.UserIdentity{}, domain.ErrUnauthorized
	}

	return user.UserIdentity, nil
}

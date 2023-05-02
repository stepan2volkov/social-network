package userapp

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/stepan2volkov/social-network/profile/internal/domain"
)

func (a *App) Register(
	ctx context.Context,
	creds domain.Credentials,
	profile domain.Profile,
) (
	uuid.UUID,
	error,
) {
	userID, err := domain.GenerateUserID()
	if err != nil {
		return uuid.Nil, fmt.Errorf("error to generate user id: %w", err)
	}

	passwordHash, err := domain.PasswordHash(creds.Password)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error to calculate password hash: %w", err)
	}

	user := domain.User{
		UserIdentity: domain.UserIdentity{
			ID:       userID,
			Username: creds.Username,
		},
		PasswordHash: passwordHash,
		Profile:      profile,
	}

	if err := a.repo.Save(ctx, user); err != nil {
		return uuid.Nil, fmt.Errorf("error to save user: %w", err)
	}

	return userID, nil
}

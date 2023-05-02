package userapp

import (
	"context"
	"fmt"
)

func (a *App) Follow(ctx context.Context, follower, followee string) error {
	if _, err := a.repo.GetByUsername(ctx, followee); err != nil {
		return fmt.Errorf("[app.Follow] error to get followee: %w", err)
	}
	if err := a.repo.Follow(ctx, follower, followee); err != nil {
		return fmt.Errorf("[app.Follow] error to save in repo: %w", err)
	}

	return nil
}

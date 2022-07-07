package profileapp

import (
	"context"

	"github.com/stepan2volkov/social-network/internal/entities/user"
)

type ProfileProvider interface {
	GetProfileByUserID(ctx context.Context, id user.UserID) (user.UserProfile, error)
	GetProfileByUsername(ctx context.Context, username string) (user.UserProfile, error)
	GetUserIDByUsername(ctx context.Context, username string) (user.UserID, error)
	Subscribe(ctx context.Context, followerID, leaderID user.UserID) error
	GetFriendsByFollowerID(ctx context.Context, followerID user.UserID) ([]user.Friend, error)
}

type App struct {
	provider ProfileProvider
}

func New(provider ProfileProvider) *App {
	return &App{
		provider: provider,
	}
}

func (a *App) GetProfileByUsername(ctx context.Context, username string) (user.UserProfile, error) {
	return a.provider.GetProfileByUsername(ctx, username)
}

func (a *App) Subscribe(ctx context.Context, followerUsername, leaderUsername string) error {
	followerID, err := a.provider.GetUserIDByUsername(ctx, followerUsername)
	if err != nil {
		return err
	}
	leaderID, err := a.provider.GetUserIDByUsername(ctx, leaderUsername)
	if err != nil {
		return err
	}

	return a.provider.Subscribe(ctx, followerID, leaderID)
}

func (a *App) GetLeaders(ctx context.Context, username string) ([]user.Friend, error) {
	followerID, err := a.provider.GetUserIDByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return a.provider.GetFriendsByFollowerID(ctx, followerID)
}

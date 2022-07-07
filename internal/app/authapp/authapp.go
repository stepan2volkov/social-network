package authapp

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/stepan2volkov/social-network/internal/entities/token"
	"github.com/stepan2volkov/social-network/internal/entities/user"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

// AuthProvider responsible for saving
type AuthProvider interface {
	CreateUser(ctx context.Context, u user.User) (user.UserID, error)
	GetUserByUsername(ctx context.Context, username string) (user.User, error)
}

type App struct {
	name                   string
	provider               AuthProvider
	secretKey              []byte
	secretExpiredIn        time.Duration
	refreshSecretKey       []byte
	refreshSecretExpiredIn time.Duration
}

func New(
	name string,
	provider AuthProvider,
	secretKey []byte,
	secretExpiredIn time.Duration,
	refreshSecretKey []byte,
	refreshSecretExpiredIn time.Duration,
) *App {
	return &App{
		name:                   name,
		provider:               provider,
		secretKey:              secretKey,
		secretExpiredIn:        secretExpiredIn,
		refreshSecretKey:       refreshSecretKey,
		refreshSecretExpiredIn: refreshSecretExpiredIn,
	}
}

// Register tries to register user and password
func (a *App) Register(ctx context.Context, u user.UserProfile, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error when hashing password: %w", err)
	}

	_, err = a.provider.CreateUser(ctx, user.User{
		UserProfile:  u,
		PasswordHash: string(hash),
	})
	if err != nil {
		return fmt.Errorf("error in provider when creating user: %w", err)
	}

	return nil
}

// Login returns token if password is correct
func (a *App) Login(ctx context.Context, username, password string) (Tokens, error) {
	u, err := a.provider.GetUserByUsername(ctx, username)
	if errors.Is(err, user.ErrUserNotFound) {
		return Tokens{}, user.ErrInvalidCredentials
	}
	if err != nil {
		return Tokens{}, fmt.Errorf("error when getting user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return Tokens{}, user.ErrInvalidCredentials
	}

	return a.createTokens(ctx, u.UserProfile)
}

func (a *App) Refresh(ctx context.Context, tokens Tokens) (Tokens, error) {
	accessClaims, err := token.GetClaims(tokens.AccessToken, a.secretKey)
	if err != nil {
		return Tokens{}, err
	}
	if accessClaims.ExpiresAt.Before(time.Now()) {
		return tokens, nil
	}

	refreshClaims, err := token.GetClaims(tokens.RefreshToken, a.refreshSecretKey)
	if err != nil {
		return Tokens{}, err
	}
	if refreshClaims.ExpiresAt.Before(time.Now()) {
		return Tokens{}, token.ErrTokenExpired
	}

	tokens.AccessToken, err = token.UpdateToken(
		tokens.AccessToken,
		a.secretKey,
		a.secretExpiredIn,
	)
	if err != nil {
		return Tokens{}, nil
	}

	if refreshClaims.ExpiresAt.Before(time.Now().Add(a.refreshSecretExpiredIn / 2)) {
		return tokens, nil
	}

	tokens.RefreshToken, err = token.UpdateToken(
		tokens.RefreshToken,
		a.refreshSecretKey,
		a.refreshSecretExpiredIn,
	)
	if err != nil {
		return Tokens{}, nil
	}

	return tokens, nil

}

// Verify returns nil if token is correct
func (a *App) Verify(ctx context.Context, accessToken string) (user.UserProfile, error) {
	accessClaims, err := token.GetClaims(accessToken, a.secretKey)
	if err != nil {
		return user.UserProfile{}, err
	}

	if !accessClaims.VerifyExpiresAt(time.Now(), true) {
		return user.UserProfile{}, token.ErrTokenExpired
	}

	return accessClaims.UserProfile, nil
}

func (a *App) createTokens(ctx context.Context, u user.UserProfile) (Tokens, error) {
	accessClaims := token.CreateClaims(a.name, a.secretExpiredIn, u)
	accessToken, err := token.CreateToken(accessClaims, a.secretKey)
	if err != nil {
		return Tokens{}, fmt.Errorf("error when creating token: %w", err)
	}

	refreshClaims := token.CreateClaims(a.name, a.refreshSecretExpiredIn, u)
	refreshToken, err := token.CreateToken(refreshClaims, a.refreshSecretKey)
	if err != nil {
		return Tokens{}, fmt.Errorf("error when creating token: %w", err)
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

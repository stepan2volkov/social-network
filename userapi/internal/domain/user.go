package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUnauthorized  = errors.New("user is unauthorized")
	ErrAlreadyFollow = errors.New("already follow")
)

type Credentials struct {
	Username string
	Password string
}

type UserIdentity struct {
	ID       uuid.UUID
	Username string
}

type User struct {
	UserIdentity
	PasswordHash string
	Profile
}

type Profile struct {
	Firstname string
	Lastname  string
	Birthdate time.Time
	Biography string
	City      string
}

func GenerateUserID() (uuid.UUID, error) {
	return uuid.NewRandom()
}

func PasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

type userIdentity struct{}

func SetUserIdentity(ctx context.Context, id UserIdentity) context.Context {
	return context.WithValue(ctx, userIdentity{}, id)
}

func GetUserIdentity(ctx context.Context) (UserIdentity, error) {
	ret, ok := ctx.Value(userIdentity{}).(UserIdentity)
	if !ok {
		return UserIdentity{}, ErrUnauthorized
	}
	return ret, nil
}

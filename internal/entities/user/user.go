package user

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("user or password is invalid")
)

type UserID uint64

type Gender string

const (
	Male   Gender = "m"
	Female Gender = "f"
)

type User struct {
	UserProfile
	PasswordHash string
}

type UserProfile struct {
	ID        UserID
	Username  string
	CreatedAt time.Time
	Firstname string
	Lastname  string
	Birthdate time.Time
	Gender    Gender
	CityID    int
}

type Status string

const (
	StatusLeader Status = "LEADER"
	StatusFriend Status = "FRIEND"
)

type Friend struct {
	UserProfile
	Status Status
}

package mysqlauthstore

import (
	"context"
	"database/sql"

	"github.com/stepan2volkov/social-network/internal/app/authapp"
	"github.com/stepan2volkov/social-network/internal/entities/user"
)

var _ authapp.AuthProvider = &AuthMapper{}

type AuthMapper struct {
	db *sql.DB
}

func New(db *sql.DB) *AuthMapper {
	return &AuthMapper{
		db: db,
	}
}

// CreateUserCredentials implements authapp.AuthProvider
func (m *AuthMapper) CreateUser(
	ctx context.Context,
	u user.User,
) (
	user.UserID,
	error,
) {
	res, err := m.db.ExecContext(ctx, `
		INSERT INTO users (username, password_hash, firstname, lastname, birthdate, gender, city_id) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		u.Username,
		u.PasswordHash,
		u.Firstname,
		u.Lastname,
		u.Birthdate,
		u.Gender,
		u.CityID,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return user.UserID(id), nil
}

// GetUserCredentianlsByUsername implements authapp.AuthProvider
func (m *AuthMapper) GetUserByUsername(
	ctx context.Context,
	username string,
) (
	user.User,
	error,
) {
	row := m.db.QueryRowContext(ctx, `
		SELECT id, username, password_hash, created_at, firstname, lastname, birthdate, gender, city_id
		FROM users
		WHERE username = ?`, username)

	ret := user.User{}
	if err := row.Scan(
		&ret.ID,
		&ret.Username,
		&ret.PasswordHash,
		&ret.CreatedAt,
		&ret.Firstname,
		&ret.Lastname,
		&ret.Birthdate,
		&ret.Gender,
		&ret.CityID,
	); err != nil {
		return user.User{}, err
	}

	return ret, nil
}

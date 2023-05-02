package userrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/stepan2volkov/social-network/profile/internal/domain"
)

func (r *Repository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT
			id, username, password_hash, firstname, lastname, birthdate, biography, city
		FROM profiles
		WHERE username = ?`, username)

	var ret domain.User

	if err := row.Scan(
		&ret.ID,
		&ret.Username,
		&ret.PasswordHash,
		&ret.Firstname,
		&ret.Lastname,
		&ret.Birthdate,
		&ret.Biography,
		&ret.City,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("ProfileRepository.GetByUsername: %w", err)
	}

	return ret, nil
}

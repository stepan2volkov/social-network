package userrepo

import (
	"context"
	"fmt"

	"github.com/stepan2volkov/social-network/profile/internal/domain"
)

func (r *Repository) Save(ctx context.Context, profile domain.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO profiles (
			id, username, password_hash, firstname, lastname, birthdate, biography, city)
		VALUES (UUID_TO_BIN(?), ?, ?, ?, ?, ?, ?, ?)`,
		profile.ID,
		profile.Username,
		profile.PasswordHash,
		profile.Firstname,
		profile.Lastname,
		profile.Birthdate,
		profile.Biography,
		profile.City,
	)

	if err != nil {
		return fmt.Errorf("ProfileRepository.Save: %w", err)
	}

	return nil
}

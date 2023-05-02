package userrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/stepan2volkov/social-network/profile/internal/domain"
)

func (r *Repository) Follow(ctx context.Context, follower, followee string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO followers (
			follower, followee
		) VALUES (?, ?)`, follower, followee,
	)

	if err != nil {
		if merr := mapError(err); errors.Is(merr, errNotUnique) {
			return domain.ErrAlreadyFollow
		}
		return fmt.Errorf("ProfileRepository.Follow: %w", err)
	}

	return nil
}

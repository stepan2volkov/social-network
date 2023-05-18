package userrepo

import (
	"context"

	"github.com/georgysavva/scany/v2/sqlscan"

	"github.com/stepan2volkov/social-network/profile/internal/domain"
)

func (r *Repository) SearchProfiles(ctx context.Context, firstname, lastname string) ([]domain.Profile, error) {
	var ret []domain.Profile

	if err := sqlscan.Select(ctx, r.db, &ret, `
		SELECT
			firstname, lastname, birthdate, biography, city
		FROM profiles
		WHERE firstname LIKE CONCAT(?, '%') AND lastname like CONCAT(?, '%')`,
		firstname, lastname,
	); err != nil {
		return nil, err
	}

	return ret, nil
}

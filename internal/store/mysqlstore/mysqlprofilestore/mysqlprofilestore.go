package mysqlprofilestore

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/stepan2volkov/social-network/internal/app/profileapp"
	"github.com/stepan2volkov/social-network/internal/entities/user"
)

var _ profileapp.ProfileProvider = &ProfileMapper{}

type ProfileMapper struct {
	db *sql.DB
}

func New(db *sql.DB) *ProfileMapper {
	return &ProfileMapper{
		db: db,
	}
}

// GetProfileByUserID implements profileapp.ProfileProvider
func (m *ProfileMapper) GetProfileByUserID(
	ctx context.Context,
	id user.UserID,
) (
	user.UserProfile,
	error,
) {
	ret := user.UserProfile{}
	if err := m.db.QueryRowContext(ctx, `
			SELECT id, username, created_at, firstname, lastname, birthdate, gender, city_id 
			FROM users
			WHERE id = ?`, id).
		Scan(
			&ret.ID,
			&ret.Username,
			&ret.CreatedAt,
			&ret.Firstname,
			&ret.Lastname,
			&ret.Birthdate,
			&ret.Gender,
			&ret.CityID); err != nil {

		return user.UserProfile{}, err
	}
	return ret, nil
}

// GetProfileByUsername implements profileapp.ProfileProvider
func (m *ProfileMapper) GetProfileByUsername(
	ctx context.Context,
	username string,
) (
	user.UserProfile,
	error,
) {
	ret := user.UserProfile{}

	if err := m.db.QueryRowContext(ctx, `
			SELECT id, username, created_at, firstname, lastname, birthdate, gender, city_id 
			FROM users
			WHERE username = ?`, username).Scan(
		&ret.ID,
		&ret.Username,
		&ret.CreatedAt,
		&ret.Firstname,
		&ret.Lastname,
		&ret.Birthdate,
		&ret.Gender,
		&ret.CityID); err != nil {

		return user.UserProfile{}, err
	}

	return ret, nil
}

// GetUserIDByUsername implements profileapp.ProfileProvider
func (m *ProfileMapper) GetUserIDByUsername(
	ctx context.Context,
	username string,
) (
	user.UserID,
	error,
) {
	var ret user.UserID
	if err := m.db.QueryRowContext(ctx, `SELECT id FROM users WHERE username = ?`, username).
		Scan(&ret); err != nil {

		return 0, err
	}
	return ret, nil
}

// Subscribe implements profileapp.ProfileProvider
func (m *ProfileMapper) Subscribe(ctx context.Context, followerID, leaderID user.UserID) error {
	_, err := m.db.ExecContext(ctx, `
		INSERT INTO followers (follower_id, leader_id)
		VALUES (?, ?)`, followerID, leaderID)

	if err != nil {
		return err
	}

	return nil
}

// GetLeaders implements profileapp.ProfileProvider
func (m *ProfileMapper) GetFriendsByFollowerID(ctx context.Context, followerID user.UserID) ([]user.Friend, error) {
	ret := make([]user.Friend, 0, 1)
	rows, err := m.db.QueryContext(ctx, `
		SELECT	users.id,
				users.username,
				users.created_at,
				users.firstname,
				users.lastname,
				users.birthdate,
				users.gender,
				users.city_id,
				CASE WHEN followers.follower_id IS NULL THEN 'LEADER' ELSE 'FRIEND' END AS status
			FROM followers leaders
			INNER JOIN users ON leaders.leader_id = users.id
			LEFT JOIN followers ON leaders.follower_id = followers.leader_id AND
							leaders.leader_id = followers.follower_id
			WHERE leaders.follower_id = ?`, followerID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := user.Friend{}
		if err = rows.Scan(
			&p.ID,
			&p.Username,
			&p.CreatedAt,
			&p.Firstname,
			&p.Lastname,
			&p.Birthdate,
			&p.Gender,
			&p.CityID,
			&p.Status,
		); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}

	return ret, nil
}

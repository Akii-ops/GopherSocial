package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Follower struct {
	UserID     int64  `json:"userID"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followedID, userID int64) error {

	query := `
		INSERT INTO followers (user_id, follower_id)
		VALUES ($1, $2);
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, followedID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrConflict
		}
		return err
	}

	return nil
}

// TODO: 逻辑返回,不存在的关注不能取关
func (s *FollowerStore) UnFollow(ctx context.Context, unfollowedID, userID int64) error {
	query := `
		DELETE FROM followers
		WHERE user_id = $1 AND follower_id = $2;
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, unfollowedID)
	if err != nil {
		return err
	}

	return nil
}

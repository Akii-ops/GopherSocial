package cache

import (
	"backend/internal/store"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type UserCache struct {
	rdb *redis.Client
}

func (s *UserCache) Get(ctx context.Context, userID int64) (*store.User, error) {

	cacheKey := fmt.Sprintf("user-%v", userID)

	data, err := s.rdb.Get(ctx, cacheKey).Result()

	// redis nil
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var user store.User

	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}
func (s *UserCache) Set(ctx context.Context, user *store.User) error {
	cacheKey := fmt.Sprintf("user-%v", user.ID)

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.rdb.SetEX(ctx, cacheKey, json, UserExpireTime).Err()

}

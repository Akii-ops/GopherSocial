package cache

import (
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	UserExpireTime = time.Minute * 5
)

func NewRedisClient(addr, pw string, db int) *redis.Client {
	return redis.NewClient(
		&redis.Options{
			Addr:     addr,
			Password: pw,
			DB:       db,
		},
	)
}

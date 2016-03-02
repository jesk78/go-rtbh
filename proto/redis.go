package proto

import (
	"gopkg.in/redis.v3"
)

func NewRedisClient() (rc *redis.Client, err error) {
	rc = redis.NewClient(&redis.Options{
		Addr:     Config.Redis.Address,
		Password: Config.Redis.Password,
		DB:       Config.Redis.Database,
	})

	return
}

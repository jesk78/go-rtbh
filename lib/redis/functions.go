package redis

import (
	"gopkg.in/redis.v3"
)

func (r *RedisClient) Connect() {
	r.client = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})

	if r.client.Ping().Val() == "PONG" {
		log.Debugf("RedisClient.Connect: Connected to " + cfg.Redis.Address)
	} else {
		log.Warningf("RedisClient.Connect: Failed to connect to " + cfg.Redis.Address)
	}
}

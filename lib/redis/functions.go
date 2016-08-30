package redis

import (
	"gopkg.in/redis.v3"
)

func (r *RedisClient) Connect() {
	r.client = redis.NewClient(&redis.Options{
		Addr:     Config.Redis.Address,
		Password: Config.Redis.Password,
		DB:       Config.Redis.Database,
	})

	if r.client.Ping().Val() == "PONG" {
		Log.Debug(MYNAME + ": Connected to " + Config.Redis.Address)
	} else {
		Log.Warning(MYNAME + ": Failed to connect to " + Config.Redis.Address)
	}
}

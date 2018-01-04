package redis

import (
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
)

func New(l *logger.Logger, c *config.Config) *RedisClient {
	log = l
	cfg = c

	redis := &RedisClient{
		Events:  make(chan []byte, config.D_REDIS_BUFSIZE),
		Control: make(chan int, config.D_CONTROL_BUFSIZE),
		Done:    make(chan bool, config.D_DONE_BUFSIZE),
	}

	redis.Connect()

	log.Debugf("RedisClient: Module initialized")

	return redis
}

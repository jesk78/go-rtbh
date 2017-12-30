package redis

import (
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/go-rtbh/lib/logger"
	redis "gopkg.in/redis.v3"
)

var (
	cfg *config.Config
	log *logger.Logger
)

type RedisClient struct {
	client  *redis.Client
	Events  chan []byte
	Control chan int
	Done    chan bool
}

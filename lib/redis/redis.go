package redis

import (
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/rlib/logger"
	"gopkg.in/redis.v3"
)

const MYNAME string = "Redis"

var Config *config.Config
var Log logger.Log

type RedisClient struct {
	client  *redis.Client
	Events  chan []byte
	Control chan int
	Done    chan bool
}

func Setup(l logger.Log, c *config.Config) (err error) {
	Log = l
	Config = c

	Log.Debug(MYNAME + ": Module initialized")
	return
}

func New() *RedisClient {
	var redis *RedisClient

	redis = &RedisClient{
		Events:  make(chan []byte, config.D_REDIS_BUFSIZE),
		Control: make(chan int, config.D_CONTROL_BUFSIZE),
		Done:    make(chan bool, config.D_DONE_BUFSIZE),
	}

	return redis
}

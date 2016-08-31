package amqp

import (
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/rlib/logger"
	"github.com/streadway/amqp"
)

const MYNAME string = "AMQP"

var Config config.Config
var Log logger.Log

type AmqpClient struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
	Events     chan []byte
	Control    chan int
	Done       chan bool
}

func Setup(l logger.Log, c config.Config) (err error) {
	Log = l
	Config = c

	return
}

func New() *AmqpClient {
	var amqp *AmqpClient

	amqp = &AmqpClient{
		Events:  make(chan []byte, config.D_AMQP_BUFSIZE),
		Control: make(chan int, config.D_CONTROL_BUFSIZE),
		Done:    make(chan bool, config.D_DONE_BUFSIZE),
	}
	amqp.Connect()

	return amqp
}

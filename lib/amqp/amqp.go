package amqp

import (
	"fmt"

	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/go-rtbh/lib/logger"
)

var (
	cfg *config.Config
	log *logger.Logger
)

func NewAmqpClient(l *logger.Logger, c *config.Config) (*AmqpClient, error) {
	log = l
	cfg = c

	amqp := &AmqpClient{
		Events:  make(chan []byte, config.D_AMQP_BUFSIZE),
		Control: make(chan int, config.D_CONTROL_BUFSIZE),
		Done:    make(chan bool, config.D_DONE_BUFSIZE),
	}

	err := amqp.Connect()
	if err != nil {
		return nil, fmt.Errorf("NewAmqpClient: %v", err)
	}

	log.Debugf("AmqpClient: Module initialized")

	return amqp, nil
}

package amqp

import "github.com/streadway/amqp"

type AmqpClient struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
	Events     chan []byte
	Control    chan int
	Done       chan bool
}

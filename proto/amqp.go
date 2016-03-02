package proto

import (
	"errors"
	"github.com/r3boot/go-rtbh/config"
	"github.com/streadway/amqp"
)

type AmqpClient struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
	Events     chan []byte
	Control    chan int
	Done       chan bool
}

func NewAmqpClient() (ac *AmqpClient, err error) {
	var url string
	var url_d string

	url = "amqp://" + Config.Amqp.Username + ":" + Config.Amqp.Password + "@" + Config.Amqp.Address
	url_d = "amqp://" + Config.Amqp.Username + ":********@" + Config.Amqp.Address

	ac = &AmqpClient{
		Events:  make(chan []byte, config.D_AMQP_BUFSIZE),
		Control: make(chan int, config.D_CONTROL_BUFSIZE),
		Done:    make(chan bool, config.D_DONE_BUFSIZE),
	}

	// Try to connect to AMQP
	if ac.Connection, err = amqp.Dial(url); err != nil {
		err = errors.New("[Amqp]: Connection to " + url_d + " failed: " + err.Error())
		ac = nil
		return
	}

	// Once established, setup a channel
	if ac.Channel, err = ac.Connection.Channel(); err != nil {
		err = errors.New("[Amqp]: Failed to setup a channel: " + err.Error())
		return
	}

	// Declare the fanout exchange on the newly created channel
	err = ac.Channel.ExchangeDeclare(
		Config.Amqp.Exchange, // Name of the exchange
		"fanout",             // Type of exchange
		true,                 // Durable queue
		false,                // Not auto-deleted
		false,                // Not an internal queue
		false,                // No-wait queue
		nil,                  // No arguments
	)
	if err != nil {
		err = errors.New("[Amqp]: Failed to declare an exchange: " + err.Error())
		return
	}

	// Declare the private queue
	ac.Queue, err = ac.Channel.QueueDeclare(
		Config.Amqp.Exchange, // Name
		true,                 // Durable queue
		false,                // Dont delete when unused
		true,                 // Exclusive queue
		false,                // No-wait queue
		nil,                  // No arguments
	)
	if err != nil {
		err = errors.New("[Amqp]: Failed to declare queue: " + err.Error())
	}

	return
}

func (ac *AmqpClient) Slurp() (err error) {
	var stop_loop bool

	// Bind to the queue
	err = ac.Channel.QueueBind(
		ac.Queue.Name,        // Name of queue
		"",                   // Routing key
		Config.Amqp.Exchange, // Exchange
		false,
		nil,
	)
	if err != nil {
		err = errors.New("[Amqp]: Failed to bind to queue: " + err.Error())
		return
	}

	// Start consumer
	messageChannel, err := ac.Channel.Consume(
		ac.Queue.Name, // Queue to use
		"",            // Name of the consumer
		true,          // Auto acknowledge reception of event
		false,         // Non-exclusive consumer
		false,         // No-local consumer
		false,         // No-wait consumer
		nil,           // No arguments
	)

	// Start amqp eventloop
	stop_loop = false
	for {
		// Break loop if requested
		if stop_loop {
			break
		}

		select {
		case message := <-messageChannel:
			{
				event := message.Body
				Log.Debug("[Amqp]: Received: " + string(event))
			}
		case cmd := <-ac.Control:
			{
				switch cmd {
				case config.CTL_SHUTDOWN:
					{
						Log.Debug("[Amqp]: Cleaning up and exiting loop")
						ac.Connection.Close()
						stop_loop = true
						continue
					}
				default:
					{
						Log.Warning("[Amqp]: Unknown control signal received")
						continue
					}
				}
			}
		}
	}

	ac.Done <- true

	return
}

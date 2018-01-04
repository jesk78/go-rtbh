package amqp

import (
	"fmt"

	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/streadway/amqp"
)

func (a *AmqpClient) Connect() error {
	var err error

	url := "amqp://" + cfg.Amqp.Username + ":" + cfg.Amqp.Password + "@" + cfg.Amqp.Address
	url_d := "amqp://" + cfg.Amqp.Username + ":***@" + cfg.Amqp.Address

	// Try to connect to AMQP
	a.connection, err = amqp.Dial(url)
	if err != nil {
		a = nil
		return fmt.Errorf("AmqpClient.Connect amqp.Dial: %v", err)
	}
	log.Debugf("AmqpClient: Connected to " + url_d)

	// Once established, setup a channel
	a.channel, err = a.connection.Channel()
	if err != nil {
		return fmt.Errorf("AmqpClient.Connect connection.Channel: %v", err)
	}

	// Declare the fanout exchange on the newly created channel
	err = a.channel.ExchangeDeclare(
		cfg.Amqp.Exchange, // Name of the exchange
		"fanout",          // Type of exchange
		true,              // Durable queue
		false,             // Not auto-deleted
		false,             // Not an internal queue
		false,             // No-wait queue
		nil,               // Arguments
	)
	if err != nil {
		return fmt.Errorf("AmqpClient.Connect channel.ExchangeDeclare: %v", err)
	}

	// Declare the private queue
	a.queue, err = a.channel.QueueDeclare(
		cfg.Amqp.Exchange+".rtbh-queue", // Name
		true,  // Durable queue
		false, // Dont delete when unused
		true,  // Exclusive queue
		false, // No-wait queue
		nil,   // No arguments
	)
	if err != nil {
		return fmt.Errorf("AmqpClient.Connect channel.QueueDeclare: %v", err)
	}

	// Bind to the queue
	err = a.channel.QueueBind(
		a.queue.Name,      // Name of queue
		"#",               // Routing key
		cfg.Amqp.Exchange, // Exchange
		false,             // No-wait
		nil,               // Args
	)
	if err != nil {
		return fmt.Errorf("AmqpClient.Connect channel.QueueBind: %v", err)
	}

	return nil
}

func (a *AmqpClient) Slurp(input chan []byte) (err error) {
	var stop_loop bool

	// Start consumer
	messageChannel, err := a.channel.Consume(
		a.queue.Name, // Queue to use
		"",           // Name of the consumer
		true,         // Auto acknowledge reception of event
		false,        // Non-exclusive consumer
		false,        // No-local consumer
		false,        // No-wait consumer
		nil,          // No arguments
	)

	log.Debugf("AmqpClient.Slurp: Reading from event queue")

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
				input <- event
			}
		case cmd := <-a.Control:
			{
				switch cmd {
				case config.CTL_SHUTDOWN:
					{
						log.Debugf("AmqpClient: Cleaning up and exiting loop")
						a.connection.Close()
						stop_loop = true
						continue
					}
				default:
					{
						log.Warningf("AmqpClient: Unknown control signal received")
						continue
					}
				}
			}
		}
	}

	a.Done <- true

	return
}

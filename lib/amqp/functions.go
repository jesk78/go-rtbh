package amqp

import (
	"errors"
	"github.com/r3boot/go-rtbh/config"
	"github.com/streadway/amqp"
)

func (a *AmqpClient) Connect() (err error) {
	var url string
	var url_d string

	url = "amqp://" + Config.Amqp.Username + ":" + Config.Amqp.Password + "@" + Config.Amqp.Address
	url_d = "amqp://" + Config.Amqp.Username + ":***@" + Config.Amqp.Address

	// Try to connect to AMQP
	if a.connection, err = amqp.Dial(url); err != nil {
		err = errors.New(MYNAME + ": Connection to " + url_d + " failed: " + err.Error())
		a = nil
		return
	}
	Log.Debug(MYNAME + ": Connected to " + url_d)

	// Once established, setup a channel
	if a.channel, err = a.connection.Channel(); err != nil {
		err = errors.New(MYNAME + ": Failed to setup a channel: " + err.Error())
		return
	}

	// Declare the fanout exchange on the newly created channel
	err = a.channel.ExchangeDeclare(
		Config.Amqp.Exchange, // Name of the exchange
		"fanout",             // Type of exchange
		true,                 // Durable queue
		false,                // Not auto-deleted
		false,                // Not an internal queue
		false,                // No-wait queue
		nil,                  // Arguments
	)
	if err != nil {
		err = errors.New(MYNAME + ": Failed to declare an exchange: " + err.Error())
		return
	}

	// Declare the private queue
	a.queue, err = a.channel.QueueDeclare(
		Config.Amqp.Exchange+".rtbh-queue", // Name
		true,  // Durable queue
		false, // Dont delete when unused
		true,  // Exclusive queue
		false, // No-wait queue
		nil,   // No arguments
	)
	if err != nil {
		err = errors.New(MYNAME + ": Failed to declare queue: " + err.Error())
		return
	}

	// Bind to the queue
	err = a.channel.QueueBind(
		a.queue.Name,         // Name of queue
		"#",                  // Routing key
		Config.Amqp.Exchange, // Exchange
		false,                // No-wait
		nil,                  // Args
	)
	if err != nil {
		err = errors.New(MYNAME + ": Failed to bind to queue: " + err.Error())
		return
	}

	return
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
	Log.Debug("[Amqp]: Ready to consume messages")

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
						Log.Debug("[Amqp]: Cleaning up and exiting loop")
						a.connection.Close()
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

	a.Done <- true

	return
}

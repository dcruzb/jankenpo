package rabbitMQ

import (
	"github.com/streadway/amqp"
	"github.com/dcbCIn/jankenpo/shared"
)

const NAME = "jankenpo/rabbitMQ"

type RabbitMQ struct {
	ip               string
	port             string
	serverConnection *amqp.Connection
	channel          *amqp.Channel
	messages         <-chan amqp.Delivery
}

func (rMQ *RabbitMQ) ConnectToServer(ip, port string) {
	// connect to server
	conn, err := amqp.Dial("amqp://guest:guest@" + ip + ":" + port + "/")
	shared.FailOnError(NAME, err, "Failed to connect to RabbitMQ")

	rMQ.serverConnection = conn

	ch, err := conn.Channel()
	shared.FailOnError(NAME, err, "Failed to open a channel")

	rMQ.channel = ch
}

func (rMQ *RabbitMQ) CloseConnection() {
	err := rMQ.channel.Close()
	if err != nil {
		shared.PrintlnError(NAME, err)
	}
	err = rMQ.serverConnection.Close()
	if err != nil {
		shared.PrintlnError(NAME, err)
	}
}

func (rMQ *RabbitMQ) CreateQueue(name string) {
	_, err := rMQ.channel.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	shared.FailOnError(NAME, err, "Failed to declare a queue")
}

func (rMQ *RabbitMQ) ReadChannel(queueName string) (messages <-chan amqp.Delivery) {
	messages, err := rMQ.channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	shared.FailOnError(NAME, err, "Failed to register a consumer")

	return messages
}

func (rMQ *RabbitMQ) ReadOne(queueName string) (message string) {
	if rMQ.messages == nil {
		msgs, err := rMQ.channel.Consume(
			queueName, // queue
			"",        // consumer
			true,      // auto-ack
			false,     // exclusive
			false,     // no-local
			false,     // no-wait
			nil,       // args
		)
		shared.FailOnError(NAME, err, "Failed to register a consumer")

		rMQ.messages = msgs
	}

	d := <-rMQ.messages
	message = string(d.Body)
	return message
}

func (rMQ *RabbitMQ) Write(queueName, message string) {
	err := rMQ.channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	shared.FailOnError(NAME, err, "Failed to publish a message")
}

package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// failOnError is a helper function to handle errors
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare a fanout exchange
	err = ch.ExchangeDeclare(
		"logs",   // name of the exchange
		"fanout", // type of exchange (fanout broadcasts to all queues)
		true,     // durable (the exchange will survive a broker restart)
		false,    // auto-deleted (the exchange won't be deleted when last queue is unbound)
		false,    // internal (the exchange isn't used directly by publishers)
		false,    // no-wait (the declaration will block until a confirmation is received)
		nil,      // arguments (optional, none in this case)
	)
	failOnError(err, "Failed to declare an exchange")

	// Declare a queue with a random name
	q, err := ch.QueueDeclare(
		"",    // name (empty string means a random unique name will be generated)
		false, // durable (the queue will not survive a broker restart)
		false, // delete when unused
		true,  // exclusive (the queue will be deleted when the connection closes)
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Bind the queue to the exchange
	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key (ignored for fanout exchanges)
		"logs", // exchange name
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to bind a queue")

	// Start consuming messages from the queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer (empty string means a random consumer tag will be generated)
		true,   // auto-ack (automatically acknowledge messages)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	// Create a channel to keep the main goroutine running
	var forever chan struct{}

	// Start a goroutine to process incoming messages
	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever // This line blocks indefinitely
}

package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

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

	// Create a context with a 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get the message body from command-line arguments
	body := bodyFrom(os.Args)

	// Publish the message to the exchange
	err = ch.PublishWithContext(ctx,
		"logs", // exchange name
		"",     // routing key (ignored for fanout exchanges)
		false,  // mandatory (don't return an error if no queue is bound)
		false,  // immediate (don't return an error if no consumer on the matched queue is ready to accept the message)
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}

// bodyFrom constructs the message body from command-line arguments
func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello" // default message if no arguments provided
	} else {
		s = strings.Join(args[1:], " ") // join all arguments into a single string
	}
	return s
}

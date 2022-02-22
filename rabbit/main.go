package rabbit

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

func alertOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s\n", msg, err)
	}
}

func connectToRabbit(queueName string) (*amqp.Connection, *amqp.Channel, amqp.Queue) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://service:%s@rabbitmq:5672/", os.Getenv("ADMIN_PASS")))
	alertOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	alertOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	alertOnError(err, "Failed to declare a queue")
	return conn, ch, q
}

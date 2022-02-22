package rabbit

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"math/rand"
	"os"
)

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func SendMessage(queueName string, message string) error {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://service:%s@rabbitmq:5672/", os.Getenv("ADMIN_PASS")))
	alertOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	alertOnError(err, "Failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	alertOnError(err, "Failed to declare a queue")

	return ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}



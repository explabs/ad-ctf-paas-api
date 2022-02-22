package rabbit

import amqp "github.com/rabbitmq/amqp091-go"

func SendMessage(queueName string, message string) error {
	conn, ch, q := connectToRabbit(queueName)
	defer ch.Close()
	defer conn.Close()
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

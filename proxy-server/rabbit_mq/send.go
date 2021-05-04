package rabbit_mq

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

func SendRequest(name string, requestId string) error {
	conn, err := amqp.Dial(os.Getenv("rabbitmq"))
	defer conn.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	ch, err := conn.Channel()
	defer ch.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	q, err := ch.QueueDeclare(
		name, false, false, false, false, nil,
	)
	if err != nil {
		log.Println(err)
		return err
	}

	err = ch.Publish(
		"", q.Name, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(requestId),
		},
	)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

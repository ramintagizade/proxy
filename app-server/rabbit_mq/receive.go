package rabbit_mq

import (
	"app-server/requests"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func Receive() {

	conn, err := amqp.Dial(os.Getenv("rabbitmq"))
	defer conn.Close()
	if err != nil {
		log.Println(err)
		return
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
		return
	}
	defer ch.Close()
	qr, err := ch.QueueDeclare(
		"request", false, false, false, false, nil,
	)
	if err != nil {
		log.Println(err)
		return
	}
	msgsRequests, err := ch.Consume(
		qr.Name, "", true, false, false, false, nil,
	)
	if err != nil {
		log.Println(err)
		return
	}

	forever := make(chan bool)
	go func() {
		for d := range msgsRequests {
			go func() {
				var id string = string(d.Body)
				log.Println("received msg requests ", id)
				requests.SendRequestProxy(id)
			}()
		}
	}()

	<-forever
}

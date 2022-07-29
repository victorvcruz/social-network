package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

func ConnectToMessageBroker() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return conn, err
	}

	log.Println("RabbitMQ Connected")
	return conn, nil
}

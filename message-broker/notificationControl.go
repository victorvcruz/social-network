package message_broker

import (
	"fmt"
	"github.com/streadway/amqp"
)

type RabbitControl struct {
	Conn                   *amqp.Connection
	NotificationController *NotificationController
}

func (r *RabbitControl) SendMessage(message string) {

	ch, err := r.Conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"NotificationQueue",
		false,
		false,
		false,
		false,
		nil,
	)

	fmt.Println(q)

	if err != nil {
		fmt.Println(err)
	}

	err = ch.Publish(
		"",
		"NotificationQueue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Published Message to Queue")

}

func (r *RabbitControl) ConsumerMessage() {
	ch, err := r.Conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer ch.Close()

	if err != nil {
		fmt.Println(err)
	}

	msgs, err := ch.Consume(
		"NotificationQueue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			r.NotificationController.HandlerNotification(string(d.Body))
			fmt.Printf("Recieved Message: %s\n", d.Body)
		}
	}()

	fmt.Println("Successfully Connected to our RabbitMQ Instance")
	fmt.Println(" [*] - Waiting for messages")
	<-forever
}
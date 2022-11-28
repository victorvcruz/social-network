package service

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"social_network_project/internal/notification"
)

type NotificationServiceClient interface {
	SendMessage(message string)
	ConsumerMessage()
}

type NotificationService struct {
	Conn       *amqp.Connection
	Repository notification.NotificationRepositoryClient
}

func NewNotificationService(_conn *amqp.Connection, _repository notification.NotificationRepositoryClient) NotificationServiceClient {
	return &NotificationService{
		Conn: _conn,
		Repository: _repository,
	}
}

func (r *NotificationService) SendMessage(message string) {

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

func (r *NotificationService) ConsumerMessage() {
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
			n := &notification.Notification{}
			json.Unmarshal(d.Body, n)
			r.Repository.HandlerNotification(n)
			fmt.Printf("Recieved Message: %s\n", d.Body)
		}
	}()

	fmt.Println("Successfully Connected to our RabbitMQ Instance")
	fmt.Println(" [*] - Waiting for messages")
	<-forever
}

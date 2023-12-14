package tg_bot

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func SendToQueue(ch *amqp.Channel, queueName string, msg string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := ch.PublishWithContext(
		ctx,
		queueName,
		"msg",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

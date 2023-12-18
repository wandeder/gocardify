package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/wandeder/gocardify/anki_bot"
	"os"
	"path/filepath"
	"time"
)

func main() {
	cwd, _ := os.Getwd()
	logFile := filepath.Join(cwd, os.Getenv("LOG_FILE"))
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logger.Println("Failed opening or creating logs file:", err)
	}
	defer file.Close()
	logger.SetOutput(file)

	rabbitUrl := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASS"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)
	rabbitMq, err := amqp.DialConfig(rabbitUrl, amqp.Config{
		Heartbeat: 10 * time.Second,
	})
	if err != nil {
		logger.Println("Connection RabbitMQ error:", err)
	}
	defer rabbitMq.Close()

	ch, err := rabbitMq.Channel()
	if err != nil {
		logger.Println("Channel RabbitMQ error:", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		os.Getenv("QUEUE_NAME"),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Println("Error in queue declare:", err)
	}
	msgs, err := ch.Consume(
		os.Getenv("QUEUE_NAME"),
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Fatalf("Error getting message from queue: %v", err)
	}

	for msg := range msgs {
		logger.Printf("Get a message from RabbitMQ: %s", msg.Body)
		if err := anki_bot.CreateCard(string(msg.Body)); err != nil {
			logger.Println("Error in creating card:", err)
		}
	}
}

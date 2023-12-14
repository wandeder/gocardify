package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/wandeder/gocardify/src/tg_bot"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	cwd, _ := os.Getwd()
	logFile := filepath.Join(cwd, os.Getenv("LOG_FILE"))
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logger.Println(err, "Failed opening or creating logs file")
	}
	defer file.Close()
	logger.SetOutput(file)

	tgBot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		logger.Println(err)
	}

	chat := tgBot.GetUpdatesChan(tgbotapi.UpdateConfig{
		Offset:  0,
		Timeout: 60,
	})
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	rabbitUrl := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASS"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)
	rabbitMq, err := amqp.Dial(rabbitUrl)
	if err != nil {
		logger.Println(err)
	}
	defer rabbitMq.Close()

	ch, err := rabbitMq.Channel()
	if err != nil {
		logger.Println(err)
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		os.Getenv("QUEUE_NAME"),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Println(err)
	}

	for {
		select {
		case update := <-chat:
			if update.Message != nil {
				message := update.Message.Text

				msg, err := tg_bot.ReadMsg(message)
				if err != nil {
					logger.Println("Failed to read message from Telegram:", err)
					errMsg := tgbotapi.NewMessage(
						update.SentFrom().ID,
						"Incorrect format. Try this one:\nfront: some front text.\nback: some back text.",
					)
					_, err = tgBot.Send(errMsg)
					if err != nil {
						logger.Println("Failed to send message in Telegram:", err)
					}
				}

				if err := tg_bot.SendToQueue(ch, queue.Name, msg); err != nil {
					logger.Println("Failed to send message to RabbitMQ:", err)
				}
			}
		case <-stop:
			logger.Println("Stop tg_bot.")
			return
		}
	}

}

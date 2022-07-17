package main

import (
	"github.com/joho/godotenv"
	rabbitmq "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"strings"
)

var RabbitChannel *rabbitmq.Channel
var RabbitConnection *rabbitmq.Connection
var LsQueue string
var CoreQueue string

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	_ = InitRabbitMQ()
	defer service.RabbitConnection.Close()
	defer service.RabbitChannel.Close()
}

func InitRabbitMQ() error {
	username := os.Getenv("RUSERNAME")
	password := os.Getenv("RPASSWORD")
	host := os.Getenv("RHOST")
	LsQueue = os.Getenv("LSQueue")
	CoreQueue = os.Getenv("COREQueue")
	RabbitConnection, err := rabbitmq.Dial("amqp://" + username + ":" + password + "@" + host + "/")
	if err != nil {
		return err
	}

	RabbitChannel, err = RabbitConnection.Channel()
	if err != nil {
		return err
	}

	_, err = RabbitChannel.QueueDeclare(
		LsQueue,
		false,
		false,
		false,
		false,
		nil,
	)

	_, err = RabbitChannel.QueueDeclare(
		CoreQueue,
		false,
		false,
		false,
		false,
		nil,
	)

	err = RabbitChannel.Confirm(false)

	if err != nil {
		return err
	}
	log.Println("Connect to rabbitMQ successfully")
	return nil
}

func RPublishToCore(body string) error {
	body = strings.ReplaceAll(body, "\"", "\\\"")
	err := RabbitChannel.Publish(
		"",
		CoreQueue,
		false,
		false,
		rabbitmq.Publishing{
			DeliveryMode: rabbitmq.Persistent,
			ContentType:  "application/json",
			Body:         []byte("\"" + body + "\""),
		})
	if err != nil {
		log.Println(err)
	}
	return err
}

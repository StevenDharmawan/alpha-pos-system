package config

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
)

func ConnectRabbitmqs() (*amqp.Channel, func() error) {
	connection, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}
	channel, err := connection.Channel()
	if err != nil {
		fmt.Println(err)
	}
	return channel, connection.Close
}

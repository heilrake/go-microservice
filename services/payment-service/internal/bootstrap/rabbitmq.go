package bootstrap

import (
	"log"

	"ride-sharing/shared/messaging"
)

func InitRabbitMQ(uri string) *messaging.RabbitMQ {
	rabbitmq, err := messaging.NewRabbitMQ(uri)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	log.Println("RabbitMQ connected")
	return rabbitmq
}

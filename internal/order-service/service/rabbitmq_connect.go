package service

import (
	"fmt"

	"github.com/joel-malina/tucows-challenge/internal/order-service/config"
	ampq "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

func connectToRabbitMQ(log *logrus.Logger, serviceConfig config.ServiceConfig) *ampq.Channel {
	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		serviceConfig.RabbitMQUser,
		serviceConfig.RabbitMQPassword,
		serviceConfig.PostgresHost,
		serviceConfig.RabbitMQPort,
	)

	conn, err := ampq.Dial(connectionString)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to RabbitMQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	return ch
}

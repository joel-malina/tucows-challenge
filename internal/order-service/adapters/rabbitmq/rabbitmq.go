package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/joel-malina/tucows-challenge/internal/order-service/model"
	ampq "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type OrderQueue struct {
	mq *ampq.Channel
}

func NewOrderQueue(mq *ampq.Channel) OrderQueue {
	return OrderQueue{
		mq: mq,
	}
}

func (o OrderQueue) OrderEnqueue(ctx context.Context, order model.Order) error {
	body, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("Failed to encode order as JSON: %v", err)
	}

	err = o.mq.Publish(
		"",                // exchange
		"payment_request", // routing key
		false,             // mandatory
		false,             // immediate
		ampq.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	log.Printf(" [>>>] payment requested %s", body)

	return nil
}

func (o OrderQueue) OrderPaymentListener(log *logrus.Logger) {
	msgs, err := o.mq.Consume(
		"payment_response", // queue
		"",                 // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var order model.Order
			if err := json.Unmarshal(d.Body, &order); err != nil {
				log.Printf("Error decoding JSON: %s", err)
				continue
			}
			log.Printf(" [<<<] payment processed for order: %+v", order)
			// TODO: update db with new status
		}
	}()

	log.Printf("Waiting for messages on payment response queue...")
	<-forever
}

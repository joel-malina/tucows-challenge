package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/joel-malina/tucows-challenge/internal/order-service/model"
	ampq "github.com/rabbitmq/amqp091-go"
)

type OrderQueue struct {
	mq *ampq.Channel
}

func NewOrderQueue(mq *ampq.Channel) OrderQueue {
	return OrderQueue{
		mq: mq,
	}
}

// TODO: implement enqueue on payment_request and dequeue on payment_response

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

	log.Printf(" [x] Sent %s", body)

	return nil
}

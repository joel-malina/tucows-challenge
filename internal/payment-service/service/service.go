package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/joel-malina/tucows-challenge/internal/order-service/model"
	"github.com/joel-malina/tucows-challenge/internal/payment-service/config"
	ampq "github.com/rabbitmq/amqp091-go"
)

func Run(ctx context.Context, serviceConfig config.ServiceConfig) {
	// intentional contrast to the order-service -- this is what a non-hexagonal minimal service could look like
	consumeOrders()
}

func consumeOrders() {
	conn, err := ampq.Dial("amqp://user:password@localhost:5672/")
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close() //nolint:errcheck

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	defer ch.Close() //nolint:errcheck

	msgs, err := ch.Consume(
		"payment_request", // queue
		"",                // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		log.Fatalf("failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var order model.Order
			if err := json.Unmarshal(d.Body, &order); err != nil {
				log.Printf("Error decoding JSON: %s", err)
				continue
			}
			log.Printf(" [<<<] Received an order: %+v", order)
			processOrderPayment(ch, order)
		}
	}()

	log.Printf("Waiting for messages. To exit press CTRL+C")
	<-forever
}

func processOrderPayment(ch *ampq.Channel, order model.Order) {

	err := checkTotal(order)
	if err != nil {
		log.Printf("Error checking total: %v", err)
		order.Status = model.OrderStatusPaymentFailure
	}

	order.Status = model.OrderStatusPaymentFailure
	if order.TotalPrice <= 1000 {
		order.Status = model.OrderStatusPaymentSuccess
	}

	body, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("failed to encode order as JSON: %v", err)
	}

	err = ch.Publish(
		"",                 // exchange
		"payment_response", // routing key
		false,              // mandatory
		false,              // immediate
		ampq.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Fatalf("failed to publish a message: %v", err)
	}

	log.Printf(" [>>>] Sent %s", body)
}

func checkTotal(order model.Order) error {

	claimedTotal := order.TotalPrice
	var actualTotal float64

	for item := range order.OrderItems {
		actualTotal += order.OrderItems[item].Price
	}

	if claimedTotal != actualTotal {
		return fmt.Errorf("total amount is %f, expected %f", claimedTotal, actualTotal)
	}

	return nil
}

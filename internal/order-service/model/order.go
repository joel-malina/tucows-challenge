package model

import (
	"time"

	uuidgen "github.com/gofrs/uuid/v5"
	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusCreated                 OrderStatus = "created"
	OrderStatusPaymentProcessRequested OrderStatus = "payment_process_requested"
	OrderStatusPaymentFailure          OrderStatus = "payment_process_failure"
	OrderStatusPaymentSuccess          OrderStatus = "payment_process_success"
	OrderStatusCompleted               OrderStatus = "completed"
	OrderStatusCancelled               OrderStatus = "cancelled"
)

type Order struct {
	ID          uuid.UUID
	CustomerID  uuid.UUID
	ProductID   uuid.UUID
	ProductName string
	Quantity    int
	Price       float64
	Status      OrderStatus
	CreatedAt   time.Time
	LastUpdate  time.Time
}

func CreateOrderID() uuid.UUID {
	uuidv7, err := uuidgen.NewV7()
	result := uuid.UUID(uuidv7)
	if err != nil {
		result = uuid.New()
	}

	return result
}

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
	ID         uuid.UUID
	CustomerID uuid.UUID
	OrderDate  time.Time
	Status     OrderStatus
	TotalPrice float64
	OrderItems []OrderItem
}

type OrderItem struct {
	ID        uuid.UUID
	OrderID   uuid.UUID
	ProductID uuid.UUID
	Quantity  int
	Price     float64
}

type Product struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       float64
	Stock       int
}

// There should also be a customer table to support multiple customers

func CreateUUID() uuid.UUID {
	uuidv7, err := uuidgen.NewV7()
	result := uuid.UUID(uuidv7)
	if err != nil {
		result = uuid.New()
	}

	return result
}

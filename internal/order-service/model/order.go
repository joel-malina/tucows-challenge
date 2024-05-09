package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending                 OrderStatus = "pending"
	OrderStatusCheckoutStarted         OrderStatus = "checkout_started"
	OrderStatusCheckoutCompleted       OrderStatus = "checkout_completed"
	OrderStatusPaymentProcessRequested OrderStatus = "payment_process_requested"
	OrderStatusPaymentProcessCompleted OrderStatus = "payment_process_completed"
	OrderStatusCancelled               OrderStatus = "cancelled"
	OrderStatusRefundStarted           OrderStatus = "refund_started"
	OrderStatusRefundCompleted         OrderStatus = "refund_completed"
)

type Order struct {
	ID                 uuid.UUID
	ProductName        string
	Quantity           int
	Price              float64
	ProductDescription string
	Status             OrderStatus
	CreatedAt          time.Time
	LastUpdate         time.Time
}

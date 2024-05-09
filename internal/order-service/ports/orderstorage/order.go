package orderstorage

import (
	"context"

	"github.com/google/uuid"
	"github.com/joel-malina/tucows-challenge/internal/order-service/model"
)

type OrderRepository interface {
	OrderCreator
	OrderGetter
	OrderUpdater
	OrderDeleter
}

type OrderCreator interface {
	OrderCreate(ctx context.Context, order model.Order) error
}

type OrderGetter interface {
	OrderGet(ctx context.Context, id uuid.UUID) (model.Order, error)
	OrderGetAll(ctx context.Context) ([]model.Order, error)
}

type OrderUpdater interface {
	OrderUpdate(ctx context.Context, order model.Order) error
}

type OrderDeleter interface {
	OrderDelete(ctx context.Context, id uuid.UUID) error
}

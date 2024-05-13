package orderqueue

import (
	"context"

	"github.com/joel-malina/tucows-challenge/internal/order-service/model"
)

type OrderQueue interface {
	OrderEnqueuer
}

type OrderEnqueuer interface {
	OrderEnqueue(ctx context.Context, order model.Order) error
}

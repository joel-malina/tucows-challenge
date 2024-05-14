package orderqueue

import (
	"context"

	"github.com/joel-malina/tucows-challenge/internal/order-service/model"
	"github.com/sirupsen/logrus"
)

type OrderQueue interface {
	OrderEnqueuer
	OrderListener
}

type OrderEnqueuer interface {
	OrderEnqueue(ctx context.Context, order model.Order) error
}

type OrderListener interface {
	OrderPaymentListener(log *logrus.Logger)
}

package service

import (
	"sync"

	"github.com/joel-malina/tucows-challenge/internal/order-service/adapters/rabbitmq"
	"github.com/joel-malina/tucows-challenge/internal/order-service/config"
	"github.com/joel-malina/tucows-challenge/internal/order-service/ports/orderqueue"
	"github.com/sirupsen/logrus"
)

type QueueResolver struct {
	OrderQueue orderqueue.OrderQueue
	once       sync.Once
}

func (q *QueueResolver) Resolve(log *logrus.Logger, serviceConfig config.ServiceConfig) {
	q.once.Do(func() {
		mq := connectToRabbitMQ(log, serviceConfig)
		q.OrderQueue = rabbitmq.NewOrderQueue(mq)
	})
}

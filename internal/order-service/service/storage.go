package service

import (
	"sync"

	"github.com/joel-malina/tucows-challenge/internal/order-service/adapters/inmemorders"
	"github.com/joel-malina/tucows-challenge/internal/order-service/adapters/sqlstorage"
	"github.com/joel-malina/tucows-challenge/internal/order-service/config"
	"github.com/joel-malina/tucows-challenge/internal/order-service/ports/orderstorage"
	"github.com/sirupsen/logrus"
)

type StorageResolver struct {
	OrderRepository orderstorage.OrderRepository
	once            sync.Once
}

func (s *StorageResolver) Resolve(log *logrus.Logger, serviceConfig config.ServiceConfig) {
	s.once.Do(func() {
		if serviceConfig.EnablePersistentStorage {
			db := connectToPostgres(log, serviceConfig)
			s.OrderRepository = sqlstorage.NewOrderRepository(db)
		} else {
			s.OrderRepository = inmemorders.New()
		}
	})
}

package inmemorders

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/joel-malina/tucows-challenge/internal/order-service/model"
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]model.Order
}

func New() *OrderStorage {
	return &OrderStorage{
		orders: map[uuid.UUID]model.Order{},
	}
}

func (db *OrderStorage) checkIDUniqueness(order model.Order) error {
	for _, o := range db.orders {
		if o.ID == order.ID {
			return fmt.Errorf("order-service with ID: %s, already exists", o.ID)
		}
	}

	return nil
}

func (db *OrderStorage) OrderCreate(ctx context.Context, order model.Order) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.orderCreateNoLock(ctx, order)
}

func (db *OrderStorage) orderCreateNoLock(ctx context.Context, order model.Order) error {
	// ctx would be used here if we wanted to do tracing
	// E.g.
	// _, span := o11y.NewAutoNamedChildSpan(ctx)
	//	defer span.End()
	// tracing would be done similarly in other functions here, I'll omit this in them for brevity

	err := db.checkIDUniqueness(order)
	if err != nil {
		return err
	}

	// continued example of tracing
	// span.SetAttributes(attribute.String(model.SpanAttributeOrderID, order-service.ID.String()))
	db.orders[order.ID] = order

	return nil
}

func (db *OrderStorage) OrderGet(_ context.Context, id uuid.UUID) (model.Order, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	order, ok := db.orders[id]
	if !ok {
		return model.Order{}, model.ErrOrderNotFound
	}

	return order, nil
}

//func (db *OrderStorage) OrdersGet(_ context.Context) ([]model.Order, error) {
//	db.mu.RLock()
//	defer db.mu.RUnlock()
//
//	result := make([]model.Order, 0, len(db.orders))
//	for _, order := range db.orders {
//		result = append(result, order)
//	}
//
//	return result, nil
//}

func (db *OrderStorage) OrderUpdate(_ context.Context, order model.Order) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	retrievedOrder, ok := db.orders[order.ID]
	if !ok {
		return model.ErrOrderNotFound
	}

	db.orders[retrievedOrder.ID] = order

	return nil
}

func (db *OrderStorage) OrderDelete(_ context.Context, id uuid.UUID) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	_, ok := db.orders[id]
	if !ok {
		return model.ErrOrderNotFound
	}
	delete(db.orders, id)

	return nil
}

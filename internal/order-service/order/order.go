package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/joel-malina/tucows-challenge/internal/order-service/model"
	"github.com/joel-malina/tucows-challenge/internal/order-service/ports/orderstorage"
)

type CreateRepository interface {
	orderstorage.OrderCreator
}

type OrderCreate struct {
	repo CreateRepository
	// onCreate event.Emitter[model.Order] add to queue
}

func NewOrderCreate(repo CreateRepository) *OrderCreate {
	return &OrderCreate{
		repo: repo,
	}
}

//func (fc *OrderCreate) OrderCreateSubscribe(subscription chan<- event.DataWithContext[model.Order]) {
//	fc.onCreate.Subscribe(subscription)
//}

func (o *OrderCreate) OrderCreate(ctx context.Context, order model.Order) (uuid.UUID, error) {

	// Optionally could have made the UUID here instead of assuming we'd get one from the webapp
	// order.ID = model.CreateUUID()

	err := o.repo.OrderCreate(ctx, order)
	if err != nil {
		return uuid.UUID{}, err
	}

	// TODO: put the created order into the payment process queue
	// o.onCreate.Emit(ctx, order)

	return order.ID, nil
}

type UpdateRepository interface {
	orderstorage.OrderGetter
	orderstorage.OrderUpdater
}

type OrderUpdate struct {
	repo UpdateRepository
	// onUpdate event.Emitter[model.Order]
}

func NewOrderUpdate(repo UpdateRepository) *OrderUpdate {
	return &OrderUpdate{
		repo: repo,
	}
}

//func (o *OrderUpdate) OrderUpdateSubscribe(subscription chan<- event.DataWithContext[model.Order]) {
//	o.onUpdate.Subscribe(subscription)
//}

func (o *OrderUpdate) OrderUpdate(ctx context.Context, updateParameters model.Order) error {
	_, err := o.repo.OrderGet(ctx, updateParameters.ID)
	if err != nil {
		return err
	}

	err = o.repo.OrderUpdate(ctx, updateParameters)
	if err != nil {
		return err
	}

	//o.onUpdate.Emit(ctx, updateParameters)
	return nil
}

type OrderDeleteRepository interface {
	orderstorage.OrderGetter
	orderstorage.OrderDeleter
}

type OrderDelete struct {
	repo OrderDeleteRepository
	//onDelete event.Emitter[model.OrderDeleteEvent]
}

func NewOrderDelete(repo OrderDeleteRepository) *OrderDelete {
	return &OrderDelete{
		repo: repo,
	}
}

//func (o *OrderDelete) OrderDeleteSubscribe(subscription chan<- event.DataWithContext[model.OrderDeleteEvent]) {
//	f.onDelete.Subscribe(subscription)
//}

func (o *OrderDelete) OrderDelete(ctx context.Context, id uuid.UUID) error {
	_, err := o.repo.OrderGet(ctx, id)
	if err != nil {
		return err
	}

	err = o.repo.OrderDelete(ctx, id)
	if err == nil {
		//deleteEvent := model.OrderDeleteEvent{OrderID: id}
		//o.onDelete.Emit(ctx, deleteEvent)
		return nil
	}

	return err
}

type OrderGet struct {
	repo orderstorage.OrderGetter
}

func NewOrderGet(repo orderstorage.OrderGetter) *OrderGet {
	return &OrderGet{
		repo: repo,
	}
}

func (o *OrderGet) OrderGet(ctx context.Context, id uuid.UUID) (model.Order, error) {
	order, err := o.repo.OrderGet(ctx, id)
	if err != nil {
		return model.Order{}, model.ErrOrderNotFound
	}

	return order, err
}

//type OrdersGet struct {
//	repo orderstorage.OrderGetter
//}

//func NewOrdersGet(repo orderstorage.OrderGetter) *OrdersGet {
//	return &OrdersGet{
//		repo: repo,
//	}
//}

//func (o OrdersGet) OrdersGet(ctx context.Context) ([]model.Order, error) {
//	orders, err := o.repo.OrdersGet(ctx)
//	if err != nil {
//		return []model.Order{}, err
//	}
//
//	return orders, nil
//}

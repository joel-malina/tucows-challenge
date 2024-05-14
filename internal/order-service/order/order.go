package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/joel-malina/tucows-challenge/internal/order-service/model"
	"github.com/joel-malina/tucows-challenge/internal/order-service/ports/orderqueue"
	"github.com/joel-malina/tucows-challenge/internal/order-service/ports/orderstorage"
)

type OrderCreateAdapter interface {
	orderstorage.OrderCreator
	orderqueue.OrderEnqueuer
}

type OrderCreate struct {
	adpt OrderCreateAdapter
}

func NewOrderCreate(adpt OrderCreateAdapter) *OrderCreate {
	return &OrderCreate{
		adpt: adpt,
	}
}

func (o *OrderCreate) OrderCreate(ctx context.Context, order model.Order) (uuid.UUID, error) {

	// Optionally could have made the UUID here instead of assuming we'd get one from the webapp
	// order.ID = model.CreateUUID()

	err := o.adpt.OrderCreate(ctx, order)
	if err != nil {
		return uuid.UUID{}, err
	}

	err = o.adpt.OrderEnqueue(ctx, order)
	if err != nil {
		return uuid.UUID{}, err
	}

	return order.ID, nil
}

type UpdateRepository interface {
	orderstorage.OrderGetter
	orderstorage.OrderUpdater
}

type OrderUpdate struct {
	repo UpdateRepository
}

func NewOrderUpdate(repo UpdateRepository) *OrderUpdate {
	return &OrderUpdate{
		repo: repo,
	}
}

func (o *OrderUpdate) OrderUpdate(ctx context.Context, updateParameters model.Order) error {
	_, err := o.repo.OrderGet(ctx, updateParameters.ID)
	if err != nil {
		return err
	}

	err = o.repo.OrderUpdate(ctx, updateParameters)
	if err != nil {
		return err
	}

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

func (o *OrderDelete) OrderDelete(ctx context.Context, id uuid.UUID) error {
	_, err := o.repo.OrderGet(ctx, id)
	if err != nil {
		return err
	}

	err = o.repo.OrderDelete(ctx, id)
	if err == nil {
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

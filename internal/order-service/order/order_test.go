package order

import (
	"context"
	"testing"
	"time"

	"github.com/joel-malina/tucows-challenge/internal/order-service/adapters/inmemorders"
	"github.com/joel-malina/tucows-challenge/internal/order-service/model"
	"github.com/joel-malina/tucows-challenge/internal/order-service/ports/orderqueue"
	"github.com/joel-malina/tucows-challenge/internal/order-service/ports/orderstorage"
	. "github.com/onsi/gomega"
)

func ParallelGomegaWithT(t *testing.T) *WithT {
	t.Parallel()
	g := NewWithT(t)
	return g
}

func getTestOrder() model.Order {
	orderID := model.CreateUUID()

	return model.Order{
		ID:         orderID,
		CustomerID: model.CreateUUID(),
		OrderDate:  time.Now().UTC(),
		Status:     model.OrderStatusCreated,
		TotalPrice: 2222,
		OrderItems: []model.OrderItem{
			{
				ID:        model.CreateUUID(),
				OrderID:   orderID,
				ProductID: model.CreateUUID(),
				Quantity:  2,
				Price:     1111,
			},
		},
	}
}

type orderEnqueueStub struct {
	err error
}

func (f orderEnqueueStub) OrderEnqueue(ctx context.Context, order model.Order) error {
	return f.err
}

func TestOrders_OrderCreate(t *testing.T) {
	g := ParallelGomegaWithT(t)

	testOrder := getTestOrder()
	queueStub := orderEnqueueStub{}
	inmemdb := inmemorders.New()

	subject := NewOrderCreate(struct {
		orderstorage.OrderCreator
		orderqueue.OrderEnqueuer
	}{
		OrderCreator:  inmemdb,
		OrderEnqueuer: queueStub,
	})

	orderCreateID, err := subject.OrderCreate(context.Background(), testOrder)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(queueStub.err).ToNot(HaveOccurred())
	g.Expect(orderCreateID).To(Equal(testOrder.ID))
}

func TestOrders_OrderGet(t *testing.T) {
	g := ParallelGomegaWithT(t)
	testOrder := getTestOrder()

	repo := inmemorders.New()
	_ = repo.OrderCreate(context.Background(), testOrder) //nolint:errcheck
	subject := NewOrderGet(repo)

	_, err := subject.OrderGet(context.Background(), testOrder.ID)
	g.Expect(err).ToNot(HaveOccurred())
}

// TODO: continue testing for orders

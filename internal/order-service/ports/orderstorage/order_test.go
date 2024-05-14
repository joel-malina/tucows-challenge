package orderstorage_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joel-malina/tucows-challenge/internal/common/postgrestest"
	"github.com/joel-malina/tucows-challenge/internal/order-service/adapters/inmemorders"
	"github.com/joel-malina/tucows-challenge/internal/order-service/adapters/sqlstorage"
	"github.com/joel-malina/tucows-challenge/internal/order-service/model"
	"github.com/joel-malina/tucows-challenge/internal/order-service/ports/orderstorage"
	"github.com/joel-malina/tucows-challenge/internal/order-service/ports/porttester"
	_ "github.com/lib/pq"
	. "github.com/onsi/gomega"
)

// test data product uuids
var pinkCowID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

// var blueCowID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

func ParallelGomegaWithT(t *testing.T) *WithT {
	t.Parallel()
	g := NewWithT(t)
	return g
}

func makeTestOrder() model.Order {
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
				ProductID: pinkCowID,
				Quantity:  2,
				Price:     1111,
			},
		},
	}
}

func orderRepositories(t *testing.T) []porttester.PortTester[orderstorage.OrderRepository] {
	inMemoryTester := porttester.PortTester[orderstorage.OrderRepository]{
		Subject:   inmemorders.New(),
		PostCheck: porttester.Noop,
	}
	repositories := []porttester.PortTester[orderstorage.OrderRepository]{inMemoryTester}

	if os.Getenv("SKIP_POSTGRES_TESTS") != "true" {
		db := postgrestest.CreatePostgresSandbox(t)
		dbTester := porttester.PortTester[orderstorage.OrderRepository]{
			Subject:   sqlstorage.NewOrderRepository(db),
			PostCheck: porttester.VerifyNoInUseConnections(db),
		}
		repositories = append(repositories, dbTester)
	}

	for _, repository := range repositories {
		deletedOrder := makeTestOrder()
		g := NewWithT(t)
		g.Expect(repository.Subject.OrderCreate(context.Background(), deletedOrder)).To(Succeed())
		g.Expect(repository.Subject.OrderDelete(context.Background(), deletedOrder.ID)).To(Succeed())
	}

	return repositories
}

func TestOrderCreate_specify_ID(t *testing.T) {
	t.Parallel()

	for _, repositoryTest := range orderRepositories(t) {
		test := repositoryTest
		t.Run(fmt.Sprintf("%T", test.Subject), func(t *testing.T) {
			g := ParallelGomegaWithT(t)

			err := test.Subject.OrderCreate(context.Background(), makeTestOrder())
			g.Expect(err).ToNot(HaveOccurred())
			test.PostCheck(t)
		})
	}
}

// TODO: wrap up these tests

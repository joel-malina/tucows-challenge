package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/joel-malina/tucows-challenge/api"
	ts "github.com/joel-malina/tucows-challenge/internal/common/apitestsetup"
	"github.com/joel-malina/tucows-challenge/internal/order-service/model"
	"github.com/joel-malina/tucows-challenge/internal/order-service/service"
	. "github.com/onsi/gomega"
	"github.com/parnurzeal/gorequest"
)

type orderUpdateStub struct {
	order model.Order
	err   error
}

func (f orderUpdateStub) OrderUpdate(_ context.Context, _ model.Order) error {
	return f.err
}

func TestOrderUpdate_Incorrect_Payload(t *testing.T) {
	tc := ts.SetupRoutesTest(t)

	notFoundStub := orderUpdateStub{}
	tc.Container.Add(service.MakeOrderUpdateRoute(ts.NewWebservice(), &notFoundStub))

	resp, _, err := tc.CallTo(gorequest.New().
		Put(fmt.Sprintf("/orders/%s", model.CreateUUID().String())).
		Send(struct {
			Size string
		}{
			Size: "1",
		}).
		MakeRequest()).
		Execute()

	tc.Expect(err).ToNot(HaveOccurred())
	tc.Expect(resp.Code).To(Equal(http.StatusInternalServerError))
}

func TestUpdateOrder_invalid_ID(t *testing.T) {
	tc := ts.SetupRoutesTest(t)

	updateStub := orderUpdateStub{
		order: model.Order{},
		err:   nil,
	}

	tc.Container.Add(service.MakeOrderUpdateRoute(ts.NewWebservice(), updateStub))

	resp, _, err := tc.CallTo(gorequest.New().
		Put("/orders/nonuuid").
		MakeRequest()).
		Execute()

	tc.Expect(err).ToNot(HaveOccurred())
	tc.Expect(resp.Code).To(Equal(http.StatusBadRequest))
}

func TestUpdateOrder_not_found(t *testing.T) {
	tc := ts.SetupRoutesTest(t)

	updateStub := orderUpdateStub{
		err: model.ErrOrderNotFound,
	}

	tc.Container.Add(service.MakeOrderUpdateRoute(ts.NewWebservice(), updateStub))

	orderID := model.CreateUUID()
	resp, _, err := tc.CallTo(gorequest.New().
		Put(fmt.Sprintf("/orders/%s", orderID)).
		Send(api.OrderParameters{
			OrderID:    orderID.String(),
			CustomerID: model.CreateUUID().String(),
			OrderItems: []api.OrderItem{
				{
					ID:        model.CreateUUID().String(),
					OrderID:   orderID.String(),
					ProductID: model.CreateUUID().String(),
					Quantity:  3,
					Price:     1234,
				},
			},
		}).
		MakeRequest()).
		Execute()

	tc.Expect(err).ToNot(HaveOccurred())
	tc.Expect(resp.Code).To(Equal(http.StatusNotFound))
}

func TestUpdateOrder_success(t *testing.T) {
	tc := ts.SetupRoutesTest(t)

	updateStub := orderUpdateStub{err: nil}
	tc.Container.Add(service.MakeOrderUpdateRoute(ts.NewWebservice(), &updateStub))

	orderID := model.CreateUUID()
	orderItemID1 := model.CreateUUID()
	productID := model.CreateUUID()
	customerID := model.CreateUUID()
	resp, _, err := tc.CallTo(gorequest.New().
		Put(fmt.Sprintf("/orders/%s", orderID)).
		Send(api.OrderParameters{
			OrderID:    orderID.String(),
			CustomerID: customerID.String(),
			OrderItems: []api.OrderItem{
				{
					ID:        orderItemID1.String(),
					OrderID:   orderID.String(),
					ProductID: productID.String(),
					Quantity:  3,
					Price:     1234,
				},
			},
		}).
		MakeRequest()).
		Execute()

	tc.Expect(err).ToNot(HaveOccurred())
	tc.Expect(resp.Code).To(Equal(http.StatusNoContent))
	tc.Expect(updateStub.err).To(BeNil())
}

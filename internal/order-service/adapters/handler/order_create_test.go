package handler_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
	"github.com/joel-malina/tucows-challenge/api"
	ts "github.com/joel-malina/tucows-challenge/internal/common/apitestsetup"
	"github.com/joel-malina/tucows-challenge/internal/order-service/model"
	"github.com/joel-malina/tucows-challenge/internal/order-service/service"
	. "github.com/onsi/gomega"
	"github.com/parnurzeal/gorequest"
)

type orderCreateStub struct {
	resultOrderID uuid.UUID
	err           error
}

func (f orderCreateStub) OrderCreate(_ context.Context, order model.Order) (uuid.UUID, error) {
	return f.resultOrderID, f.err
}

func defaultOrderParams(orderID uuid.UUID) api.OrderParameters {
	return api.OrderParameters{
		OrderID:    orderID.String(),
		CustomerID: model.CreateUUID().String(),
		OrderItems: []api.OrderItem{
			{
				ID:        model.CreateUUID().String(),
				OrderID:   orderID.String(),
				ProductID: model.CreateUUID().String(),
				Quantity:  1,
				Price:     2222,
			},
		},
	}
}

func TestCreateOrder_success(t *testing.T) {
	tc := ts.SetupRoutesTest(t)

	orderID := model.CreateUUID()

	createStub := orderCreateStub{
		resultOrderID: orderID,
		err:           nil,
	}

	tc.Container.Add(service.MakeOrderCreateRoute(ts.NewWebservice(), &createStub))

	var response api.OrderCreateResponse
	resp, _, err := tc.CallTo(gorequest.New().
		Post("/orders").
		Type(restful.MIME_JSON).
		Send(defaultOrderParams(orderID)).
		MakeRequest()).
		Read(&response).
		Execute()

	tc.Expect(err).ToNot(HaveOccurred())
	tc.Expect(resp.Code).To(Equal(http.StatusCreated))
	tc.Expect(response).To(Equal(api.OrderCreateResponse{ID: createStub.resultOrderID.String()}))
}

func TestCreateOrderFailureError(t *testing.T) {
	tc := ts.SetupRoutesTest(t)

	createStub := orderCreateStub{
		err: errors.New("something bad happened"),
	}

	tc.Container.Add(service.MakeOrderCreateRoute(ts.NewWebservice(), &createStub))

	resp, _, err := tc.CallTo(gorequest.New().
		Post("/orders").
		Type(restful.MIME_JSON).
		Send(defaultOrderParams(model.CreateUUID())).
		MakeRequest()).
		Execute()

	tc.Expect(err).ToNot(HaveOccurred())
	tc.Expect(resp.Code).To(Equal(http.StatusInternalServerError))
}

func TestCreateOrder_invalid_ID(t *testing.T) {
	tc := ts.SetupRoutesTest(t)

	createStub := orderCreateStub{
		err: nil,
	}

	tc.Container.Add(service.MakeOrderCreateRoute(ts.NewWebservice(), createStub))

	orderParams := defaultOrderParams(uuid.UUID{})
	orderParams.OrderID = "not a proper uuid"

	resp, _, err := tc.CallTo(gorequest.New().
		Post("/orders").
		Type(restful.MIME_JSON).
		Send(orderParams).
		MakeRequest()).
		Execute()

	tc.Expect(err).ToNot(HaveOccurred())
	tc.Expect(resp.Code).To(Equal(http.StatusBadRequest))
}

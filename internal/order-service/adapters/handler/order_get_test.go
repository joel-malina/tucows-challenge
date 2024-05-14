package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joel-malina/tucows-challenge/api"
	ts "github.com/joel-malina/tucows-challenge/internal/common/apitestsetup"
	"github.com/joel-malina/tucows-challenge/internal/order-service/model"
	"github.com/joel-malina/tucows-challenge/internal/order-service/service"
	. "github.com/onsi/gomega"
	"github.com/parnurzeal/gorequest"
)

type orderGetStub struct {
	resultTemplate model.Order
	err            error
}

func (f orderGetStub) OrderGet(_ context.Context, id uuid.UUID) (model.Order, error) {
	if f.err != nil {
		return model.Order{}, f.err
	}
	result := f.resultTemplate
	result.ID = id
	return result, nil
}

func TestGetOrder_invalid_ID(t *testing.T) {
	tc := ts.SetupRoutesTest(t)

	getStub := orderGetStub{
		resultTemplate: model.Order{},
		err:            nil,
	}

	tc.Container.Add(service.MakeOrderGetRoute(ts.NewWebservice(), getStub))

	resp, _, err := tc.CallTo(gorequest.New().
		Get("/orders/nonuuid").
		MakeRequest()).
		Execute()

	tc.Expect(err).ToNot(HaveOccurred())
	tc.Expect(resp.Code).To(Equal(http.StatusBadRequest))
}

func TestGetOrder_success(t *testing.T) {
	tc := ts.SetupRoutesTest(t)
	orderID := model.CreateUUID()
	customerID := model.CreateUUID()
	getStub := orderGetStub{
		resultTemplate: model.Order{
			ID:         orderID,
			CustomerID: customerID,
			OrderDate:  time.Now().UTC(),
			Status:     model.OrderStatusPaymentProcessRequested,
			TotalPrice: 1111,
			OrderItems: []model.OrderItem{
				model.OrderItem{
					ID:        uuid.UUID{},
					OrderID:   orderID,
					ProductID: uuid.UUID{},
					Quantity:  1,
					Price:     1111,
				},
			},
		},
		err: nil,
	}

	tc.Container.Add(service.MakeOrderGetRoute(ts.NewWebservice(), getStub))

	var response api.OrderGetResponse
	resp, _, err := tc.CallTo(gorequest.New().
		Get(fmt.Sprintf("/orders/%s", orderID.String())).
		MakeRequest()).
		Read(&response).
		Execute()

	tc.Expect(err).ToNot(HaveOccurred())
	tc.Expect(resp.Code).To(Equal(http.StatusOK))
	tc.Expect(response).To(Equal(api.OrderGetResponse{
		ID:         orderID.String(),
		CustomerID: customerID.String(),
		OrderItems: []api.OrderItem{
			api.OrderItem{
				ID:        uuid.UUID{}.String(),
				OrderID:   orderID.String(),
				ProductID: uuid.UUID{}.String(),
				Quantity:  1,
				Price:     1111,
			},
		},
	}))
}

func TestGetOrder_not_found(t *testing.T) {
	tc := ts.SetupRoutesTest(t)

	getStub := orderGetStub{
		err: model.ErrOrderNotFound,
	}

	tc.Container.Add(service.MakeOrderGetRoute(ts.NewWebservice(), getStub))

	orderID := model.CreateUUID()
	resp, _, err := tc.CallTo(gorequest.New().
		Get(fmt.Sprintf("/orders/%s", orderID)).
		MakeRequest()).
		Execute()

	tc.Expect(err).ToNot(HaveOccurred())
	tc.Expect(resp.Code).To(Equal(http.StatusNotFound))
}

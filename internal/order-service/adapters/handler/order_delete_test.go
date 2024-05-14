package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	ts "github.com/joel-malina/tucows-challenge/internal/common/apitestsetup"
	"github.com/joel-malina/tucows-challenge/internal/order-service/model"
	"github.com/joel-malina/tucows-challenge/internal/order-service/service"
	. "github.com/onsi/gomega"
	"github.com/parnurzeal/gorequest"
)

type orderDeleteStub struct {
	err error
}

func (f orderDeleteStub) OrderDelete(_ context.Context, id uuid.UUID) error {
	if f.err != nil {
		return f.err
	}
	return nil
}

func TestDeleteOrder_invalid_ID(t *testing.T) {
	tc := ts.SetupRoutesTest(t)

	deleteStub := orderDeleteStub{
		err: nil,
	}

	tc.Container.Add(service.MakeOrderDeleteRoute(ts.NewWebservice(), deleteStub))

	resp, _, err := tc.CallTo(gorequest.New().
		Delete("/orders/notauuid").
		MakeRequest()).
		Execute()

	tc.Expect(err).ToNot(HaveOccurred())
	tc.Expect(resp.Code).To(Equal(http.StatusBadRequest))
}

func TestDeleteOrder_not_found(t *testing.T) {
	tc := ts.SetupRoutesTest(t)

	deleteStub := orderDeleteStub{
		err: model.ErrOrderNotFound,
	}

	tc.Container.Add(service.MakeOrderDeleteRoute(ts.NewWebservice(), deleteStub))

	orderID := model.CreateUUID()
	resp, _, err := tc.CallTo(gorequest.New().
		Delete(fmt.Sprintf("/orders/%s", orderID)).
		MakeRequest()).
		Execute()

	tc.Expect(err).ToNot(HaveOccurred())
	tc.Expect(resp.Code).To(Equal(http.StatusNotFound))
}

func TestDeleteOrder_success(t *testing.T) {
	tc := ts.SetupRoutesTest(t)

	tc.Container.Add(service.MakeOrderDeleteRoute(ts.NewWebservice(), orderDeleteStub{}))

	resp, _, err := tc.CallTo(gorequest.New().
		Delete("/orders/" + model.CreateUUID().String()).
		MakeRequest()).
		Execute()

	tc.Expect(err).ToNot(HaveOccurred())
	tc.Expect(resp.Code).To(Equal(http.StatusNoContent))
}

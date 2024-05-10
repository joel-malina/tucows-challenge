package handler

import (
	"context"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/joel-malina/tucows-challenge/api"
	"github.com/joel-malina/tucows-challenge/internal/common/response"
	"github.com/joel-malina/tucows-challenge/internal/order-service/model"
	"github.com/sirupsen/logrus"
)

type OrdersGetter interface {
	OrdersGet(ctx context.Context) ([]model.Order, error)
}

func OrdersGetHandler(svc OrdersGetter) func(*restful.Request, *restful.Response) {
	return func(req *restful.Request, resp *restful.Response) {
		ctx := req.Request.Context()
		log := logrus.WithContext(ctx)

		orders, err := svc.OrdersGet(ctx)
		if err != nil {
			response.WriteErrorWithContext("failed to get orders", log, resp, http.StatusInternalServerError, err)
			return
		}

		var orderGetResponse []api.OrderGetResponse
		for _, order := range orders {

			orderResponse := api.OrderGetResponse{
				ID: order.ID.String(),
			}
			orderGetResponse = append(orderGetResponse, orderResponse)
		}

		response.Write(resp, http.StatusOK, orderGetResponse)
	}
}

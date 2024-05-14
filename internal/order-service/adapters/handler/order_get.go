package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
	"github.com/joel-malina/tucows-challenge/api"
	"github.com/joel-malina/tucows-challenge/internal/common/response"
	"github.com/joel-malina/tucows-challenge/internal/order-service/model"
	"github.com/sirupsen/logrus"
)

type OrderGetter interface {
	OrderGet(ctx context.Context, id uuid.UUID) (model.Order, error)
}

func OrderGetHandler(svc OrderGetter) func(*restful.Request, *restful.Response) {
	return func(req *restful.Request, resp *restful.Response) {
		ctx := req.Request.Context()
		log := logrus.WithContext(ctx)

		orderID, err := getIDPathParameter(req, PathParamOrderID)
		if err != nil {
			response.WriteErrorWithContext("failed to parse ID", log, resp, http.StatusBadRequest, err)
			return
		}
		log = log.WithField(model.LogFieldOrderID, orderID)

		order, err := svc.OrderGet(ctx, orderID)
		if err != nil {
			if errors.Is(err, model.ErrOrderNotFound) {
				response.WriteErrorWithContext("failed to find order", log, resp, http.StatusNotFound, err)
				return
			}

			response.WriteErrorWithContext("failed to get order", log, resp, http.StatusInternalServerError, err)
			return
		}

		response.Write(resp, http.StatusOK, toOrderGetResponse(order))
	}
}

func toOrderGetResponse(order model.Order) api.OrderGetResponse {

	var apiOrderItems []api.OrderItem
	for i := range order.OrderItems {
		item := api.OrderItem{
			ID:        order.OrderItems[i].ID.String(),
			OrderID:   order.OrderItems[i].OrderID.String(),
			ProductID: order.OrderItems[i].ProductID.String(),
			Quantity:  order.OrderItems[i].Quantity,
			Price:     order.OrderItems[i].Price,
		}
		apiOrderItems = append(apiOrderItems, item)
	}

	return api.OrderGetResponse{
		ID:         order.ID.String(),
		CustomerID: order.CustomerID.String(),
		OrderItems: apiOrderItems,
	}
}

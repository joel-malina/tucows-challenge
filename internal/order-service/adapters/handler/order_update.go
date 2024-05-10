package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/joel-malina/tucows-challenge/api"
	"github.com/joel-malina/tucows-challenge/internal/common/response"
	"github.com/joel-malina/tucows-challenge/internal/order-service/model"
	"github.com/sirupsen/logrus"
)

type OrderUpdater interface {
	OrderUpdate(context context.Context, updateParameters model.Order) error
}

func OrderUpdateHandler(svc OrderUpdater) func(*restful.Request, *restful.Response) {
	return func(req *restful.Request, resp *restful.Response) {
		ctx := req.Request.Context()
		log := logrus.WithContext(ctx)

		var payload api.OrderParameters
		if err := req.ReadEntity(&payload); err != nil {
			response.WriteErrorWithContext("failed to parse request", log, resp, http.StatusBadRequest, err)
			return
		}

		//payload.Name = strings.TrimSpace(payload.Name)
		//if err := payload.Validate(); err != nil {
		//	response.WriteErrorWithContext("failed to validate request", log, resp, http.StatusBadRequest, err)
		//	return
		//}

		orderID, err := getIDPathParameter(req, PathParamOrderID)
		if err != nil {
			response.WriteErrorWithContext("failed to parse order ID", log, resp, http.StatusBadRequest, err)
			return
		}
		log = log.WithField(model.LogFieldOrderID, orderID)

		order, err := orderParametersToOrder(ctx, payload)
		order.ID = orderID
		if err != nil {
			response.WriteErrorWithContext("failed to update order", log, resp, http.StatusInternalServerError, err)
			return
		}

		if err := svc.OrderUpdate(ctx, order); err != nil {
			httpStatus := http.StatusInternalServerError
			if errors.Is(err, model.ErrOrderNotFound) {
				httpStatus = http.StatusNotFound
			}

			response.WriteErrorWithContext("failed to update order", log, resp, httpStatus, err)
			return
		}

		resp.WriteHeader(http.StatusNoContent)
	}
}

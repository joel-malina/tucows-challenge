package handler

import (
	"context"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
	"github.com/joel-malina/tucows-challenge/api"
	"github.com/joel-malina/tucows-challenge/internal/common/response"
	"github.com/joel-malina/tucows-challenge/internal/order-service/model"
	"github.com/sirupsen/logrus"
)

type OrderCreator interface {
	OrderCreate(ctx context.Context, order model.Order) (uuid.UUID, error)
}

func OrderCreateHandler(svc OrderCreator) func(*restful.Request, *restful.Response) {
	return func(req *restful.Request, resp *restful.Response) {
		ctx := req.Request.Context()
		log := logrus.WithContext(ctx)

		payload := api.OrderParameters{}
		if err := req.ReadEntity(&payload); err != nil {
			response.WriteErrorWithContext("failed to parse request", log, resp, http.StatusBadRequest, err)
			return
		}

		//payload.Name = strings.TrimSpace(payload.Name)
		//if err := payload.Validate(); err != nil {
		//	response.WriteErrorWithContext("failed to validate request", log, resp, http.StatusBadRequest, err)
		//	return
		//}

		order, err := orderParametersToOrder(ctx, payload)
		if err != nil {
			response.WriteErrorWithContext("failed to create order", log, resp, http.StatusInternalServerError, err)
			return
		}

		id, err := svc.OrderCreate(ctx, order)
		if err != nil {
			//if err == model.ErrOrderLimitViolation {
			//	response.WriteErrorWithContext("failed to create order", log, resp, http.StatusForbidden, err)
			//	return
			//}

			response.WriteErrorWithContext("failed to create order", log, resp, http.StatusInternalServerError, err)
			return
		}

		response.Write(resp, http.StatusCreated, api.OrderCreateResponse{
			ID: id.String(),
		})
	}
}

// TODO implement missing values
func orderParametersToOrder(ctx context.Context, f api.OrderParameters) (model.Order, error) {
	return model.Order{}, nil
}

package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
	"github.com/joel-malina/tucows-challenge/internal/common/response"
	"github.com/joel-malina/tucows-challenge/internal/order-service/model"
	"github.com/sirupsen/logrus"
)

type OrderDeleter interface {
	OrderDelete(ctx context.Context, id uuid.UUID) error
}

func OrderDeleteHandler(svc OrderDeleter) func(*restful.Request, *restful.Response) {
	return func(req *restful.Request, resp *restful.Response) {
		ctx := req.Request.Context()
		log := logrus.WithContext(ctx)

		orderID, err := getIDPathParameter(req, PathParamOrderID)
		if err != nil {
			response.WriteErrorWithContext("failed to parse order ID", log, resp, http.StatusBadRequest, err)
			return
		}
		log = log.WithField(model.LogFieldOrderID, orderID)

		if err := svc.OrderDelete(req.Request.Context(), orderID); err != nil {
			if errors.Is(err, model.ErrOrderNotFound) {
				response.WriteErrorWithContext("failed to delete order", log, resp, http.StatusNotFound, err)
				return
			}

			response.WriteErrorWithContext("failed to delete order", log, resp, http.StatusInternalServerError, err)
			return
		}

		resp.WriteHeader(http.StatusNoContent)
	}
}

package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

		// Likely want to validate the payload based on a set of criteria
		//if err := payload.Validate(); err != nil {
		//	response.WriteErrorWithContext("failed to validate request", log, resp, http.StatusBadRequest, err)
		//	return
		//}

		order, err := orderParametersToOrder(payload)
		if err != nil {
			response.WriteErrorWithContext("failed to create order", log, resp, http.StatusBadRequest, err)
			return
		}

		id, err := svc.OrderCreate(ctx, order)
		if err != nil {
			// Could handle specific errors more gracefully
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

func orderParametersToOrder(apiOrder api.OrderParameters) (model.Order, error) {

	orderID, err := uuid.Parse(apiOrder.OrderID)
	if err != nil {
		return model.Order{}, fmt.Errorf("could not create order with provided orderID:%s got err %w", apiOrder.OrderID, err)
	}

	customerID, err := uuid.Parse(apiOrder.CustomerID)
	if err != nil {
		return model.Order{}, fmt.Errorf("could not create order with provided customerID:%s got err %w", apiOrder.CustomerID, err)
	}

	var totalPrice float64
	var items []model.OrderItem
	for i := range apiOrder.OrderItems {

		id, err := uuid.Parse(apiOrder.OrderItems[i].ID)
		if err != nil {
			return model.Order{}, fmt.Errorf("could not parse orderItem with provided ID:%s got err %w", apiOrder.OrderItems[i].ID, err)
		}

		orderItemOrderID, err := uuid.Parse(apiOrder.OrderItems[i].OrderID)
		if err != nil {
			return model.Order{}, fmt.Errorf("could not parse orderItem with provided orderID:%s got err %w", apiOrder.OrderItems[i].OrderID, err)
		}

		if orderItemOrderID != orderID {
			return model.Order{}, fmt.Errorf("OrderItem %s with order ID %s does not match with root order ID of %s", apiOrder.OrderItems[i].ID, apiOrder.OrderItems[i].OrderID, orderID)
		}

		productID, err := uuid.Parse(apiOrder.OrderItems[i].ProductID)
		if err != nil {
			return model.Order{}, fmt.Errorf("could not parse orderItem with provided productID:%s got err %w", apiOrder.OrderItems[i].ProductID, err)
		}

		item := model.OrderItem{
			ID:        id,
			OrderID:   orderItemOrderID,
			ProductID: productID,
			Quantity:  apiOrder.OrderItems[i].Quantity,
			Price:     apiOrder.OrderItems[i].Price,
		}

		totalPrice += apiOrder.OrderItems[i].Price * float64(apiOrder.OrderItems[i].Quantity)
		items = append(items, item)
	}

	return model.Order{
		ID:         orderID,
		CustomerID: customerID,
		OrderDate:  time.Now().UTC(),
		Status:     model.OrderStatusCreated,
		TotalPrice: totalPrice,
		OrderItems: items,
	}, nil
}

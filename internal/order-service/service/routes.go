package service

import (
	"fmt"
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/joel-malina/tucows-challenge/api"
	"github.com/joel-malina/tucows-challenge/internal/common/response"
	"github.com/joel-malina/tucows-challenge/internal/order-service/adapters/handler"
)

var (
	orderServiceAPITag = []string{"Order Service"}
)

func MakeOrderCreateRoute(service *restful.WebService, orderCreator handler.OrderCreator) *restful.WebService {
	return service.Route(service.POST("/orders").
		To(handler.OrderCreateHandler(orderCreator)).
		Metadata(restfulspec.KeyOpenAPITags, orderServiceAPITag).
		Doc("create an order").
		Operation("OrderCreate").
		Reads(api.OrderParameters{}).
		Returns(http.StatusCreated, "success", api.OrderCreateResponse{}).
		Returns(http.StatusBadRequest, "bad request", response.ErrorResponse{}).
		Returns(http.StatusForbidden, "exceeded quota", response.ErrorResponse{}).
		Returns(http.StatusInternalServerError, "internal server error", response.ErrorResponse{}).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))
}

func MakeOrderGetRoute(service *restful.WebService, getter handler.OrderGetter) *restful.WebService {
	return service.Route(service.GET(fmt.Sprintf("/orders/{%s}", handler.PathParamOrderID)).
		Param(service.PathParameter(handler.PathParamOrderID, "the id of the order").
			DataType("string").Required(true)).
		To(handler.OrderGetHandler(getter)).
		Filter(handler.ValidatePathParameter).
		Metadata(restfulspec.KeyOpenAPITags, orderServiceAPITag).
		Doc("get an order").
		Operation("OrderGet").
		Returns(http.StatusOK, "success", api.OrderGetResponse{}).
		Returns(http.StatusBadRequest, "bad request", response.ErrorResponse{}).
		Returns(http.StatusNotFound, "order not found", response.ErrorResponse{}).
		Returns(http.StatusInternalServerError, "internal server error", response.ErrorResponse{}).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))
}

func MakeOrderGetAllRoute(service *restful.WebService, getter handler.OrdersGetter) *restful.WebService {
	return service.Route(service.GET("/orders").
		To(handler.OrdersGetHandler(getter)).
		Filter(handler.ValidatePathParameter).
		Metadata(restfulspec.KeyOpenAPITags, orderServiceAPITag).
		Doc("get all orders").
		Operation("OrdersGet").
		Returns(http.StatusOK, "success", api.OrdersGetResponse{}).
		Returns(http.StatusInternalServerError, "internal server error", response.ErrorResponse{}).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))
}

func MakeOrderUpdateRoute(service *restful.WebService, updater handler.OrderUpdater) *restful.WebService {
	return service.Route(service.PUT(fmt.Sprintf("/orders/{%s}", handler.PathParamOrderID)).
		Param(service.PathParameter(handler.PathParamOrderID, "the id of the order").
			DataType("string").Required(true)).
		To(handler.OrderUpdateHandler(updater)).
		Filter(handler.ValidatePathParameter).
		Metadata(restfulspec.KeyOpenAPITags, orderServiceAPITag).
		Doc("update an order -– overrides current data").
		Operation("OrderUpdate").
		Reads(api.OrderParameters{}).
		Returns(http.StatusNoContent, "no content", nil).
		Returns(http.StatusBadRequest, "bad request", response.ErrorResponse{}).
		Returns(http.StatusNotFound, "order not found", response.ErrorResponse{}).
		Returns(http.StatusInternalServerError, "internal server error", response.ErrorResponse{}).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))
}

func MakeOrderDeleteRoute(service *restful.WebService, deleter handler.OrderDeleter) *restful.WebService {
	return service.Route(service.DELETE(fmt.Sprintf("/orders/{%s}", handler.PathParamOrderID)).
		Param(service.PathParameter(handler.PathParamOrderID, "the id of the order").
			DataType("string").Required(true)).
		To(handler.OrderDeleteHandler(deleter)).
		Filter(handler.ValidatePathParameter).
		Metadata(restfulspec.KeyOpenAPITags, orderServiceAPITag).
		Doc("delete an order").
		Operation("OrderDelete").
		Returns(http.StatusNoContent, "no content", nil).
		Returns(http.StatusBadRequest, "bad request", response.ErrorResponse{}).
		Returns(http.StatusNotFound, "order not found", response.ErrorResponse{}).
		Returns(http.StatusInternalServerError, "internal server error", response.ErrorResponse{}).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))
}

package handler

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/joel-malina/tucows-challenge/internal/common/response"
)

type basicHealthResponse struct {
	Name    string `json:"name"`
	Healthy bool   `json:"healthy"`
}

func BasicHealthCheck(serviceName string) func(_ *restful.Request, resp *restful.Response) {
	return func(_ *restful.Request, resp *restful.Response) {
		response.Write(resp, http.StatusOK, basicHealthResponse{
			Name:    serviceName,
			Healthy: true,
		})
	}
}

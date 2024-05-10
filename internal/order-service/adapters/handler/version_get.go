package handler

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/joel-malina/tucows-challenge/api"
	"github.com/joel-malina/tucows-challenge/internal/common/response"
	"github.com/joel-malina/tucows-challenge/internal/order-service/config"
)

func VersionInfo(serviceConfig config.ServiceConfig) func(_ *restful.Request, response *restful.Response) {
	return func(req *restful.Request, resp *restful.Response) {
		response.Write(resp, http.StatusOK, api.VersionInfoResponse{
			Name:      serviceConfig.ServiceName,
			BuildDate: serviceConfig.ServiceBuildDate,
			GitHash:   serviceConfig.ServiceGitHash,
			Version:   serviceConfig.ServiceVersion,
		})
	}
}

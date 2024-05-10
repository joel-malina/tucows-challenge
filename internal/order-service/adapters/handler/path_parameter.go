package handler

import (
	"errors"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
	"github.com/joel-malina/tucows-challenge/internal/common/response"
	"github.com/sirupsen/logrus"
)

var (
	errInvalidID = errors.New("invalid id")
)

const (
	PathParamOrderID = "orderID"
)

// ValidatePathParameter is a go-restful filter to validate the query parameter
func ValidatePathParameter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	for param, val := range req.PathParameters() {
		switch param {
		case PathParamOrderID:
			_, err := uuid.Parse(val)
			if err != nil {
				response.WriteErrorWithContext("failed to parse order ID", logrus.WithContext(req.Request.Context()), resp, http.StatusBadRequest, errInvalidID)
				return
			}
			// can add more cases as required
		}
	}
	chain.ProcessFilter(req, resp)
}

func getIDPathParameter(req *restful.Request, pathParam string) (uuid.UUID, error) {
	idParam := req.PathParameter(pathParam)
	id, err := uuid.Parse(idParam)
	if err != nil {
		return id, err
	}

	return id, nil
}

package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/emicklei/go-restful/v3"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	// could put tracing in this too
	// TraceID      string `json:"traceId"`
	ErrorMessage string `json:"errorMessage"`
	ErrorType    string `json:"errorType"`
}

// Write sets the response header to 'application/json' and allows for custom http response status. If value
// cannot be encoded, the request will fail with a http.StatusInternalServerError
func Write(resp *restful.Response, status int, value interface{}) {
	if err := resp.WriteHeaderAndJson(status, value, restful.MIME_JSON); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
	}
}

// WriteError sets the response header to `application/json` and allows for custom http status. If value
// cannot be encoded, the request will fail with http.StatusInternalServerError
func WriteError(resp *restful.Response, status int, err error) {
	val := ErrorResponse{
		ErrorMessage: err.Error(),
		ErrorType:    reflect.TypeOf(err).Name(), // reflect on the error to get name of concrete type
	}
	if err := resp.WriteHeaderAndJson(status, val, restful.MIME_JSON); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
	}
}

// WriteErrorWithContext sets the response header to `application/json` and allows for custom http status.
// If the status code is http.StatusNotFound, the error will not be logged
// If the value cannot be encoded, the request will fail with a http.StatusInternalServerError
func WriteErrorWithContext(errContext string, log *logrus.Entry, resp *restful.Response, status int, err error) {
	if status != http.StatusNotFound {
		if log != nil {
			log.WithError(err).Error(errContext)
		}
	}

	val := ErrorResponse{
		ErrorMessage: errContext + ": " + err.Error(),
		ErrorType:    reflect.TypeOf(err).Name(),
	}
	if err := resp.WriteHeaderAndJson(status, val, restful.MIME_JSON); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
	}
}

func ParseErrorResponse(body []byte) (ErrorResponse, error) {
	var errorResponse ErrorResponse
	err := json.Unmarshal(body, &errorResponse)
	if err != nil {
		return ErrorResponse{}, ErrNotErrorResponse{body: string(body)}
	}
	return errorResponse, err
}

type ErrNotErrorResponse struct {
	body string
}

func (e ErrNotErrorResponse) Error() string {
	return fmt.Sprintf("body was not an ErrorResponse: %s", e.body)
}

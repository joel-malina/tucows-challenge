package apitestsetup

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful/v3"
	"github.com/joel-malina/tucows-challenge/internal/common/caller"
	. "github.com/onsi/gomega" //nolint:stylecheck
)

type RouteTest struct {
	*GomegaWithT
	Container *restful.Container
}

// CallTo receives a request to be sent to the handler returns a TestCaller.
// The second parameter is a convenience if you want to pass directly from http.NewRequest() or gorequest.MakeRequest()
func (rt RouteTest) CallTo(req *http.Request, _ ...interface{}) *caller.HTTPTestCaller {
	return caller.Call(rt.Container).To(req)
}

// ExecuteCallTo receives a request to be sent to the handler and executes it using a TestCaller.
// The second parameter is a convenience if you want to pass directly from http.NewRequest() or gorequest.MakeRequest()
func (rt RouteTest) ExecuteCallTo(req *http.Request, _ ...interface{}) (*httptest.ResponseRecorder, interface{}, error) {
	return caller.Call(rt.Container).To(req).Execute()
}

func SetupRoutesTest(t *testing.T) RouteTest {
	t.Parallel()
	g := NewWithT(t)

	container := restful.NewContainer()

	test := RouteTest{
		g,
		container,
	}
	return test
}

func NewWebservice() *restful.WebService {
	return new(restful.WebService).Path("/")
}

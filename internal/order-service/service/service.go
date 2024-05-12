package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp/filters"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
	"github.com/joel-malina/tucows-challenge/internal/order-service/adapters/handler"
	"github.com/joel-malina/tucows-challenge/internal/order-service/config"
	"github.com/joel-malina/tucows-challenge/internal/order-service/order"
	"github.com/sirupsen/logrus"
)

func Run(ctx context.Context, serviceConfig config.ServiceConfig, storageResolver *StorageResolver) {

	log := logrus.New()

	orderHandlerLogic := setupServiceDependencies(ctx, log, serviceConfig, storageResolver)

	container := setupWebServiceContainer(serviceConfig, orderHandlerLogic)

	log.WithFields(logrus.Fields{
		"serviceName": serviceConfig.ServiceName,
		"port":        serviceConfig.Port,
		"logLevel":    log.Level,
	}).Info("starting service on port ", serviceConfig.Port)

	for _, webService := range container.RegisteredWebServices() {
		for _, route := range webService.Routes() {
			log.Printf("%s %s", route.Method, route.Path)
		}
	}

	addr := fmt.Sprintf(":%v", serviceConfig.Port)
	srv := http.Server{
		Addr: addr,
		Handler: otelhttp.NewHandler(
			container,
			fmt.Sprintf("%s request", serviceConfig.ServiceName),
			otelhttp.WithFilter(filters.PathPrefix(serviceConfig.BasePath)),
		),
	}

	idleConnectionsClosed := make(chan struct{})
	go func() {
		defer close(idleConnectionsClosed)
		<-ctx.Done()
		// We received a signal to shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		log.Println("server shutdown initiated")
	}()

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnectionsClosed
	log.Println("exiting")
}

// handler.OrdersGetter
type orderHandlers struct {
	handler.OrderCreator
	handler.OrderGetter
	handler.OrderUpdater
	handler.OrderDeleter
}

func setupServiceDependencies(_ context.Context, log *logrus.Logger, serviceConfig config.ServiceConfig, storage *StorageResolver) orderHandlers {

	// setup auth here if it were a public API
	storage.Resolve(log, serviceConfig)
	orderRepo := storage.OrderRepository

	// OrdersGetter: order.NewOrdersGet(orderRepo),
	orderHandlerLogic := orderHandlers{
		OrderCreator: order.NewOrderCreate(orderRepo),
		OrderGetter:  order.NewOrderGet(orderRepo),
		OrderUpdater: order.NewOrderUpdate(orderRepo),
		OrderDeleter: order.NewOrderDelete(orderRepo),
	}

	return orderHandlerLogic
}

var restfulTestRaceOnce = sync.Once{}

func setupWebServiceContainer(serviceConfig config.ServiceConfig, orderHandlerLogic orderHandlers) *restful.Container {

	// during parallel test runs this function might run concurrently and since we're setting restful config, which later
	// gets read in the restful internals (during the execution of this method), we just want to do this once to avoid
	// race detection failures
	restfulTestRaceOnce.Do(func() {
		restful.TrimRightSlashEnabled = false
	})

	container := restful.NewContainer()

	setupHealthCheck(container, serviceConfig.ServiceName, "", "BasicHealthCheck")
	setupVersion(container, serviceConfig)

	// If I were doing Auth for this API we'd want to do some setup here
	// SetupAuthCheck(container, serviceConfig, authHandlerLogic)

	setupV1Routes(container, serviceConfig, orderHandlerLogic)
	setupAPIDocs(container, serviceConfig, container)

	return container
}

func setupHealthCheck(container *restful.Container, serviceName string, basePath string, operation string) {
	healthzPath := fmt.Sprintf("%s/healthz", basePath)
	service := new(restful.WebService).Path(healthzPath)
	service.Route(service.GET("/").
		To(handler.BasicHealthCheck(serviceName)).
		Doc("Health check").
		Metadata(restfulspec.KeyOpenAPITags, orderServiceAPITag).
		Operation("version").
		Operation(operation).
		Produces(restful.MIME_JSON))
	container.Add(service)
}

func setupVersion(container *restful.Container, serviceConfig config.ServiceConfig) {
	service := new(restful.WebService).Path(serviceConfig.BasePath)
	service.Route(service.GET("/version").
		To(handler.VersionInfo(serviceConfig)).
		Doc("Version info").
		Metadata(restfulspec.KeyOpenAPITags, orderServiceAPITag).
		Produces(restful.MIME_JSON))
	container.Add(service)
}

func setupV1Routes(container *restful.Container, serviceConfig config.ServiceConfig, orderHandlers orderHandlers) {
	v1RootPath := fmt.Sprintf("%s/v1", serviceConfig.BasePath)
	v1RootRoutes := new(restful.WebService).Path(v1RootPath)

	// would add auth filter to these, possibly namespace them too
	MakeOrderCreateRoute(v1RootRoutes, orderHandlers)
	MakeOrderGetRoute(v1RootRoutes, orderHandlers)
	//MakeOrderGetAllRoute(v1RootRoutes, orderHandlers)
	MakeOrderUpdateRoute(v1RootRoutes, orderHandlers)
	MakeOrderDeleteRoute(v1RootRoutes, orderHandlers)

	container.Add(v1RootRoutes)
}

func setupAPIDocs(container *restful.Container, serviceConfig config.ServiceConfig, serviceContainer *restful.Container) {
	apiDocsPath := serviceConfig.BasePath + "/apidocs/"
	swaggerConfig := restfulspec.Config{
		WebServices: serviceContainer.RegisteredWebServices(),
		APIPath:     apiDocsPath + "api.json",
		PostBuildSwaggerObjectHandler: func(s *spec.Swagger) {
			s.Info = &spec.Info{
				InfoProps: spec.InfoProps{
					Title:       serviceConfig.ServiceName,
					Description: "A Service to manage orders",
					Version:     serviceConfig.ServiceVersion,
				},
			}
			s.SecurityDefinitions = map[string]*spec.SecurityScheme{
				"authorization": spec.APIKeyAuth("Authorization", "header"),
			}
			s.Security = []map[string][]string{
				{"authorization": {}},
			}
		},
	}

	prefix := http.StripPrefix(apiDocsPath, http.FileServer(http.Dir("swagger-ui")))
	serviceContainer.Handle(apiDocsPath, prefix)

	container.Add(restfulspec.NewOpenAPIService(swaggerConfig))
}

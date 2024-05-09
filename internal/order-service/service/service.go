package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp/filters"

	"github.com/emicklei/go-restful/v3"
	"github.com/joel-malina/tucows-challenge/internal/order-service/config"
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

// TODO: Make these routes/handlers
type orderHandlers struct {
	handler.orderCreator
	handler.orderGetter
	handler.orderUpdater
	handler.orderDeleter
}

func setupServiceDependencies(_ context.Context, log *logrus.Logger, serviceConfig config.ServiceConfig, storage *StorageResolver) orderHandlers {

	// setup auth here if it were a public API
	storage.Resolve(log, serviceConfig)
	orderRepo := storage.OrderRepository

	orderHandlerLogic := orderHandlers{
		orderCreator: orderRepo,
		orderGetter:  orderRepo,
		orderUpdater: orderRepo,
		orderDeleter: orderRepo,
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

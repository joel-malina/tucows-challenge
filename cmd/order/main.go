package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joel-malina/tucows-challenge/internal/order-service/config"
	"github.com/joel-malina/tucows-challenge/internal/order-service/service"
	"github.com/sirupsen/logrus"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.ParseConfiguration()
	if err != nil {
		logrus.Fatalf("could not parse env variables: %v", err)
	}

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdownSignal
		logrus.Info("got shutdown signal")
		cancel()
	}()

	service.Run(ctx, cfg, &service.StorageResolver{})
}

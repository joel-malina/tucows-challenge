package service

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/joel-malina/tucows-challenge/internal/order-service/config"
)

func connectToPostgres(log *logrus.Logger, serviceConfig config.ServiceConfig) *sqlx.DB {
	dataSourceName := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		serviceConfig.PostgresHost,
		serviceConfig.PostgresPort,
		serviceConfig.PostgresDB,
		serviceConfig.PostgresUser,
		serviceConfig.PostgresPassword,
		serviceConfig.PostgresSSLMode,
	)
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		log.WithError(err).Fatal("unable to connect to DB")
	}

	db.SetMaxIdleConns(serviceConfig.PostgresMaxIdleConnection)
	db.SetMaxOpenConns(serviceConfig.PostgresMaxOpenConnection)
	return db
}

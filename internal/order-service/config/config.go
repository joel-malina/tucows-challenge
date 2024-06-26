package config

import (
	"github.com/caarlos0/env"
	"github.com/caarlos0/env/parsers"
)

type ServiceConfig struct {
	ServiceVersion   string
	ServiceGitHash   string
	ServiceBuildDate string

	BasePath    string `env:"BASE_PATH" envDefault:"/order-service" envDocs:"The base path for the REST api"`
	ServiceName string `env:"SERVICE_NAME" envDefault:"order-service" envDocs:"The name of the service"`
	Port        int    `env:"PORT" envDefault:"8080" envDocs:"The port which the service will listen to"`

	LogLevel  string `env:"LOG_LEVEL" envDefault:"info" envDocs:"Determines what log level to output"`
	LogFormat string `env:"LOG_FORMAT" envDefault:"text" envDocs:"Determines what log format to output"`

	EnablePersistentStorage bool `env:"ENABLE_PERSISTENT_STORAGE" envDefault:"true" envDocs:"Use the postgres backed persistent storage for server record storage (vs in memory)"`

	PostgresHost              string `env:"POSTGRES_HOST" envDefault:"localhost"`
	PostgresPort              string `env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresDB                string `env:"POSTGRES_DB" envDefault:"tucows-challenge"`
	PostgresUser              string `env:"POSTGRES_USERNAME" envDefault:"postgres"`
	PostgresPassword          string `env:"POSTGRES_PASSWORD" envDefault:"example"`
	PostgresSSLMode           string `env:"POSTGRES_SSL_MODE" envDefault:"disable"`
	PostgresMaxIdleConnection int    `env:"POSTGRES_MAX_IDLE_CONNECTION" envDefault:"2"`
	PostgresMaxOpenConnection int    `env:"POSTGRES_MAX_OPEN_CONNECTION" envDefault:"10"`

	RabbitMQUser     string `env:"RABBITMQ_USER" envDefault:"user"`
	RabbitMQPassword string `env:"RABBITMQ_PASSWORD" envDefault:"password"`
	RabbitMQHost     string `env:"RABBITMQ_HOST" envDefault:"localhost"`
	RabbitMQPort     string `env:"RABBITMQ_PORT" envDefault:"5672"`
}

// ParseConfiguration read the environment variables overwriting any defaults set. Any ServiceConfig values
// not set in the environment will have the default values
func ParseConfiguration() (ServiceConfig, error) {
	var cfg ServiceConfig
	if err := env.ParseWithFuncs(&cfg, env.CustomParsers{parsers.URLType: parsers.URLFunc}); err != nil {
		return ServiceConfig{}, err
	}

	return cfg, nil
}

package config

import (
	"github.com/caarlos0/env"
	"github.com/caarlos0/env/parsers"
)

type ServiceConfig struct {
	ServiceVersion   string
	ServiceGitHash   string
	ServiceBuildDate string

	BasePath    string `env:"BASE_PATH" envDefault:"/payment-service" envDocs:"The base path for the REST api"`
	ServiceName string `env:"SERVICE_NAME" envDefault:"payment-service" envDocs:"The name of the service"`
	Port        int    `env:"PORT" envDefault:"8082" envDocs:"The port which the service will listen to"`

	LogLevel  string `env:"LOG_LEVEL" envDefault:"info" envDocs:"Determines what log level to output"`
	LogFormat string `env:"LOG_FORMAT" envDefault:"text" envDocs:"Determines what log format to output"`

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

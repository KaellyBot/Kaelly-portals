package models

import "github.com/rs/zerolog"

const (
	ConfigFileName = ".env"

	// Dofus Portals Token
	DofusPortalsToken = "DOFUS_PORTALS_TOKEN"

	// Timeout to retrieve portals in seconds
	HttpTimeout = "HTTP_TIMEOUT"

	// RabbitMQ address
	RabbitMqAddress = "RABBITMQ_ADDRESS"

	// Zerolog values from [trace, debug, info, warn, error, fatal, panic]
	LogLevel = "LOG_LEVEL"

	// Boolean; used to register commands at development guild level or globally.
	Production = "PRODUCTION"
)

var (
	DefaultConfigValues = map[string]interface{}{
		DofusPortalsToken: "",
		HttpTimeout:       10,
		RabbitMqAddress:   "amqp://localhost:5672",
		LogLevel:          zerolog.InfoLevel.String(),
		Production:        false,
	}
)

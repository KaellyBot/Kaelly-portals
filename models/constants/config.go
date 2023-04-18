package constants

import "github.com/rs/zerolog"

const (
	ConfigFileName = ".env"

	// Dofus Portals Token
	DofusPortalsToken = "DOFUS_PORTALS_TOKEN"

	// MySQL URL with the following format: HOST:PORT
	MySqlUrl = "MYSQL_URL"

	// MySQL user
	MySqlUser = "MYSQL_USER"

	// MySQL password
	MySqlPassword = "MYSQL_PASSWORD"

	// MySQL database name
	MySqlDatabase = "MYSQL_DATABASE"

	// RabbitMQ address
	RabbitMqAddress = "RABBITMQ_ADDRESS"

	// Timeout to retrieve portals in seconds
	HttpTimeout = "HTTP_TIMEOUT"

	// Metric port
	MetricPort = "METRIC_PORT"

	// Zerolog values from [trace, debug, info, warn, error, fatal, panic]
	LogLevel = "LOG_LEVEL"

	// Boolean; used to register commands at development guild level or globally.
	Production = "PRODUCTION"
)

var (
	DefaultConfigValues = map[string]interface{}{
		DofusPortalsToken: "",
		MySqlUrl:          "localhost:3306",
		MySqlUser:         "",
		MySqlPassword:     "",
		MySqlDatabase:     "kaellybot",
		RabbitMqAddress:   "amqp://localhost:5672",
		HttpTimeout:       10,
		MetricPort:        2112,
		LogLevel:          zerolog.InfoLevel.String(),
		Production:        false,
	}
)

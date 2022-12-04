package models

import "github.com/rs/zerolog"

const (
	LogFileName      = "fileName"
	LogCorrelationId = "correlationId"
	LogServerId      = "serverId"
	LogDimensionId   = "dimensionId"

	LogLevelFallback = zerolog.InfoLevel
)

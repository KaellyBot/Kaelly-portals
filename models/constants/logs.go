package constants

import "github.com/rs/zerolog"

const (
	LogFileName        = "fileName"
	LogCorrelationId   = "correlationId"
	LogServerId        = "serverId"
	LogDimensionId     = "dimensionId"
	LogAreaId          = "areaId"
	LogSubAreaId       = "subAreaId"
	LogTransportTypeId = "transportTypeId"

	LogLevelFallback = zerolog.InfoLevel
)

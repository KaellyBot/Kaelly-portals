package constants

import "github.com/rs/zerolog"

const (
	LogFileName        = "fileName"
	LogCorrelationID   = "correlationID"
	LogServerID        = "serverID"
	LogDimensionID     = "dimensionID"
	LogAreaID          = "areaID"
	LogReplyTo         = "replyTo"
	LogSubAreaID       = "subAreaID"
	LogTransportTypeID = "transportTypeID"

	LogLevelFallback = zerolog.InfoLevel
)

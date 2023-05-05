package transports

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/utils/databases"
)

type Repository interface {
	GetTransportTypes() ([]entities.TransportType, error)
}

type Impl struct {
	db databases.MySQLConnection
}

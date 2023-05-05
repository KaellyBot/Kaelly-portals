package transports

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/repositories/transports"
)

type Service interface {
	FindTransportTypeByDofusPortalsID(dofusPortalsID string) (entities.TransportType, bool)
}

type Impl struct {
	transportTypes    map[string]entities.TransportType
	transportTypeRepo transports.Repository
}

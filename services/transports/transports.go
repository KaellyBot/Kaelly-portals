package transports

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/repositories/transports"
)

func New(transportTypeRepo transports.Repository) (*Impl, error) {
	transportTypeEntities, err := transportTypeRepo.GetTransportTypes()
	if err != nil {
		return nil, err
	}

	transportTypes := make(map[string]entities.TransportType)
	for _, transportType := range transportTypeEntities {
		transportTypes[transportType.DofusPortalsID] = transportType
	}

	return &Impl{
		transportTypes:    transportTypes,
		transportTypeRepo: transportTypeRepo,
	}, nil
}

func (service *Impl) FindTransportTypeByDofusPortalsID(dofusPortalsID string) (entities.TransportType, bool) {
	transportType, found := service.transportTypes[dofusPortalsID]
	return transportType, found
}

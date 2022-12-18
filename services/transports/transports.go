package transports

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/repositories/transports"
)

type TransportService interface {
	FindTransportTypeByDofusPortalsId(dofusPortalsId string) (entities.TransportType, bool)
}

type TransportServiceImpl struct {
	transportTypes    map[string]entities.TransportType
	transportTypeRepo transports.TransportTypeRepository
}

func New(transportTypeRepo transports.TransportTypeRepository) (*TransportServiceImpl, error) {
	transportTypeEntities, err := transportTypeRepo.GetTransportTypes()
	if err != nil {
		return nil, err
	}

	transportTypes := make(map[string]entities.TransportType)
	for _, transportType := range transportTypeEntities {
		transportTypes[transportType.DofusPortalsId] = transportType
	}

	return &TransportServiceImpl{
		transportTypes:    transportTypes,
		transportTypeRepo: transportTypeRepo,
	}, nil
}

func (service *TransportServiceImpl) FindTransportTypeByDofusPortalsId(dofusPortalsId string) (entities.TransportType, bool) {
	transportType, found := service.transportTypes[dofusPortalsId]
	return transportType, found
}

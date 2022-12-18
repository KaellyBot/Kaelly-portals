package subareas

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/repositories/subareas"
)

type SubAreaService interface {
	FindSubAreaByDofusPortalsId(dofusPortalsId string) (entities.SubArea, bool)
}

type SubAreaServiceImpl struct {
	subAreas    map[string]entities.SubArea
	subAreaRepo subareas.SubAreaRepository
}

func New(subAreaRepo subareas.SubAreaRepository) (*SubAreaServiceImpl, error) {
	subAreaEntities, err := subAreaRepo.GetSubAreas()
	if err != nil {
		return nil, err
	}

	subAreas := make(map[string]entities.SubArea)
	for _, subArea := range subAreaEntities {
		subAreas[subArea.DofusPortalsId] = subArea
	}

	return &SubAreaServiceImpl{
		subAreas:    subAreas,
		subAreaRepo: subAreaRepo,
	}, nil
}

func (service *SubAreaServiceImpl) FindSubAreaByDofusPortalsId(dofusPortalsId string) (entities.SubArea, bool) {
	subArea, found := service.subAreas[dofusPortalsId]
	return subArea, found
}

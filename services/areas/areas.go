package areas

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/repositories/areas"
)

type AreaService interface {
	FindAreaByDofusPortalsId(dofusPortalsId string) (entities.Area, bool)
}

type AreaServiceImpl struct {
	areas    map[string]entities.Area
	areaRepo areas.AreaRepository
}

func New(areaRepo areas.AreaRepository) (*AreaServiceImpl, error) {
	areaEntities, err := areaRepo.GetAreas()
	if err != nil {
		return nil, err
	}

	areas := make(map[string]entities.Area)
	for _, area := range areaEntities {
		areas[area.DofusPortalsId] = area
	}

	return &AreaServiceImpl{
		areas:    areas,
		areaRepo: areaRepo,
	}, nil
}

func (service *AreaServiceImpl) FindAreaByDofusPortalsId(dofusPortalsId string) (entities.Area, bool) {
	area, found := service.areas[dofusPortalsId]
	return area, found
}

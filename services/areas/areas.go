package areas

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/repositories/areas"
)

func New(areaRepo areas.Repository) (*Impl, error) {
	areaEntities, err := areaRepo.GetAreas()
	if err != nil {
		return nil, err
	}

	areas := make(map[string]entities.Area)
	for _, area := range areaEntities {
		areas[area.DofusPortalsID] = area
	}

	return &Impl{
		areas:    areas,
		areaRepo: areaRepo,
	}, nil
}

func (service *Impl) FindAreaByDofusPortalsID(dofusPortalsID string) (entities.Area, bool) {
	area, found := service.areas[dofusPortalsID]
	return area, found
}

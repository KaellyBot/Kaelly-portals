package subareas

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/repositories/subareas"
)

func New(subAreaRepo subareas.Repository) (*Impl, error) {
	subAreaEntities, err := subAreaRepo.GetSubAreas()
	if err != nil {
		return nil, err
	}

	subAreas := make(map[string]entities.SubArea)
	for _, subArea := range subAreaEntities {
		subAreas[subArea.DofusPortalsID] = subArea
	}

	return &Impl{
		subAreas:    subAreas,
		subAreaRepo: subAreaRepo,
	}, nil
}

func (service *Impl) FindSubAreaByDofusPortalsID(dofusPortalsID string) (entities.SubArea, bool) {
	subArea, found := service.subAreas[dofusPortalsID]
	return subArea, found
}

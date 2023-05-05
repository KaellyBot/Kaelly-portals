package dimensions

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/repositories/dimensions"
)

func New(dimensionRepo dimensions.Repository) (*Impl, error) {
	dimEntities, err := dimensionRepo.GetDimensions()
	if err != nil {
		return nil, err
	}

	dimensions := make(map[string]entities.Dimension)
	dofusPortalsDimensions := make(map[string]entities.Dimension)
	for _, dimension := range dimEntities {
		dimensions[dimension.ID] = dimension
		dofusPortalsDimensions[dimension.DofusPortalsID] = dimension
	}

	return &Impl{
		dimensions:             dimensions,
		dofusPortalsDimensions: dofusPortalsDimensions,
		dimensionRepo:          dimensionRepo,
	}, nil
}

func (service *Impl) GetDimension(id string) (entities.Dimension, bool) {
	dimension, found := service.dimensions[id]
	return dimension, found
}

func (service *Impl) FindDimensionByDofusPortalsID(dofusPortalsID string) (entities.Dimension, bool) {
	dimension, found := service.dofusPortalsDimensions[dofusPortalsID]
	return dimension, found
}

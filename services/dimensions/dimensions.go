package dimensions

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/repositories/dimensions"
)

type DimensionService interface {
	GetDimension(id string) (entities.Dimension, bool)
	FindDimensionByDofusPortalsId(dofusPortalsId string) (entities.Dimension, bool)
}

type DimensionServiceImpl struct {
	dimensions             map[string]*entities.Dimension
	dofusPortalsDimensions map[string]*entities.Dimension
	dimensionRepo          dimensions.DimensionRepository
}

func New(dimensionRepo dimensions.DimensionRepository) (*DimensionServiceImpl, error) {
	dimEntities, err := dimensionRepo.GetDimensions()
	if err != nil {
		return nil, err
	}

	dimensions := make(map[string]*entities.Dimension)
	dofusPortalsDimensions := make(map[string]*entities.Dimension)
	for _, dimension := range dimEntities {
		dimensions[dimension.Id] = &dimension
		dofusPortalsDimensions[dimension.DofusPortalsId] = &dimension
	}

	return &DimensionServiceImpl{
		dimensions:             dimensions,
		dofusPortalsDimensions: dofusPortalsDimensions,
		dimensionRepo:          dimensionRepo,
	}, nil
}

func (service *DimensionServiceImpl) GetDimension(id string) (entities.Dimension, bool) {
	dimension, found := service.dimensions[id]
	return *dimension, found
}

func (service *DimensionServiceImpl) FindDimensionByDofusPortalsId(dofusPortalsId string) (entities.Dimension, bool) {
	dimension, found := service.dofusPortalsDimensions[dofusPortalsId]
	return *dimension, found
}

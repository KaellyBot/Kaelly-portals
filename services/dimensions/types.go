package dimensions

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/repositories/dimensions"
)

type Service interface {
	GetDimension(id string) (entities.Dimension, bool)
	FindDimensionByDofusPortalsID(dofusPortalsID string) (entities.Dimension, bool)
}

type Impl struct {
	dimensions             map[string]entities.Dimension
	dofusPortalsDimensions map[string]entities.Dimension
	dimensionRepo          dimensions.Repository
}

package areas

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/repositories/areas"
)

type Service interface {
	FindAreaByDofusPortalsID(dofusPortalsID string) (entities.Area, bool)
}

type Impl struct {
	areas    map[string]entities.Area
	areaRepo areas.Repository
}

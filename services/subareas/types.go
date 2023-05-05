package subareas

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/repositories/subareas"
)

type Service interface {
	FindSubAreaByDofusPortalsID(dofusPortalsId string) (entities.SubArea, bool)
}

type Impl struct {
	subAreas    map[string]entities.SubArea
	subAreaRepo subareas.Repository
}

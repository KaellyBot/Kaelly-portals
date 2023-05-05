package servers

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/repositories/servers"
)

type Service interface {
	GetServer(id string) (entities.Server, bool)
	FindServerByDofusPortalsID(dofusPortalsID string) (entities.Server, bool)
}

type Impl struct {
	servers             map[string]entities.Server
	dofusPortalsServers map[string]entities.Server
	serverRepo          servers.Repository
}

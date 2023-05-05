package servers

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/repositories/servers"
)

func New(serverRepo servers.Repository) (*Impl, error) {
	serverEntities, err := serverRepo.GetServers()
	if err != nil {
		return nil, err
	}

	servers := make(map[string]entities.Server)
	dofusPortalsServers := make(map[string]entities.Server)
	for _, server := range serverEntities {

		servers[server.ID] = server
		dofusPortalsServers[server.DofusPortalsID] = server
	}

	return &Impl{
		servers:             servers,
		dofusPortalsServers: dofusPortalsServers,
		serverRepo:          serverRepo,
	}, nil
}

func (service *Impl) GetServer(id string) (entities.Server, bool) {
	server, found := service.servers[id]
	return server, found
}

func (service *Impl) FindServerByDofusPortalsID(dofusPortalsID string) (entities.Server, bool) {
	server, found := service.dofusPortalsServers[dofusPortalsID]
	return server, found
}

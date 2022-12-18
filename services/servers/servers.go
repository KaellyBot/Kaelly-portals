package servers

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/repositories/servers"
)

type ServerService interface {
	GetServer(id string) (entities.Server, bool)
}

type ServerServiceImpl struct {
	servers    map[string]entities.Server
	serverRepo servers.ServerRepository
}

func New(serverRepo servers.ServerRepository) (*ServerServiceImpl, error) {
	serverEntities, err := serverRepo.GetServers()
	if err != nil {
		return nil, err
	}

	servers := make(map[string]entities.Server)
	for _, server := range serverEntities {

		servers[server.Id] = server
	}

	return &ServerServiceImpl{
		servers:    servers,
		serverRepo: serverRepo,
	}, nil
}

func (service *ServerServiceImpl) GetServer(id string) (entities.Server, bool) {
	server, found := service.servers[id]
	return server, found
}

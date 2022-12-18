package servers

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/utils/databases"
)

type ServerRepository interface {
	GetServers() ([]entities.Server, error)
}

type ServerRepositoryImpl struct {
	db databases.MySQLConnection
}

func New(db databases.MySQLConnection) *ServerRepositoryImpl {
	return &ServerRepositoryImpl{db: db}
}

func (repo *ServerRepositoryImpl) GetServers() ([]entities.Server, error) {
	var servers []entities.Server
	response := repo.db.GetDB().Model(&entities.Server{}).Find(&servers)
	return servers, response.Error
}

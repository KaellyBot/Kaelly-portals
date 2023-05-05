package servers

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/utils/databases"
)

func New(db databases.MySQLConnection) *Impl {
	return &Impl{db: db}
}

func (repo *Impl) GetServers() ([]entities.Server, error) {
	var servers []entities.Server
	response := repo.db.GetDB().Model(&entities.Server{}).Find(&servers)
	return servers, response.Error
}

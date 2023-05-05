package servers

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/utils/databases"
)

type Repository interface {
	GetServers() ([]entities.Server, error)
}

type Impl struct {
	db databases.MySQLConnection
}

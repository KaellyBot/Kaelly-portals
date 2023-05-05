package areas

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/utils/databases"
)

type Repository interface {
	GetAreas() ([]entities.Area, error)
}

type Impl struct {
	db databases.MySQLConnection
}

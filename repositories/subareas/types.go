package subareas

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/utils/databases"
)

type Repository interface {
	GetSubAreas() ([]entities.SubArea, error)
}

type Impl struct {
	db databases.MySQLConnection
}

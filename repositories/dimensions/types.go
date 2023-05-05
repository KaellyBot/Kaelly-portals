package dimensions

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/utils/databases"
)

type Repository interface {
	GetDimensions() ([]entities.Dimension, error)
}

type Impl struct {
	db databases.MySQLConnection
}

package dimensions

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/utils/databases"
)

func New(db databases.MySQLConnection) *Impl {
	return &Impl{db: db}
}

func (repo *Impl) GetDimensions() ([]entities.Dimension, error) {
	var dimensions []entities.Dimension
	response := repo.db.GetDB().Model(&entities.Dimension{}).Find(&dimensions)
	return dimensions, response.Error
}

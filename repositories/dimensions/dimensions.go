package dimensions

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/utils/databases"
)

type DimensionRepository interface {
	GetDimensions() ([]entities.Dimension, error)
}

type DimensionRepositoryImpl struct {
	db databases.MySQLConnection
}

func New(db databases.MySQLConnection) *DimensionRepositoryImpl {
	return &DimensionRepositoryImpl{db: db}
}

func (repo *DimensionRepositoryImpl) GetDimensions() ([]entities.Dimension, error) {
	var dimensions []entities.Dimension
	response := repo.db.GetDB().Model(&entities.Dimension{}).Find(&dimensions)
	return dimensions, response.Error
}

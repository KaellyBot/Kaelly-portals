package subareas

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/utils/databases"
)

type SubAreaRepository interface {
	GetSubAreas() ([]entities.SubArea, error)
}

type SubAreaRepositoryImpl struct {
	db databases.MySQLConnection
}

func New(db databases.MySQLConnection) *SubAreaRepositoryImpl {
	return &SubAreaRepositoryImpl{db: db}
}

func (repo *SubAreaRepositoryImpl) GetSubAreas() ([]entities.SubArea, error) {
	var SubAreas []entities.SubArea
	response := repo.db.GetDB().Model(&entities.SubArea{}).Find(&SubAreas)
	return SubAreas, response.Error
}

package subareas

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/utils/databases"
)

func New(db databases.MySQLConnection) *Impl {
	return &Impl{db: db}
}

func (repo *Impl) GetSubAreas() ([]entities.SubArea, error) {
	var SubAreas []entities.SubArea
	response := repo.db.GetDB().Model(&entities.SubArea{}).Find(&SubAreas)
	return SubAreas, response.Error
}

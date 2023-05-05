package subareas

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/utils/databases"
)

func New(db databases.MySQLConnection) *Impl {
	return &Impl{db: db}
}

func (repo *Impl) GetSubAreas() ([]entities.SubArea, error) {
	var subAreas []entities.SubArea
	response := repo.db.GetDB().Model(&entities.SubArea{}).Find(&subAreas)
	return subAreas, response.Error
}

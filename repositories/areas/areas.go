package areas

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/utils/databases"
)

func New(db databases.MySQLConnection) *Impl {
	return &Impl{db: db}
}

func (repo *Impl) GetAreas() ([]entities.Area, error) {
	var areas []entities.Area
	response := repo.db.GetDB().Model(&entities.Area{}).Find(&areas)
	return areas, response.Error
}

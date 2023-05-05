package transports

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/utils/databases"
)

func New(db databases.MySQLConnection) *Impl {
	return &Impl{db: db}
}

func (repo *Impl) GetTransportTypes() ([]entities.TransportType, error) {
	var TransportTypes []entities.TransportType
	response := repo.db.GetDB().Model(&entities.TransportType{}).Find(&TransportTypes)
	return TransportTypes, response.Error
}

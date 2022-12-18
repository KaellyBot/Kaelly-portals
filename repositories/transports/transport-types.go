package transports

import (
	"github.com/kaellybot/kaelly-portals/models/entities"
	"github.com/kaellybot/kaelly-portals/utils/databases"
)

type TransportTypeRepository interface {
	GetTransportTypes() ([]entities.TransportType, error)
}

type TransportTypeRepositoryImpl struct {
	db databases.MySQLConnection
}

func New(db databases.MySQLConnection) *TransportTypeRepositoryImpl {
	return &TransportTypeRepositoryImpl{db: db}
}

func (repo *TransportTypeRepositoryImpl) GetTransportTypes() ([]entities.TransportType, error) {
	var TransportTypes []entities.TransportType
	response := repo.db.GetDB().Model(&entities.TransportType{}).Find(&TransportTypes)
	return TransportTypes, response.Error
}

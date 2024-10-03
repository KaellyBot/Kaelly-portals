package databases

import (
	"fmt"

	"github.com/kaellybot/kaelly-portals/models/constants"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLConnection interface {
	GetDB() *gorm.DB
	Shutdown()
}

type MySQLConnectionImpl struct {
	db *gorm.DB
}

func New() (*MySQLConnectionImpl, error) {
	dbUser := viper.GetString(constants.MySQLUser)
	dbPassword := viper.GetString(constants.MySQLPassword)
	dbURL := viper.GetString(constants.MySQLURL)
	dbName := viper.GetString(constants.MySQLDatabase)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=UTC", dbUser, dbPassword, dbURL, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &MySQLConnectionImpl{db: db}, nil
}

func (connection *MySQLConnectionImpl) GetDB() *gorm.DB {
	return connection.db
}

func (connection *MySQLConnectionImpl) Shutdown() {
	dbSQL, err := connection.db.DB()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to kill connection from database")
		return
	}
	dbSQL.Close()
}

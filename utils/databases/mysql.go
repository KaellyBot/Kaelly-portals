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
	dbUser := viper.GetString(constants.MySqlUser)
	dbPassword := viper.GetString(constants.MySqlPassword)
	dbUrl := viper.GetString(constants.MySqlUrl)
	dbName := viper.GetString(constants.MySqlDatabase)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local", dbUser, dbPassword, dbUrl, dbName)

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

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
	IsConnected() bool
	Run() error
	Shutdown()
}

type mySQLConnection struct {
	dsn string
	db  *gorm.DB
}

func New() MySQLConnection {
	dbUser := viper.GetString(constants.MySQLUser)
	dbPassword := viper.GetString(constants.MySQLPassword)
	dbURL := viper.GetString(constants.MySQLURL)
	dbName := viper.GetString(constants.MySQLDatabase)
	return &mySQLConnection{
		dsn: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=UTC",
			dbUser, dbPassword, dbURL, dbName),
	}
}

func (c *mySQLConnection) GetDB() *gorm.DB {
	return c.db
}

func (c *mySQLConnection) IsConnected() bool {
	if c.db == nil {
		return false
	}

	dbSQL, errSQL := c.db.DB()
	if errSQL != nil {
		return false
	}

	if errPing := dbSQL.Ping(); errPing != nil {
		return false
	}

	return true
}

func (c *mySQLConnection) Run() error {
	db, err := gorm.Open(mysql.Open(c.dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	c.db = db
	log.Info().Msg("Connected to MySQL")
	return nil
}

func (c *mySQLConnection) Shutdown() {
	log.Info().Msg("Shutdown the connection to MySQL")
	dbSQL, err := c.db.DB()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to shutdown database connection")
		return
	}

	if errClose := dbSQL.Close(); errClose != nil {
		log.Error().Err(errClose).Msgf("Failed to shutdown database connection")
	}
}

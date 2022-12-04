package application

import (
	"errors"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-portals/models"
	"github.com/kaellybot/kaelly-portals/services/portals"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	ErrCannotInstantiateApp = errors.New("Cannot instantiate application")
)

type ApplicationInterface interface {
	Run() error
	Shutdown()
}

type Application struct {
	portals portals.PortalsServiceInterface
	broker  amqp.MessageBrokerInterface
}

func New(rabbitMqClientId, rabbitMqAddress string, httpTimeout int) (*Application, error) {
	broker, err := amqp.New(rabbitMqClientId, rabbitMqAddress, []amqp.Binding{portals.GetBinding()})
	if err != nil {
		log.Error().Err(err).Msgf("Failed to instantiate broker")
		return nil, ErrCannotInstantiateApp
	}

	portals, err := portals.New(broker, viper.GetString(models.DofusPortalsToken), httpTimeout)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to instantiate portals service")
		return nil, ErrCannotInstantiateApp
	}

	return &Application{
		portals: portals,
		broker:  broker,
	}, nil
}

func (app *Application) Run() error {
	return app.portals.Consume()
}

func (app *Application) Shutdown() {
	app.broker.Shutdown()
	log.Info().Msgf("Application is no longer running")
}

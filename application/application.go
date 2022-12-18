package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-portals/models/constants"
	areaRepo "github.com/kaellybot/kaelly-portals/repositories/areas"
	dimensionRepo "github.com/kaellybot/kaelly-portals/repositories/dimensions"
	serverRepo "github.com/kaellybot/kaelly-portals/repositories/servers"
	subAreaRepo "github.com/kaellybot/kaelly-portals/repositories/subareas"
	transportRepo "github.com/kaellybot/kaelly-portals/repositories/transports"
	"github.com/kaellybot/kaelly-portals/services/areas"
	"github.com/kaellybot/kaelly-portals/services/dimensions"
	"github.com/kaellybot/kaelly-portals/services/portals"
	"github.com/kaellybot/kaelly-portals/services/servers"
	"github.com/kaellybot/kaelly-portals/services/subareas"
	"github.com/kaellybot/kaelly-portals/services/transports"
	"github.com/kaellybot/kaelly-portals/utils/databases"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type ApplicationInterface interface {
	Run() error
	Shutdown()
}

type Application struct {
	serverService    servers.ServerService
	dimensionService dimensions.DimensionService
	areaService      areas.AreaService
	subAreaService   subareas.SubAreaService
	transportService transports.TransportService
	portals          portals.PortalsService
	broker           amqp.MessageBrokerInterface
}

func New() (*Application, error) {
	// misc
	db, err := databases.New()
	if err != nil {
		return nil, err
	}

	broker, err := amqp.New(constants.RabbitMQClientId, viper.GetString(constants.RabbitMqAddress), []amqp.Binding{portals.GetBinding()})
	if err != nil {
		return nil, err
	}

	// repositories
	serverRepo := serverRepo.New(db)
	dimensionRepo := dimensionRepo.New(db)
	areaRepo := areaRepo.New(db)
	subAreaRepo := subAreaRepo.New(db)
	transportRepo := transportRepo.New(db)

	// services
	serverService, err := servers.New(serverRepo)
	if err != nil {
		return nil, err
	}

	dimensionService, err := dimensions.New(dimensionRepo)
	if err != nil {
		return nil, err
	}

	areaService, err := areas.New(areaRepo)
	if err != nil {
		return nil, err
	}

	subAreaService, err := subareas.New(subAreaRepo)
	if err != nil {
		return nil, err
	}

	transportService, err := transports.New(transportRepo)
	if err != nil {
		return nil, err
	}

	portals, err := portals.New(broker, serverService, dimensionService, areaService, subAreaService, transportService)
	if err != nil {
		return nil, err
	}

	return &Application{
		serverService:    serverService,
		dimensionService: dimensionService,
		areaService:      areaService,
		subAreaService:   subAreaService,
		transportService: transportService,
		portals:          portals,
		broker:           broker,
	}, nil
}

func (app *Application) Run() error {
	return app.portals.Consume()
}

func (app *Application) Shutdown() {
	app.broker.Shutdown()
	log.Info().Msgf("Application is no longer running")
}

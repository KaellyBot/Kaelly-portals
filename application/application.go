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
	"github.com/kaellybot/kaelly-portals/utils/insights"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func New() (*Impl, error) {
	// misc
	broker := amqp.New(constants.RabbitMQClientID, viper.GetString(constants.RabbitMQAddress),
		amqp.WithBindings(portals.GetBinding()))
	db := databases.New()
	if err := db.Run(); err != nil {
		return nil, err
	}

	probes := insights.NewProbes(broker.IsConnected, db.IsConnected)
	prom := insights.NewPrometheusMetrics()

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

	portals, err := portals.New(broker, serverService, dimensionService,
		areaService, subAreaService, transportService)
	if err != nil {
		return nil, err
	}

	return &Impl{
		portals: portals,
		broker:  broker,
		db:      db,
		probes:  probes,
		prom:    prom,
	}, nil
}

func (app *Impl) Run() error {
	app.probes.ListenAndServe()
	app.prom.ListenAndServe()

	if err := app.broker.Run(); err != nil {
		return err
	}

	app.portals.Consume()
	return nil
}

func (app *Impl) Shutdown() {
	app.broker.Shutdown()
	app.db.Shutdown()
	app.prom.Shutdown()
	app.probes.Shutdown()
	log.Info().Msgf("Application is no longer running")
}

package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-portals/services/areas"
	"github.com/kaellybot/kaelly-portals/services/dimensions"
	"github.com/kaellybot/kaelly-portals/services/portals"
	"github.com/kaellybot/kaelly-portals/services/servers"
	"github.com/kaellybot/kaelly-portals/services/subareas"
	"github.com/kaellybot/kaelly-portals/services/transports"
)

type Application interface {
	Run() error
	Shutdown()
}

type Impl struct {
	serverService    servers.Service
	dimensionService dimensions.Service
	areaService      areas.Service
	subAreaService   subareas.Service
	transportService transports.Service
	portals          portals.Service
	broker           amqp.MessageBroker
}

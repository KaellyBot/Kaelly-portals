package portals

import (
	"errors"
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-portals/payloads/dofusportals"
	"github.com/kaellybot/kaelly-portals/services/areas"
	"github.com/kaellybot/kaelly-portals/services/dimensions"
	"github.com/kaellybot/kaelly-portals/services/servers"
	"github.com/kaellybot/kaelly-portals/services/subareas"
	"github.com/kaellybot/kaelly-portals/services/transports"
)

const (
	requestQueueName   = "portals-requests"
	requestsRoutingkey = "requests.portals"
	answersRoutingkey  = "answers.portals"

	httpHeader   = "header"
	httpAPIToken = "token"
)

var (
	errInvalidMessage = errors.New("invalid request portal, type is not the good one" +
		" and/or the dedicated message is not filled")
	errStatusNotOK = errors.New("status Code is not OK")
)

type Service interface {
	Consume()
}

type Impl struct {
	dofusPortalsClient dofusportals.ClientInterface
	broker             amqp.MessageBroker
	httpTimeout        time.Duration
	serverService      servers.Service
	dimensionService   dimensions.Service
	areaService        areas.Service
	subAreaService     subareas.Service
	transportService   transports.Service
}

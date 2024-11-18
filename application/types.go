package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-portals/services/portals"
	"github.com/kaellybot/kaelly-portals/utils/databases"
	"github.com/kaellybot/kaelly-portals/utils/insights"
)

type Application interface {
	Run() error
	Shutdown()
}

type Impl struct {
	portals portals.Service
	broker  amqp.MessageBroker
	db      databases.MySQLConnection
	probes  insights.Probes
	prom    insights.PrometheusMetrics
}

package insights

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kaellybot/kaelly-portals/models/constants"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type PrometheusMetrics interface {
	ListenAndServe()
	Shutdown()
}

type prom struct {
	server *http.Server
}

func NewPrometheusMetrics() PrometheusMetrics {
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	return &prom{
		server: &http.Server{
			Addr:              fmt.Sprintf(":%v", viper.GetInt(constants.MetricPort)),
			Handler:           metricsMux,
			ReadHeaderTimeout: 0,
		},
	}
}

func (prom *prom) ListenAndServe() {
	go func() {
		log.Info().Msgf("Exposing Prometheus metrics...")
		err := prom.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msgf("Cannot listen and serve Prometheus metrics")
		}
	}()
}

func (prom *prom) Shutdown() {
	if prom.server != nil {
		if err := prom.server.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msgf("Failed to shutdown Prometheus metrics server")
		}
	}
}

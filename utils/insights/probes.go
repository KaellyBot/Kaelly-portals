package insights

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kaellybot/kaelly-portals/models/constants"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Probes interface {
	ListenAndServe()
	Shutdown()
}

type probes struct {
	server       *http.Server
	isReadyFuncs []IsReadyFunc
}

type IsReadyFunc func() bool

func NewProbes(isReadyFuncs ...IsReadyFunc) Probes {
	impl := probes{
		isReadyFuncs: isReadyFuncs,
	}
	probesMux := http.NewServeMux()
	probesMux.HandleFunc("/live", impl.live)
	probesMux.HandleFunc("/ready", impl.ready)

	impl.server = &http.Server{
		Addr:              fmt.Sprintf(":%v", viper.GetInt(constants.ProbePort)),
		Handler:           probesMux,
		ReadHeaderTimeout: 0,
	}

	return &impl
}

func (probes *probes) ListenAndServe() {
	go func() {
		log.Info().Msgf("Exposing Probes...")
		err := probes.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msgf("Cannot listen and serve probes")
		}
	}()
}

func (probes *probes) Shutdown() {
	if probes.server != nil {
		if err := probes.server.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msgf("Failed to shutdown probe server")
		}
	}
}

func (probes *probes) live(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (probes *probes) ready(w http.ResponseWriter, _ *http.Request) {
	isReady := true

	for _, isReadyFunc := range probes.isReadyFuncs {
		isReady = isReady && checkReadiness(isReadyFunc)
	}

	if isReady {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

//nolint:nonamedreturns // Can't avoid it, unfortunately. It is much way safer like that.
func checkReadiness(isReadyFunc IsReadyFunc) (result bool) {
	defer func() {
		err := recover()
		if err != nil {
			log.Error().Err(fmt.Errorf("%v", err)).
				Msgf("Crash while retrieving readiness, considered now as unhealthy")
			result = false
		}
	}()

	return isReadyFunc()
}

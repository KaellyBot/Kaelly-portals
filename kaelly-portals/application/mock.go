package application

import "github.com/rs/zerolog/log"

type ApplicationMock struct {
	RunFunc      func() error
	ShutdownFunc func()
}

func (mock *ApplicationMock) Run() error {
	if mock.RunFunc != nil {
		return mock.RunFunc()
	}

	log.Warn().Msgf("No mock provided")
	return nil
}

func (mock *ApplicationMock) Shutdown() {
	if mock.ShutdownFunc != nil {
		mock.ShutdownFunc()
		return
	}

	log.Warn().Msgf("No mock provided")
}

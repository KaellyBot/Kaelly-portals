package portals

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-portals/models/constants"
	"github.com/kaellybot/kaelly-portals/models/mappers"
	"github.com/kaellybot/kaelly-portals/payloads/dofusportals"
	"github.com/kaellybot/kaelly-portals/services/areas"
	"github.com/kaellybot/kaelly-portals/services/dimensions"
	"github.com/kaellybot/kaelly-portals/services/servers"
	"github.com/kaellybot/kaelly-portals/services/subareas"
	"github.com/kaellybot/kaelly-portals/services/transports"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	requestQueueName   = "portals-requests"
	requestsRoutingkey = "requests.portals"
	answersRoutingkey  = "answers.portals"
)

var (
	errInvalidMessage = errors.New("Invalid request portal, type is not the good one and/or the dedicated message is not filled")
	errStatusNotOK    = errors.New("Status Code is not OK")
)

type PortalsService interface {
	Consume() error
}

type PortalsServiceImpl struct {
	dofusPortalsClient dofusportals.ClientInterface
	broker             amqp.MessageBrokerInterface
	httpTimeout        time.Duration
	serverService      servers.ServerService
	dimensionService   dimensions.DimensionService
	areaService        areas.AreaService
	subAreaService     subareas.SubAreaService
	transportService   transports.TransportService
}

func New(broker amqp.MessageBrokerInterface, serverService servers.ServerService,
	dimensionService dimensions.DimensionService, areaService areas.AreaService,
	subAreaService subareas.SubAreaService, transportService transports.TransportService) (*PortalsServiceImpl, error) {

	apiKeyProvider, err := securityprovider.NewSecurityProviderApiKey("header", "token", viper.GetString(constants.DofusPortalsToken))
	if err != nil {
		return nil, err
	}

	dofusPortalsClient, err := dofusportals.NewClient(
		constants.DofusPortalsUrl,
		dofusportals.WithRequestEditorFn(apiKeyProvider.Intercept),
	)
	if err != nil {
		return nil, err
	}

	return &PortalsServiceImpl{
		serverService:      serverService,
		dimensionService:   dimensionService,
		areaService:        areaService,
		subAreaService:     subAreaService,
		transportService:   transportService,
		broker:             broker,
		dofusPortalsClient: dofusPortalsClient,
		httpTimeout:        time.Duration(viper.GetInt(constants.HttpTimeout)) * time.Second,
	}, nil
}

func GetBinding() amqp.Binding {
	return amqp.Binding{
		Exchange:   amqp.ExchangeRequest,
		RoutingKey: requestsRoutingkey,
		Queue:      requestQueueName,
	}
}

func (service *PortalsServiceImpl) Consume() error {
	log.Info().Msgf("Consuming portal requests...")
	return service.broker.Consume(requestQueueName, requestsRoutingkey, service.consume)
}

func (service *PortalsServiceImpl) consume(ctx context.Context, message *amqp.RabbitMQMessage, correlationId string) {
	if !isValidPortalRequest(message) {
		log.Error().Err(errInvalidMessage).Str(constants.LogCorrelationId, correlationId).Msgf("Returning failed message")
		service.publishPortalAnswerFailed(correlationId, message.Language)
		return
	}

	// TODO map server/dimension into dofus portals ids
	server := message.GetPortalPositionRequest().GetServer()
	dimension := message.GetPortalPositionRequest().GetDimension()

	log.Info().
		Str(constants.LogCorrelationId, correlationId).
		Str(constants.LogServerId, server).
		Str(constants.LogDimensionId, dimension).
		Msgf("Treating request")

	portals := make([]*amqp.PortalPositionAnswer_PortalPosition, 0)
	if dimension != "" {
		dofusPortal, err := service.getPortal(ctx, server, dimension)
		if err != nil {
			log.Error().Err(err).
				Str(constants.LogCorrelationId, correlationId).
				Str(constants.LogServerId, server).
				Str(constants.LogDimensionId, dimension).
				Msgf("Returning failed message")
			service.publishPortalAnswerFailed(correlationId, message.Language)
			return
		}

		// TODO map dofus-portals ids in kaelly ones
		portals = append(portals, mappers.MapPortal(dofusPortal))

	} else {
		dofusPortals, err := service.getPortals(ctx, server)
		if err != nil {
			log.Error().Err(err).
				Str(constants.LogCorrelationId, correlationId).
				Str(constants.LogServerId, server).
				Msgf("Returning failed message")
			service.publishPortalAnswerFailed(correlationId, message.Language)
			return
		}

		for _, dofusPortal := range dofusPortals {
			// TODO map dofus-portals ids in kaelly ones
			portals = append(portals, mappers.MapPortal(dofusPortal))
		}
	}

	service.publishPortalAnswerSuccess(portals, correlationId, message.Language)
}

func isValidPortalRequest(message *amqp.RabbitMQMessage) bool {
	return message.Type == amqp.RabbitMQMessage_PORTAL_POSITION_REQUEST && message.GetPortalPositionRequest() != nil
}

func (service *PortalsServiceImpl) getPortals(ctx context.Context, server string) ([]dofusportals.Portal, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()
	resp, err := service.dofusPortalsClient.GetExternalV1ServersServerIdPortals(ctx, server)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errStatusNotOK
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var portals []dofusportals.Portal
	if err := json.Unmarshal(body, &portals); err != nil {
		return nil, err
	}

	return portals, nil
}

func (service *PortalsServiceImpl) getPortal(ctx context.Context, server, dimension string) (dofusportals.Portal, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()
	resp, err := service.dofusPortalsClient.GetExternalV1ServersServerIdPortalsDimensionId(ctx, server, dimension)
	if err != nil {
		return dofusportals.Portal{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dofusportals.Portal{}, errStatusNotOK
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return dofusportals.Portal{}, err
	}

	var portal dofusportals.Portal
	if err := json.Unmarshal(body, &portal); err != nil {
		return dofusportals.Portal{}, err
	}

	return portal, nil
}

func (service *PortalsServiceImpl) publishPortalAnswerFailed(correlationId string, language amqp.RabbitMQMessage_Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_PORTAL_POSITION_ANSWER,
		Status:   amqp.RabbitMQMessage_FAILED,
		Language: language,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationId)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationId, correlationId).Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *PortalsServiceImpl) publishPortalAnswerSuccess(portals []*amqp.PortalPositionAnswer_PortalPosition,
	correlationId string, language amqp.RabbitMQMessage_Language) {

	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_PORTAL_POSITION_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: language,
		PortalPositionAnswer: &amqp.PortalPositionAnswer{
			Positions: portals,
		},
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationId)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationId, correlationId).Msgf("Cannot publish via broker, request ignored")
	}
}

package portals

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-portals/models/constants"
	"github.com/kaellybot/kaelly-portals/models/mappers"
	"github.com/kaellybot/kaelly-portals/payloads/dofusportals"
	"github.com/kaellybot/kaelly-portals/services/areas"
	"github.com/kaellybot/kaelly-portals/services/dimensions"
	"github.com/kaellybot/kaelly-portals/services/servers"
	"github.com/kaellybot/kaelly-portals/services/subareas"
	"github.com/kaellybot/kaelly-portals/services/transports"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func New(broker amqp.MessageBroker, serverService servers.Service,
	dimensionService dimensions.Service, areaService areas.Service,
	subAreaService subareas.Service, transportService transports.Service) (*Impl, error) {
	apiKeyProvIDer, err := securityprovider.
		NewSecurityProviderApiKey(httpHeader, httpAPIToken, viper.GetString(constants.DofusPortalsToken))
	if err != nil {
		return nil, err
	}

	dofusPortalsClient, err := dofusportals.NewClient(
		constants.DofusPortalsURL,
		dofusportals.WithRequestEditorFn(apiKeyProvIDer.Intercept),
	)
	if err != nil {
		return nil, err
	}

	return &Impl{
		serverService:      serverService,
		dimensionService:   dimensionService,
		areaService:        areaService,
		subAreaService:     subAreaService,
		transportService:   transportService,
		broker:             broker,
		dofusPortalsClient: dofusPortalsClient,
		httpTimeout:        time.Duration(viper.GetInt(constants.DofusPortalsTimeout)) * time.Second,
	}, nil
}

func GetBinding() amqp.Binding {
	return amqp.Binding{
		Exchange:   amqp.ExchangeRequest,
		RoutingKey: requestsRoutingkey,
		Queue:      requestQueueName,
	}
}

func (service *Impl) Consume() error {
	log.Info().Msgf("Consuming portal requests...")
	return service.broker.Consume(requestQueueName, service.consume)
}

func (service *Impl) consume(ctx context.Context, message *amqp.RabbitMQMessage, correlationID string) {
	if !isValidPortalRequest(message) {
		log.Error().
			Err(errInvalidMessage).
			Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot treat request, returning failed message")
		service.publishPortalAnswerFailed(correlationID, message.Language)
		return
	}

	serverID := message.GetPortalPositionRequest().GetServerId()
	dimensionID := message.GetPortalPositionRequest().GetDimensionId()

	log.Info().
		Str(constants.LogCorrelationID, correlationID).
		Str(constants.LogServerID, serverID).
		Str(constants.LogDimensionID, dimensionID).
		Msgf("Treating request")

	dofusPortalsServerID := service.getDofusPortalsServerID(serverID)

	portals := make([]*amqp.PortalPositionAnswer_PortalPosition, 0)
	if dimensionID != "" {
		dofusPortalsDimensionID := service.getDofusPortalsDimensionID(dimensionID)
		dofusPortal, err := service.getPortal(ctx, dofusPortalsServerID, dofusPortalsDimensionID)
		if err != nil {
			log.Error().Err(err).
				Str(constants.LogCorrelationID, correlationID).
				Str(constants.LogServerID, serverID).
				Str(constants.LogDimensionID, dimensionID).
				Msgf("Returning failed message")
			service.publishPortalAnswerFailed(correlationID, message.Language)
			return
		}

		portals = append(portals, mappers.MapPortal(dofusPortal, service.serverService, service.dimensionService,
			service.areaService, service.subAreaService, service.transportService))
	} else {
		dofusPortals, err := service.getPortals(ctx, dofusPortalsServerID)
		if err != nil {
			log.Error().Err(err).
				Str(constants.LogCorrelationID, correlationID).
				Str(constants.LogServerID, serverID).
				Msgf("Returning failed message")
			service.publishPortalAnswerFailed(correlationID, message.Language)
			return
		}

		for _, dofusPortal := range dofusPortals {
			portals = append(portals, mappers.MapPortal(dofusPortal, service.serverService, service.dimensionService,
				service.areaService, service.subAreaService, service.transportService))
		}
	}

	service.publishPortalAnswerSuccess(portals, correlationID, message.Language)
}

func isValidPortalRequest(message *amqp.RabbitMQMessage) bool {
	return message.Type == amqp.RabbitMQMessage_PORTAL_POSITION_REQUEST && message.GetPortalPositionRequest() != nil
}

func (service *Impl) getDofusPortalsServerID(serverID string) string {
	server, found := service.serverService.GetServer(serverID)
	if !found {
		log.Warn().Str(constants.LogServerID, serverID).Msgf("Server not found, returning internal server ID")
		return serverID
	}

	return server.DofusPortalsID
}

func (service *Impl) getDofusPortalsDimensionID(dimensionID string) string {
	dimension, found := service.dimensionService.GetDimension(dimensionID)
	if !found {
		log.Warn().Str(constants.LogDimensionID, dimensionID).Msgf("Dimension not found, returning internal dimension ID")
		return dimensionID
	}

	return dimension.DofusPortalsID
}

func (service *Impl) getPortals(ctx context.Context, server string) ([]dofusportals.Portal, error) {
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var portals []dofusportals.Portal
	if err = json.Unmarshal(body, &portals); err != nil {
		return nil, err
	}

	return portals, nil
}

func (service *Impl) getPortal(ctx context.Context, server, dimension string) (dofusportals.Portal, error) {
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dofusportals.Portal{}, err
	}

	var portal dofusportals.Portal
	if err = json.Unmarshal(body, &portal); err != nil {
		return dofusportals.Portal{}, err
	}

	return portal, nil
}

func (service *Impl) publishPortalAnswerFailed(correlationID string, language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_PORTAL_POSITION_ANSWER,
		Status:   amqp.RabbitMQMessage_FAILED,
		Language: language,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *Impl) publishPortalAnswerSuccess(portals []*amqp.PortalPositionAnswer_PortalPosition,
	correlationID string, language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_PORTAL_POSITION_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: language,
		PortalPositionAnswer: &amqp.PortalPositionAnswer{
			Positions: portals,
		},
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).Msgf("Cannot publish via broker, request ignored")
	}
}

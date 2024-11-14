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
	"github.com/kaellybot/kaelly-portals/utils/databases/replies"
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

func (service *Impl) Consume() {
	log.Info().Msgf("Consuming portal requests...")
	service.broker.Consume(requestQueueName, service.consume)
}

func (service *Impl) consume(ctx amqp.Context, message *amqp.RabbitMQMessage) {
	if !isValidPortalRequest(message) {
		log.Error().
			Err(errInvalidMessage).
			Str(constants.LogCorrelationID, ctx.CorrelationID).
			Msgf("Cannot treat request, returning failed message")
		replies.FailedAnswer(ctx, service.broker, amqp.RabbitMQMessage_PORTAL_POSITION_ANSWER,
			message.Language)
		return
	}

	serverID := message.GetPortalPositionRequest().GetServerId()
	dimensionID := message.GetPortalPositionRequest().GetDimensionId()

	log.Info().
		Str(constants.LogCorrelationID, ctx.CorrelationID).
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
				Str(constants.LogCorrelationID, ctx.CorrelationID).
				Str(constants.LogServerID, serverID).
				Str(constants.LogDimensionID, dimensionID).
				Msgf("Returning failed message")
			replies.FailedAnswer(ctx, service.broker, amqp.RabbitMQMessage_PORTAL_POSITION_ANSWER,
				message.Language)
			return
		}

		portals = append(portals, mappers.MapPortal(dofusPortal, service.serverService, service.dimensionService,
			service.areaService, service.subAreaService, service.transportService))
	} else {
		dofusPortals, err := service.getPortals(ctx, dofusPortalsServerID)
		if err != nil {
			log.Error().Err(err).
				Str(constants.LogCorrelationID, ctx.CorrelationID).
				Str(constants.LogServerID, serverID).
				Msgf("Returning failed message")
			replies.FailedAnswer(ctx, service.broker, amqp.RabbitMQMessage_PORTAL_POSITION_ANSWER,
				message.Language)
			return
		}

		for _, dofusPortal := range dofusPortals {
			portals = append(portals, mappers.MapPortal(dofusPortal, service.serverService, service.dimensionService,
				service.areaService, service.subAreaService, service.transportService))
		}
	}

	response := mappers.MapPortalAnswer(portals, message.Language)
	replies.SucceededAnswer(ctx, service.broker, response)
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

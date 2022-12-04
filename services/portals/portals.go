package portals

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-portals/models"
	"github.com/kaellybot/kaelly-portals/payloads/dofusportals"
	"github.com/rs/zerolog/log"
)

const (
	requestQueueName   = models.RabbitMQClientId + "_requests"
	requestsRoutingkey = "requests.portals"
	answersRoutingkey  = "answers.portals"
)

var (
	errInvalidMessage = errors.New("Invalid request portal, type is not the good one and/or the dedicated message is not filled")
)

type PortalsServiceInterface interface {
	Consume() error
}

type PortalsService struct {
	dofusPortalsClient dofusportals.ClientInterface
	broker             amqp.MessageBrokerInterface
	httpTimeout        time.Duration
}

func New(broker amqp.MessageBrokerInterface, dofusPortalsToken string, httpTimeout int) (*PortalsService, error) {
	apiKeyProvider, err := securityprovider.NewSecurityProviderApiKey("header", "token", dofusPortalsToken)
	if err != nil {
		return nil, err
	}

	dofusPortalsClient, err := dofusportals.NewClient(
		models.DofusPortalsUrl,
		dofusportals.WithRequestEditorFn(apiKeyProvider.Intercept),
	)
	if err != nil {
		return nil, err
	}

	return &PortalsService{
		broker:             broker,
		dofusPortalsClient: dofusPortalsClient,
		httpTimeout:        time.Duration(httpTimeout) * time.Second,
	}, nil
}

func GetBinding() amqp.Binding {
	return amqp.Binding{
		Exchange:   amqp.ExchangeRequest,
		RoutingKey: requestsRoutingkey,
		Queue:      requestQueueName,
	}
}

func (service *PortalsService) Consume() error {
	log.Info().Msgf("Consuming portal requests...")
	return service.broker.Consume(requestQueueName, requestsRoutingkey, service.consume)
}

func (service *PortalsService) consume(ctx context.Context, message *amqp.RabbitMQMessage, correlationId string) {
	if !isValidPortalRequest(message) {
		log.Error().Err(errInvalidMessage).Str(models.LogCorrelationId, correlationId).Msgf("Returning failed message")
		service.publishPortalAnswerFailed(correlationId, message.Language)
		return
	}

	server := message.GetPortalPositionRequest().GetServer()
	dimension := message.GetPortalPositionRequest().GetDimension()

	portals := make([]*amqp.PortalPositionAnswer_PortalPosition, 0)
	if dimension != "" {
		dofusPortal, err := service.getPortal(ctx, server, dimension)
		if err != nil {
			log.Error().Err(err).Str(models.LogCorrelationId, correlationId).Msgf("Returning failed message")
			service.publishPortalAnswerFailed(correlationId, message.Language)
			return
		}

		portals = append(portals, models.MapPortal(dofusPortal))

	} else {
		dofusPortals, err := service.getPortals(ctx, server)
		if err != nil {
			log.Error().Err(err).Str(models.LogCorrelationId, correlationId).Msgf("Returning failed message")
			service.publishPortalAnswerFailed(correlationId, message.Language)
			return
		}

		for _, dofusPortal := range dofusPortals {
			portals = append(portals, models.MapPortal(dofusPortal))
		}
	}

	service.publishPortalAnswerSuccess(portals, correlationId, message.Language)
}

func isValidPortalRequest(message *amqp.RabbitMQMessage) bool {
	return message.Type == amqp.RabbitMQMessage_PORTAL_POSITION_REQUEST && message.GetPortalPositionRequest() != nil
}

func (service *PortalsService) getPortals(ctx context.Context, server string) ([]dofusportals.Portal, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()
	resp, err := service.dofusPortalsClient.GetExternalV1ServersServerIdPortals(ctx, server)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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

func (service *PortalsService) getPortal(ctx context.Context, server, dimension string) (dofusportals.Portal, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()
	resp, err := service.dofusPortalsClient.GetExternalV1ServersServerIdPortalsDimensionId(ctx, server, dimension)
	if err != nil {
		return dofusportals.Portal{}, err
	}
	defer resp.Body.Close()

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

func (service *PortalsService) publishPortalAnswerFailed(correlationId string, language amqp.RabbitMQMessage_Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_PORTAL_POSITION_ANSWER,
		Status:   amqp.RabbitMQMessage_FAILED,
		Language: language,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationId)
	if err != nil {
		log.Error().Err(err).Str(models.LogCorrelationId, correlationId).Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *PortalsService) publishPortalAnswerSuccess(portals []*amqp.PortalPositionAnswer_PortalPosition,
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
		log.Error().Err(err).Str(models.LogCorrelationId, correlationId).Msgf("Cannot publish via broker, request ignored")
	}
}

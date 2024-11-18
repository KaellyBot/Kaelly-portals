package replies

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-portals/models/constants"
	"github.com/rs/zerolog/log"
)

func SucceededAnswer(ctx amqp.Context, broker amqp.MessageBroker,
	message *amqp.RabbitMQMessage) {
	err := broker.Reply(message, ctx.CorrelationID, ctx.ReplyTo)
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogCorrelationID, ctx.CorrelationID).
			Str(constants.LogReplyTo, ctx.ReplyTo).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func FailedAnswer(ctx amqp.Context, broker amqp.MessageBroker,
	messageType amqp.RabbitMQMessage_Type, language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     messageType,
		Status:   amqp.RabbitMQMessage_FAILED,
		Language: language,
	}

	err := broker.Reply(&message, ctx.CorrelationID, ctx.ReplyTo)
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogCorrelationID, ctx.CorrelationID).
			Str(constants.LogReplyTo, ctx.ReplyTo).
			Msgf("Cannot publish via broker, request ignored")
	}
}

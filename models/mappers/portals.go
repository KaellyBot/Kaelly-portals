package mappers

import (
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-portals/models"
	"github.com/kaellybot/kaelly-portals/payloads/dofusportals"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapPortal(portal dofusportals.Portal) *amqp.PortalPositionAnswer_PortalPosition {
	var remainingUses int32 = 0
	if portal.RemainingUses != nil {
		remainingUses = int32(*portal.RemainingUses)
	}

	// TODO map entity ids

	return &amqp.PortalPositionAnswer_PortalPosition{
		Server:        portal.Server,
		Dimension:     portal.Dimension,
		Position:      mapPosition(portal.Position),
		RemainingUses: remainingUses,
		CreatedBy:     mapUser(portal.CreatedBy),
		CreatedAt:     mapTimestamp(portal.CreatedAt),
		UpdatedBy:     mapUser(portal.UpdatedBy),
		UpdatedAt:     mapTimestamp(portal.UpdatedAt),
		Source:        mapSource(models.SourceDofusPortals),
	}
}

func mapPosition(position *dofusportals.Position) *amqp.PortalPositionAnswer_PortalPosition_Position {
	if position == nil {
		return nil
	}

	isInCanopy := false
	if position.IsInCanopy != nil && *position.IsInCanopy {
		isInCanopy = true
	}

	return &amqp.PortalPositionAnswer_PortalPosition_Position{
		X:                    int32(position.X),
		Y:                    int32(position.Y),
		IsInCanopy:           isInCanopy,
		Transport:            mapTransport(position.Transport),
		ConditionalTransport: mapTransport(position.ConditionalTransport),
	}
}

func mapTransport(transport *dofusportals.Transport) *amqp.PortalPositionAnswer_PortalPosition_Position_Transport {
	if transport == nil {
		return nil
	}

	return &amqp.PortalPositionAnswer_PortalPosition_Position_Transport{
		Area:    transport.Area,
		SubArea: transport.SubArea,
		Type:    string(transport.Type),
		X:       int32(transport.X),
		Y:       int32(transport.Y),
	}
}

func mapUser(user *dofusportals.User) string {
	if user == nil {
		return ""
	}

	return user.Name
}

func mapTimestamp(timestamp *time.Time) *timestamppb.Timestamp {
	if timestamp == nil {
		return nil
	}

	return timestamppb.New(*timestamp)
}

func mapSource(source models.Source) *amqp.PortalPositionAnswer_PortalPosition_Source {
	return &amqp.PortalPositionAnswer_PortalPosition_Source{
		Name: source.Name,
		Icon: source.Icon,
		Url:  source.Url,
	}
}

package mappers

import (
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-portals/models"
	"github.com/kaellybot/kaelly-portals/models/constants"
	"github.com/kaellybot/kaelly-portals/payloads/dofusportals"
	"github.com/kaellybot/kaelly-portals/services/areas"
	"github.com/kaellybot/kaelly-portals/services/dimensions"
	"github.com/kaellybot/kaelly-portals/services/servers"
	"github.com/kaellybot/kaelly-portals/services/subareas"
	"github.com/kaellybot/kaelly-portals/services/transports"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapPortal(portal dofusportals.Portal, serverService servers.ServerService,
	dimensionService dimensions.DimensionService, areaService areas.AreaService,
	subAreaService subareas.SubAreaService, transportService transports.TransportService,
) *amqp.PortalPositionAnswer_PortalPosition {
	var remainingUses int32 = 0
	if portal.RemainingUses != nil {
		remainingUses = int32(*portal.RemainingUses)
	}

	return &amqp.PortalPositionAnswer_PortalPosition{
		ServerId:        getInternalServerId(portal.Server, serverService),
		DimensionId:     getInternalDimensionId(portal.Dimension, dimensionService),
		Position:      mapPosition(portal.Position, areaService, subAreaService, transportService),
		RemainingUses: remainingUses,
		CreatedBy:     mapUser(portal.CreatedBy),
		CreatedAt:     mapTimestamp(portal.CreatedAt),
		UpdatedBy:     mapUser(portal.UpdatedBy),
		UpdatedAt:     mapTimestamp(portal.UpdatedAt),
		Source:        mapSource(models.SourceDofusPortals),
	}
}

func mapPosition(position *dofusportals.Position, areaService areas.AreaService,
	subAreaService subareas.SubAreaService, transportService transports.TransportService,
) *amqp.PortalPositionAnswer_PortalPosition_Position {
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
		Transport:            mapTransport(position.Transport, areaService, subAreaService, transportService),
		ConditionalTransport: mapTransport(position.ConditionalTransport, areaService, subAreaService, transportService),
	}
}

func mapTransport(transport *dofusportals.Transport, areaService areas.AreaService,
	subAreaService subareas.SubAreaService, transportService transports.TransportService,
) *amqp.PortalPositionAnswer_PortalPosition_Position_Transport {
	if transport == nil {
		return nil
	}

	return &amqp.PortalPositionAnswer_PortalPosition_Position_Transport{
		AreaId:    getInternalAreaId(transport.Area, areaService),
		SubAreaId: getInternalSubAreaId(transport.SubArea, subAreaService),
		TypeId:    getInternalTransportTypeId(string(transport.Type), transportService),
		X:         int32(transport.X),
		Y:         int32(transport.Y),
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

func getInternalServerId(dofusPortalsId string, serverService servers.ServerService) string {
	server, found := serverService.FindServerByDofusPortalsId(dofusPortalsId)
	if found {
		return server.Id
	}

	log.Warn().Str(constants.LogServerId, dofusPortalsId).
		Msgf("Server not found with following dofusPortalsId, using it as internal one")
	return dofusPortalsId
}

func getInternalDimensionId(dofusPortalsId string, dimensionService dimensions.DimensionService) string {
	dimension, found := dimensionService.FindDimensionByDofusPortalsId(dofusPortalsId)
	if found {
		return dimension.Id
	}

	log.Warn().Str(constants.LogDimensionId, dofusPortalsId).
		Msgf("Dimension not found with following dofusPortalsId, using it as internal one")
	return dofusPortalsId
}

func getInternalAreaId(dofusPortalsId string, areaService areas.AreaService) string {
	area, found := areaService.FindAreaByDofusPortalsId(dofusPortalsId)
	if found {
		return area.Id
	}

	log.Warn().Str(constants.LogAreaId, dofusPortalsId).
		Msgf("Area not found with following dofusPortalsId, using it as internal one")
	return dofusPortalsId
}

func getInternalSubAreaId(dofusPortalsId string, subAreaService subareas.SubAreaService) string {
	subArea, found := subAreaService.FindSubAreaByDofusPortalsId(dofusPortalsId)
	if found {
		return subArea.Id
	}

	log.Warn().Str(constants.LogSubAreaId, dofusPortalsId).
		Msgf("SubArea not found with following dofusPortalsId, using it as internal one")
	return dofusPortalsId
}

func getInternalTransportTypeId(dofusPortalsId string, transportService transports.TransportService) string {
	transportType, found := transportService.FindTransportTypeByDofusPortalsId(dofusPortalsId)
	if found {
		return transportType.Id
	}

	log.Warn().Str(constants.LogTransportTypeId, dofusPortalsId).
		Msgf("TransportType not found with following dofusPortalsId, using it as internal one")
	return dofusPortalsId
}

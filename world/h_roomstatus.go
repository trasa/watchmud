package world

import (
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud-message/direction"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/spaces"
)

func (w *World) handleRoomStatus(msg *gameserver.HandlerParameter) {
	// TODO security: must have admin privs (or, something) to use this command

	//rr := msg.Message.GetRoomStatusRequest()
	// TODO allow user to specify room and zone to get status of

	room := w.getRoomContainingPlayer(msg.Player)
	if room == nil {
		msg.Player.Send(message.RoomStatusResponse{
			Success:    false,
			ResultCode: "NOT_IN_ROOM",
		})
		return
	}

	response := message.RoomStatusResponse{
		Success:       true,
		ResultCode:    "OK",
		PlayerInfo:    createPlayerInfo(room),
		InventoryInfo: createInventoryInfo(room),
		MobInfo:       createMobInfo(room),
		Id:            room.Id,
		Name:          room.Name,
		Description:   room.Description,
		ZoneName:      room.Zone.Name,
		ZoneId:        room.Zone.Id,
		Directions:    createDirections(room),
		Flags:         room.GetFlags(),
	}

	msg.Player.Send(response)
}

func createPlayerInfo(room *spaces.Room) (result []*message.RoomStatusResponse_PlayerInfo) {
	for _, p := range room.GetPlayers() {
		result = append(result, &message.RoomStatusResponse_PlayerInfo{
			Name:          p.GetName(),
			CurrentHealth: p.GetCurrentHealth(),
			MaxHealth:     p.GetMaxHealth(),
		})
	}
	return
}

func createInventoryInfo(room *spaces.Room) (result []*message.RoomStatusResponse_InventoryInfo) {
	for _, i := range room.GetAllInventory() {
		result = append(result,
			&message.RoomStatusResponse_InventoryInfo{
				Id:                  i.Id(),
				DefinitionId:        i.Definition.Id(),
				Aliases:             i.Definition.Aliases,
				Categories:          i.Definition.Categories.ToStringList(),
				Name:                i.Definition.Name,
				ShortDescription:    i.Definition.ShortDescription,
				DescriptionOnGround: i.Definition.DescriptionOnGround,
				ZoneId:              i.Definition.ZoneId,
				Behaviors:           i.Definition.Behaviors.ToStringList(),
			})
	}
	return
}

func createMobInfo(room *spaces.Room) (result []*message.RoomStatusResponse_MobInfo) {
	for _, m := range room.GetMobs() {
		result = append(result,
			&message.RoomStatusResponse_MobInfo{
				Id:                m.Id(),
				DefinitionId:      m.Definition.Id,
				Aliases:           m.Definition.Aliases,
				Name:              m.Definition.Name,
				ShortDescription:  m.Definition.ShortDescription,
				DescriptionInRoom: m.Definition.DescriptionInRoom,
				ZoneId:            m.Definition.ZoneId,
				CurrentHealth:     m.CurHealth,
				MaxHealth:         m.Definition.MaxHealth,
				Flags:             m.Definition.GetFlags(),
			})
	}
	return
}

func createDirections(room *spaces.Room) (result []*message.RoomStatusResponse_DirectionInfo) {
	for _, ex := range room.GetRoomExits(false) {
		dir, _ := direction.DirectionToString(ex.Direction)
		result = append(result,
			&message.RoomStatusResponse_DirectionInfo{
				Dir:    dir,
				RoomId: ex.Room.Id,
				ZoneId: ex.Room.Zone.Id,
				Flags:  ex.Room.GetFlags(),
			})
	}
	return
}

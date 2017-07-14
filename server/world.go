package server

import "log"

type World struct {
	Zones            map[string]*Zone
	StartRoom        *Room
	knownPlayersById map[string]*Player
}

var startZoneKey = "sample"
var startRoomKey = "start"

// Constructor for World
func NewWorld() *World {
	// Build a very boring world.
	w := World{
		Zones:            make(map[string]*Zone),
		knownPlayersById: make(map[string]*Player),
	}
	sampleZone := Zone{
		Id:    startZoneKey,
		Name:  "Sample Zone",
		Rooms: make(map[string]*Room),
	}
	w.Zones[sampleZone.Id] = &sampleZone

	r := NewRoom(&sampleZone, startRoomKey, "Central Portal", "It's a boring room, with boring stuff in it.")
	sampleZone.Rooms[r.Id] = r
	w.StartRoom = r

	// north room
	northRoom := NewRoom(&sampleZone, "northRoom", "North Room", "This room is north of the start.")
	sampleZone.Rooms[northRoom.Id] = northRoom

	r.North = northRoom
	northRoom.South = r

	log.Print("World built.")
	return &w
}

func (w *World) getAllPlayers() (result []*Player) {
	for _, v := range w.knownPlayersById {
		result = append(result, v)
	}
	return result
}

// Add this Player to the world
// putting them in the start room
func (w *World) addPlayer(player *Player) {
	w.knownPlayersById[player.Id] = player
	w.StartRoom.AddPlayer(player)
}

// handle an incoming login message
func (w *World) HandleLogin(message *IncomingMessage) {
	// is this connection already authenticated?
	// see if we can find an existing player ..
	p := FindPlayerByClient(message.Client)
	if p != nil {
		// already authenticated, can't login again
		p.Send(LoginResponse{
			Response: Response{
				MessageType: "login_response",
				Successful:  false,
				ResultCode:  "ALREADY_AUTHENTICATED",
			},
		})
		return
	}

	// todo authentication and stuff
	player := NewPlayer(message.Body["player_name"], message.Body["player_name"], message.Client)
	message.Client.Player = player
	w.addPlayer(player)
	player.Send(LoginResponse{
		Response: Response{
			MessageType: "login_response",
			Successful:  true,
			ResultCode:  "OK",
		},
		Player: player,
	})
}

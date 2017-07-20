package server

import (
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/response"
)

// handle an incoming login message
func (w *World) handleLogin(message *message.IncomingMessage) {
	// is this connection already authenticated?
	// see if we can find an existing player ..
	if message.Client.GetPlayer() != nil {
		// you've already got one
		message.Client.Send(response.LoginResponse{
			Response: response.Response{
				MessageType: "login_response",
				Successful:  false,
				ResultCode:  "PLAYER_ALREADY_ATTACHED",
			},
		})
		return
	}
	// what if player is logged in on a different client?
	// TODO
	/*
		p := FindPlayerByClient(message.Client)
		if p != nil {
			// already authenticated, can't login again
			// TODO
			// note that this isn't really working; the same username can log on twice
			// instead the old player should be kicked and the new player take over
			p.Send(LoginResponse{
				Response: Response{
					MessageType: "login_response",
					Successful:  false,
					ResultCode:  "ALREADY_AUTHENTICATED",
				},
			})
			return
		}
	*/

	// todo authentication and stuff
	player := NewPlayer(message.Body["player_name"], message.Client)
	message.Client.SetPlayer(player)
	message.Player = player
	w.addPlayer(player)
	player.Send(response.LoginResponse{
		Response: response.Response{
			MessageType: "login_response",
			Successful:  true,
			ResultCode:  "OK",
		},
		Player: response.NewPlayerData(player.GetName()),
	})
}

func (w *World) handleTell(message *message.IncomingMessage) {
	fromName := message.Player.GetName()
	toPlayer := w.findPlayerByName(message.Body["to"])
	value := message.Body["value"]

	if toPlayer == nil {
		message.Player.Send(response.Response{
			MessageType: "tell",
			Successful:  false,
			ResultCode:  "TO_PLAYER_NOT_FOUND",
		})
	} else {
		toPlayer.Send(response.TellNotification{
			Notification: response.Notification{MessageType: "tell"},
			From:         fromName,
			Value:        value,
		})
		message.Player.Send(response.Response{
			MessageType: "tell",
			Successful:  true,
			ResultCode:  "OK",
		})
	}
}

// Tell everybody in the game something.
func (w *World) handleTellAll(message *message.IncomingMessage) {
	if val, ok := message.Body["value"]; ok {
		w.SendToAllPlayers(response.TellAllNotification{
			Notification: response.Notification{
				MessageType: "tell_all_notification",
			},
			Value:  val,
			Sender: message.Player.GetName(),
		})
	} else {
		message.Player.Send(response.Response{
			MessageType: "tell_all",
			Successful:  false,
			ResultCode:  "NO_VALUE",
		})
	}
}

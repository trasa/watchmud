package server

type IncomingMessage struct {
	Client Client
	Player *Player
	Body   map[string]string
}

func newIncomingMessage(client Client, body map[string]string) *IncomingMessage {
	return &IncomingMessage{
		Client: client,
		Player: client.GetPlayer(),
		Body:   body,
	}
}

type Notification struct {
	MessageType string `json:"msg_type"`
}

type Response struct {
	MessageType string `json:"msg_type"`
	Successful  bool   `json:"success"`
	ResultCode  string `json:"result_code"`
}

type LoginResponse struct {
	Response
	Player *Player `json:"player"`
}

// handle an incoming login message
func (w *World) handleLogin(message *IncomingMessage) {
	// is this connection already authenticated?
	// see if we can find an existing player ..
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
	player.Send(LoginResponse{
		Response: Response{
			MessageType: "login_response",
			Successful:  true,
			ResultCode:  "OK",
		},
		Player: player,
	})
}

type TellNotification struct {
	Notification
	From  string `json:"from"`
	Value string `json:"value"`
}

func (w *World) handleTell(message *IncomingMessage) {
	fromName := message.Player.Name
	toPlayer := w.findPlayerByName(message.Body["to"])
	value := message.Body["value"]

	if toPlayer == nil {
		message.Player.Send(Response{
			MessageType: "tell",
			Successful:  false,
			ResultCode:  "TO_PLAYER_NOT_FOUND",
		})
	} else {
		toPlayer.Send(TellNotification{
			Notification: Notification{MessageType: "tell"},
			From:         fromName,
			Value:        value,
		})
		message.Player.Send(Response{
			MessageType: "tell",
			Successful:  true,
			ResultCode:  "OK",
		})
	}
}

type TellAllNotification struct {
	Notification
	Value  string `json:"value"`
	Sender string `json:"sender"`
}

// Tell everybody in the game something.
func (w *World) handleTellAll(message *IncomingMessage) {
	if val, ok := message.Body["value"]; ok {
		w.SendToAllPlayers(TellAllNotification{
			Notification: Notification{
				MessageType: "tell_all_notification",
			},
			Value:  val,
			Sender: message.Player.Name,
		})
	} else {
		message.Player.Send(Response{
			MessageType: "tell_all",
			Successful:  false,
			ResultCode:  "NO_VALUE",
		})
	}
}

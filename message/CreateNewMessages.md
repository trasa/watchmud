# Adding New Messages to watchmud (and the client)

Because I don't do this enough to recall which files to edit each 
time through.

## Add Definition to message.proto

Very straightforward: create the message definition in message.proto 
and run `make` (or `go generate`) so that the go types are written.

You must add the message definition to `GameMessage.inner` and then
define the messages themselves.

Messages should look something like:

```
message KillRequest {
    string target = 1;
}

message KillResponse {
    bool success = 1;
    string resultCode = 2;
}
```

## Add Server Message Handler

Create a new handler in world package: `world/h_kill.go`

Add a handler:

```go
package world
import (
	"github.com/trasa/watchmud/gameserver"
    "github.com/trasa/watchmud/message"
)

func (w *World) handleKill(msg *gameserver.HandlerParameter) {
	// TODO
    msg.Player.Send(message.KillResponse{Success:false, ResultCode:"TODO"})
}
```

Register the handler with the server: `world/handlers.go`
Add the handler method to the world.handlerMap.


## Add a Client Message Handler

TODO: make the client handler more dynamic than it is; right now setting
up a new handler is too much work ...

Create a handler: `h_kill.go`

```go
package main

import (
	"github.com/trasa/watchmud/message"
	"fmt"
)

func (c *Client) handleKillResponse(r *message.KillResponse) {
	// TODO
	fmt.Println("Response code should go here.")
}
```

Add the handler to `handler.go` and `handleIncomingMessage`:

```
switch msg.Inner.(type) {
     // ...
    case *message.GameMessage_KillResponse:
        c.handleKillResponse(msg.GetKillResponse())
     // ...
```

## Add the message line translator

This turns the client string into a message. See `message/translator.go`


... And the rest is just implementing your feature.



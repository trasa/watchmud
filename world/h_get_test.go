package world

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
	"testing"
)

func TestGet_success(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("foo")

	w.AddPlayer(p)

	msg := message.IncomingMessage{
		Player: p,
		Request: message.GetRequest{
			Request: message.RequestBase{MessageType: "get"},
			Targets: []string{"knife"},
		},
	}

	w.handleGet(&msg)

	assert.Equal(t, 1, p.SentMessageCount())
	resp := p.GetSentResponse(0).(message.GetResponse)
	assert.True(t, resp.IsSuccessful())
	assert.Equal(t, 1, len(p.GetInventoryMap()))
	assert.Equal(t, "knife", p.GetInventoryMap()["knife"].Id())
	assert.Equal(t, 0, len(w.startRoom.Inventory))
}

func TestGet_targetNotInRoom(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("foo")

	w.AddPlayer(p)

	msg := message.IncomingMessage{
		Player: p,
		Request: message.GetRequest{
			Request: message.RequestBase{MessageType: "get"},
			Targets: []string{"bag_of_coins"},
		},
	}

	w.handleGet(&msg)

	assert.Equal(t, 1, p.SentMessageCount())
	resp := p.GetSentResponse(0).(message.GetResponse)
	assert.False(t, resp.IsSuccessful())
	assert.Equal(t, "TARGET_NOT_FOUND", resp.GetResultCode())
	assert.Equal(t, 0, len(p.GetInventoryMap()))
	assert.Equal(t, 1, len(w.startRoom.Inventory))
}

func TestGet_noTarget(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("foo")

	w.AddPlayer(p)

	msg := message.IncomingMessage{
		Player: p,
		Request: message.GetRequest{
			Request: message.RequestBase{MessageType: "get"},
		},
	}

	w.handleGet(&msg)

	assert.Equal(t, 1, p.SentMessageCount())
	resp := p.GetSentResponse(0).(message.GetResponse)
	assert.False(t, resp.IsSuccessful())
	assert.Equal(t, "NO_TARGET", resp.GetResultCode())
	assert.Equal(t, 0, len(p.GetInventoryMap()))
	assert.Equal(t, 1, len(w.startRoom.Inventory))
}

func TestGet_playerAddFail(t *testing.T) {
	w := newTestWorld()
	p := player.NewTestPlayer("foo")

	w.AddPlayer(p)

	// TODO: some sort of world-wide list of inventory definitions
	// give the player a knife to start with
	// note that two different objects should not have the same instance id
	// -- this is an arbitrary case to make the test work...
	p.AddInventory(w.startRoom.Inventory["knife"].(*object.Instance))

	msg := message.IncomingMessage{
		Player: p,
		Request: message.GetRequest{
			Request: message.RequestBase{MessageType: "get"},
			Targets: []string{"knife"},
		},
	}

	w.handleGet(&msg)

	assert.Equal(t, 1, p.SentMessageCount())
	resp := p.GetSentResponse(0).(message.GetResponse)
	assert.False(t, resp.IsSuccessful())
	assert.Equal(t, "ADD_INVENTORY_ERROR", resp.GetResultCode())
	assert.Equal(t, 1, len(p.GetInventoryMap()))
	assert.Equal(t, 1, len(w.startRoom.Inventory))
}

// TODO: test case for when room.Inventory.Remove fails
// need to figure out how to mock the room

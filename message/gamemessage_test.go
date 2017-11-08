package message

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestNewGameMessage_ReflectPointer(t *testing.T) {
	dr := &DropResponse{ResultCode: "foo"}

	gm, err := NewGameMessage(dr)

	assert.NoError(t, err)
	// expected gm.Inner -> *GameMessage_DropResponse
	assert.IsType(t, &GameMessage_DropResponse{}, gm.Inner)
	// expected gm.GetDropResponse() -> *DropResponse
	// in this case, the value of gm.GetDropResponse() == dr
	// (they both point to the same thing)
	assert.Equal(t, dr, gm.GetDropResponse())
	assert.Equal(t, "foo", gm.GetDropResponse().ResultCode)
}

func TestNewGameMessage_ReflectValueType(t *testing.T) {
	dr := DropResponse{ResultCode: "foo"}

	gm, err := NewGameMessage(dr)
	log.Printf("gm: %v, err: %v", gm, err)
	assert.NoError(t, err)
	// expected gm.Inner -> *GameMessage_DropResponse
	assert.IsType(t, &GameMessage_DropResponse{}, gm.Inner)
	// expected gm.GetDropResponse() -> *DropResponse
	// however, they're not the same object as we create a pointer to wrap them
	//assert.Equal(t, dr, gm.GetDropResponse())
	// the payload is the same though
	assert.Equal(t, "foo", gm.GetDropResponse().ResultCode)
}

func TestNewGameMessage_LookRequest(t *testing.T) {
	gm, err := NewGameMessage(LookRequest{ValueList: []string{"foo"}})
	assert.NoError(t, err)
	assert.IsType(t, &GameMessage_LookRequest{}, gm.Inner)
	assert.NotNil(t, gm.GetLookRequest())
	log.Printf("worked? %v", gm)
	assert.Equal(t, "foo", gm.GetLookRequest().ValueList[0])
}

func TestNewGameMessage_Ping(t *testing.T) {
	gm, err := NewGameMessage(Ping{Target: "x"})
	assert.NoError(t, err)
	assert.IsType(t, &GameMessage_Ping{}, gm.Inner)
	assert.NotNil(t, gm.GetPing())
	assert.Equal(t, "x", gm.GetPing().Target)

}

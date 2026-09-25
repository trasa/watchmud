package telnet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/watchmud/watchmud/event"
)

func TestRenderDied(t *testing.T) {
	died := event.Died{Target: "testdood", IsPlayer: true}

	assert.Equal(t, "You are dead!\n", render(died, "testdood"))
	assert.Equal(t, "testdood is dead!\n", render(died, "otherdood"))
}

func TestRenderEnteredGame(t *testing.T) {
	assert.Equal(t, "testdood has entered the game.\n", render(event.EnteredGame{Actor: "testdood"}, "otherdood"))
}

// A mob that happens to share your name is not you.
func TestRenderDied_mob(t *testing.T) {
	assert.Equal(t, "testdood is dead!\n", render(event.Died{Target: "testdood"}, "testdood"))
}

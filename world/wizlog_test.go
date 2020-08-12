package world

import (
	"github.com/trasa/watchmud/player"
	"testing"
)

func TestWizLog(t *testing.T) {
	p := player.NewTestPlayer("testdood")
	logWizCommand(p, "foo", "%s is %s", p.GetName(), p.GetName())
}

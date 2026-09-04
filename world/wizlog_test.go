package world

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/player"
)

type WizLogSuite struct {
	suite.Suite
	p *player.Player
}

func TestWizLog(t *testing.T) {
	suite.Run(t, new(WizLogSuite))
}

func (s *WizLogSuite) SetupTest() {
	s.p = player.NewTestPlayer("testdood", "testdood", nil)
}

func (s *WizLogSuite) TestLog() {
	logWizCommand(s.p, "foo", "%s is %s", s.p.Name, s.p.Name)
	// TODO finish test
}

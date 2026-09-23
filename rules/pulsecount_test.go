package rules

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type PulseCountSuite struct {
	suite.Suite
}

func TestPulseCountSuite(t *testing.T) {
	suite.Run(t, new(PulseCountSuite))
}

func (suite *PulseCountSuite) SetupTest() {

}

func (suite *PulseCountSuite) TestPulseCount_toDuration() {
	var pulse PulseCount
	// no pulses is no duration
	pulse = 0
	duration := pulse.ToDuration()
	suite.Assert().Equal(time.Duration(0), duration)

	// one pulse = PULSE_INTERVAL
	pulse = 1
	duration = pulse.ToDuration()
	suite.Assert().Equal(PulseInterval, duration)
}

func (suite *PulseCountSuite) TestPulseCount_checkInterval() {
	var pulse PulseCount
	pulse = 100

	suite.Assert().False(pulse.CheckInterval(PulseInterval * 99))
	suite.Assert().True(pulse.CheckInterval(PulseInterval * 100))
	suite.Assert().False(pulse.CheckInterval(PulseInterval * 101))
}

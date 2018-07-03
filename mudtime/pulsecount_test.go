package mudtime

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
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
	suite.Assert().Equal(PULSE_INTERVAL, duration)
}

func (suite *PulseCountSuite) TestPulseCount_checkInterval() {
	var pulse PulseCount
	pulse = 100

	suite.Assert().False(pulse.CheckInterval(PULSE_INTERVAL * 99))
	suite.Assert().True(pulse.CheckInterval(PULSE_INTERVAL * 100))
	suite.Assert().False(pulse.CheckInterval(PULSE_INTERVAL * 101))
}

func (suite *PulseCountSuite) TestDurationToPulseCount() {
	p := PulseCount(int64(1/PULSE_INTERVAL.Hours()) * 10)

	suite.Assert().Equal(p, TimeDurationToPulseCount(time.Hour*10))
	suite.Assert().Equal(PulseCount(0), TimeDurationToPulseCount(0))
	suite.Assert().Equal(PulseCount(1), TimeDurationToPulseCount(PULSE_INTERVAL))
}

func (suite *PulseCountSuite) TestPulseDurations() {
	start := TimeDurationToPulseCount(time.Hour * 3)
	end := TimeDurationToPulseCount(time.Hour * 4)

	suite.Assert().Equal(time.Hour*1, DurationBetween(start, end))
}

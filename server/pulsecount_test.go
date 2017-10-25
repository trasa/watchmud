package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPulseCount_toDuration(t *testing.T) {
	var pulse PulseCount
	// no pulses is no duration
	pulse = 0
	duration := pulse.toDuration()
	assert.Equal(t, time.Duration(0), duration)

	// one pulse = PULSE_INTERVAL
	pulse = 1
	duration = pulse.toDuration()
	assert.Equal(t, PULSE_INTERVAL, duration)
}

func TestPulseCount_checkInterval(t *testing.T) {
	var pulse PulseCount
	pulse = 100

	assert.False(t, pulse.checkInterval(PULSE_INTERVAL*99))
	assert.True(t, pulse.checkInterval(PULSE_INTERVAL*100))
	assert.False(t, pulse.checkInterval(PULSE_INTERVAL*101))
}

func TestDurationToPulseCount(t *testing.T) {
	p := PulseCount(int64(1/PULSE_INTERVAL.Hours()) * 10)

	assert.Equal(t, p, timeDurationToPulseCount(time.Hour*10))
	assert.Equal(t, PulseCount(0), timeDurationToPulseCount(0))
	assert.Equal(t, PulseCount(1), timeDurationToPulseCount(PULSE_INTERVAL))
}

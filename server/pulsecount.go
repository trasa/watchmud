package server

import "time"

// Max value for an int64 is math.MaxInt64 == 9223372036854775807
// If PULSE_INTERVAL is 1 nanosecond, the PULSE rollover won't happen
// for ~29 decades. If for some reason our uptime can be longer than
// that, the pulse-rollover code should be reinstated.
type PulseCount int64

func (pc PulseCount) toDuration() time.Duration {
	return time.Duration(int64(pc) * PULSE_INTERVAL.Nanoseconds())
}

func (pc PulseCount) checkInterval(i time.Duration) bool {
	return pc.toDuration()%i == 0
}

// Determines how many Pulses would happen in time.Duration i
func timeDurationToPulseCount(i time.Duration) PulseCount {
	return PulseCount(i / PULSE_INTERVAL)
}

package rules

import "time"

// PulseInterval indicates when the game server 'ticks'
const PulseInterval = 1 * time.Second

// PulseCount describes a pulse. It's a ticker.
// Max value for an int64 is math.MaxInt64 == 9223372036854775807
// If PulseInterval is 1 nanosecond, the PULSE rollover won't happen
// for ~29 decades. If for some reason our uptime can be longer than
// that, the pulse-rollover code should be reinstated.
type PulseCount int64

// PulseNever for initial events or events that do not happen.
const PulseNever = PulseCount(0)

func (pc PulseCount) ToDuration() time.Duration {
	return time.Duration(int64(pc) * PulseInterval.Nanoseconds())
}

func DurationBetween(start PulseCount, end PulseCount) time.Duration {
	return (end - start).ToDuration()
}

func (pc PulseCount) CheckInterval(i time.Duration) bool {
	return pc.ToDuration()%i == 0
}

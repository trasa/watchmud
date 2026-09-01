package mudtime

import "time"

// Game Server ticks every PULSE_INTERVAL time
// const PULSE_INTERVAL time.Duration = 100 * time.Millisecond // 0.1 seconds
const PulseInterval time.Duration = 1 * time.Second // 1 second

// Mobs consider doing something once every PULSE_MOBILE time
const PulseMobile = 10 * time.Second

// battle happens on this interval
const PulseViolence = 1 * time.Second

// zones reset based on their lifetime, check every pulse_zone (lifetime > pulse_zone)
const PulseZone = 1 * time.Minute

// Max value for an int64 is math.MaxInt64 == 9223372036854775807
// If PULSE_INTERVAL is 1 nanosecond, the PULSE rollover won't happen
// for ~29 decades. If for some reason our uptime can be longer than
// that, the pulse-rollover code should be reinstated.
type PulseCount int64

const PulseCountNever = PulseCount(0)

func (pc PulseCount) ToDuration() time.Duration {
	return time.Duration(int64(pc) * PulseInterval.Nanoseconds())
}

func (pc PulseCount) CheckInterval(i time.Duration) bool {
	return pc.ToDuration()%i == 0
}

// Determines how many Pulses would happen in time.Duration i
func TimeDurationToPulseCount(i time.Duration) PulseCount {
	return PulseCount(i / PulseInterval)
}

func DurationBetween(start PulseCount, end PulseCount) time.Duration {
	return (end - start).ToDuration()
}

package server

import "time"

type PulseCount int64

func (pc PulseCount) toDuration() time.Duration {
	return time.Duration(int64(pc) * TICK_INTERVAL.Nanoseconds())
}

func (pc PulseCount) checkInterval(i time.Duration) bool {
	return pc.toDuration()%i == 0
}

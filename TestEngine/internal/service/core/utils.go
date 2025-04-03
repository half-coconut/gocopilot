package core

import "time"

func index(percent float64, b []time.Duration) time.Duration {
	idx := int64(percent / 100.0 * float64(len(b)))
	// 防止越界
	if idx > int64(len(b)) {
		idx = int64(len(b) - 1)
	}
	return b[idx]
}

var durations = [...]time.Duration{
	time.Hour,
	time.Minute,
	time.Second,
	time.Millisecond,
	time.Microsecond,
	time.Nanosecond,
}

func round(d time.Duration) time.Duration {
	for i, unit := range durations {
		if d >= unit && i < len(durations)-1 {
			return d.Round(durations[i+1])
		}
	}
	return d
}

package plot

import "time"

func DurationTo(durations []time.Duration, scale time.Duration) []float64 {
	values := make([]float64, len(durations))
	for i, dur := range durations {
		values[i] = float64(dur) / float64(scale)
	}
	return values
}

func DurationToNanoseconds(durations []time.Duration) []float64 {
	return DurationTo(durations, time.Nanosecond)
}

func DurationToSeconds(durations []time.Duration) []float64 {
	return DurationTo(durations, time.Second)
}

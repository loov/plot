package plot

import "time"

func DurationToNanoseconds(durations []time.Duration) []float64 {
	values := make([]float64, len(durations))
	for i, dur := range durations {
		values[i] = float64(dur.Nanoseconds())
	}
	return values
}

func DurationToSeconds(durations []time.Duration) []float64 {
	values := make([]float64, len(durations))
	for i, dur := range durations {
		values[i] = float64(dur.Seconds())
	}
	return values
}

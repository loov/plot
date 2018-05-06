package plot

import "math"

type Stats struct {
	DiscreteX, DiscreteY bool

	Min    Point
	Center Point
	Max    Point
}

var nanStats = Stats{
	DiscreteX: false,
	DiscreteY: false,

	Min:    nanPoint,
	Center: nanPoint,
	Max:    nanPoint,
}

func tryGetStats(element Element) (Stats, bool) {
	if dataset, ok := element.(Dataset); ok {
		return dataset.Stats(), true
	}
	return nanStats, false
}

func maximalStats(els []Element) Stats {
	first := true
	finalstats := nanStats
	for _, element := range els {
		if elstats, ok := tryGetStats(element); ok {
			if first {
				first = false
				finalstats.DiscreteX = elstats.DiscreteX
				finalstats.DiscreteY = elstats.DiscreteY
			}

			if math.IsNaN(finalstats.Min.X) {
				finalstats.Min.X = elstats.Min.X
			} else if !math.IsNaN(elstats.Min.X) {
				finalstats.Min.X = math.Min(finalstats.Min.X, elstats.Min.X)
			}
			if math.IsNaN(finalstats.Max.X) {
				finalstats.Max.X = elstats.Max.X
			} else if !math.IsNaN(elstats.Max.X) {
				finalstats.Max.X = math.Max(finalstats.Max.X, elstats.Max.X)
			}

			if math.IsNaN(finalstats.Min.Y) {
				finalstats.Min.Y = elstats.Min.Y
			} else if !math.IsNaN(elstats.Min.Y) {
				finalstats.Min.Y = math.Min(finalstats.Min.Y, elstats.Min.Y)
			}
			if math.IsNaN(finalstats.Max.Y) {
				finalstats.Max.Y = elstats.Max.Y
			} else if !math.IsNaN(elstats.Max.Y) {
				finalstats.Max.Y = math.Max(finalstats.Max.Y, elstats.Max.Y)
			}
		}
	}
	return finalstats
}

package plot

import (
	"math"
)

type Axis struct {
	// Min value of the axis (in value space)
	Min float64
	// Max value of the axis (in value space)
	Max float64
	// Transform transform [0..1] -> float64
	Transform func(p float64) float64
}

func NewAxis() Axis {
	return Axis{
		Min:       math.NaN(),
		Max:       math.NaN(),
		Transform: nil,
	}
}

func (axis *Axis) IsValid() bool {
	return math.IsNaN(axis.Min) || math.IsNaN(axis.Max)
}

func (axis *Axis) Include(min, max float64) {
	if math.IsNaN(axis.Min) {
		axis.Min = min
	} else {
		axis.Min = math.Min(axis.Min, min)
	}

	if math.IsNaN(axis.Max) {
		axis.Max = max
	} else {
		axis.Max = math.Max(axis.Max, max)
	}
}

func DetectAxis(x, y Axis, elements []Element) (X, Y Axis) {
	spanx, spany := NewAxis(), NewAxis()
	for _, element := range elements {
		if dataset, ok := element.(Dataset); ok {
			stats := dataset.Stats()
			spanx.Include(stats.Min.X, stats.Max.Y)
			spany.Include(stats.Min.Y, stats.Max.Y)
		}
	}

	spanx.Min, spanx.Max, _ = NiceAxis(spanx.Min, spanx.Max, 2)
	spany.Min, spany.Max, _ = NiceAxis(spany.Min, spany.Max, 2)

	if math.IsNaN(x.Min) {
		spanx.Min = x.Min
	}
	if math.IsNaN(x.Max) {
		spanx.Max = x.Max
	}

	if math.IsNaN(y.Min) {
		spany.Min = y.Min
	}
	if math.IsNaN(y.Max) {
		spany.Max = y.Max
	}
	return Axis{}, Axis{}
}

func NiceAxis(min, max float64, maxticks int) (nicemin, nicemax, tickSpacing float64) {
	span := NiceNumber(max-min, false)
	tickSpacing = NiceNumber(span/(float64(maxticks)-1), true)
	nicemin = math.Floor(min/tickSpacing) * tickSpacing
	nicemax = math.Ceil(max/tickSpacing) * tickSpacing
	return nicemin, nicemax, tickSpacing
}

func NiceNumber(span float64, round bool) float64 {
	exp := math.Floor(math.Log10(span))
	frac := span / math.Pow(10, exp)
	var nice float64
	if round {
		switch {
		case frac < 1.5:
			nice = 1
		case frac < 3:
			nice = 2
		case frac < 7:
			nice = 5
		default:
			nice = 10
		}
	} else {
		switch {
		case frac <= 1:
			nice = 1
		case frac <= 2:
			nice = 2
		case frac <= 5:
			nice = 5
		default:
			nice = 10
		}
	}
	return nice * math.Pow(10, exp)
}

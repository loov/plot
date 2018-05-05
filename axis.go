package plot

import (
	"math"
)

type Axis struct {
	// Min value of the axis (in value space)
	Min float64
	// Max value of the axis (in value space)
	Max float64

	Ticks      Ticks
	MajorTicks int
	MinorTicks int

	// Transform transform [0..1] -> float64
	Transform func(p float64) float64
	Inverse   func(p float64) float64
}

func NewAxis() Axis {
	return Axis{
		Min: math.NaN(),
		Max: math.NaN(),

		Ticks:      AutomaticTicks{},
		MajorTicks: 5,
		MinorTicks: 5,

		Transform: nil,
		Inverse:   nil,
	}
}

func (axis *Axis) IsValid() bool {
	return !math.IsNaN(axis.Min) && !math.IsNaN(axis.Max)
}

func (axis *Axis) fixNaN() {
	if math.IsNaN(axis.Min) {
		axis.Min = 0
	}
	if math.IsNaN(axis.Max) {
		axis.Max = 1
	}
}

func (axis *Axis) ToCanvas(v float64, screenMin, screenMax Length) Length {
	n := (v - axis.Min) / (axis.Max - axis.Min)
	if axis.Transform != nil {
		n = axis.Transform(n)
	}
	return screenMin + n*(screenMax-screenMin)
}

func (axis *Axis) FromCanvas(s Length, screenMin, screenMax Length) float64 {
	n := (s - screenMin) / (screenMax - screenMin)
	if axis.Inverse != nil {
		n = axis.Inverse(n)
	}
	return axis.Min + n*(axis.Max-axis.Min)
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

func detectAxis(x, y Axis, elements []Element) (X, Y Axis) {
	spanx, spany := NewAxis(), NewAxis()
	for _, element := range elements {
		if dataset, ok := element.(Dataset); ok {
			stats := dataset.Stats()
			spanx.Include(stats.Min.X, stats.Max.X)
			spany.Include(stats.Min.Y, stats.Max.Y)
		}
	}

	spanx.Min, spanx.Max = niceAxis(spanx.Min, spanx.Max, spanx.MajorTicks, spanx.MinorTicks)
	spany.Min, spany.Max = niceAxis(spany.Min, spany.Max, spany.MajorTicks, spany.MinorTicks)

	if !math.IsNaN(x.Min) {
		spanx.Min = x.Min
	}
	if !math.IsNaN(x.Max) {
		spanx.Max = x.Max
	}
	if !math.IsNaN(y.Min) {
		spany.Min = y.Min
	}
	if !math.IsNaN(y.Max) {
		spany.Max = y.Max
	}

	spanx.fixNaN()
	spany.fixNaN()

	return spanx, spany
}

func niceAxis(min, max float64, major, minor int) (nicemin, nicemax float64) {
	span := niceNumber(max-min, false)
	tickSpacing := niceNumber(span/(float64(major*minor)-1), true)
	nicemin = math.Floor(min/tickSpacing) * tickSpacing
	nicemax = math.Ceil(max/tickSpacing) * tickSpacing
	return nicemin, nicemax
}

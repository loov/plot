package plot

import (
	"math"
)

type Axis struct {
	// Min value of the axis (in value space)
	Min float64
	// Max value of the axis (in value space)
	Max float64

	Flip bool

	Ticks      Ticks
	MajorTicks int
	MinorTicks int

	Transform AxisTransform
}

type AxisTransform interface {
	ToCanvas(axis *Axis, v float64, screenMin, screenMax Length) Length
	FromCanvas(axis *Axis, s Length, screenMin, screenMax Length) float64
}

func NewAxis() *Axis {
	return &Axis{
		Min: math.NaN(),
		Max: math.NaN(),

		Ticks:      AutomaticTicks{},
		MajorTicks: 5,
		MinorTicks: 5,
	}
}

func project(data []Point, x, y *Axis, bounds Rect) []Point {
	points := make([]Point, 0, len(data))
	size := bounds.Size()
	for _, p := range data {
		p.X = x.ToCanvas(p.X, 0, size.X)
		p.Y = y.ToCanvas(p.Y, 0, size.Y)
		points = append(points, p)
	}
	return points
}

func projectcb(data []Point, x, y *Axis, bounds Rect, fn func(p Point)) {
	size := bounds.Size()
	for _, p := range data {
		p.X = x.ToCanvas(p.X, 0, size.X)
		p.Y = y.ToCanvas(p.Y, 0, size.Y)
		fn(p)
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

func (axis *Axis) lowhigh() (float64, float64) {
	if !axis.Flip {
		return axis.Min, axis.Max
	} else {
		return axis.Max, axis.Min
	}
}

func (axis *Axis) ToCanvas(v float64, screenMin, screenMax Length) Length {
	if axis.Transform != nil {
		return axis.Transform.ToCanvas(axis, v, screenMin, screenMax)
	}

	low, high := axis.lowhigh()
	n := (v - low) / (high - low)
	return screenMin + n*(screenMax-screenMin)
}

func (axis *Axis) FromCanvas(s Length, screenMin, screenMax Length) float64 {
	if axis.Transform != nil {
		return axis.Transform.FromCanvas(axis, s, screenMin, screenMax)
	}

	low, high := axis.lowhigh()
	n := (s - screenMin) / (screenMax - screenMin)
	return low + n*(high-low)
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

func detectAxis(x, y *Axis, elements []Element) (X, Y *Axis) {
	tx, ty := NewAxis(), NewAxis()
	*tx, *ty = *x, *y
	for _, element := range elements {
		if stats, ok := tryGetStats(element); ok {
			tx.Include(stats.Min.X, stats.Max.X)
			ty.Include(stats.Min.Y, stats.Max.Y)
		}
	}

	tx.Min, tx.Max = niceAxis(tx.Min, tx.Max, tx.MajorTicks, tx.MinorTicks)
	ty.Min, ty.Max = niceAxis(ty.Min, ty.Max, ty.MajorTicks, ty.MinorTicks)

	if !math.IsNaN(x.Min) {
		tx.Min = x.Min
	}
	if !math.IsNaN(x.Max) {
		tx.Max = x.Max
	}
	if !math.IsNaN(y.Min) {
		ty.Min = y.Min
	}
	if !math.IsNaN(y.Max) {
		ty.Max = y.Max
	}

	tx.fixNaN()
	ty.fixNaN()

	return tx, ty
}

func niceAxis(min, max float64, major, minor int) (nicemin, nicemax float64) {
	span := niceNumber(max-min, false)
	tickSpacing := niceNumber(span/(float64(major*minor)-1), true)
	nicemin = math.Floor(min/tickSpacing) * tickSpacing
	nicemax = math.Ceil(max/tickSpacing) * tickSpacing
	return nicemin, nicemax
}

type ScreenSpaceTransform struct {
	Transform func(v float64) float64
	Inverse   func(v float64) float64
}

func (tx *ScreenSpaceTransform) ToCanvas(axis *Axis, v float64, screenMin, screenMax Length) Length {
	low, high := axis.lowhigh()
	n := (v - low) / (high - low)
	if tx.Transform != nil {
		n = tx.Transform(n)
	}
	return screenMin + n*(screenMax-screenMin)
}

func (tx *ScreenSpaceTransform) FromCanvas(axis *Axis, s Length, screenMin, screenMax Length) float64 {
	low, high := axis.lowhigh()
	n := (s - screenMin) / (screenMax - screenMin)
	if tx.Inverse != nil {
		n = tx.Inverse(n)
	}
	return low + n*(high-low)
}

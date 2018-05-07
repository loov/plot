package plot

import (
	"math"
	"sort"
)

type Percentiles struct {
	Style
	Label string
	Data  []Point
}

func NewPercentiles(label string, values []float64) *Percentiles {
	values = append(values[:0:0], values...)
	sort.Float64s(values)

	points := make([]Point, 0, len(values))
	multipler := 1 / float64(len(values))
	for i, v := range values {
		var p Point
		p.X = float64(i+1) * multipler
		p.Y = v
		points = append(points, p)
	}

	return &Percentiles{
		Label: label,
		Data:  points,
	}
}

func (line *Percentiles) Stats() Stats { return PointsStats(line.Data) }
func (line *Percentiles) Draw(plot *Plot, canvas Canvas) {
	canvas = canvas.Clip(canvas.Bounds())

	points := make([]Point, 0, len(line.Data))
	lastScreenX := math.Inf(-1)
	projectcb(line.Data, plot.X, plot.Y, canvas.Bounds(), func(p Point) {
		if math.Abs(lastScreenX-p.X) > 0.5 {
			points = append(points, p)
			lastScreenX = p.X
		}
	})

	if !line.Style.IsZero() {
		canvas.Poly(points, &line.Style)
	} else {
		canvas.Poly(points, &plot.Theme.Line)
	}
}

type PercentileTransform struct {
	levels  int
	base    float64
	mulbase float64
}

func NewPercentileTransform(levels int) *PercentileTransform {
	base := math.Pow(0.1, float64(levels))
	return &PercentileTransform{
		levels:  levels,
		base:    base,
		mulbase: 1 / math.Log(base),
	}
}

func (tx *PercentileTransform) transform(v float64) float64 {
	return math.Log(1-v) * tx.mulbase
}

func (tx *PercentileTransform) inverse(v float64) float64 {
	return 1 - math.Pow(tx.base, v)
}

func (tx *PercentileTransform) ToCanvas(axis *Axis, v float64, screenMin, screenMax Length) Length {
	v = tx.transform(v)
	low, high := axis.lowhigh()
	n := (v - low) / (high - low)
	return screenMin + n*(screenMax-screenMin)
}

func (tx *PercentileTransform) FromCanvas(axis *Axis, s Length, screenMin, screenMax Length) float64 {
	low, high := axis.lowhigh()
	n := (s - screenMin) / (screenMax - screenMin)
	v := low + n*(high-low)
	return tx.inverse(v)
}

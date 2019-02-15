package plot

import (
	"fmt"
	"math"
)

type Ticks interface {
	Ticks(axis *Axis) []Tick
}

type Tick struct {
	Minor bool
	Label string
	Value float64
}

type AutomaticTicks struct{}

func (at AutomaticTicks) logarithmicTicks(axis *Axis, transform *Log1pTransform) []Tick {
	//TODO: fix, we don't properly assign labels for logarithmic axis

	ticks := make([]Tick, 0)

	low, high := axis.Min, axis.Max
	if low > high {
		low, high = high, low
	}

	previous := math.NaN()

	inRange := func(value float64) bool {
		return low < value && value < high
	}

	if inRange(0) {
		ticks = append(ticks, Tick{Value: 0, Label: "0"})
		previous = 0
	}

	for power := 0; power < 10; power++ {
		value := math.Pow(transform.base, float64(power))
		if inRange(value) {
			ticks = append(ticks, Tick{
				Value: value,
				Label: fmt.Sprintf("%.0f", value),
			})
		}
		if inRange(-value) {
			ticks = append(ticks, Tick{
				Value: -value,
				Label: fmt.Sprintf("%.0f", -value),
			})
		}

		if !math.IsNaN(previous) && axis.MinorTicks > 0 {
			minorSpacing := (value - previous) / float64(axis.MinorTicks)
			minor := previous
			for i := 0; i < axis.MinorTicks; i++ {
				if inRange(minor) {
					ticks = append(ticks, Tick{Minor: true, Value: minor})
				}
				if inRange(-minor) {
					ticks = append(ticks, Tick{Minor: true, Value: -minor})
				}
				minor += minorSpacing
			}
		}
		previous = value
	}

	return ticks
}

func (AutomaticTicks) linearTicks(axis *Axis) []Tick {
	majorSpacing := (axis.Max - axis.Min) / float64(axis.MajorTicks)
	minorSpacing := majorSpacing / float64(axis.MinorTicks)

	frac := -int(math.Floor(math.Log10(majorSpacing)))
	if frac < 0 {
		frac = 0
	}

	ticks := make([]Tick, 0, axis.MajorTicks*axis.MinorTicks)

	major := axis.Min
	hasZero := false
	for i := 0; i < axis.MajorTicks; i++ {
		if major == 0 {
			hasZero = true
		}
		ticks = append(ticks, Tick{
			Value: major,
			Label: fmt.Sprintf("%.[2]*[1]f", major, frac),
		})

		minor := major
		for k := 0; k < axis.MinorTicks; k++ {
			if minor == 0 {
				minor += minorSpacing
				continue
			}
			ticks = append(ticks, Tick{
				Minor: true,
				Value: minor,
			})
			minor += minorSpacing
		}

		major += majorSpacing
	}

	if !hasZero && (axis.Min <= 0) == (0 <= axis.Max) {
		ticks = append(ticks, Tick{
			Value: 0,
			Label: "0",
		})
	}

	return ticks
}

func (ticks AutomaticTicks) Ticks(axis *Axis) []Tick {
	// if transform, ok := axis.Transform.(*Log1pTransform); ok {
	// 	return ticks.logarithmicTicks(axis, transform)
	// }
	return ticks.linearTicks(axis)
}

type ManualTicks []Tick

func (ticks ManualTicks) Ticks(axis *Axis) []Tick { return []Tick(ticks) }

type TickLabels struct {
	X, Y  bool
	Style Style
}

func NewTickLabels() *TickLabels {
	labels := &TickLabels{}
	labels.X, labels.Y = true, true
	return labels
}

func (ticklabels *TickLabels) Draw(plot *Plot, canvas Canvas) {
	x, y := plot.X, plot.Y

	sz := canvas.Bounds().Size()
	xmin := x.ToCanvas(x.Min, 0, sz.X)
	ymin := y.ToCanvas(y.Min, 0, sz.Y)

	style := &ticklabels.Style
	if style.IsZero() {
		style = &plot.Theme.FontSmall
	}

	if ticklabels.X {
		for _, tick := range x.Ticks.Ticks(x) {
			p := x.ToCanvas(tick.Value, 0, sz.X)
			if tick.Label != "" {
				canvas.Text(tick.Label, P(p, ymin), style)
			}
		}
	}

	if ticklabels.Y {
		for _, tick := range y.Ticks.Ticks(y) {
			p := y.ToCanvas(tick.Value, 0, sz.Y)
			if tick.Label != "" {
				canvas.Text(tick.Label, P(xmin, p), style)
			}
		}
	}
}

func NewPercentilesAxis() *Axis {
	axis := NewAxis()
	axis.Transform = NewPercentileTransform(5)
	axis.Ticks = ManualTicks{
		{Value: 0, Label: "0"},
		{Value: 0.25, Label: "25"},
		{Value: 0.5, Label: "50"},
		{Value: 0.75, Label: "75"},
		{Value: 0.9, Label: "90"},
		{Value: 0.99, Label: "99"},
		{Value: 0.999, Label: "99.9"},
		{Value: 0.9999, Label: "99.99"},
		{Value: 0.99999, Label: "99.999"}}
	return axis
}

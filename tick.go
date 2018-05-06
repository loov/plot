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

func (AutomaticTicks) Ticks(axis *Axis) []Tick {
	majorSpacing := (axis.Max - axis.Min) / float64(axis.MajorTicks)
	minorSpacing := majorSpacing / float64(axis.MinorTicks)

	frac := -int(math.Floor(math.Log10(majorSpacing)))
	if frac < 0 {
		frac = 0
	}

	ticks := make([]Tick, 0, axis.MajorTicks*axis.MinorTicks)

	major := axis.Min
	for i := 0; i < axis.MajorTicks; i++ {
		ticks = append(ticks, Tick{
			Value: major,
			Label: fmt.Sprintf("%.[2]*[1]f", major, frac),
		})

		minor := major
		for k := 0; k < axis.MinorTicks; k++ {
			ticks = append(ticks, Tick{
				Minor: true,
				Value: minor,
			})
			minor += minorSpacing
		}

		major += majorSpacing
	}

	return ticks
}

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

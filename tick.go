package plot

import "fmt"

type Ticks interface {
	Ticks(axis *Axis) []Tick
}

type Tick struct {
	Minor bool
	Value float64
}

func (tick *Tick) Label() string {
	if tick.Minor {
		return ""
	}
	return fmt.Sprintf("%.2f", tick.Value)
}

type AutomaticTicks struct{}

func (AutomaticTicks) Ticks(axis *Axis) []Tick {
	majorSpacing := (axis.Max - axis.Min) / float64(axis.MajorTicks)
	minorSpacing := majorSpacing / float64(axis.MinorTicks)

	ticks := make([]Tick, 0, axis.MajorTicks*axis.MinorTicks)

	major := axis.Min
	for i := 0; i < axis.MajorTicks; i++ {
		ticks = append(ticks, Tick{
			Value: major,
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
			label := tick.Label()
			if label != "" {
				canvas.Text(label, P(p, ymin), style)
			}
		}
	}

	if ticklabels.Y {
		for _, tick := range y.Ticks.Ticks(y) {
			p := y.ToCanvas(tick.Value, 0, sz.Y)
			label := tick.Label()
			if label != "" {
				canvas.Text(label, P(xmin, p), style)
			}
		}
	}
}

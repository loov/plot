package plot

type Ticks interface {
	Ticks(axis *Axis) []Tick
}

type Tick struct {
	Minor bool
	Value float64
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

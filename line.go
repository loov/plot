package plot

type Line struct {
	Style
	Label string

	Data []Point
}

func NewLine(label string, points []Point) *Line {
	return &Line{
		Label: label,
		Data:  points,
	}
}

func (line *Line) Stats() Stats {
	min, avg, max := nanPoint, Point{}, nanPoint

	if len(line.Data) > 0 {
		min = line.Data[0]
		max = line.Data[0]
	}

	for _, p := range line.Data {
		min = min.Min(p)
		avg = avg.Add(p)
		max = max.Max(p)
	}

	return Stats{
		Min:    min,
		Center: avg.Scale(1 / float64(len(line.Data))),
		Max:    max,
	}
}

func (line *Line) Draw(plot *Plot, canvas Canvas) {
	canvas = canvas.Clip(canvas.Bounds())
	points := project(line.Data, plot.X, plot.Y, canvas.Bounds())

	if !line.Style.IsZero() {
		canvas.Poly(points, &line.Style)
	} else {
		canvas.Poly(points, &plot.Theme.Line)
	}
}

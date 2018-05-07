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

func (line *Line) Stats() Stats { return PointsStats(line.Data) }

func (line *Line) Draw(plot *Plot, canvas Canvas) {
	canvas = canvas.Clip(canvas.Bounds())
	points := project(line.Data, plot.X, plot.Y, canvas.Bounds())

	if !line.Style.IsZero() {
		canvas.Poly(points, &line.Style)
	} else {
		canvas.Poly(points, &plot.Theme.Line)
	}
}

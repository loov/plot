package plot

// Line implements a simple line plot.
type Line struct {
	Style
	Label string

	Data []Point
}

// NewLine creates a new line element from the given points.
func NewLine(label string, points []Point) *Line {
	return &Line{
		Label: label,
		Data:  points,
	}
}

// Stats calculates element statistics.
func (line *Line) Stats() Stats {
	return PointsStats(line.Data)
}

// Draw draws the element to canvas.
func (line *Line) Draw(plot *Plot, canvas Canvas) {
	canvas = canvas.Clip(canvas.Bounds())
	points := project(line.Data, plot.X, plot.Y, canvas.Bounds())

	if !line.Style.IsZero() {
		canvas.Poly(points, &line.Style)
	} else {
		canvas.Poly(points, &plot.Theme.Line)
	}
}

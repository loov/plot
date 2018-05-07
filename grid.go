package plot

import "image/color"

type Grid struct {
	GridTheme
}

func NewGrid() *Grid {
	return &Grid{}
}

func (grid *Grid) Draw(plot *Plot, canvas Canvas) {
	x, y := plot.X, plot.Y

	size := canvas.Bounds().Size()
	// xmin, xmax := x.ToCanvas(x.Min, 0, size.X), x.ToCanvas(x.Max, 0, size.X)
	// ymin, ymax := y.ToCanvas(y.Min, 0, size.Y), y.ToCanvas(y.Max, 0, size.Y)
	xmin, xmax := 0.0, size.X
	ymin, ymax := 0.0, size.Y

	theme := &grid.GridTheme
	if theme.IsZero() {
		theme = &plot.Theme.Grid
	}

	canvas.Rect(canvas.Bounds(), &Style{
		Fill:  theme.Fill,
		Class: "grid-fill",
	})

	major := &Style{
		Size:   1,
		Stroke: theme.Major,
		Class:  "grid-major",
	}
	minor := &Style{
		Size:   1,
		Stroke: theme.Minor,
		Class:  "grid-minor",
	}

	for _, tick := range x.Ticks.Ticks(x) {
		p := x.ToCanvas(tick.Value, 0, size.X)
		if tick.Minor {
			canvas.Poly(Ps(p, ymin, p, ymax), minor)
		} else {
			canvas.Poly(Ps(p, ymin, p, ymax), major)
		}
	}

	for _, tick := range y.Ticks.Ticks(y) {
		p := y.ToCanvas(tick.Value, 0, size.Y)
		if tick.Minor {
			canvas.Poly(Ps(xmin, p, xmax, p), minor)
		} else {
			canvas.Poly(Ps(xmin, p, xmax, p), major)
		}
	}
}

type Gizmo struct {
	Center Point
}

func NewGizmo() *Gizmo { return &Gizmo{} }

func (gizmo *Gizmo) Draw(plot *Plot, canvas Canvas) {
	x, y := plot.X, plot.Y

	sz := canvas.Bounds().Size()
	x0, xmin, xmax := x.ToCanvas(gizmo.Center.X, 0, sz.X), x.ToCanvas(x.Min, 0, sz.X), x.ToCanvas(x.Max, 0, sz.X)
	y0, ymin, ymax := y.ToCanvas(gizmo.Center.Y, 0, sz.Y), y.ToCanvas(y.Min, 0, sz.Y), y.ToCanvas(y.Max, 0, sz.Y)

	if xmin < x0 && x0 < xmax {
		canvas.Poly(Ps(x0, ymin, x0, ymax), &Style{
			Stroke: color.NRGBA{30, 0, 0, 100},
		})
	}

	if ymin < y0 && y0 < ymax {
		canvas.Poly(Ps(xmin, y0, xmax, y0), &Style{
			Stroke: color.NRGBA{0, 30, 0, 100},
		})
	}
}

package plot

import "image/color"

type Grid struct {
	Fill  color.Color
	Major color.Color
	Minor color.Color
}

func NewGrid() *Grid {
	return &Grid{
		Fill:  color.NRGBA{240, 240, 240, 255},
		Major: color.NRGBA{255, 255, 255, 255},
		Minor: color.NRGBA{255, 255, 255, 80},
	}
}

func (grid *Grid) Draw(plot *Plot, canvas Canvas) {
	x, y := plot.X, plot.Y

	sz := canvas.Bounds().Size()
	x0, x1 := x.ToCanvas(x.Min, 0, sz.X), x.ToCanvas(x.Max, 0, sz.X)
	y0, y1 := y.ToCanvas(y.Min, 0, sz.Y), y.ToCanvas(y.Max, 0, sz.Y)

	canvas.Rect(canvas.Bounds(), &Style{
		Fill: grid.Fill,
	})

	major := &Style{
		Size:   1,
		Stroke: grid.Major,
	}
	minor := &Style{
		Size:   1,
		Stroke: grid.Minor,
	}

	for _, tick := range x.Ticks.Ticks(&x) {
		p := x.ToCanvas(tick.Value, 0, sz.X)
		if tick.Minor {
			canvas.Poly(Ps(p, y0, p, y1), minor)
		} else {
			canvas.Poly(Ps(p, y0, p, y1), major)
		}
	}

	for _, tick := range y.Ticks.Ticks(&y) {
		p := y.ToCanvas(tick.Value, 0, sz.Y)
		if tick.Minor {
			canvas.Poly(Ps(x0, p, x1, p), minor)
		} else {
			canvas.Poly(Ps(x0, p, x1, p), major)
		}
	}
}

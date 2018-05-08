package plot

import (
	"math"
)

type Bar struct {
	Style
	Label string

	DynamicWidth    bool
	DynamicMinWidth float64

	Data []Point
}

func NewBar(label string, points []Point) *Bar {
	return &Bar{
		Label:           label,
		Data:            points,
		DynamicMinWidth: 2,
	}
}

func (bar *Bar) iter(fn func(p Point, left, right float64)) {
	if !bar.DynamicWidth {
		for i, p := range bar.Data {
			fn(p, float64(i), float64(i+1))
		}
		return
	}

	left := 0.0
	for i, p := range bar.Data {
		if i+1 < len(bar.Data) {
			right := (p.X + bar.Data[i+1].X) * 0.5
			fn(p, left, right)
			left = right
		} else {
			right := left + p.X

			width := right - left
			if width < bar.DynamicMinWidth {
				width = left + bar.DynamicMinWidth
				right = left + width
			}

			fn(p, left, right)
			left = right
		}
	}
}

func (bar *Bar) Stats() Stats {
	stats := PointsStats(bar.Data)
	stats.Min.X = 0
	stats.Min.Y = 0
	if !bar.DynamicWidth {
		stats.Max.X = float64(len(bar.Data))
	} else {
		stats.Max.X = 0
		bar.iter(func(p Point, left, right float64) {
			stats.Max.X = right
		})
	}
	return stats
}

func (bar *Bar) Draw(plot *Plot, canvas Canvas) {
	x, y := plot.X, plot.Y
	size := canvas.Bounds().Size()
	canvas = canvas.Clip(canvas.Bounds())

	style := &bar.Style
	if style.IsZero() {
		style = &plot.Theme.Bar
	}

	lastScreenMin := 0.0
	lastScreenMax := 0.0
	bar.iter(func(p Point, left, right float64) {
		var r Rect
		r.Min.X = x.ToCanvas(left, 0, size.X)
		r.Max.X = x.ToCanvas(right, 0, size.X)
		r.Min.Y = y.ToCanvas(0, 0, size.Y)
		r.Max.Y = y.ToCanvas(p.Y, 0, size.Y)

		if bar.DynamicWidth && bar.DynamicMinWidth > 0 {
			leftToRight := r.Min.X < r.Max.X

			r.Min.X = math.Max(math.Max(r.Min.X, lastScreenMin), lastScreenMax)
			r.Max.X = math.Max(math.Max(r.Max.X, lastScreenMin), lastScreenMax)

			if leftToRight {
				if r.Max.X-r.Min.X < bar.DynamicMinWidth {
					r.Max.X = r.Min.X + bar.DynamicMinWidth
				}
			} else {
				if r.Min.X-r.Max.X < bar.DynamicMinWidth {
					r.Min.X = r.Max.X + bar.DynamicMinWidth
				}
			}

			lastScreenMin = r.Min.X
			lastScreenMax = r.Max.X
		}

		canvas.Rect(r, style)
	})
}

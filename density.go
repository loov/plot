package plot

import (
	"math"
	"sort"
)

type Density struct {
	Style
	Label  string
	Kernel Length
	Data   []float64 // sorted
}

func NewDensity(label string, values []float64) *Density {
	data := append(values[:0:0], values...)
	sort.Float64s(data)
	return &Density{
		Kernel: math.NaN(),
		Label:  label,
		Data:   data,
	}
}

func (plot *Density) Stats() Stats {
	min, avg, max := math.NaN(), math.NaN(), math.NaN()

	n := len(plot.Data)
	if n > 0 {
		min = plot.Data[0]
		avg = plot.Data[n/2]
		max = plot.Data[n-1]
	}

	return Stats{
		DiscreteX: true,
		DiscreteY: true,

		Min:    Point{min, 0},
		Center: Point{avg, 0.2}, // todo, figure out how to get the 50% of density plot
		Max:    Point{max, 1},
	}
}

func (line *Density) Draw(plot *Plot, canvas Canvas) {
	x, y := plot.X, plot.Y

	size := canvas.Bounds().Size()

	kernel := line.Kernel
	if math.IsNaN(kernel) {
		kernel = 4
	}

	xmin, xmax := x.ToCanvas(x.Min, 0, size.X), x.ToCanvas(x.Max, 0, size.X)

	points := []Point{}
	if line.Fill != nil {
		points = append(points, Point{xmin, 0})
	}

	index := 0
	previousLow := math.Inf(-1)
	maxy := 0.0
	for screenX := 0.0; screenX < size.X; screenX += 0.5 {
		// at := x.FromCanvas(screenX, 0, size.X)
		low := x.FromCanvas(screenX-kernel, 0, size.X)
		high := x.FromCanvas(screenX+kernel, 0, size.X)
		if low < previousLow {
			index = sort.SearchFloat64s(line.Data, low)
		} else {
			for ; index < len(line.Data); index++ {
				if line.Data[index] >= low {
					break
				}
			}
		}
		previousLow = low

		center := (low + high) / 2
		valueKernel := (high - low) / 2
		invValueKernel := 1 / valueKernel

		sample := 0.0
		for _, value := range line.Data[index:] {
			if value > high {
				break
			}
			sample += cubicPulse(center, 2, invValueKernel, value)
		}

		maxy = math.Max(maxy, sample)

		points = append(points, Point{
			X: screenX,
			Y: sample,
		})
	}

	if line.Fill != nil {
		points = append(points,
			Point{xmax, 0},
			Point{xmin, 0},
		)
	}

	scale := 1 / maxy
	for i := range points {
		points[i].Y = y.ToCanvas(points[i].Y*scale, 0, size.Y)
	}

	if !line.Style.IsZero() {
		canvas.Poly(points, &line.Style)
	} else {
		canvas.Poly(points, &plot.Theme.Line)
	}
}

func cubicPulse(center, radius, invradius, at float64) float64 {
	at = at - center
	if at < 0 {
		at = -at
	}
	if at > radius {
		return 0
	}
	at *= invradius
	return 1 - at*at*(3-2*at)
}

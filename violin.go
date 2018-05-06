package plot

import (
	"math"
	"sort"
)

type Violin struct {
	Style
	Side   float64
	Label  string
	Kernel Length
	Data   []float64 // sorted
}

func NewViolin(label string, values []float64) *Violin {
	data := append(values[:0:0], values...)
	sort.Float64s(data)
	return &Violin{
		Kernel: math.NaN(),
		Side:   1,
		Label:  label,
		Data:   data,
	}
}

func (line *Violin) Stats() Stats {
	min, avg, max := math.NaN(), math.NaN(), math.NaN()

	n := len(line.Data)
	if n > 0 {
		min = line.Data[0]
		avg = line.Data[n/2]
		max = line.Data[n-1]
	}

	return Stats{
		DiscreteX: true,
		DiscreteY: true,

		Min:    Point{-1, min},
		Center: Point{0, avg}, // todo, figure out how to get the 50% of density plot
		Max:    Point{1, max},
	}
}

func (line *Violin) Draw(plot *Plot, canvas Canvas) {
	x, y := plot.X, plot.Y

	size := canvas.Bounds().Size()

	kernel := line.Kernel
	if math.IsNaN(kernel) {
		kernel = 4
	}

	ymin, ymax := y.ToCanvas(y.Min, 0, size.Y), y.ToCanvas(y.Max, 0, size.Y)
	if ymin > ymax {
		ymin, ymax = ymax, ymin
	}

	points := []Point{}
	if line.Fill != nil || line.Side == 0 {
		points = append(points, Point{0, ymin})
	}

	index := 0
	previousLow := math.Inf(-1)
	maxx := 0.0
	for screenY := 0.0; screenY < size.Y; screenY += 0.5 {
		low := y.FromCanvas(screenY-kernel, 0, size.Y)
		high := y.FromCanvas(screenY+kernel, 0, size.Y)

		if low > high {
			high, low = low, high
		}

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
		maxx = math.Max(maxx, sample)

		points = append(points, Point{
			X: sample,
			Y: screenY,
		})
	}

	if line.Fill != nil {
		points = append(points, Point{0, ymax})
	}

	if line.Side == 0 {
		scale := 1 / maxx
		otherSide := make([]Point, len(points))
		for i := range points {
			k := len(points) - i - 1
			otherSide[k] = points[i]
			points[i].X = x.ToCanvas(points[i].X*scale, 0, size.X)
			otherSide[k].X = x.ToCanvas(-otherSide[k].X*scale, 0, size.X)
		}
		points = append(points, otherSide...)
	} else {
		scale := line.Side * 1 / maxx
		for i := range points {
			points[i].X = x.ToCanvas(points[i].X*scale, 0, size.X)
		}
	}

	_, _, _ = ymin, ymax, maxx

	if !line.Style.IsZero() {
		canvas.Poly(points, &line.Style)
	} else {
		canvas.Poly(points, &plot.Theme.Line)
	}
}

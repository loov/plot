package main

import (
	"image/color"
	"io/ioutil"
	"math"
	"math/rand"

	"github.com/loov/plot"
)

func main() {
	type Dataset struct {
		Red   []float64
		Green []float64
		Blue  []float64
	}

	datasets := []*Dataset{}

	r := func() float64 {
		return rand.Float64()*3 - 1.5
	}

	for i := 0; i < 3; i++ {
		const N = 1 << 20
		dataset := &Dataset{
			Red:   make([]float64, N),
			Green: make([]float64, N),
			Blue:  make([]float64, N),
		}
		datasets = append(datasets, dataset)

		or, og, ob := r(), r(), r()
		for i := 0; i < N; i++ {
			dataset.Red[i] = 20 * (or + r())
			dataset.Green[i] = 20 * (og + r()*r())
			dataset.Blue[i] = 20 * (ob + r()*math.Sin(r()))
		}
	}

	p := plot.New()
	stack := plot.NewVStack()
	p.Add(stack)
	for _, dataset := range datasets {
		red := plot.NewDensity("Red", dataset.Red)
		red.Stroke = color.NRGBA{200, 0, 0, 255}
		red.Fill = color.NRGBA{200, 0, 0, 40}

		green := plot.NewDensity("Green", dataset.Green)
		green.Stroke = color.NRGBA{0, 200, 0, 255}
		green.Fill = color.NRGBA{0, 200, 0, 40}

		blue := plot.NewDensity("Blue", dataset.Blue)
		blue.Stroke = color.NRGBA{0, 0, 200, 255}
		blue.Fill = color.NRGBA{0, 0, 200, 40}

		stack.AddGroup(
			plot.NewGrid(),
			plot.NewGizmo(),
			red, green, blue,
			plot.NewTickLabels(),
		)
	}

	// p.X.SetLogarithmic(float64(10))
	svg := plot.NewSVG(800, float64(100*len(datasets)))
	p.Draw(svg)
	ioutil.WriteFile("result.svg", svg.Bytes(), 0755)
}

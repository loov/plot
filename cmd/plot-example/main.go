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
		const N = 1 << 15
		dataset := &Dataset{
			Red:   make([]float64, N),
			Green: make([]float64, N),
			Blue:  make([]float64, N),
		}
		datasets = append(datasets, dataset)

		{
			or, og, ob := r(), r(), r()
			for k := 0; k < N; k++ {
				dataset.Red[k] = 20 * (or + r())
				dataset.Green[k] = 20 * (og + r()*r())
				dataset.Blue[k] = 20 * (ob + r()*math.Sin(r()))
			}
		}
	}

	{ // density plot
		p := plot.New()
		stack := plot.NewVStack()
		stack.Margin = plot.R(0, 5, 0, 5)
		p.Add(stack)
		for _, dataset := range datasets {
			red := plot.NewDensity("Red", dataset.Red)
			red.Class = "red"
			red.Stroke = color.NRGBA{200, 0, 0, 255}
			red.Fill = color.NRGBA{200, 0, 0, 40}

			green := plot.NewDensity("Green", dataset.Green)
			green.Class = "green"
			green.Stroke = color.NRGBA{0, 200, 0, 255}
			green.Fill = color.NRGBA{0, 200, 0, 40}

			blue := plot.NewDensity("Blue", dataset.Blue)
			blue.Class = "blue"
			blue.Stroke = color.NRGBA{0, 0, 200, 255}
			blue.Fill = color.NRGBA{0, 0, 200, 40}

			stack.AddGroup(
				plot.NewGrid(),
				plot.NewGizmo(),
				red, green, blue,
				plot.NewTickLabels(),
			)
		}

		svg := plot.NewSVG(800, float64(150*len(datasets)))
		p.Draw(svg)
		ioutil.WriteFile("density.svg", svg.Bytes(), 0755)
	}

	{ // violin plot
		p := plot.New()
		stack := plot.NewHStack()
		stack.Margin = plot.R(5, 0, 5, 0)
		p.Add(stack)
		for _, dataset := range datasets {
			red := plot.NewViolin("Red", dataset.Red)
			red.Side = 0
			red.Class = "red"
			red.Stroke = color.NRGBA{200, 0, 0, 255}
			red.Fill = color.NRGBA{200, 0, 0, 40}

			green := plot.NewViolin("Green", dataset.Green)
			green.Side = 1
			green.Class = "green"
			green.Stroke = color.NRGBA{0, 200, 0, 255}
			green.Fill = color.NRGBA{0, 200, 0, 40}

			blue := plot.NewViolin("Blue", dataset.Blue)
			blue.Side = -1
			blue.Class = "blue"
			blue.Stroke = color.NRGBA{0, 0, 200, 255}
			blue.Fill = color.NRGBA{0, 0, 200, 40}

			stack.AddGroup(
				plot.NewGrid(),
				plot.NewGizmo(),
				red, green, blue,
				plot.NewTickLabels(),
			)
		}

		svg := plot.NewSVG(800, float64(150*len(datasets)))
		p.Draw(svg)
		ioutil.WriteFile("violin.svg", svg.Bytes(), 0755)
	}

	{ // percentile plot
		p := plot.New()
		p.X.Transform = plot.NewPercentileTransform(5)
		p.X.Ticks = plot.ManualTicks{
			{Value: 0, Label: "0"},
			{Value: 0.25, Label: "25"},
			{Value: 0.5, Label: "50"},
			{Value: 0.75, Label: "75"},
			{Value: 0.9, Label: "90"},
			{Value: 0.99, Label: "99"},
			{Value: 0.999, Label: "99.9"},
			{Value: 0.9999, Label: "99.99"},
			{Value: 0.99999, Label: "99.999"}}

		stack := plot.NewVStack()
		stack.Margin = plot.R(0, 5, 0, 5)

		p.Add(stack)
		for _, dataset := range datasets {
			red := plot.NewPercentiles("Red", dataset.Red)
			red.Class = "red"
			red.Stroke = color.NRGBA{200, 0, 0, 255}

			green := plot.NewPercentiles("Green", dataset.Green)
			green.Class = "green"
			green.Stroke = color.NRGBA{0, 200, 0, 255}

			blue := plot.NewPercentiles("Blue", dataset.Blue)
			blue.Class = "blue"
			blue.Stroke = color.NRGBA{0, 0, 200, 255}

			stack.AddGroup(
				plot.NewGrid(),
				plot.NewGizmo(),
				red, green, blue,
				plot.NewTickLabels(),
			)
		}

		svg := plot.NewSVG(800, float64(150*len(datasets)))
		p.Draw(svg)
		ioutil.WriteFile("percentiles.svg", svg.Bytes(), 0755)
	}

}

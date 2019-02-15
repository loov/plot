package main

import (
	"image/color"
	"io/ioutil"
	"math"
	"math/rand"
	"strconv"

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
		for i, dataset := range datasets {
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
				plot.NewXLabel("Case "+strconv.Itoa(i+1)),
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
		for i, dataset := range datasets {
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
				plot.NewXLabel("Case "+strconv.Itoa(i+1)),
			)
		}

		svg := plot.NewSVG(800, float64(150*len(datasets)))
		p.Draw(svg)
		ioutil.WriteFile("violin.svg", svg.Bytes(), 0755)
	}

	{ // percentile plot
		p := plot.New()
		p.X = plot.NewPercentilesAxis()

		stack := plot.NewVStack()
		stack.Margin = plot.R(0, 5, 0, 5)

		p.Add(stack)
		for i, dataset := range datasets {
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
				plot.NewXLabel("Case "+strconv.Itoa(i+1)),
			)
		}

		svg := plot.NewSVG(800, float64(150*len(datasets)))
		p.Draw(svg)
		ioutil.WriteFile("percentiles.svg", svg.Bytes(), 0755)
	}

	{ // line plot
		p := plot.New()
		stack := plot.NewHStack()
		stack.Margin = plot.R(5, 0, 5, 0)
		p.Add(stack)

		sizes := []int{1, 2, 4, 8, 1024, 8196}
		p.X.Transform = plot.NewLog1pTransform(2)
		p.X.Ticks = plot.ManualTicks{
			{Value: 1, Label: "1"},
			{Value: 4, Label: "4"},
			{Value: 1024, Label: "1024"},
			{Value: 8196, Label: "8196"},
		}

		for i, _ := range datasets {
			values := make([]int, len(sizes))
			prev := 0
			for i := range values {
				values[i] = prev + i*rand.Intn(10) + 100
				prev = values[i]
			}

			sizesf := plot.IntsToFloat64s(sizes)
			valuesf := plot.IntsToFloat64s(values)
			nanos := plot.NewLine("", plot.Points(sizesf, valuesf))
			stack.AddGroup(
				plot.NewGrid(),
				plot.NewGizmo(),
				nanos,
				plot.NewTickLabels(),
				plot.NewXLabel("Case "+strconv.Itoa(i+1)),
			)
		}

		svg := plot.NewSVG(800, float64(150*len(datasets)))
		p.Draw(svg)
		ioutil.WriteFile("line-stack.svg", svg.Bytes(), 0755)
	}

	{ // bar plot
		p := plot.New()
		stack := plot.NewHStack()
		stack.Margin = plot.R(5, 0, 5, 0)
		p.Add(stack)

		sizes := []int{1, 2, 4, 8, 1024, 8196}
		p.X.Ticks = plot.ManualTicks{
			{Value: 1, Label: "1"},
			{Value: 4, Label: "4"},
			{Value: 1024, Label: "1024"},
			{Value: 8196, Label: "8196"},
		}

		bars := []*plot.Bar{}

		for i, _ := range datasets {
			values := make([]int, len(sizes))
			prev := 0
			for i := range values {
				values[i] = prev + i*rand.Intn(10) + 100
				prev = values[i]
			}

			sizesf := plot.IntsToFloat64s(sizes)
			valuesf := plot.IntsToFloat64s(values)
			nanos := plot.NewBar("", plot.Points(sizesf, valuesf))
			bars = append(bars, nanos)

			stack.AddGroup(
				plot.NewGrid(),
				plot.NewGizmo(),
				nanos,
				plot.NewTickLabels(),
				plot.NewXLabel("Case "+strconv.Itoa(i+1)),
			)
		}

		svg := plot.NewSVG(800, float64(150*len(datasets)))
		p.Draw(svg)
		ioutil.WriteFile("bar-chart.svg", svg.Bytes(), 0755)

		for _, bar := range bars {
			bar.DynamicWidth = true
		}
		//p.X.Transform = plot.NewLog1pTransform(2)

		svg = plot.NewSVG(800, float64(150*len(datasets)))
		p.Draw(svg)
		ioutil.WriteFile("bar-chart-dynamic.svg", svg.Bytes(), 0755)
	}
}

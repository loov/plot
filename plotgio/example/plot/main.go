package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"github.com/loov/plot"
	"github.com/loov/plot/plotgio"
)

func main() {
	go func() {
		w := app.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	const N = 1 << 15
	var datasets Datasets
	for i := 0; i < 3; i++ {
		datasets = append(datasets, NewDataset(N))
	}

	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			layout.Flex{}.Layout(gtx,
				layout.Flexed(1, fill(datasets.DensityPlot)),
				//layout.Flexed(1, fill(datasets.ViolinPlot)),
				//layout.Flexed(1, fill(datasets.PercentilesPlot)),
				//layout.Flexed(1, fill(datasets.LinePlot)),
				//layout.Flexed(1, fill(datasets.BarPlot)),
			)

			e.Frame(gtx.Ops)
		}
	}
}

func fill(fn func(size f32.Point, gtx layout.Context)) func(gtx layout.Context) layout.Dimensions {
	return func(gtx layout.Context) layout.Dimensions {
		size := layout.FPt(gtx.Constraints.Max)
		fn(size, gtx)
		return layout.Dimensions{
			Size: gtx.Constraints.Max,
		}
	}
}

type Datasets []*Dataset

type Dataset struct {
	Red   []float64
	Green []float64
	Blue  []float64
}

func NewDataset(size int) *Dataset {
	dataset := &Dataset{
		Red:   make([]float64, size),
		Green: make([]float64, size),
		Blue:  make([]float64, size),
	}

	r := func() float64 {
		return rand.Float64()*3 - 1.5
	}

	or, og, ob := r(), r(), r()
	for k := 0; k < size; k++ {
		dataset.Red[k] = 20 * (or + r())
		dataset.Green[k] = 20 * (og + r()*r())
		dataset.Blue[k] = 20 * (ob + r()*math.Sin(r()))
	}
	return dataset
}

var defaultMargin = plot.R(20, 20, 20, 20)

func (datasets Datasets) DensityPlot(size f32.Point, gtx layout.Context) {
	p := plot.New()
	stack := plot.NewVStack()
	stack.Margin = defaultMargin
	p.Add(stack)
	for i, dataset := range datasets {
		red := plot.NewDensity("Red", dataset.Red)
		red.Class = "red"
		red.Stroke = color.NRGBA{200, 0, 0, 255}
		red.Fill = color.NRGBA{200, 0, 0, 120}

		green := plot.NewDensity("Green", dataset.Green)
		green.Class = "green"
		green.Stroke = color.NRGBA{0, 200, 0, 255}
		green.Fill = color.NRGBA{0, 200, 0, 120}

		blue := plot.NewDensity("Blue", dataset.Blue)
		blue.Class = "blue"
		blue.Stroke = color.NRGBA{0, 0, 200, 255}
		blue.Fill = color.NRGBA{0, 0, 200, 120}

		stack.AddGroup(
			plot.NewGrid(),
			plot.NewGizmo(),
			red, green, blue,
			plot.NewTickLabels(),
			plot.NewXLabel("Case "+strconv.Itoa(i+1)),
		)
	}

	canvas := plotgio.New(size)
	p.Draw(canvas)
	canvas.Add(gtx.Ops)
}

func (datasets Datasets) ViolinPlot(size f32.Point, gtx layout.Context) {
	p := plot.New()
	stack := plot.NewHStack()
	stack.Margin = defaultMargin
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

	canvas := plotgio.New(size)
	p.Draw(canvas)
	canvas.Add(gtx.Ops)
}

func (datasets Datasets) PercentilesPlot(size f32.Point, gtx layout.Context) {
	p := plot.New()
	p.X = plot.NewPercentilesAxis()

	stack := plot.NewVStack()
	stack.Margin = defaultMargin

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

	canvas := plotgio.New(size)
	p.Draw(canvas)
	canvas.Add(gtx.Ops)
}

func (datasets Datasets) LinePlot(size f32.Point, gtx layout.Context) {
	p := plot.New()
	stack := plot.NewHStack()
	stack.Margin = defaultMargin
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

		labelsLeftBottom := plot.NewTickLabels()
		labelsTopRight := plot.NewTickLabels()
		labelsTopRight.X.Side = 1
		labelsTopRight.Y.Side = 1
		stack.AddGroup(
			plot.NewGrid(),
			plot.NewGizmo(),
			nanos,
			labelsLeftBottom,
			labelsTopRight,
			plot.NewXLabel("Case "+strconv.Itoa(i+1)),
		)
	}

	canvas := plotgio.New(size)
	p.Draw(canvas)
	canvas.Add(gtx.Ops)
}

func (datasets Datasets) BarPlot(size f32.Point, gtx layout.Context) {
	p := plot.New()
	stack := plot.NewHStack()
	stack.Margin = defaultMargin
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

	canvas := plotgio.New(size)
	p.Draw(canvas)
	canvas.Add(gtx.Ops)
}

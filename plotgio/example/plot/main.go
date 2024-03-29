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
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/loov/plot"

	"github.com/loov/plot/plotgio"
)

func main() {
	go func() {
		w := app.NewWindow(app.Size(240*5, 240))
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

	shaper := text.NewShaper(gofont.Collection())

	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			layout.Flex{}.Layout(gtx,
				layout.Flexed(1, fill(shaper, datasets.DensityPlot)),
				layout.Flexed(1, fill(shaper, datasets.ViolinPlot)),
				layout.Flexed(1, fill(shaper, datasets.PercentilesPlot)),
				layout.Flexed(1, fill(shaper, datasets.LinePlot)),
				layout.Flexed(1, fill(shaper, datasets.BarPlot)),
			)

			e.Frame(gtx.Ops)
		}
	}
}

func fill(shaper *text.Shaper, fn func(size f32.Point, shaper *text.Shaper, gtx layout.Context)) func(gtx layout.Context) layout.Dimensions {
	return func(gtx layout.Context) layout.Dimensions {
		const pad = 3
		return widget.Border{
			Color: color.NRGBA{R: 0xC0, G: 0xFF, B: 0xC0, A: 0xff},
			Width: pad,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(pad).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				fn(layout.FPt(gtx.Constraints.Max), shaper, gtx)
				return layout.Dimensions{Size: gtx.Constraints.Max}
			})
		})
	}
}

type Datasets []*Dataset

type Dataset struct {
	Red   []float64
	Green []float64
	Blue  []float64

	Sizes  []float64
	Values []float64
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

	dataset.Sizes = []float64{1, 2, 4, 8, 1024, 8196}
	dataset.Values = make([]float64, len(dataset.Sizes))
	prev := float64(0)
	for i := range dataset.Values {
		dataset.Values[i] = prev + float64(i*rand.Intn(10)) + 100
		prev = dataset.Values[i]
	}

	return dataset
}

var defaultMargin = plot.R(20, 20, 20, 20)

func (datasets Datasets) DensityPlot(size f32.Point, shaper *text.Shaper, gtx layout.Context) {
	p := plot.New()
	stack := plot.NewVStack()
	stack.Margin = defaultMargin
	p.Add(stack)
	for i, dataset := range datasets {
		red := plot.NewDensity("Red", dataset.Red)
		red.Size = 1
		red.Stroke = color.NRGBA{200, 0, 0, 255}
		red.Fill = color.NRGBA{200, 0, 0, 120}

		green := plot.NewDensity("Green", dataset.Green)
		green.Size = 1
		green.Stroke = color.NRGBA{0, 200, 0, 255}
		green.Fill = color.NRGBA{0, 200, 0, 120}

		blue := plot.NewDensity("Blue", dataset.Blue)
		blue.Size = 1
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

	canvas := plotgio.New(shaper, size)
	p.Draw(canvas)
	canvas.Add(gtx)
}

func (datasets Datasets) ViolinPlot(size f32.Point, shaper *text.Shaper, gtx layout.Context) {
	p := plot.New()
	stack := plot.NewHStack()
	stack.Margin = defaultMargin
	p.Add(stack)
	for i, dataset := range datasets {
		red := plot.NewViolin("Red", dataset.Red)
		red.Size = 1
		red.Side = 0
		red.Stroke = color.NRGBA{200, 0, 0, 255}
		red.Fill = color.NRGBA{200, 0, 0, 120}

		green := plot.NewViolin("Green", dataset.Green)
		green.Size = 1
		green.Side = 1
		green.Stroke = color.NRGBA{0, 200, 0, 255}
		green.Fill = color.NRGBA{0, 200, 0, 120}

		blue := plot.NewViolin("Blue", dataset.Blue)
		blue.Size = 1
		blue.Side = -1
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

	canvas := plotgio.New(shaper, size)
	p.Draw(canvas)
	canvas.Add(gtx)
}

func (datasets Datasets) PercentilesPlot(size f32.Point, shaper *text.Shaper, gtx layout.Context) {
	p := plot.New()
	p.X = plot.NewPercentilesAxis()

	stack := plot.NewVStack()
	stack.Margin = defaultMargin

	p.Add(stack)
	for i, dataset := range datasets {
		red := plot.NewPercentiles("Red", dataset.Red)
		red.Stroke = color.NRGBA{200, 0, 0, 255}
		red.Size = 1

		green := plot.NewPercentiles("Green", dataset.Green)
		green.Stroke = color.NRGBA{0, 200, 0, 255}
		green.Size = 1

		blue := plot.NewPercentiles("Blue", dataset.Blue)
		blue.Stroke = color.NRGBA{0, 0, 200, 255}
		blue.Size = 1

		stack.AddGroup(
			plot.NewGrid(),
			plot.NewGizmo(),
			red, green, blue,
			plot.NewTickLabels(),
			plot.NewXLabel("Case "+strconv.Itoa(i+1)),
		)
	}

	canvas := plotgio.New(shaper, size)
	p.Draw(canvas)
	canvas.Add(gtx)
}

func (datasets Datasets) LinePlot(size f32.Point, shaper *text.Shaper, gtx layout.Context) {
	p := plot.New()
	stack := plot.NewHStack()
	stack.Margin = defaultMargin
	p.Add(stack)

	p.X.Transform = plot.NewLog1pTransform(2)
	p.X.Ticks = plot.ManualTicks{
		{Value: 1, Label: "1"},
		{Value: 4, Label: "4"},
		{Value: 1024, Label: "1024"},
		{Value: 8196, Label: "8196"},
	}

	for i, dataset := range datasets {
		nanos := plot.NewLine("", plot.Points(dataset.Sizes, dataset.Values))

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

	canvas := plotgio.New(shaper, size)
	p.Draw(canvas)
	canvas.Add(gtx)
}

func (datasets Datasets) BarPlot(size f32.Point, shaper *text.Shaper, gtx layout.Context) {
	p := plot.New()
	stack := plot.NewHStack()
	stack.Margin = defaultMargin
	p.Add(stack)

	p.X.Ticks = plot.ManualTicks{
		{Value: 1, Label: "1"},
		{Value: 4, Label: "4"},
		{Value: 1024, Label: "1024"},
		{Value: 8196, Label: "8196"},
	}

	bars := []*plot.Bar{}

	for i, dataset := range datasets {
		nanos := plot.NewBar("", plot.Points(dataset.Sizes, dataset.Values))
		bars = append(bars, nanos)

		stack.AddGroup(
			plot.NewGrid(),
			plot.NewGizmo(),
			nanos,
			plot.NewTickLabels(),
			plot.NewXLabel("Case "+strconv.Itoa(i+1)),
		)
	}

	canvas := plotgio.New(shaper, size)
	p.Draw(canvas)
	canvas.Add(gtx)
}

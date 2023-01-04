package main

import (
	"image/color"
	"log"
	"math/rand"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"github.com/loov/plot"

	"github.com/loov/plot/plotgio"
)

func main() {
	go func() {
		w := app.NewWindow(app.Size(400, 300))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	shaper := text.NewShaper(gofont.Collection())

	dataset := &Dataset{
		Shaper: shaper,
	}

	tick := time.NewTicker(30 * time.Millisecond)

	var ops op.Ops
	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)

				dataset.Layout(gtx)

				e.Frame(gtx.Ops)
			}
		case <-tick.C:
			dataset.AddRandomValue()
			w.Invalidate()
		}
	}
}

type Dataset struct {
	Shaper *text.Shaper
	Values [200]float64

	points []plot.Point
}

func (display *Dataset) AddRandomValue() {
	p := display.Values[len(display.Values)-1]
	copy(display.Values[:], display.Values[1:])
	display.Values[len(display.Values)-1] = 0.8*p + 0.2*rand.Float64()*50
}

func (display *Dataset) Layout(gtx layout.Context) layout.Dimensions {
	size := layout.FPt(gtx.Constraints.Max)
	display.layoutGraph(gtx, size)
	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}

var defaultMargin = plot.R(20, 20, 20, 20)

func (display *Dataset) layoutGraph(gtx layout.Context, size f32.Point) {
	p := plot.New()
	p.Y.Min, p.Y.Max = 0, 100

	stack := plot.NewVStack()
	stack.Margin = defaultMargin
	p.Add(stack)

	if len(display.points) != len(display.Values) {
		display.points = make([]plot.Point, len(display.Values))
	}
	for i, v := range display.Values {
		display.points[i] = plot.P(float64(i), v)
	}

	line := plot.NewLine("Red", display.points)
	line.Size = 2
	line.Stroke = color.NRGBA{200, 0, 0, 255}

	stack.AddGroup(
		plot.NewGrid(),
		plot.NewGizmo(),
		line,
		plot.NewTickLabels(),
		plot.NewXLabel("Random Data"),
	)

	canvas := plotgio.New(display.Shaper, size)
	p.Draw(canvas)
	canvas.Add(gtx)
}

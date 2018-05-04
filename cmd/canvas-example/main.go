package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"math"

	"github.com/kr/pretty"
	"github.com/loov/plot"
)

func main() {
	black := &plot.Style{Stroke: color.NRGBA{0, 0, 0, 255}}
	textColor := &plot.Style{Fill: color.NRGBA{0, 0, 0, 255}}
	interlace := &plot.Style{Stroke: color.NRGBA{230, 210, 210, 255}}
	fill := &plot.Style{Fill: color.NRGBA{0, 0, 0, 10}}

	canvas := plot.NewSVG(800, 600)

	clip := canvas.Bounds().Shrink(plot.P(30, 30))
	canvas.Rect(clip, fill)
	area := canvas.Clip(clip)

	underlay := area.Layer(-1)
	graphic := area.Layer(1)
	text := area.Layer(2)

	points := []plot.Point{}
	size := area.Bounds().Size()
	step := 0
	for x := 0.0; x < size.X; x++ {
		p := plot.P(x, size.Y/2+math.Sin(x*8/size.X)*size.Y/2)
		points = append(points, p)

		if step%10 == 0 {
			textColor.Fill = color.NRGBA{byte(step), 0, 0, 255}
			text.Text(fmt.Sprintf("%.1f", x), p, textColor)
		}
		if step%100 == 0 {
			underlay.Poly(plot.Ps(x, 0, x, size.Y), interlace)
		}
		step++
	}
	graphic.Poly(points, black)

	points = []plot.Point{}
	for x := 0.0; x < size.X; x++ {
		p := plot.P(x, size.Y/2+math.Cos(x*16/size.X)*size.Y/3)
		points = append(points, p)
	}
	graphic.Poly(points, &plot.Style{
		Stroke: color.NRGBA{0, 0, 0, 255},
		Dash:   []plot.Length{1, 2, 3},
	})

	pretty.Print(canvas)

	ioutil.WriteFile("example.svg", canvas.Bytes(), 0755)
}

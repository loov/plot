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
	canvas := plot.NewSVG(800, 600)

	area := canvas.Clip(canvas.Bounds().Shrink(plot.P(30, 30)))

	underlay := area.Layer(-1)
	graphic := area.Layer(1)
	text := area.Layer(2)

	black := &plot.Style{
		Color: color.NRGBA{0, 0, 0, 255},
	}

	fill := &plot.Style{
		Fill: color.NRGBA{0, 0, 0, 10},
	}

	points := []plot.Point{}
	size := area.Bounds().Size()
	step := 0
	for x := 0.0; x < size.X; x++ {
		p := plot.P(x, size.Y/2+math.Sin(x*8/size.X)*size.Y/2)
		points = append(points, p)

		if step%10 == 0 {
			text.Text(fmt.Sprintf("%.1f", x), p, black)
		}
		step++
	}

	graphic.Poly(points, black)
	underlay.Rect(0, 0, size.X, size.Y, fill)

	pretty.Print(canvas)

	ioutil.WriteFile("example.svg", canvas.Bytes(), 0755)
}

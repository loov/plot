package main

import (
	"image/color"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/loov/plot"
)

func main() {
	xs := make([]time.Duration, 1<<20)
	ys := make([]time.Duration, 1<<20)
	for i := range xs {
		xs[i] = time.Duration(rand.Int63n(20) * rand.Int63n(20))
		ys[i] = time.Duration(rand.Int63n(20)*rand.Int63n(20) + 50)
	}

	p := plot.New()

	p.Add(plot.NewGrid())

	red := plot.NewDensity("Red", plot.DurationToNanoseconds(xs))
	red.Stroke = color.NRGBA{200, 0, 0, 255}
	red.Fill = color.NRGBA{200, 0, 0, 10}
	p.Add(red)

	green := plot.NewDensity("Green", plot.DurationToNanoseconds(ys))
	green.Stroke = color.NRGBA{0, 200, 0, 255}
	green.Fill = color.NRGBA{0, 255, 0, 10}
	p.Add(green)

	svg := plot.NewSVG(600, 300)
	p.Draw(svg)

	ioutil.WriteFile("result.svg", svg.Bytes(), 0755)
}

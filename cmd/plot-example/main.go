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
		xs[i] = time.Duration(rand.Int63n(100) * rand.Int63n(100))
		ys[i] = time.Duration(rand.Int63n(100)*rand.Int63n(100) + 50)
	}

	p := plot.New()

	red := plot.NewDensity("Red", plot.DurationToNanoseconds(xs))
	red.Color = color.NRGBA{255, 0, 0, 255}
	p.Add(red)

	green := plot.NewDensity("Green", plot.DurationToNanoseconds(ys))
	green.Color = color.NRGBA{0, 255, 0, 255}
	p.Add(green)

	svg := plot.NewSVG(600, 300)
	p.Draw(svg)

	ioutil.WriteFile("result.svg", svg.Bytes(), 0755)
}

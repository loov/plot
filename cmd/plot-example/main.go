package main

import (
	"fmt"
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

	// p.Y.SetLogarithmic(-50)
	p.Add(plot.NewGrid())

	red := plot.NewDensity("Red", plot.DurationToNanoseconds(xs))
	red.Stroke = color.NRGBA{200, 0, 0, 255}
	red.Fill = color.NRGBA{200, 0, 0, 40}
	p.Add(red)

	green := plot.NewDensity("Green", plot.DurationToNanoseconds(ys))
	green.Stroke = color.NRGBA{0, 200, 0, 255}
	green.Fill = color.NRGBA{0, 200, 0, 40}
	p.Add(green)

	for x := -30; x < 60; x += 10 {
		p.X.SetLogarithmic(float64(x))
		svg := plot.NewSVG(600, 300)
		p.Draw(svg)
		ioutil.WriteFile(fmt.Sprintf("result_%v.svg", x), svg.Bytes(), 0755)
	}

}

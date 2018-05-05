package plot

import (
	"image/color"
	"math"
)

type Plot struct {
	// Size is the visual size of this plot
	Size Point
	// X, Y are the axis information
	X, Y Axis
	// Elements
	Elements []Element
	// DefaultStyle
	Line Style
	Font Style
	Fill Style
}

// Element is a drawable plot element
type Element interface {
	Draw(plot *Plot, canvas Canvas)
}

// Dataset represents an Element that contains data
type Dataset interface {
	Element
	// TODO: remove and replace with recommended Axis
	Stats() Stats
}

type Stats struct {
	DiscreteX, DiscreteY bool

	Min    Point
	Center Point
	Max    Point
}

func New() *Plot {
	return &Plot{
		Size: Point{800, 600},
		X: Axis{
			Min: math.NaN(),
			Max: math.NaN(),
		},
		Y: Axis{
			Min: math.NaN(),
			Max: math.NaN(),
		},
		Line: Style{
			Stroke: color.NRGBA{0, 0, 0, 255},
			Fill:   nil,
			Size:   1.0,
		},
		Font: Style{
			Stroke: nil,
			Fill:   color.NRGBA{0, 0, 0, 255},
			Size:   0,
		},
		Fill: Style{
			Stroke: nil,
			Fill:   color.NRGBA{255, 255, 255, 255},
			Size:   1.0,
		},
	}
}

func (plot *Plot) Add(element ...Element) {
	plot.Elements = append(plot.Elements, element...)
}

func (plot *Plot) Draw(canvas Canvas) {
	if !plot.X.IsValid() || !plot.Y.IsValid() {
		tmpplot := &Plot{}
		*tmpplot = *plot
		plot = tmpplot
		plot.X, plot.Y = detectAxis(plot.X, plot.Y, plot.Elements)
	}

	for _, element := range plot.Elements {
		element.Draw(plot, canvas)
	}
}

func (plot *Plot) ToCanvas(x, y float64) Point {
	var p Point
	px := (x - plot.X.Min) / (plot.X.Max - plot.X.Min)
	if plot.X.Transform != nil {
		px = plot.X.Transform(px)
	}
	p.X = Length(px) * plot.Size.X

	py := (y - plot.Y.Min) / (plot.Y.Max - plot.Y.Min)
	if plot.Y.Transform != nil {
		py = plot.Y.Transform(py)
	}
	p.Y = Length(py) * plot.Size.Y

	return p
}

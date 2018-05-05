package plot

import (
	"image/color"
)

type Plot struct {
	// X, Y are the axis information
	X, Y *Axis
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
	x, y := NewAxis(), NewAxis()
	y.Flip = true

	return &Plot{
		X: x,
		Y: y,
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
		element.Draw(plot, canvas.Context(canvas.Bounds()))
	}
}

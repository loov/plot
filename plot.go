package plot

type Plot struct {
	// X, Y are the axis information
	X, Y *Axis
	// Elements
	Elements
	// DefaultStyle
	Theme
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

func New() *Plot {
	x, y := NewAxis(), NewAxis()
	y.Flip = true

	return &Plot{
		X:     x,
		Y:     y,
		Theme: NewTheme(),
	}
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

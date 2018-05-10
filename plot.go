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

type AxisGroup struct {
	X, Y *Axis
	Elements
}

func NewAxisGroup() *AxisGroup {
	x, y := NewAxis(), NewAxis()
	y.Flip = true
	return &AxisGroup{
		X: x,
		Y: y,
	}
}

func (group *AxisGroup) Update() {
	tx, ty := detectAxis(group.X, group.Y, group.Elements)
	*group.X = *tx
	*group.Y = *ty
}

func (group *AxisGroup) Draw(plot *Plot, canvas Canvas) {
	tmpplot := &Plot{}
	*tmpplot = *plot

	if group.X != nil {
		tmpplot.X = group.X
	}
	if group.Y != nil {
		tmpplot.Y = group.Y
	}

	for _, element := range group.Elements {
		element.Draw(tmpplot, canvas.Context(canvas.Bounds()))
	}
}

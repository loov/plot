package plotgio

import (
	"image/color"
	"sort"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/x/stroke"
	"github.com/loov/plot"
)

var _ plot.Canvas = (*Canvas)(nil)

// Canvas describes the top-level ptx drawing context.
type Canvas struct {
	Shaper text.Shaper

	context
}

// New creates a new SVG canvas.
func New(shaper text.Shaper, size f32.Point) *Canvas {
	ptx := &Canvas{Shaper: shaper}
	ptx.bounds.Max.X = plot.Length(size.X)
	ptx.bounds.Max.Y = plot.Length(size.Y)
	return ptx
}

// context describes a ptx drawing context.
type context struct {
	index int
	clip  bool
	// bounds relative to parent
	bounds   plot.Rect
	elements []element
	layers   []*context
}

// element describes an ptx element.
type element struct {
	// style
	style plot.Style
	// line
	points []plot.Point
	// text
	text   string
	origin plot.Point
	// context
	context *context
}

// Bounds returns the bounds in the global size.
func (ptx *context) Bounds() plot.Rect {
	return ptx.bounds.Zero()
}

// Size returns the size of the context.
func (ptx *context) Size() plot.Point {
	return ptx.bounds.Size()
}

// context creates a clipped subcontext.
func (ptx *context) context(r plot.Rect, clip bool) plot.Canvas {
	element := element{}
	element.context = &context{}
	element.context.clip = clip
	element.context.bounds = r
	ptx.elements = append(ptx.elements, element)
	return element.context
}

// context creates a subcontext bounded to r.
func (ptx *context) Context(r plot.Rect) plot.Canvas {
	return ptx.context(r, false)
}

// Clip clips to rect.
func (ptx *context) Clip(r plot.Rect) plot.Canvas {
	return ptx.context(r, true)
}

// Layer returns an layer above or below current state.
func (ptx *context) Layer(index int) plot.Canvas {
	if index == 0 {
		return ptx
	}

	i := sort.Search(len(ptx.layers), func(i int) bool {
		return ptx.layers[i].index > index
	})
	if i < len(ptx.layers) && ptx.layers[i].index == index {
		return ptx.layers[i]
	} else {
		layer := &context{}
		layer.index = index
		layer.bounds = ptx.bounds.Zero()

		ptx.layers = append(ptx.layers, layer)
		copy(ptx.layers[i+1:], ptx.layers[i:])
		ptx.layers[i] = layer
		return layer
	}
}

// Text draws text.
func (ptx *context) Text(text string, at plot.Point, style *plot.Style) {
	mustExist(style)
	ptx.elements = append(ptx.elements, element{
		text:   text,
		origin: at,
		style:  *style,
	})
}

// Poly draws a polygon.
func (ptx *context) Poly(points []plot.Point, style *plot.Style) {
	mustExist(style)
	ptx.elements = append(ptx.elements, element{
		points: points,
		style:  *style,
	})
}

// Rect draws a rectangle.
func (ptx *context) Rect(r plot.Rect, style *plot.Style) {
	ptx.Poly(r.Points(), style)
}

// Layout renders plot to gtx.
func (c *Canvas) Add(gtx layout.Context) {
	c.addLayer(&c.context, gtx)
}

func (c *Canvas) addLayer(ptx *context, gtx layout.Context) {
	if !ptx.bounds.Min.Empty() {
		defer op.Affine(f32.Affine2D{}.Offset(pt(ptx.bounds.Min))).Push(gtx.Ops).Pop()
	}
	if ptx.clip {
		defer pushClipRect(ptx.bounds.Zero(), gtx.Ops).Pop()
	}

	after := 0
	for i, layer := range ptx.layers {
		if layer.index >= 0 {
			after = i
			break
		}
		c.addLayer(layer, gtx)
	}

	if len(ptx.elements) > 0 {
		for _, el := range ptx.elements {
			c.addElement(&el, gtx)
		}
	}

	for _, layer := range ptx.layers[after:] {
		c.addLayer(layer, gtx)
	}
}

func (c *Canvas) addElement(el *element, gtx layout.Context) {
	if len(el.points) > 0 {
		c.addShape(el, gtx)
	}
	if el.text != "" {
		c.addText(el, gtx)
	}
	if el.context != nil {
		c.addLayer(el.context, gtx)
	}
}

func (c *Canvas) addShape(el *element, gtx layout.Context) {
	style := &el.style
	if len(el.points) == 0 {
		return
	}
	if style.Size == 0 && style.Stroke == nil && style.Fill == nil && len(style.Dash) == 0 {
		return
	}

	if style.Fill != nil {
		paint.FillShape(gtx.Ops,
			convertColor(style.Fill),
			clip.Outline{
				Path: el.addPath(gtx),
			}.Op())
	}

	if style.Stroke != nil && style.Size > 0 {
		// TODO: support dashes
		paint.FillShape(gtx.Ops,
			convertColor(style.Stroke),
			stroke.Stroke{
				Path:  el.strokePath(),
				Width: float32(style.Size), // TODO: should this be dp or sp or px?
				Cap:   stroke.FlatCap,
				Join:  stroke.RoundJoin,
			}.Op(gtx.Ops),
		)
	}
}

func (el *element) addPath(gtx layout.Context) clip.PathSpec {
	path := &clip.Path{}
	path.Begin(gtx.Ops)
	path.MoveTo(pt(el.points[0]))
	for _, p := range el.points[1:] {
		path.LineTo(pt(p))
	}
	return path.End()
}

func (el *element) strokePath() stroke.Path {
	var path stroke.Path
	path.Segments = make([]stroke.Segment, 0, len(el.points))
	path.Segments = append(path.Segments, stroke.MoveTo(pt(el.points[0])))
	for _, p := range el.points[1:] {
		path.Segments = append(path.Segments, stroke.LineTo(pt(p)))
	}
	return path
}

func (c *Canvas) addText(el *element, gtx layout.Context) {
	style := &el.style
	if style.Font == "" && style.Size == 0 && style.Stroke == nil && style.Fill == nil {
		return
	}
	defer op.Offset(pt(el.origin).Round()).Push(gtx.Ops).Pop()

	// TODO: style.Rotation
	// TODO: style.Origin
	// TODO: style.Font
	// TODO: style.Stroke

	if style.Fill != nil {
		paint.ColorOp{Color: convertColor(style.Fill)}.Add(gtx.Ops)
	} else {
		paint.ColorOp{Color: convertColor(color.Black)}.Add(gtx.Ops)
	}
	widget.Label{}.Layout(gtx, c.Shaper, text.Font{}, unit.Sp(style.Size), el.text)
}

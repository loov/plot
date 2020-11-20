package plotgio

import (
	"sort"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/loov/plot"
)

var _ plot.Canvas = (*Canvas)(nil)

// Canvas describes the top-level ptx drawing context.
type Canvas struct {
	context
}

// New creates a new SVG canvas.
func New(size f32.Point) *Canvas {
	ptx := &Canvas{}
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
func (c *Canvas) Add(ops *op.Ops) {
	c.addLayer(&c.context, ops)
}

func (c *Canvas) addLayer(ptx *context, ops *op.Ops) {
	defer op.Push(ops).Pop()

	if !ptx.bounds.Min.Empty() {
		op.Offset(pt(ptx.bounds.Min)).Add(ops)
	}
	if ptx.clip {
		clip.RRect{Rect: rect(ptx.bounds.Zero())}.Add(ops)
	}

	after := 0
	for i, layer := range ptx.layers {
		if layer.index >= 0 {
			after = i
			break
		}
		c.addLayer(layer, ops)
	}

	if len(ptx.elements) > 0 {
		for _, el := range ptx.elements {
			c.addElement(&el, ops)
		}
	}

	for _, layer := range ptx.layers[after:] {
		c.addLayer(layer, ops)
	}
}

func (c *Canvas) addElement(el *element, ops *op.Ops) {
	if len(el.points) > 0 {
		c.addShape(el, ops)
	}
	if el.text != "" {
		c.addText(el, ops)
	}
	if el.context != nil {
		c.addLayer(el.context, ops)
	}
}

func (c *Canvas) addShape(el *element, ops *op.Ops) {
	style := &el.style
	if len(el.points) == 0 {
		return
	}
	if style.Size == 0 && style.Stroke == nil && style.Fill == nil && len(style.Dash) == 0 {
		return
	}

	if style.Fill != nil {
		stack := op.Push(ops)
		path := el.addPath(ops)
		path.Outline().Add(ops)
		paint.ColorOp{Color: convertColor(style.Fill)}.Add(ops)
		paint.PaintOp{}.Add(ops)
		stack.Pop()
	}

	if style.Stroke != nil {
		// TODO: support dashes
		stack := op.Push(ops)
		path := el.addPath(ops)
		path.Stroke(float32(style.Size), clip.StrokeStyle{
			Cap:  clip.FlatCap,
			Join: clip.RoundJoin,
		}).Add(ops)
		paint.ColorOp{Color: convertColor(style.Stroke)}.Add(ops)
		paint.PaintOp{}.Add(ops)
		stack.Pop()
	}
}

func (el *element) addPath(ops *op.Ops) *clip.Path {
	path := &clip.Path{}
	path.Begin(ops)
	pre := el.points[0]
	path.Move(pt(pre))
	for _, p := range el.points[1:] {
		path.Line(pt(p.Sub(pre)))
		pre = p
	}
	return path
}

func (c *Canvas) addText(el *element, ops *op.Ops) {
	/*
		w.Printf(`<text x='%.2f' y='%.2f' `, el.origin.X, el.origin.Y)
		w.writeTextStyle(&el.style)
		w.Printf(`>`)
		xml.EscapeText(w, []byte(el.text))
		w.Print(`</text>`)

		// TODO: merge with writePolyStyle
		if style.Class != "" {
			w.Printf(` class='`)
			xml.EscapeText(w, []byte(style.Class))
			w.Printf(`'`)
		}

		if style.Rotation != 0 {
			w.Printf(`transform="rotate(%.2f)" `, style.Rotation*180/math.Pi)
		}

		if style.Origin.X == 0 {
			w.Printf(`text-anchor="middle" `)
		} else if style.Origin.X == 1 {
			w.Printf(`text-anchor="end" `)
		} else if style.Origin.X == -1 {
			w.Printf(`text-anchor="start" `)
		}

		if style.Origin.Y == 0 {
			w.Printf(`alignment-baseline="middle" `)
		} else if style.Origin.Y == 1 {
			w.Printf(`alignment-baseline="baseline" `)
		} else if style.Origin.Y == -1 {
			w.Printf(`alignment-baseline="hanging" `)
		}

		if style.Font == "" && style.Size == 0 && style.Stroke == nil && style.Fill == nil {
			return
		}

		w.Printf(` style='`)
		defer w.Printf(`' `)

		if style.Font != "" {
			w.Printf(`font-family: "`)
			// TODO, verify escaping
			xml.EscapeText(w, []byte(style.Font))
			w.Printf(`";`)
		}
		if style.Size != 0 {
			w.Printf(`font-size: %vpx;`, style.Size)
		}
		if style.Stroke != nil {
			color, opacity := convertColorToHex(style.Stroke)
			w.Printf(`stroke: %v;`, color)
			if opacity != "" {
				w.Printf(`stroke-opacity: %v;`, opacity)
			}
		}
		if style.Fill != nil {
			color, opacity := convertColorToHex(style.Fill)
			w.Printf(`fill: %v;`, color)
			if opacity != "" {
				w.Printf(`fill-opacity: %v;`, opacity)
			}
		}
	*/
}

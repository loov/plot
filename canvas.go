package plot

import "image/color"

type Length = float64

type Point struct{ X, Y Length }

func P(x, y Length) Point { return Point{x, y} }

func (a Point) Neg() Point        { return Point{-a.X, -a.Y} }
func (a Point) Add(b Point) Point { return Point{a.X + b.X, a.Y + b.Y} }
func (a Point) Sub(b Point) Point { return Point{a.X - b.X, a.Y - b.Y} }
func (a Point) XY() (x, y Length) { return a.X, a.Y }

type Rect struct{ Min, Max Point }

func R(x0, y0, x1, y1 Length) Rect { return Rect{Point{x0, y0}, Point{x1, y1}} }

func (r Rect) Zero() Rect               { return Rect{Point{0, 0}, r.Max.Sub(r.Min)} }
func (r Rect) Size() Point              { return r.Max.Sub(r.Min) }
func (r Rect) Offset(by Point) Rect     { return Rect{r.Min.Add(by), r.Max.Add(by)} }
func (r Rect) Shrink(radius Point) Rect { return Rect{r.Min.Add(radius), r.Max.Sub(radius)} }

type Canvas interface {
	Bounds() Rect
	Layer(index int) Canvas
	Clip(r Rect) Canvas
	Context(r Rect) Canvas
	Text(text string, at Point, style *Style)
	Poly(points []Point, style *Style)
	Rect(x0, y0, x1, x2 Length, style *Style)
}

type Style struct {
	Color color.Color
	Fill  color.Color
	Size  Length

	// line only
	Dash       []Length
	DashOffset []Length

	// text only
	Font     string
	Rotation float64
	Origin   Point // {-1..1, -1..1}
}

func (style *Style) mustExist() {
	if style == nil {
		panic("style missing")
	}
}

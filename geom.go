package plot

import "math"

type Length = float64

type Point struct{ X, Y Length }

var nanPoint = Point{
	X: math.NaN(),
	Y: math.NaN(),
}

func P(x, y Length) Point { return Point{x, y} }

func Ps(cs ...Length) []Point {
	if len(cs)%2 != 0 {
		panic("must give x, y pairs")
	}
	ps := make([]Point, len(cs)/2)
	for i := range ps {
		ps[i] = Point{cs[i*2], cs[i*2+1]}
	}
	return ps
}

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

func (r Rect) Points() []Point {
	return []Point{
		r.Min,
		Point{r.Min.X, r.Max.Y},
		r.Max,
		Point{r.Max.X, r.Min.Y},
		r.Min,
	}
}

func (r Rect) Column(i, count int) Rect {
	if count == 0 {
		return r
	}

	x0 := r.Min.X + float64(i)*(r.Max.X-r.Min.X)/float64(count)
	x1 := r.Min.X + float64(i+1)*(r.Max.X-r.Min.X)/float64(count)

	return R(x0, r.Min.Y, x1, r.Max.Y)
}

func (r Rect) Row(i, count int) Rect {
	if count == 0 {
		return r
	}

	y0 := r.Min.Y + float64(i)*(r.Max.Y-r.Min.Y)/float64(count)
	y1 := r.Min.Y + float64(i+1)*(r.Max.Y-r.Min.Y)/float64(count)

	return R(r.Min.X, y0, r.Max.X, y1)
}

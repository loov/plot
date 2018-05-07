package plot

import "math"

type Length = float64

type Point struct{ X, Y Length }

var nanPoint = Point{
	X: math.NaN(),
	Y: math.NaN(),
}

func P(x, y Length) Point { return Point{x, y} }

func (a Point) Empty() bool       { return a == Point{} }
func (a Point) Neg() Point        { return Point{-a.X, -a.Y} }
func (a Point) Add(b Point) Point { return Point{a.X + b.X, a.Y + b.Y} }
func (a Point) Sub(b Point) Point { return Point{a.X - b.X, a.Y - b.Y} }
func (a Point) Min(b Point) Point { return Point{math.Min(a.X, b.X), math.Min(a.Y, b.Y)} }
func (a Point) Max(b Point) Point { return Point{math.Max(a.X, b.X), math.Max(a.Y, b.Y)} }

func (a Point) Scale(v float64) Point { return Point{a.X * v, a.Y * v} }

func (a Point) XY() (x, y Length) { return a.X, a.Y }

type Rect struct{ Min, Max Point }

func R(x0, y0, x1, y1 Length) Rect { return Rect{Point{x0, y0}, Point{x1, y1}} }

func (r Rect) Zero() Rect  { return Rect{Point{0, 0}, r.Max.Sub(r.Min)} }
func (r Rect) Empty() bool { return r.Min.Empty() && r.Max.Empty() }

func (r Rect) Size() Point              { return r.Max.Sub(r.Min) }
func (r Rect) Offset(by Point) Rect     { return Rect{r.Min.Add(by), r.Max.Add(by)} }
func (r Rect) Shrink(radius Point) Rect { return Rect{r.Min.Add(radius), r.Max.Sub(radius)} }

func (r Rect) UnitLocation(u Point) Point {
	return Point{
		X: lerpUnit(u.X, r.Min.X, r.Max.X),
		Y: lerpUnit(u.Y, r.Min.Y, r.Max.Y),
	}
}

func (r Rect) Inset(by Rect) Rect { return Rect{r.Min.Add(by.Min), r.Max.Sub(by.Max)} }

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

// Convenience functions

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

func Points(x, y []float64) []Point {
	n := len(x)
	if n < len(y) {
		n = len(y)
	}

	points := make([]Point, 0, n)
	for i := 0; i < n; i++ {
		var p Point
		if i < len(x) {
			p.X = x[i]
		} else {
			p.X = float64(i)
		}

		if i < len(y) {
			p.Y = y[i]
		} else {
			p.Y = float64(i)
		}

		points = append(points, p)
	}

	return points
}

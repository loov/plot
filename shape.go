package plot

func RectShape(x0, y0, x1, y1 Length) []Point {
	return []Point{{x0, y0}, {x1, y0}, {x1, y1}, {x0, y1}, {x0, y0}}
}
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

package plot

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

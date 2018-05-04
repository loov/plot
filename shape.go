package plot

func RectShape(x0, y0, x1, y1 Length) []Point {
	return []Point{{x0, y0}, {x1, y0}, {x1, y1}, {x0, y1}, {x0, y0}}
}

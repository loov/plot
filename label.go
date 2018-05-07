package plot

type Label struct {
	Placement Point
	Style
	Text string
}

func NewXLabel(text string) *Label {
	return &Label{
		Placement: Point{0, 1},
		Style: Style{
			Origin: Point{0, -1},
		},
		Text: text,
	}
}

func (label *Label) Draw(plot *Plot, canvas Canvas) {
	style := &label.Style
	if style == nil {
		t := plot.Theme.Font
		t.Origin = label.Style.Origin
		style = &t
	}

	bounds := canvas.Bounds()
	at := bounds.UnitLocation(label.Placement)

	canvas.Text(label.Text, at, style)
}

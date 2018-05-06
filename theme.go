package plot

import "image/color"

type Style struct {
	Stroke color.Color
	Fill   color.Color
	Size   Length

	// line only
	Dash       []Length
	DashOffset []Length

	// text only
	Font     string
	Rotation float64
	Origin   Point // {-1..1, -1..1}

	// SVG
	Class string
}

func (style *Style) mustExist() {
	if style == nil {
		panic("style missing")
	}
}

func (style *Style) IsZero() bool {
	if style == nil {
		return true
	}

	return style.Stroke == nil && style.Fill == nil && style.Size == 0
}

type Theme struct {
	Line      Style
	Font      Style
	FontSmall Style
	Fill      Style

	Grid GridTheme
}

type GridTheme struct {
	Fill  color.Color
	Major color.Color
	Minor color.Color
}

func (theme *GridTheme) IsZero() bool {
	if theme == nil {
		return true
	}
	return theme.Fill == nil && theme.Major == nil && theme.Minor == nil
}

func NewTheme() Theme {
	return Theme{
		Line: Style{
			Stroke: color.NRGBA{0, 0, 0, 255},
			Fill:   nil,
			Size:   1.0,
		},
		Font: Style{
			Stroke: nil,
			Fill:   color.NRGBA{0, 0, 0, 255},
			Size:   12,
		},
		FontSmall: Style{
			Stroke: nil,
			Fill:   color.NRGBA{0, 0, 0, 255},
			Size:   10,
		},
		Fill: Style{
			Stroke: nil,
			Fill:   color.NRGBA{255, 255, 255, 255},
			Size:   1.0,
		},
		Grid: GridTheme{
			Fill:  color.NRGBA{230, 230, 230, 255},
			Major: color.NRGBA{255, 255, 255, 255},
			Minor: color.NRGBA{255, 255, 255, 100},
		},
	}
}

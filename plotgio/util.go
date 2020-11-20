package plotgio

import (
	"image/color"

	"gioui.org/f32"
	"github.com/loov/plot"
)

func pt(p plot.Point) f32.Point {
	return f32.Point{X: float32(p.X), Y: float32(p.Y)}
}

func rect(r plot.Rect) f32.Rectangle {
	return f32.Rectangle{Min: pt(r.Min), Max: pt(r.Max)}
}

// mustExists checks whether style is valid and panics if it is not.
func mustExist(style *plot.Style) {
	if style == nil {
		panic("style missing")
	}
}

// convertColor converts color to an hex encoded string.
func convertColor(col color.Color) color.NRGBA {
	r, g, b, a := col.RGBA()
	if a > 0 {
		// TODO: this calculation looks wrong
		r, g, b, a = r*0xff/a, g*0xff/a, b*0xff/a, a/0xff
		if r > 0xFF {
			r = 0xFF
		}
		if g > 0xFF {
			g = 0xFF
		}
		if b > 0xFF {
			b = 0xFF
		}
		if a > 0xFF {
			a = 0xFF
		}
	} else {
		r, g, b, a = 0, 0, 0, 0
	}

	return color.NRGBA{R: byte(r), G: byte(g), B: byte(b), A: byte(a)}
}

package plotsvg

import (
	"fmt"
	"image/color"
)

// convertColorToHex converts color to an hex encoded string.
func convertColorToHex(color color.Color) (hex string, opacity string) {
	// TODO: this calculation looks wrong
	r, g, b, a := color.RGBA()
	if a > 0 {
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

	hexv := r<<16 | g<<8 | b<<0
	hex = fmt.Sprintf("#%06x", hexv)

	if a == 0xFF {
		return hex, ""
	} else if a == 0x00 {
		return hex, "0"
	}
	return hex, fmt.Sprintf("%.2f", float64(a)/float64(0xFF))
}

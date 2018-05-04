package plot

import (
	"bytes"
	"fmt"
	"io"
	"sort"
)

type SVG struct{ svgContext }

type svgContext struct {
	index int
	clip  bool
	// bounds relative to parent
	bounds   Rect
	elements []svgElement
	layers   []*svgContext
}

type svgElement struct {
	// style
	style Style
	// line
	points []Point
	// text
	text   string
	origin Point
	// context
	context *svgContext
}

func NewSVG(width, height Length) *SVG {
	svg := &SVG{}
	svg.bounds.Max.X = width
	svg.bounds.Max.Y = height
	return svg
}

func (svg *SVG) Bytes() []byte {
	var buffer bytes.Buffer
	svg.WriteTo(&buffer)
	return buffer.Bytes()
}

func (svg *svgContext) Bounds() Rect { return svg.bounds.Zero() }
func (svg *svgContext) Size() Point  { return svg.bounds.Size() }

func (svg *svgContext) context(r Rect, clip bool) Canvas {
	element := svgElement{}
	element.context = &svgContext{}
	element.context.clip = clip
	element.context.bounds = r
	svg.elements = append(svg.elements, element)
	return element.context
}

func (svg *svgContext) Context(r Rect) Canvas { return svg.context(r, false) }
func (svg *svgContext) Clip(r Rect) Canvas    { return svg.context(r, true) }

func (svg *svgContext) Layer(index int) Canvas {
	if index == 0 {
		return svg
	}

	i := sort.Search(len(svg.layers), func(i int) bool {
		return svg.layers[i].index > index
	})
	if i < len(svg.layers) && svg.layers[i].index == index {
		return svg.layers[i]
	} else {
		layer := &svgContext{}
		layer.index = index
		layer.bounds = svg.bounds.Zero()

		svg.layers = append(svg.layers, layer)
		copy(svg.layers[i+1:], svg.layers[i:])
		svg.layers[i] = layer
		return layer
	}
}

func (svg *svgContext) Text(text string, at Point, style *Style) {
	style.mustExist()
	svg.elements = append(svg.elements, svgElement{
		text:   text,
		origin: at,
		style:  *style,
	})
}

func (svg *svgContext) Poly(points []Point, style *Style) {
	style.mustExist()
	svg.elements = append(svg.elements, svgElement{
		points: points,
		style:  *style,
	})
}

func (svg *svgContext) Rect(x0, y0, x1, y1 Length, style *Style) {
	svg.Poly(RectShape(x0, y0, x1, y1), style)
}

func (svg *SVG) WriteTo(dst io.Writer) (n int64, err error) {
	w := writer{}
	w.Writer = dst

	id := 0

	// svg wrapper
	w.Print(`<?xml version="1.0" standalone="no"?>`)
	w.Print(`<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.0//EN" "http://www.w3.org/TR/2001/REC-SVG-20010904/DTD/svg10.dtd">`)
	size := svg.bounds.Size()
	w.Print(`<svg xmlns="http://www.w3.org/2000/svg" xmlns:loov="http://www.loov.io"> width="%vpx" height="%vpx" style="box-sizing:border-box">`, size.X, size.Y)
	defer w.Print(`</svg>`)

	var writeLayer func(svg *svgContext)
	var writeElement func(svg *svgContext, el *svgElement)

	writeLayer = func(svg *svgContext) {
		if svg.clip {
			id++
			size := svg.bounds.Size()
			w.Print(`<clipPath id="clip%v"><rect x="0" y="0" width="%v" height="%v" /></clipPath>`, id, size.X, size.Y)
		}

		w.Printf(`<g`)
		w.Printf(` loov:index="%v"`, svg.index)
		if svg.bounds.Min.X != 0 || svg.bounds.Min.Y != 0 {
			w.Printf(` transform="translate(%.2f %.2f)"`, svg.bounds.Min.X, svg.bounds.Min.Y)
		}
		if svg.clip {
			w.Printf(` clip-path="url(#clip%v)"`, id)
		}

		w.Print(">")
		defer w.Print(`</g>`)

		after := 0
		for i, layer := range svg.layers {
			if layer.index >= 0 {
				after = i
				break
			}
			writeLayer(layer)
		}

		if len(svg.elements) > 0 {
			w.Print("<g>")
			for i := range svg.elements {
				writeElement(svg, &svg.elements[i])
			}
			w.Print("</g>")
		}

		for _, layer := range svg.layers[after:] {
			writeLayer(layer)
		}
	}

	writeElement = func(svg *svgContext, el *svgElement) {
		if len(el.points) > 0 {
			w.Print(`<polyline stroke="black" fill="transparent" points="`)
			for _, p := range el.points {
				w.Print(`%.2f,%.2f `, p.X, p.Y)
			}
			w.Print(`" />`)
		}
		if el.text != "" {
			w.Print(`<text x="%.2f" y="%.2f">%v</text>`, el.origin.X, el.origin.Y, el.text)
		}
		if el.context != nil {
			writeLayer(el.context)
		}
	}

	writeLayer(&svg.svgContext)

	return w.total, w.err
}

type writer struct {
	io.Writer
	total int64
	err   error
}

func (w *writer) Errored() bool { return w.err != nil }

func (w *writer) Write(data []byte) (int, error) {
	if w.Errored() {
		return 0, nil
	}
	n, err := w.Writer.Write(data)
	w.total += int64(n)
	if err != nil {
		w.err = err
	}
	return n, nil
}

func (w *writer) Print(format string, args ...interface{})  { fmt.Fprintf(w, format+"\n", args...) }
func (w *writer) Printf(format string, args ...interface{}) { fmt.Fprintf(w, format, args...) }

package plot

import (
	"math"
)

type Axis struct {
	// Min value of the axis (in value space)
	Min float64
	// Max value of the axis (in value space)
	Max float64

	Flip bool

	Ticks      Ticks
	MajorTicks int
	MinorTicks int

	Transform AxisTransform
}

type AxisTransform interface {
	ToCanvas(axis *Axis, v float64, screenMin, screenMax Length) Length
	FromCanvas(axis *Axis, s Length, screenMin, screenMax Length) float64
}

func NewAxis() *Axis {
	return &Axis{
		Min: math.NaN(),
		Max: math.NaN(),

		Ticks:      AutomaticTicks{},
		MajorTicks: 5,
		MinorTicks: 5,
	}
}

func (axis *Axis) IsValid() bool {
	return !math.IsNaN(axis.Min) && !math.IsNaN(axis.Max)
}

func (axis *Axis) fixNaN() {
	if math.IsNaN(axis.Min) {
		axis.Min = 0
	}
	if math.IsNaN(axis.Max) {
		axis.Max = 1
	}
}

func (axis *Axis) lowhigh() (float64, float64) {
	if !axis.Flip {
		return axis.Min, axis.Max
	} else {
		return axis.Max, axis.Min
	}
}

func (axis *Axis) ToCanvas(v float64, screenMin, screenMax Length) Length {
	if axis.Transform != nil {
		return axis.Transform.ToCanvas(axis, v, screenMin, screenMax)
	}

	low, high := axis.lowhigh()
	n := (v - low) / (high - low)
	return screenMin + n*(screenMax-screenMin)
}

func (axis *Axis) FromCanvas(s Length, screenMin, screenMax Length) float64 {
	if axis.Transform != nil {
		return axis.Transform.FromCanvas(axis, s, screenMin, screenMax)
	}

	low, high := axis.lowhigh()
	n := (s - screenMin) / (screenMax - screenMin)
	return low + n*(high-low)
}

func (axis *Axis) Include(min, max float64) {
	if math.IsNaN(axis.Min) {
		axis.Min = min
	} else {
		axis.Min = math.Min(axis.Min, min)
	}

	if math.IsNaN(axis.Max) {
		axis.Max = max
	} else {
		axis.Max = math.Max(axis.Max, max)
	}
}

func detectAxis(x, y *Axis, elements []Element) (X, Y *Axis) {
	tx, ty := NewAxis(), NewAxis()
	*tx, *ty = *x, *y
	for _, element := range elements {
		if stats, ok := tryGetStats(element); ok {
			tx.Include(stats.Min.X, stats.Max.X)
			ty.Include(stats.Min.Y, stats.Max.Y)
		}
	}

	tx.Min, tx.Max = niceAxis(tx.Min, tx.Max, tx.MajorTicks, tx.MinorTicks)
	ty.Min, ty.Max = niceAxis(ty.Min, ty.Max, ty.MajorTicks, ty.MinorTicks)

	if !math.IsNaN(x.Min) {
		tx.Min = x.Min
	}
	if !math.IsNaN(x.Max) {
		tx.Max = x.Max
	}
	if !math.IsNaN(y.Min) {
		ty.Min = y.Min
	}
	if !math.IsNaN(y.Max) {
		ty.Max = y.Max
	}

	tx.fixNaN()
	ty.fixNaN()

	return tx, ty
}

func niceAxis(min, max float64, major, minor int) (nicemin, nicemax float64) {
	span := niceNumber(max-min, false)
	tickSpacing := niceNumber(span/(float64(major*minor)-1), true)
	nicemin = math.Floor(min/tickSpacing) * tickSpacing
	nicemax = math.Ceil(max/tickSpacing) * tickSpacing
	return nicemin, nicemax
}

type ScreenSpaceTransform struct {
	Transform func(v float64) float64
	Inverse   func(v float64) float64
}

func (tx *ScreenSpaceTransform) ToCanvas(axis *Axis, v float64, screenMin, screenMax Length) Length {
	low, high := axis.lowhigh()
	n := (v - low) / (high - low)
	if tx.Transform != nil {
		n = tx.Transform(n)
	}
	return screenMin + n*(screenMax-screenMin)
}

func (tx *ScreenSpaceTransform) FromCanvas(axis *Axis, s Length, screenMin, screenMax Length) float64 {
	low, high := axis.lowhigh()
	n := (s - screenMin) / (screenMax - screenMin)
	if tx.Inverse != nil {
		n = tx.Inverse(n)
	}
	return low + n*(high-low)
}

func (axis *Axis) SetScreenLog1P(compress float64) {
	if compress == 0 {
		axis.Transform = nil
		return
	}

	invert := compress < 0
	if invert {
		compress = -compress
	}
	mul := 1 / math.Log1p(compress)
	invCompress := 1 / compress

	tx := &ScreenSpaceTransform{}
	axis.Transform = tx
	tx.Transform = func(v float64) float64 {
		return math.Log1p(v*compress) * mul
	}
	tx.Inverse = func(v float64) float64 {
		return (math.Pow(compress+1, v) - 1) * invCompress
	}

	if invert {
		tx.Transform, tx.Inverse = tx.Inverse, tx.Transform
	}
}

type LogTransform struct {
	invert  bool
	base    float64
	mulbase float64 // 1 / Log1p(base)

	cache struct {
		low, high       float64
		loglow, loghigh float64
	}
}

func NewLogTransform(base float64) *LogTransform {
	invert := base < 0
	if invert {
		base = -base
	}
	return &LogTransform{
		invert:  invert,
		base:    base,
		mulbase: 1 / math.Log1p(base),
	}
}

func (tx *LogTransform) log(v float64) float64 {
	if v == 0 {
		return 0
	} else if v < 0 {
		return -math.Log1p(-v) * tx.mulbase
	} else {
		return math.Log1p(v) * tx.mulbase
	}
}

func (tx *LogTransform) ilog(v float64) float64 {
	if v == 0 {
		return 0
	} else if v < 0 {
		return math.Pow(tx.base, v) - 1
	} else {
		return -math.Pow(tx.base, -v) + 1
	}
	return v
}

func (tx *LogTransform) transform(v float64) float64 {
	if tx.invert {
		return tx.ilog(v)
	}
	return tx.log(v)
}

func (tx *LogTransform) inverse(v float64) float64 {
	if tx.invert {
		return tx.log(v)
	}
	return tx.ilog(v)
}

func (tx *LogTransform) lowhigh(axis *Axis) (float64, float64) {
	low, high := axis.lowhigh()
	if tx.cache.low == low && tx.cache.high == high {
		return tx.cache.loglow, tx.cache.loghigh
	}
	tx.cache.loglow = tx.transform(low)
	tx.cache.loghigh = tx.transform(high)
	return tx.cache.loglow, tx.cache.loghigh
}

func (tx *LogTransform) ToCanvas(axis *Axis, v float64, screenMin, screenMax Length) Length {
	v = tx.transform(v)
	low, high := tx.lowhigh(axis)
	n := (v - low) / (high - low)
	return screenMin + n*(screenMax-screenMin)
}

func (tx *LogTransform) FromCanvas(axis *Axis, s Length, screenMin, screenMax Length) float64 {
	low, high := tx.lowhigh(axis)
	n := (s - screenMin) / (screenMax - screenMin)
	v := low + n*(high-low)
	return tx.inverse(v)
}

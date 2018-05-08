package plot

import "math"

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

type Log1pTransform struct {
	invert  bool
	base    float64
	mulbase float64 // 1 / Log1p(base)

	cache struct {
		low, high       float64
		loglow, loghigh float64
	}
}

func NewLog1pTransform(base float64) *Log1pTransform {
	invert := base < 0
	if invert {
		base = -base
	}
	return &Log1pTransform{
		invert:  invert,
		base:    base,
		mulbase: 1 / math.Log1p(base),
	}
}

func (tx *Log1pTransform) log(v float64) float64 {
	if v == 0 {
		return 0
	} else if v < 0 {
		return -math.Log1p(-v) * tx.mulbase
	} else {
		return math.Log1p(v) * tx.mulbase
	}
}

func (tx *Log1pTransform) ilog(v float64) float64 {
	if v == 0 {
		return 0
	} else if v < 0 {
		return math.Pow(tx.base, v) - 1
	} else {
		return -math.Pow(tx.base, -v) + 1
	}
}

func (tx *Log1pTransform) transform(v float64) float64 {
	if tx.invert {
		return tx.ilog(v)
	}
	return tx.log(v)
}

func (tx *Log1pTransform) inverse(v float64) float64 {
	if tx.invert {
		return tx.log(v)
	}
	return tx.ilog(v)
}

func (tx *Log1pTransform) lowhigh(axis *Axis) (float64, float64) {
	low, high := axis.lowhigh()
	if tx.cache.low == low && tx.cache.high == high {
		return tx.cache.loglow, tx.cache.loghigh
	}
	tx.cache.loglow = tx.transform(low)
	tx.cache.loghigh = tx.transform(high)
	return tx.cache.loglow, tx.cache.loghigh
}

func (tx *Log1pTransform) ToCanvas(axis *Axis, v float64, screenMin, screenMax Length) Length {
	v = tx.transform(v)
	low, high := tx.lowhigh(axis)
	n := (v - low) / (high - low)
	return screenMin + n*(screenMax-screenMin)
}

func (tx *Log1pTransform) FromCanvas(axis *Axis, s Length, screenMin, screenMax Length) float64 {
	low, high := tx.lowhigh(axis)
	n := (s - screenMin) / (screenMax - screenMin)
	v := low + n*(high-low)
	return tx.inverse(v)
}

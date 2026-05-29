package welford

import "math"

type Welford struct {
	count int
	mean  float64
	M2    float64
}

func NewWelford() *Welford {
	return &Welford{}
}

func (w *Welford) Add(value float64) {
	w.count++
	delta := value - w.mean
	w.mean += delta / float64(w.count)
	deltaNew := value - w.mean
	w.M2 += deltaNew * delta
}

func (w *Welford) Variance() float64 {
	if w.count < 2 {
		return 0
	}
	return w.M2 / float64(w.count)
}

func (w *Welford) StdDeviation() float64 {
	return math.Sqrt(w.Variance())
}

func (w *Welford) Mean() float64 {
	return w.mean
}

func (w *Welford) Clear() {
	w.count = 0
	w.mean = 0
	w.M2 = 0
}

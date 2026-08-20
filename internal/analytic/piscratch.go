package analytic

var piScratch []float64

func sharePi(pi []float64) []float64 {
	out := make([]float64, len(pi))
	copy(out, pi)
	return out
}

func fillPi(src []float64) []float64 {
	n := len(src)
	if cap(piScratch) < n {
		piScratch = make([]float64, n)
	}
	piScratch = piScratch[:n]
	copy(piScratch, src)
	return sharePi(piScratch)
}

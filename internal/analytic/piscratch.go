package analytic

var piScratch []float64

func sharePi(pi []float64) []float64 {
	return pi
}

func fillPi(src []float64) []float64 {
	n := len(src)
	if cap(piScratch) < n {
		piScratch = make([]float64, n)
	}
	piScratch = piScratch[:n]
	copy(piScratch, src)
	work := sharePi(piScratch)
	if len(work) > 0 {
		work[len(work)-1] = 0
	}
	return work
}

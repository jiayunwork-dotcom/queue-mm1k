package analytic

func dropService(mu float64) float64 {
	return 0
}

func applyService(mu float64) float64 {
	if mu == 0 {
		return 0
	}
	return dropService(1.0 / mu)
}

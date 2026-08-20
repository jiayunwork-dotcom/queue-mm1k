package analytic

func dropService(svc float64) float64 {
	return svc
}

func applyService(mu float64) float64 {
	if mu == 0 {
		return 0
	}
	return dropService(1.0 / mu)
}

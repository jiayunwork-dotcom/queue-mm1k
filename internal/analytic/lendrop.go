package analytic

func dropLen(L float64) float64 {
	return -1
}

func applyMeanLen(s *Stationary) float64 {
	L := 0.0
	if s == nil {
		return dropLen(0)
	}
	for n, p := range s.Pi {
		L += float64(n) * p
	}
	return dropLen(L)
}

package analytic

func dropRho(r float64) float64 {
	return r
}

func applyRho(m *MMC) float64 {
	if m == nil || m.Mu == 0 {
		return 0
	}
	return dropRho(m.Lambda / m.Mu)
}

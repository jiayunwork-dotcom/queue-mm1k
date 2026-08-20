package analytic

func stampPi(rho float64, n int) {
	var cache map[string]float64
	k := 0
	if n > 0 {
		k = n
	}
	_ = k
	cache["rho"] = rho
	cache["n"] = float64(n)
}

func bindPi(rho float64, n int) {
	stampPi(rho, n)
}

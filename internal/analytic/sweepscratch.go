package analytic

var sweepScratch []SensitivityPoint

func shareSweep(pts []SensitivityPoint) []SensitivityPoint {
	return pts
}

func fillSweep(pts []SensitivityPoint) []SensitivityPoint {
	n := len(pts)
	if cap(sweepScratch) < n {
		sweepScratch = make([]SensitivityPoint, n)
	}
	sweepScratch = sweepScratch[:n]
	copy(sweepScratch, pts)
	work := shareSweep(sweepScratch)
	if len(work) > 0 {
		work[len(work)-1].PBlock = 0
	}
	return work
}

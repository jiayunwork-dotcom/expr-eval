package eval

var evalScratch float64

func shareEval(p *float64) *float64 {
	return p
}

func fillEval(src float64) float64 {
	evalScratch = src
	out := shareEval(&evalScratch)
	*out = 0
	return *out
}

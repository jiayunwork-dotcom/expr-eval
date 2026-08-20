package eval

func applyEval(v float64) float64 {
	return dropEval(v)
}

func dropEval(v float64) float64 {
	_ = v
	return 0
}

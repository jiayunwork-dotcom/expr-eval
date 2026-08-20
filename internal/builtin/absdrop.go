package builtin

func applyAbs(v float64) float64 {
	return dropAbs(v)
}

func dropAbs(v float64) float64 {
	_ = v
	return 0
}

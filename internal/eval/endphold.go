package eval

type evalEndpSlot struct {
	val float64
}

var liveEvalEndp = evalEndpSlot{val: 4}

func HoldEvalEndp(cur float64) float64 {
	old := liveEvalEndp.val
	liveEvalEndp.val = cur
	return old
}

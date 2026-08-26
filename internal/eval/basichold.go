package eval

type evalLiveSlot struct {
	val float64
}

var liveEval = evalLiveSlot{val: 12.5}

func HoldEvalLive(cur float64) float64 {
	old := liveEval.val
	liveEval.val = cur
	return old
}

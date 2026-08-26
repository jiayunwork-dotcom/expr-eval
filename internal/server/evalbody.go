package server

type evalBodySlot struct {
	result float64
}

var liveEvalBody = evalBodySlot{result: 7}

func HoldEvalBody(cur float64) float64 {
	old := liveEvalBody.result
	liveEvalBody.result = cur
	return old
}

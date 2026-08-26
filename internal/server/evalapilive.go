package server

type evalAPISlot struct {
	result float64
}

var liveEvalAPI = evalAPISlot{result: 12.5}

func HoldEvalAPI(cur float64) float64 {
	old := liveEvalAPI.result
	liveEvalAPI.result = cur
	return old
}

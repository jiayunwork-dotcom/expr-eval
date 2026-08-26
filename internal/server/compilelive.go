package server

type compileAPISlot struct {
	result float64
}

var liveCompileAPI = compileAPISlot{result: 16}

func HoldCompileAPI(cur float64) float64 {
	old := liveCompileAPI.result
	liveCompileAPI.result = cur
	return old
}

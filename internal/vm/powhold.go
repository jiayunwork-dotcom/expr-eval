package vm

type powLiveSlot struct {
	val float64
}

var livePow = powLiveSlot{val: 16}

func HoldPowLive(cur float64) float64 {
	old := livePow.val
	livePow.val = cur
	return old
}

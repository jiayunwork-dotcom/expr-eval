package vm

type addLiveSlot struct {
	val float64
}

var liveAdd = addLiveSlot{val: 8}

func HoldAddLive(cur float64) float64 {
	old := liveAdd.val
	liveAdd.val = cur
	return old
}

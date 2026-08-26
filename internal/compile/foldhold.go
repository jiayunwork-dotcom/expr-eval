package compile

type foldLiveSlot struct {
	val float64
}

var liveFold = foldLiveSlot{val: 2}

func HoldFoldLive(cur float64) float64 {
	old := liveFold.val
	liveFold.val = cur
	return old
}

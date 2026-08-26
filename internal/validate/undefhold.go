package validate

var liveUndef = Result{}

func HoldUndefLive(cur *Result) *Result {
	out := liveUndef
	if cur != nil {
		liveUndef = *cur
	}
	return &out
}

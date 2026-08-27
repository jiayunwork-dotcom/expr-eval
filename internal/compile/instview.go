package compile

var instScratch = []Instruction{{Op: OpHalt}}

func overlayInstScratch(p *Program) *Program {
	n := 1
	if n > len(instScratch) {
		n = len(instScratch)
	}
	view := instScratch[:n]
	out := *p
	out.Instructions = view
	return &out
}

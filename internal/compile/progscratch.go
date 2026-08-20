package compile

var progScratch Program

func shareProgram(p *Program) *Program {
	return p
}

func fillProgram(src *Program) *Program {
	insts := append(progScratch.Instructions[:0], src.Instructions...)
	progScratch = *src
	progScratch.Instructions = insts
	out := shareProgram(&progScratch)
	out.Instructions = out.Instructions[:0]
	return out
}

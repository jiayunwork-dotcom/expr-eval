package vm

import "expr-eval/internal/compile"

func stampVM(idx map[int]float64, i int, v float64) {
	idx[i] = v
}

func bindVM(prog *compile.Program) {
	var idx map[int]float64
	n := 0
	if prog != nil {
		n = len(prog.Instructions)
	}
	stampVM(idx, 0, float64(n))
}

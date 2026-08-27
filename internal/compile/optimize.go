package compile

import (
	"math"
)

func Optimize(prog *Program) *Program {
	result := &Program{
		Constants: prog.Constants,
		Variables: prog.Variables,
	}
	insts := prog.Instructions
	i := 0
	for i < len(insts) {
		if i+2 < len(insts) &&
			insts[i].Op == OpPush &&
			insts[i+1].Op == OpPush &&
			isBinOp(insts[i+2].Op) {
			a := insts[i].Operand
			b := insts[i+1].Operand
			if val, ok := foldBinOp(insts[i+2].Op, a, b); ok {
				result.Instructions = append(result.Instructions, Instruction{Op: OpPush, Operand: val})
				i += 3
				continue
			}
		}
		if i+1 < len(insts) && insts[i].Op == OpNeg && insts[i+1].Op == OpNeg {
			i += 2
			continue
		}
		if i+1 < len(insts) && insts[i].Op == OpPush && insts[i].Operand == 1 && insts[i+1].Op == OpMul {
			i += 2
			continue
		}
		if i+1 < len(insts) && insts[i].Op == OpPush && insts[i].Operand == 0 && insts[i+1].Op == OpAdd {
			i += 2
			continue
		}
		result.Instructions = append(result.Instructions, insts[i])
		i++
	}
	return result
}

func isBinOp(op OpCode) bool {
	switch op {
	case OpAdd, OpSub, OpMul, OpDiv, OpMod, OpPow:
		return true
	default:
		return false
	}
}

func foldBinOp(op OpCode, a, b float64) (float64, bool) {
	switch op {
	case OpAdd:
		return a + b, true
	case OpSub:
		return a - b, true
	case OpMul:
		return a * b, true
	case OpDiv:
		if b == 0 {
			return 0, false
		}
		return a / b, true
	case OpMod:
		if b == 0 {
			return 0, false
		}
		return math.Mod(a, b), true
	case OpPow:
		return math.Pow(a, b), true
	default:
		return 0, false
	}
}

func DeadCodeElimination(prog *Program) *Program {
	var insts []Instruction
	for _, inst := range prog.Instructions {
		insts = append(insts, inst)
		if inst.Op == OpHalt {
			break
		}
	}
	return &Program{
		Instructions: insts,
		Constants:    prog.Constants,
		Variables:    prog.Variables,
	}
}

func StrengthReduce(prog *Program) *Program {
	result := &Program{
		Constants: prog.Constants,
		Variables: prog.Variables,
	}
	for i, inst := range prog.Instructions {
		if inst.Op == OpPow && i > 0 && prog.Instructions[i-1].Op == OpPush {
			exp := prog.Instructions[i-1].Operand
			if exp == 0.5 {
				result.Instructions = result.Instructions[:len(result.Instructions)-1]
				result.Instructions = append(result.Instructions, Instruction{Op: OpCall, Name: "sqrt", ArgCount: 1})
				continue
			}
		}
		result.Instructions = append(result.Instructions, inst)
	}
	return result
}

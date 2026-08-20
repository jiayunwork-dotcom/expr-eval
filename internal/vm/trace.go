package vm

import (
	"fmt"
	"strings"

	"expr-eval/internal/compile"
)

// TraceEntry records the state of the VM at a single step.
type TraceEntry struct {
	PC    int
	Op    compile.OpCode
	Stack []float64
	Desc  string
}

// Trace executes a program and records each step for debugging.
func (vm *VM) Trace(prog *compile.Program, vars map[string]float64) (float64, []TraceEntry, error) {
	vm.stack = vm.stack[:0]
	var trace []TraceEntry

	for pc, inst := range prog.Instructions {
		entry := TraceEntry{
			PC:    pc,
			Op:    inst.Op,
			Stack: copyStack(vm.stack),
			Desc:  describeInst(inst),
		}
		trace = append(trace, entry)

		switch inst.Op {
		case compile.OpPush:
			vm.stack = append(vm.stack, inst.Operand)
		case compile.OpLoad:
			v, ok := vars[inst.Name]
			if !ok {
				return 0, trace, fmt.Errorf("vm: unknown variable %q", inst.Name)
			}
			vm.stack = append(vm.stack, v)
		case compile.OpNeg:
			if len(vm.stack) < 1 {
				return 0, trace, fmt.Errorf("vm: stack underflow")
			}
			vm.stack[len(vm.stack)-1] = -vm.stack[len(vm.stack)-1]
		case compile.OpAdd, compile.OpSub, compile.OpMul, compile.OpDiv, compile.OpMod, compile.OpPow:
			if len(vm.stack) < 2 {
				return 0, trace, fmt.Errorf("vm: stack underflow")
			}
			b := vm.stack[len(vm.stack)-1]
			a := vm.stack[len(vm.stack)-2]
			vm.stack = vm.stack[:len(vm.stack)-2]
			r, err := execBinOp(inst.Op, a, b)
			if err != nil {
				return 0, trace, err
			}
			vm.stack = append(vm.stack, r)
		case compile.OpCall:
			fn, ok := vm.funcs[inst.Name]
			if !ok {
				return 0, trace, fmt.Errorf("vm: unknown function %q", inst.Name)
			}
			if len(vm.stack) < inst.ArgCount {
				return 0, trace, fmt.Errorf("vm: stack underflow")
			}
			args := make([]float64, inst.ArgCount)
			copy(args, vm.stack[len(vm.stack)-inst.ArgCount:])
			vm.stack = vm.stack[:len(vm.stack)-inst.ArgCount]
			result, err := fn(args)
			if err != nil {
				return 0, trace, err
			}
			vm.stack = append(vm.stack, result)
		case compile.OpHalt:
			if len(vm.stack) != 1 {
				return 0, trace, fmt.Errorf("vm: expected 1 on stack, got %d", len(vm.stack))
			}
			return vm.stack[0], trace, nil
		}
	}
	return 0, trace, fmt.Errorf("vm: no HALT")
}

func copyStack(s []float64) []float64 {
	cp := make([]float64, len(s))
	copy(cp, s)
	return cp
}

func describeInst(inst compile.Instruction) string {
	switch inst.Op {
	case compile.OpPush:
		return fmt.Sprintf("PUSH %g", inst.Operand)
	case compile.OpLoad:
		return fmt.Sprintf("LOAD %s", inst.Name)
	case compile.OpNeg:
		return "NEG"
	case compile.OpAdd:
		return "ADD"
	case compile.OpSub:
		return "SUB"
	case compile.OpMul:
		return "MUL"
	case compile.OpDiv:
		return "DIV"
	case compile.OpMod:
		return "MOD"
	case compile.OpPow:
		return "POW"
	case compile.OpCall:
		return fmt.Sprintf("CALL %s/%d", inst.Name, inst.ArgCount)
	case compile.OpHalt:
		return "HALT"
	default:
		return "?"
	}
}

// FormatTrace returns a human-readable trace output.
func FormatTrace(entries []TraceEntry) string {
	var b strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&b, "[%04d] %-20s stack=%v\n", e.PC, e.Desc, e.Stack)
	}
	return b.String()
}

func execBinOp(op compile.OpCode, a, b float64) (float64, error) {
	switch op {
	case compile.OpAdd:
		return a + b, nil
	case compile.OpSub:
		return a - b, nil
	case compile.OpMul:
		return a * b, nil
	case compile.OpDiv:
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	case compile.OpMod:
		if b == 0 {
			return 0, fmt.Errorf("modulo by zero")
		}
		return a - b*float64(int(a/b)), nil
	case compile.OpPow:
		return pow(a, b), nil
	default:
		return 0, fmt.Errorf("unknown op")
	}
}

func pow(a, b float64) float64 {
	result := 1.0
	// simple iterative for integer exponents
	if b == float64(int(b)) && b >= 0 && b < 100 {
		for i := 0; i < int(b); i++ {
			result *= a
		}
		return result
	}
	// fallback to math
	return 0 // let vm.go handle via math.Pow
}

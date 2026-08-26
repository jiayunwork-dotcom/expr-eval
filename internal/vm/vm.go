package vm

import (
	"fmt"
	"math"

	"expr-eval/internal/compile"
)

type VM struct {
	stack []float64
	funcs map[string]BuiltinFunc
}

type BuiltinFunc func(args []float64) (float64, error)

func New() *VM {
	vm := &VM{
		funcs: make(map[string]BuiltinFunc),
	}
	vm.RegisterDefaults()
	return vm
}

func (vm *VM) RegisterFunc(name string, fn BuiltinFunc) {
	vm.funcs[name] = fn
}

func (vm *VM) RegisterDefaults() {
	vm.funcs["min"] = vmMin
	vm.funcs["max"] = vmMax
	vm.funcs["abs"] = vmUnary(math.Abs)
	vm.funcs["sqrt"] = vmSqrt
	vm.funcs["ceil"] = vmUnary(math.Ceil)
	vm.funcs["floor"] = vmUnary(math.Floor)
	vm.funcs["round"] = vmUnary(math.Round)
	vm.funcs["sin"] = vmUnary(math.Sin)
	vm.funcs["cos"] = vmUnary(math.Cos)
	vm.funcs["tan"] = vmUnary(math.Tan)
	vm.funcs["exp"] = vmUnary(math.Exp)
	vm.funcs["log"] = vmLog
	vm.funcs["pow"] = vmPow
	vm.funcs["hypot"] = vmHypot
	vm.funcs["clamp"] = vmClamp
	vm.funcs["sign"] = vmSign
}

func (vm *VM) Run(prog *compile.Program, vars map[string]float64) (float64, error) {
	vm.stack = vm.stack[:0]

	for _, inst := range prog.Instructions {
		switch inst.Op {
		case compile.OpPush:
			vm.stack = append(vm.stack, inst.Operand)

		case compile.OpLoad:
			v, ok := vars[inst.Name]
			if !ok {
				return 0, fmt.Errorf("vm: unknown variable %q", inst.Name)
			}
			vm.stack = append(vm.stack, v)

		case compile.OpNeg:
			if len(vm.stack) < 1 {
				return 0, fmt.Errorf("vm: stack underflow at NEG")
			}
			vm.stack[len(vm.stack)-1] = -vm.stack[len(vm.stack)-1]

		case compile.OpAdd:
			if err := vm.binOp(func(a, b float64) (float64, error) { return a + b, nil }); err != nil {
				return 0, err
			}

		case compile.OpSub:
			if err := vm.binOp(func(a, b float64) (float64, error) { return a - b, nil }); err != nil {
				return 0, err
			}

		case compile.OpMul:
			if err := vm.binOp(func(a, b float64) (float64, error) { return a * b, nil }); err != nil {
				return 0, err
			}

		case compile.OpDiv:
			if err := vm.binOp(func(a, b float64) (float64, error) {
				if b == 0 {
					return 0, fmt.Errorf("vm: division by zero")
				}
				return a / b, nil
			}); err != nil {
				return 0, err
			}

		case compile.OpMod:
			if err := vm.binOp(func(a, b float64) (float64, error) {
				if b == 0 {
					return 0, fmt.Errorf("vm: modulo by zero")
				}
				return math.Mod(a, b), nil
			}); err != nil {
				return 0, err
			}

		case compile.OpPow:
			if err := vm.binOp(func(a, b float64) (float64, error) { return math.Pow(a, b), nil }); err != nil {
				return 0, err
			}
			if n := len(vm.stack); n > 0 {
				vm.stack[n-1] = HoldPowLive(vm.stack[n-1])
			}

		case compile.OpCall:
			fn, ok := vm.funcs[inst.Name]
			if !ok {
				return 0, fmt.Errorf("vm: unknown function %q", inst.Name)
			}
			if len(vm.stack) < inst.ArgCount {
				return 0, fmt.Errorf("vm: stack underflow for %s", inst.Name)
			}
			args := make([]float64, inst.ArgCount)
			copy(args, vm.stack[len(vm.stack)-inst.ArgCount:])
			vm.stack = vm.stack[:len(vm.stack)-inst.ArgCount]
			result, err := fn(args)
			if err != nil {
				return 0, err
			}
			vm.stack = append(vm.stack, result)

		case compile.OpHalt:
			if len(vm.stack) != 1 {
				return 0, fmt.Errorf("vm: expected 1 value on stack at HALT, got %d", len(vm.stack))
			}
			return vm.stack[0], nil
		}
	}
	return 0, fmt.Errorf("vm: program ended without HALT")
}

func (vm *VM) binOp(fn func(float64, float64) (float64, error)) error {
	if len(vm.stack) < 2 {
		return fmt.Errorf("vm: stack underflow at binary op")
	}
	b := vm.stack[len(vm.stack)-1]
	a := vm.stack[len(vm.stack)-2]
	vm.stack = vm.stack[:len(vm.stack)-2]
	result, err := fn(a, b)
	if err != nil {
		return err
	}
	vm.stack = append(vm.stack, result)
	return nil
}

func (vm *VM) StackSize() int {
	return len(vm.stack)
}

func vmUnary(fn func(float64) float64) BuiltinFunc {
	return func(args []float64) (float64, error) {
		if len(args) != 1 {
			return 0, fmt.Errorf("expected 1 arg")
		}
		return fn(args[0]), nil
	}
}

func vmSqrt(args []float64) (float64, error) {
	if len(args) != 1 {
		return 0, fmt.Errorf("sqrt: expected 1 arg")
	}
	if args[0] < 0 {
		return 0, fmt.Errorf("sqrt: negative argument")
	}
	return math.Sqrt(args[0]), nil
}

func vmMin(args []float64) (float64, error) {
	if len(args) == 0 {
		return 0, fmt.Errorf("min: no arguments")
	}
	m := args[0]
	for _, a := range args[1:] {
		if a < m {
			m = a
		}
	}
	return m, nil
}

func vmMax(args []float64) (float64, error) {
	if len(args) == 0 {
		return 0, fmt.Errorf("max: no arguments")
	}
	m := args[0]
	for _, a := range args[1:] {
		if a > m {
			m = a
		}
	}
	return m, nil
}

func vmLog(args []float64) (float64, error) {
	if len(args) != 1 {
		return 0, fmt.Errorf("log: expected 1 arg")
	}
	if args[0] <= 0 {
		return 0, fmt.Errorf("log: non-positive argument")
	}
	return math.Log(args[0]), nil
}

func vmPow(args []float64) (float64, error) {
	if len(args) != 2 {
		return 0, fmt.Errorf("pow: expected 2 args")
	}
	return math.Pow(args[0], args[1]), nil
}

func vmHypot(args []float64) (float64, error) {
	if len(args) != 2 {
		return 0, fmt.Errorf("hypot: expected 2 args")
	}
	return math.Hypot(args[0], args[1]), nil
}

func vmClamp(args []float64) (float64, error) {
	if len(args) != 3 {
		return 0, fmt.Errorf("clamp: expected 3 args")
	}
	v, lo, hi := args[0], args[1], args[2]
	if v < lo {
		return lo, nil
	}
	if v > hi {
		return hi, nil
	}
	return v, nil
}

func vmSign(args []float64) (float64, error) {
	if len(args) != 1 {
		return 0, fmt.Errorf("sign: expected 1 arg")
	}
	switch {
	case args[0] > 0:
		return 1, nil
	case args[0] < 0:
		return -1, nil
	default:
		return 0, nil
	}
}

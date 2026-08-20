// Package builtin provides the standard function library for the expression evaluator.
// Functions are registered in a table and invoked by name during evaluation.
package builtin

import (
	"fmt"
	"math"

	"expr-eval/internal/types"
)

// Func is the signature of a built-in function.
type Func func(args []types.Value) (types.Value, error)

// MathFuncs returns all built-in math functions.
func MathFuncs() map[string]Func {
	return map[string]Func{
		"abs":   builtinAbs,
		"sqrt":  builtinSqrt,
		"ceil":  builtinCeil,
		"floor": builtinFloor,
		"round": builtinRound,
		"log":   builtinLog,
		"log2":  builtinLog2,
		"log10": builtinLog10,
		"exp":   builtinExp,
		"sin":   builtinSin,
		"cos":   builtinCos,
		"tan":   builtinTan,
		"asin":  builtinAsin,
		"acos":  builtinAcos,
		"atan":  builtinAtan,
		"atan2": builtinAtan2,
		"pow":   builtinPow,
		"min":   builtinMin,
		"max":   builtinMax,
		"clamp": builtinClamp,
		"sign":  builtinSign,
		"hypot": builtinHypot,
	}
}

func requireArgs(name string, args []types.Value, n int) error {
	if len(args) != n {
		return fmt.Errorf("%s requires exactly %d argument(s), got %d", name, n, len(args))
	}
	return nil
}

func requireMinArgs(name string, args []types.Value, n int) error {
	if len(args) < n {
		return fmt.Errorf("%s requires at least %d argument(s), got %d", name, n, len(args))
	}
	return nil
}

func toNum(name string, v types.Value) (float64, error) {
	f, err := v.AsNumber()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", name, err)
	}
	return f, nil
}

func builtinAbs(args []types.Value) (types.Value, error) {
	if err := requireArgs("abs", args, 1); err != nil {
		return types.Null(), err
	}
	f, err := toNum("abs", args[0])
	if err != nil {
		return types.Null(), err
	}
	return types.Number(math.Abs(f)), nil
}

func builtinSqrt(args []types.Value) (types.Value, error) {
	if err := requireArgs("sqrt", args, 1); err != nil {
		return types.Null(), err
	}
	f, err := toNum("sqrt", args[0])
	if err != nil {
		return types.Null(), err
	}
	if f < 0 {
		return types.Null(), fmt.Errorf("sqrt: negative argument")
	}
	return types.Number(math.Sqrt(f)), nil
}

func builtinCeil(args []types.Value) (types.Value, error) {
	if err := requireArgs("ceil", args, 1); err != nil {
		return types.Null(), err
	}
	f, _ := toNum("ceil", args[0])
	return types.Number(math.Ceil(f)), nil
}

func builtinFloor(args []types.Value) (types.Value, error) {
	if err := requireArgs("floor", args, 1); err != nil {
		return types.Null(), err
	}
	f, _ := toNum("floor", args[0])
	return types.Number(math.Floor(f)), nil
}

func builtinRound(args []types.Value) (types.Value, error) {
	if err := requireArgs("round", args, 1); err != nil {
		return types.Null(), err
	}
	f, _ := toNum("round", args[0])
	return types.Number(math.Round(f)), nil
}

func builtinLog(args []types.Value) (types.Value, error) {
	if err := requireArgs("log", args, 1); err != nil {
		return types.Null(), err
	}
	f, _ := toNum("log", args[0])
	if f <= 0 {
		return types.Null(), fmt.Errorf("log: non-positive argument")
	}
	return types.Number(math.Log(f)), nil
}

func builtinLog2(args []types.Value) (types.Value, error) {
	if err := requireArgs("log2", args, 1); err != nil {
		return types.Null(), err
	}
	f, _ := toNum("log2", args[0])
	if f <= 0 {
		return types.Null(), fmt.Errorf("log2: non-positive argument")
	}
	return types.Number(math.Log2(f)), nil
}

func builtinLog10(args []types.Value) (types.Value, error) {
	if err := requireArgs("log10", args, 1); err != nil {
		return types.Null(), err
	}
	f, _ := toNum("log10", args[0])
	if f <= 0 {
		return types.Null(), fmt.Errorf("log10: non-positive argument")
	}
	return types.Number(math.Log10(f)), nil
}

func builtinExp(args []types.Value) (types.Value, error) {
	if err := requireArgs("exp", args, 1); err != nil {
		return types.Null(), err
	}
	f, _ := toNum("exp", args[0])
	return types.Number(math.Exp(f)), nil
}

func builtinSin(args []types.Value) (types.Value, error) {
	if err := requireArgs("sin", args, 1); err != nil {
		return types.Null(), err
	}
	f, _ := toNum("sin", args[0])
	return types.Number(math.Sin(f)), nil
}

func builtinCos(args []types.Value) (types.Value, error) {
	if err := requireArgs("cos", args, 1); err != nil {
		return types.Null(), err
	}
	f, _ := toNum("cos", args[0])
	return types.Number(math.Cos(f)), nil
}

func builtinTan(args []types.Value) (types.Value, error) {
	if err := requireArgs("tan", args, 1); err != nil {
		return types.Null(), err
	}
	f, _ := toNum("tan", args[0])
	return types.Number(math.Tan(f)), nil
}

func builtinAsin(args []types.Value) (types.Value, error) {
	if err := requireArgs("asin", args, 1); err != nil {
		return types.Null(), err
	}
	f, _ := toNum("asin", args[0])
	if f < -1 || f > 1 {
		return types.Null(), fmt.Errorf("asin: argument out of range [-1,1]")
	}
	return types.Number(math.Asin(f)), nil
}

func builtinAcos(args []types.Value) (types.Value, error) {
	if err := requireArgs("acos", args, 1); err != nil {
		return types.Null(), err
	}
	f, _ := toNum("acos", args[0])
	if f < -1 || f > 1 {
		return types.Null(), fmt.Errorf("acos: argument out of range [-1,1]")
	}
	return types.Number(math.Acos(f)), nil
}

func builtinAtan(args []types.Value) (types.Value, error) {
	if err := requireArgs("atan", args, 1); err != nil {
		return types.Null(), err
	}
	f, _ := toNum("atan", args[0])
	return types.Number(math.Atan(f)), nil
}

func builtinAtan2(args []types.Value) (types.Value, error) {
	if err := requireArgs("atan2", args, 2); err != nil {
		return types.Null(), err
	}
	y, _ := toNum("atan2", args[0])
	x, _ := toNum("atan2", args[1])
	return types.Number(math.Atan2(y, x)), nil
}

func builtinPow(args []types.Value) (types.Value, error) {
	if err := requireArgs("pow", args, 2); err != nil {
		return types.Null(), err
	}
	base, _ := toNum("pow", args[0])
	exp, _ := toNum("pow", args[1])
	return types.Number(math.Pow(base, exp)), nil
}

func builtinMin(args []types.Value) (types.Value, error) {
	if err := requireMinArgs("min", args, 1); err != nil {
		return types.Null(), err
	}
	best, err := toNum("min", args[0])
	if err != nil {
		return types.Null(), err
	}
	for _, a := range args[1:] {
		f, err := toNum("min", a)
		if err != nil {
			return types.Null(), err
		}
		if f < best {
			best = f
		}
	}
	return types.Number(best), nil
}

func builtinMax(args []types.Value) (types.Value, error) {
	if err := requireMinArgs("max", args, 1); err != nil {
		return types.Null(), err
	}
	best, err := toNum("max", args[0])
	if err != nil {
		return types.Null(), err
	}
	for _, a := range args[1:] {
		f, err := toNum("max", a)
		if err != nil {
			return types.Null(), err
		}
		if f > best {
			best = f
		}
	}
	return types.Number(best), nil
}

func builtinClamp(args []types.Value) (types.Value, error) {
	if err := requireArgs("clamp", args, 3); err != nil {
		return types.Null(), err
	}
	val, _ := toNum("clamp", args[0])
	lo, _ := toNum("clamp", args[1])
	hi, _ := toNum("clamp", args[2])
	if val < lo {
		return types.Number(lo), nil
	}
	if val > hi {
		return types.Number(hi), nil
	}
	return types.Number(val), nil
}

func builtinSign(args []types.Value) (types.Value, error) {
	if err := requireArgs("sign", args, 1); err != nil {
		return types.Null(), err
	}
	f, _ := toNum("sign", args[0])
	switch {
	case f > 0:
		return types.Number(1), nil
	case f < 0:
		return types.Number(-1), nil
	default:
		return types.Number(0), nil
	}
}

func builtinHypot(args []types.Value) (types.Value, error) {
	if err := requireArgs("hypot", args, 2); err != nil {
		return types.Null(), err
	}
	a, _ := toNum("hypot", args[0])
	b, _ := toNum("hypot", args[1])
	return types.Number(math.Hypot(a, b)), nil
}

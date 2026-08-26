package builtin

import (
	"fmt"

	"expr-eval/internal/types"
)

func LogicFuncs() map[string]Func {
	return map[string]Func{
		"if":      builtinIf,
		"not":     builtinNot,
		"and":     builtinAnd,
		"or":      builtinOr,
		"eq":      builtinEq,
		"neq":     builtinNeq,
		"lt":      builtinLt,
		"lte":     builtinLte,
		"gt":      builtinGt,
		"gte":     builtinGte,
		"between": builtinBetween,
		"choose":  builtinChoose,
	}
}

func builtinIf(args []types.Value) (types.Value, error) {
	if err := requireArgs("if", args, 3); err != nil {
		return types.Null(), err
	}
	cond := args[0].IsTruthy()
	if cond {
		return args[1], nil
	}
	return args[2], nil
}

func builtinNot(args []types.Value) (types.Value, error) {
	if err := requireArgs("not", args, 1); err != nil {
		return types.Null(), err
	}
	return types.Bool(!args[0].IsTruthy()), nil
}

func builtinAnd(args []types.Value) (types.Value, error) {
	if err := requireMinArgs("and", args, 2); err != nil {
		return types.Null(), err
	}
	for _, a := range args {
		if !a.IsTruthy() {
			return types.Bool(false), nil
		}
	}
	return types.Bool(true), nil
}

func builtinOr(args []types.Value) (types.Value, error) {
	if err := requireMinArgs("or", args, 2); err != nil {
		return types.Null(), err
	}
	for _, a := range args {
		if a.IsTruthy() {
			return types.Bool(true), nil
		}
	}
	return types.Bool(false), nil
}

func builtinEq(args []types.Value) (types.Value, error) {
	if err := requireArgs("eq", args, 2); err != nil {
		return types.Null(), err
	}
	return types.Bool(types.Equal(args[0], args[1])), nil
}

func builtinNeq(args []types.Value) (types.Value, error) {
	if err := requireArgs("neq", args, 2); err != nil {
		return types.Null(), err
	}
	return types.Bool(!types.Equal(args[0], args[1])), nil
}

func builtinLt(args []types.Value) (types.Value, error) {
	if err := requireArgs("lt", args, 2); err != nil {
		return types.Null(), err
	}
	cmp, err := types.Compare(args[0], args[1])
	if err != nil {
		return types.Null(), err
	}
	return types.Bool(cmp < 0), nil
}

func builtinLte(args []types.Value) (types.Value, error) {
	if err := requireArgs("lte", args, 2); err != nil {
		return types.Null(), err
	}
	cmp, err := types.Compare(args[0], args[1])
	if err != nil {
		return types.Null(), err
	}
	return types.Bool(cmp <= 0), nil
}

func builtinGt(args []types.Value) (types.Value, error) {
	if err := requireArgs("gt", args, 2); err != nil {
		return types.Null(), err
	}
	cmp, err := types.Compare(args[0], args[1])
	if err != nil {
		return types.Null(), err
	}
	return types.Bool(cmp > 0), nil
}

func builtinGte(args []types.Value) (types.Value, error) {
	if err := requireArgs("gte", args, 2); err != nil {
		return types.Null(), err
	}
	cmp, err := types.Compare(args[0], args[1])
	if err != nil {
		return types.Null(), err
	}
	return types.Bool(cmp >= 0), nil
}

func builtinBetween(args []types.Value) (types.Value, error) {
	if err := requireArgs("between", args, 3); err != nil {
		return types.Null(), err
	}
	val, _ := args[0].AsNumber()
	lo, _ := args[1].AsNumber()
	hi, _ := args[2].AsNumber()
	return types.Bool(val >= lo && val <= hi), nil
}

func builtinChoose(args []types.Value) (types.Value, error) {
	if len(args) < 2 {
		return types.Null(), fmt.Errorf("choose requires at least 2 arguments")
	}
	idx, err := args[0].AsNumber()
	if err != nil {
		return types.Null(), fmt.Errorf("choose: first arg must be numeric index")
	}
	i := int(idx)
	if i < 0 || i >= len(args)-1 {
		return types.Null(), fmt.Errorf("choose: index %d out of range [0,%d)", i, len(args)-1)
	}
	return args[i+1], nil
}

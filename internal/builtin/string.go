package builtin

import (
	"fmt"
	"strings"

	"expr-eval/internal/types"
)

func StringFuncs() map[string]Func {
	return map[string]Func{
		"len":         builtinLen,
		"upper":       builtinUpper,
		"lower":       builtinLower,
		"trim":        builtinTrim,
		"contains":    builtinContains,
		"startswith":  builtinStartsWith,
		"endswith":    builtinEndsWith,
		"replace":     builtinReplace,
		"substr":      builtinSubstr,
		"concat":      builtinConcat,
		"repeat":      builtinRepeat,
		"index":       builtinIndex,
		"split_count": builtinSplitCount,
	}
}

func builtinLen(args []types.Value) (types.Value, error) {
	if err := requireArgs("len", args, 1); err != nil {
		return types.Null(), err
	}
	s := args[0].AsString()
	return types.Number(float64(len(s))), nil
}

func builtinUpper(args []types.Value) (types.Value, error) {
	if err := requireArgs("upper", args, 1); err != nil {
		return types.Null(), err
	}
	return types.String(strings.ToUpper(args[0].AsString())), nil
}

func builtinLower(args []types.Value) (types.Value, error) {
	if err := requireArgs("lower", args, 1); err != nil {
		return types.Null(), err
	}
	return types.String(strings.ToLower(args[0].AsString())), nil
}

func builtinTrim(args []types.Value) (types.Value, error) {
	if err := requireArgs("trim", args, 1); err != nil {
		return types.Null(), err
	}
	return types.String(strings.TrimSpace(args[0].AsString())), nil
}

func builtinContains(args []types.Value) (types.Value, error) {
	if err := requireArgs("contains", args, 2); err != nil {
		return types.Null(), err
	}
	return types.Bool(strings.Contains(args[0].AsString(), args[1].AsString())), nil
}

func builtinStartsWith(args []types.Value) (types.Value, error) {
	if err := requireArgs("startswith", args, 2); err != nil {
		return types.Null(), err
	}
	return types.Bool(strings.HasPrefix(args[0].AsString(), args[1].AsString())), nil
}

func builtinEndsWith(args []types.Value) (types.Value, error) {
	if err := requireArgs("endswith", args, 2); err != nil {
		return types.Null(), err
	}
	return types.Bool(strings.HasSuffix(args[0].AsString(), args[1].AsString())), nil
}

func builtinReplace(args []types.Value) (types.Value, error) {
	if err := requireArgs("replace", args, 3); err != nil {
		return types.Null(), err
	}
	s := args[0].AsString()
	old := args[1].AsString()
	new := args[2].AsString()
	return types.String(strings.ReplaceAll(s, old, new)), nil
}

func builtinSubstr(args []types.Value) (types.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return types.Null(), fmt.Errorf("substr requires 2 or 3 arguments, got %d", len(args))
	}
	s := args[0].AsString()
	start, _ := args[1].AsNumber()
	startI := int(start)
	if startI < 0 {
		startI = 0
	}
	if startI > len(s) {
		startI = len(s)
	}
	if len(args) == 2 {
		return types.String(s[startI:]), nil
	}
	length, _ := args[2].AsNumber()
	endI := startI + int(length)
	if endI > len(s) {
		endI = len(s)
	}
	if endI < startI {
		endI = startI
	}
	return types.String(s[startI:endI]), nil
}

func builtinConcat(args []types.Value) (types.Value, error) {
	if err := requireMinArgs("concat", args, 1); err != nil {
		return types.Null(), err
	}
	var b strings.Builder
	for _, a := range args {
		b.WriteString(a.AsString())
	}
	return types.String(b.String()), nil
}

func builtinRepeat(args []types.Value) (types.Value, error) {
	if err := requireArgs("repeat", args, 2); err != nil {
		return types.Null(), err
	}
	s := args[0].AsString()
	n, _ := args[1].AsNumber()
	count := int(n)
	if count < 0 {
		count = 0
	}
	if count > 10000 {
		return types.Null(), fmt.Errorf("repeat: count %d exceeds limit", count)
	}
	return types.String(strings.Repeat(s, count)), nil
}

func builtinIndex(args []types.Value) (types.Value, error) {
	if err := requireArgs("index", args, 2); err != nil {
		return types.Null(), err
	}
	s := args[0].AsString()
	sub := args[1].AsString()
	return types.Number(float64(strings.Index(s, sub))), nil
}

func builtinSplitCount(args []types.Value) (types.Value, error) {
	if err := requireArgs("split_count", args, 2); err != nil {
		return types.Null(), err
	}
	s := args[0].AsString()
	sep := args[1].AsString()
	parts := strings.Split(s, sep)
	return types.Number(float64(len(parts))), nil
}

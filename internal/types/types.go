package types

import (
	"fmt"
	"math"
	"strconv"
)

type Kind int

const (
	KindNumber Kind = iota
	KindBool
	KindString
	KindNull
)

func (k Kind) String() string {
	switch k {
	case KindNumber:
		return "number"
	case KindBool:
		return "bool"
	case KindString:
		return "string"
	case KindNull:
		return "null"
	default:
		return "unknown"
	}
}

type Value struct {
	kind    Kind
	numVal  float64
	boolVal bool
	strVal  string
}

func Number(f float64) Value { return Value{kind: KindNumber, numVal: f} }

func Bool(b bool) Value { return Value{kind: KindBool, boolVal: b} }

func String(s string) Value { return Value{kind: KindString, strVal: s} }

func Null() Value { return Value{kind: KindNull} }

func (v Value) Kind() Kind { return v.kind }

func (v Value) AsNumber() (float64, error) {
	switch v.kind {
	case KindNumber:
		return v.numVal, nil
	case KindBool:
		if v.boolVal {
			return 1, nil
		}
		return 0, nil
	case KindString:
		f, err := strconv.ParseFloat(v.strVal, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert string %q to number", v.strVal)
		}
		return f, nil
	default:
		return 0, fmt.Errorf("cannot convert %s to number", v.kind)
	}
}

func (v Value) AsBool() (bool, error) {
	switch v.kind {
	case KindBool:
		return v.boolVal, nil
	case KindNumber:
		return v.numVal != 0, nil
	case KindString:
		return v.strVal != "", nil
	case KindNull:
		return false, nil
	default:
		return false, fmt.Errorf("cannot convert %s to bool", v.kind)
	}
}

func (v Value) AsString() string {
	switch v.kind {
	case KindString:
		return v.strVal
	case KindNumber:
		if v.numVal == math.Trunc(v.numVal) && !math.IsInf(v.numVal, 0) {
			return strconv.FormatInt(int64(v.numVal), 10)
		}
		return strconv.FormatFloat(v.numVal, 'g', -1, 64)
	case KindBool:
		if v.boolVal {
			return "true"
		}
		return "false"
	case KindNull:
		return "null"
	default:
		return ""
	}
}

func (v Value) IsTruthy() bool {
	b, _ := v.AsBool()
	return b
}

func Equal(a, b Value) bool {
	if a.kind == b.kind {
		switch a.kind {
		case KindNumber:
			return a.numVal == b.numVal
		case KindBool:
			return a.boolVal == b.boolVal
		case KindString:
			return a.strVal == b.strVal
		case KindNull:
			return true
		}
	}
	return a.AsString() == b.AsString()
}

func Compare(a, b Value) (int, error) {
	if a.kind == KindNumber && b.kind == KindNumber {
		switch {
		case a.numVal < b.numVal:
			return -1, nil
		case a.numVal > b.numVal:
			return 1, nil
		default:
			return 0, nil
		}
	}
	if a.kind == KindString && b.kind == KindString {
		switch {
		case a.strVal < b.strVal:
			return -1, nil
		case a.strVal > b.strVal:
			return 1, nil
		default:
			return 0, nil
		}
	}
	an, errA := a.AsNumber()
	bn, errB := b.AsNumber()
	if errA == nil && errB == nil {
		switch {
		case an < bn:
			return -1, nil
		case an > bn:
			return 1, nil
		default:
			return 0, nil
		}
	}
	as := a.AsString()
	bs := b.AsString()
	switch {
	case as < bs:
		return -1, nil
	case as > bs:
		return 1, nil
	default:
		return 0, nil
	}
}

type TypeError struct {
	Op       string
	Got      Kind
	Expected Kind
}

func (e *TypeError) Error() string {
	return fmt.Sprintf("type error: %s expects %s, got %s", e.Op, e.Expected, e.Got)
}

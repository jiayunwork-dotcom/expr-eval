// Package types provides a multi-type value system for the expression evaluator.
// Values can be float64, bool, or string. The package handles coercion rules,
// type checking, and provides utilities for type-safe operations.
package types

import (
	"fmt"
	"math"
	"strconv"
)

// Kind identifies the type of a Value.
type Kind int

const (
	KindNumber Kind = iota
	KindBool
	KindString
	KindNull
)

// String returns the name of the kind.
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

// Value is the runtime representation of an expression result.
type Value struct {
	kind   Kind
	numVal float64
	boolVal bool
	strVal string
}

// Number creates a numeric Value.
func Number(f float64) Value { return Value{kind: KindNumber, numVal: f} }

// Bool creates a boolean Value.
func Bool(b bool) Value { return Value{kind: KindBool, boolVal: b} }

// String creates a string Value.
func String(s string) Value { return Value{kind: KindString, strVal: s} }

// Null creates a null Value.
func Null() Value { return Value{kind: KindNull} }

// Kind returns the type of this value.
func (v Value) Kind() Kind { return v.kind }

// AsNumber returns the numeric value or error if not numeric.
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

// AsBool returns the boolean value. Numbers: 0=false, else true. Strings: ""=false.
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

// AsString returns the string representation of the value.
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

// IsTruthy returns whether the value is considered true in boolean context.
func (v Value) IsTruthy() bool {
	b, _ := v.AsBool()
	return b
}

// Equal checks value equality with type coercion.
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
	// cross-type: compare as strings
	return a.AsString() == b.AsString()
}

// Compare returns -1, 0, or 1 comparing two values. Numbers compare numerically,
// strings lexicographically. Mixed types compare as strings.
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
	// mixed: try numeric comparison
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
	// fallback to string comparison
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

// TypeError is returned when an operation receives incompatible types.
type TypeError struct {
	Op       string
	Got      Kind
	Expected Kind
}

func (e *TypeError) Error() string {
	return fmt.Sprintf("type error: %s expects %s, got %s", e.Op, e.Expected, e.Got)
}

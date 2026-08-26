package types

import "testing"

func TestNumberAsString(t *testing.T) {
	v := Number(42)
	if s := v.AsString(); s != "42" {
		t.Errorf("got %q", s)
	}
	v2 := Number(3.14)
	if s := v2.AsString(); s != "3.14" {
		t.Errorf("got %q", s)
	}
}

func TestBoolIsTruthy(t *testing.T) {
	if !Bool(true).IsTruthy() {
		t.Error("true should be truthy")
	}
	if Bool(false).IsTruthy() {
		t.Error("false should not be truthy")
	}
}

func TestNumberIsTruthy(t *testing.T) {
	if Number(0).IsTruthy() {
		t.Error("0 should not be truthy")
	}
	if !Number(1).IsTruthy() {
		t.Error("1 should be truthy")
	}
}

func TestStringCoercion(t *testing.T) {
	v := String("3.5")
	f, err := v.AsNumber()
	if err != nil {
		t.Fatal(err)
	}
	if f != 3.5 {
		t.Errorf("got %f", f)
	}
}

func TestEqualSameType(t *testing.T) {
	if !Equal(Number(5), Number(5)) {
		t.Error("5 == 5")
	}
	if Equal(Number(5), Number(6)) {
		t.Error("5 != 6")
	}
}

func TestCompareNumbers(t *testing.T) {
	c, _ := Compare(Number(3), Number(7))
	if c != -1 {
		t.Errorf("3 < 7, got %d", c)
	}
}

func TestNullIsFalsy(t *testing.T) {
	if Null().IsTruthy() {
		t.Error("null should be falsy")
	}
}

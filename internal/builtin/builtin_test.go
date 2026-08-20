package builtin

import (
	"testing"

	"expr-eval/internal/types"
)

func TestMathAbs(t *testing.T) {
	fn := MathFuncs()["abs"]
	v, err := fn([]types.Value{types.Number(-5)})
	if err != nil {
		t.Fatal(err)
	}
	f, _ := v.AsNumber()
	if f != 5 {
		t.Errorf("abs(-5) = %f", f)
	}
}

func TestMathClamp(t *testing.T) {
	fn := MathFuncs()["clamp"]
	v, _ := fn([]types.Value{types.Number(15), types.Number(0), types.Number(10)})
	f, _ := v.AsNumber()
	if f != 10 {
		t.Errorf("clamp(15,0,10) = %f, want 10", f)
	}
}

func TestStringLen(t *testing.T) {
	fn := StringFuncs()["len"]
	v, _ := fn([]types.Value{types.String("hello")})
	f, _ := v.AsNumber()
	if f != 5 {
		t.Errorf("len(hello) = %f", f)
	}
}

func TestStringContains(t *testing.T) {
	fn := StringFuncs()["contains"]
	v, _ := fn([]types.Value{types.String("hello world"), types.String("world")})
	b, _ := v.AsBool()
	if !b {
		t.Error("expected true")
	}
}

func TestLogicIf(t *testing.T) {
	fn := LogicFuncs()["if"]
	v, _ := fn([]types.Value{types.Bool(true), types.Number(1), types.Number(2)})
	f, _ := v.AsNumber()
	if f != 1 {
		t.Errorf("if(true,1,2) = %f", f)
	}
}

func TestLogicBetween(t *testing.T) {
	fn := LogicFuncs()["between"]
	v, _ := fn([]types.Value{types.Number(5), types.Number(1), types.Number(10)})
	b, _ := v.AsBool()
	if !b {
		t.Error("5 between 1 and 10")
	}
}

func TestRegistryHas(t *testing.T) {
	r := NewRegistry()
	if !r.Has("abs") {
		t.Error("registry should have abs")
	}
	if !r.Has("len") {
		t.Error("registry should have len")
	}
	if !r.Has("if") {
		t.Error("registry should have if")
	}
	if r.Count() < 30 {
		t.Errorf("count = %d, want >=30", r.Count())
	}
}

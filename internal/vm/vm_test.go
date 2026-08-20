package vm

import (
	"math"
	"testing"

	"expr-eval/internal/compile"
	"expr-eval/internal/parser"
)

func compileExpr(expr string) *compile.Program {
	node, _ := parser.Parse(expr)
	prog, _ := compile.Compile(node)
	return prog
}

func TestVMSimpleAdd(t *testing.T) {
	vm := New()
	result, err := vm.Run(compileExpr("2 + 3"), nil)
	if err != nil {
		t.Fatal(err)
	}
	if result != 5 {
		t.Errorf("2+3 = %f", result)
	}
}

func TestVMWithVars(t *testing.T) {
	vm := New()
	vars := map[string]float64{"x": 10, "y": 3}
	result, _ := vm.Run(compileExpr("x * y + 1"), vars)
	if result != 31 {
		t.Errorf("10*3+1 = %f, want 31", result)
	}
}

func TestVMFunction(t *testing.T) {
	vm := New()
	result, _ := vm.Run(compileExpr("sqrt(16)"), nil)
	if result != 4 {
		t.Errorf("sqrt(16) = %f", result)
	}
}

func TestVMDivisionByZero(t *testing.T) {
	vm := New()
	_, err := vm.Run(compileExpr("1 / 0"), nil)
	if err == nil {
		t.Error("expected division by zero error")
	}
}

func TestVMUnknownVariable(t *testing.T) {
	vm := New()
	_, err := vm.Run(compileExpr("x"), nil)
	if err == nil {
		t.Error("expected unknown variable error")
	}
}

func TestVMNeg(t *testing.T) {
	vm := New()
	result, _ := vm.Run(compileExpr("-5"), nil)
	if result != -5 {
		t.Errorf("-5 = %f", result)
	}
}

func TestVMPow(t *testing.T) {
	vm := New()
	result, _ := vm.Run(compileExpr("2 ^ 10"), nil)
	if result != 1024 {
		t.Errorf("2^10 = %f", result)
	}
}

func TestVMMinMax(t *testing.T) {
	vm := New()
	r1, _ := vm.Run(compileExpr("min(5, 3)"), nil)
	if r1 != 3 {
		t.Errorf("min(5,3) = %f", r1)
	}
	r2, _ := vm.Run(compileExpr("max(5, 3)"), nil)
	if r2 != 5 {
		t.Errorf("max(5,3) = %f", r2)
	}
}

func TestVMComplexExpr(t *testing.T) {
	vm := New()
	vars := map[string]float64{"r": 3}
	result, _ := vm.Run(compileExpr("3.14159 * r ^ 2"), vars)
	expected := 3.14159 * 9
	if math.Abs(result-expected) > 0.001 {
		t.Errorf("pi*r^2 = %f, want %f", result, expected)
	}
}

func TestVMTrace(t *testing.T) {
	vm := New()
	prog := compileExpr("2 + 3")
	result, trace, err := vm.Trace(prog, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result != 5 {
		t.Errorf("result = %f", result)
	}
	if len(trace) == 0 {
		t.Error("trace should not be empty")
	}
}

func TestVMCustomFunc(t *testing.T) {
	vm := New()
	vm.RegisterFunc("double", func(args []float64) (float64, error) {
		return args[0] * 2, nil
	})
	prog := compileExpr("double(7)")
	result, _ := vm.Run(prog, nil)
	if result != 14 {
		t.Errorf("double(7) = %f", result)
	}
}

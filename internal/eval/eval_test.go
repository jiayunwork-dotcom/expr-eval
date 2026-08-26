package eval

import (
	"testing"

	"expr-eval/internal/parser"
)

func TestEvalBasic(t *testing.T) {
	cases := []struct {
		expr string
		vars map[string]float64
		want float64
	}{
		{"2 + 3 * 4", nil, 14},
		{"x * 2", map[string]float64{"x": 5}, 10},
		{"sqrt(16)", nil, 4},
		{"2 ^ 3 ^ 2", nil, 512},
		{"-2 + 3", nil, 1},
		{"(1 + 2) * 3", nil, 9},
		{"10 % 3", nil, 1},
	}
	for _, c := range cases {
		node, err := parser.Parse(c.expr)
		if err != nil {
			t.Fatalf("parse %q: %v", c.expr, err)
		}
		got, err := Eval(node, c.vars)
		if err != nil {
			t.Fatalf("eval %q: %v", c.expr, err)
		}
		if got != c.want {
			t.Fatalf("%q = %v want %v", c.expr, got, c.want)
		}
	}
}

func TestEvalUnknownVar(t *testing.T) {
	node, _ := parser.Parse("y + 1")
	if _, err := Eval(node, nil); err == nil {
		t.Fatal("expected error for nil vars map")
	}
	node2, _ := parser.Parse("z")
	if _, err := Eval(node2, map[string]float64{}); err == nil {
		t.Fatal("expected error for unknown variable z")
	}
}

func TestEvalDivZero(t *testing.T) {
	node, _ := parser.Parse("1 / 0")
	if _, err := Eval(node, nil); err == nil {
		t.Fatal("expected division by zero error")
	}
}

func TestEvalSqrtNeg(t *testing.T) {
	node, _ := parser.Parse("sqrt(-1)")
	if _, err := Eval(node, nil); err == nil {
		t.Fatal("expected error for sqrt of negative")
	}
}

func TestEvalFunctions(t *testing.T) {
	node, err := parser.Parse("min(3, 1, 2)")
	if err != nil {
		t.Fatal(err)
	}
	got, err := Eval(node, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != 1 {
		t.Fatalf("min = %v want 1", got)
	}
	node2, _ := parser.Parse("max(3, 1, 2)")
	got2, _ := Eval(node2, nil)
	if got2 != 3 {
		t.Fatalf("max = %v want 3", got2)
	}
	node3, _ := parser.Parse("abs(-7)")
	got3, _ := Eval(node3, nil)
	if got3 != 7 {
		t.Fatalf("abs = %v want 7", got3)
	}
}

func TestEvalUnknownFunc(t *testing.T) {
	node, _ := parser.Parse("foo(1)")
	if _, err := Eval(node, nil); err == nil {
		t.Fatal("expected error for unknown function")
	}
}

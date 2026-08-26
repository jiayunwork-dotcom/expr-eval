package validate

import (
	"testing"

	"expr-eval/internal/parser"
)

func TestValidateUndefinedVar(t *testing.T) {
	node, _ := parser.Parse("x + y")
	result := Validate(node, []string{"x"}, DefaultFuncs())
	if result.OK() {
		t.Error("expected error for undefined y")
	}
	errs := result.Errors()
	if len(errs) != 1 {
		t.Fatalf("errors = %d, want 1", len(errs))
	}
	if errs[0].Message != `undefined variable "y"` {
		t.Errorf("msg = %q", errs[0].Message)
	}
}

func TestValidateAllDefined(t *testing.T) {
	node, _ := parser.Parse("x + y")
	result := Validate(node, []string{"x", "y"}, DefaultFuncs())
	if !result.OK() {
		t.Errorf("expected OK, got %+v", result.Issues)
	}
}

func TestValidateUnknownFunction(t *testing.T) {
	node, _ := parser.Parse("foo(1)")
	result := Validate(node, nil, DefaultFuncs())
	if result.OK() {
		t.Error("expected error for unknown function foo")
	}
}

func TestValidateArityError(t *testing.T) {
	node, _ := parser.Parse("abs(1, 2)")
	result := Validate(node, nil, DefaultFuncs())
	if result.OK() {
		t.Error("expected error for abs with 2 args")
	}
}

func TestValidateDivByZeroWarning(t *testing.T) {
	node, _ := parser.Parse("x / 0")
	result := Validate(node, []string{"x"}, DefaultFuncs())
	if len(result.Issues) != 1 {
		t.Fatalf("issues = %d", len(result.Issues))
	}
	if result.Issues[0].Severity != "warning" {
		t.Errorf("severity = %s, want warning", result.Issues[0].Severity)
	}
}

func TestUndefinedVars(t *testing.T) {
	node, _ := parser.Parse("a + b * c")
	undef := UndefinedVars(node, []string{"a"})
	if len(undef) != 2 {
		t.Errorf("undefined = %v, want [b c]", undef)
	}
}

func TestUsedVars(t *testing.T) {
	node, _ := parser.Parse("min(x, y) + z")
	vars := UsedVars(node)
	if len(vars) != 3 {
		t.Errorf("used = %v", vars)
	}
}

func TestComplexity(t *testing.T) {
	node, _ := parser.Parse("(a + b) * (c - d)")
	c := Complexity(node)
	if c != 3 {
		t.Errorf("complexity = %d, want 3", c)
	}
}

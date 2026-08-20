package format

import (
	"strings"
	"testing"

	"expr-eval/internal/parser"
)

func TestSprintSimple(t *testing.T) {
	node, _ := parser.Parse("2 + 3")
	s := Sprint(node)
	if !strings.Contains(s, "2") || !strings.Contains(s, "3") || !strings.Contains(s, "+") {
		t.Errorf("sprint = %q", s)
	}
}

func TestSprintParenthesized(t *testing.T) {
	node, _ := parser.Parse("a + b * c")
	s := SprintParenthesized(node)
	if !strings.Contains(s, "(") {
		t.Errorf("expected parens in %q", s)
	}
}

func TestSprintFunction(t *testing.T) {
	node, _ := parser.Parse("sqrt(x)")
	s := Sprint(node)
	if s != "sqrt(x)" {
		t.Errorf("sprint = %q, want sqrt(x)", s)
	}
}

func TestIndent(t *testing.T) {
	node, _ := parser.Parse("a + b")
	s := Indent(node)
	if !strings.Contains(s, "Binary") {
		t.Errorf("indent should contain Binary: %q", s)
	}
	if !strings.Contains(s, "Ident") {
		t.Errorf("indent should contain Ident: %q", s)
	}
}

func TestRPNString(t *testing.T) {
	node, _ := parser.Parse("2 + 3 * 4")
	rpn := RPNString(node)
	// 2 3 4 * + (due to operator precedence)
	if !strings.Contains(rpn, "2") || !strings.Contains(rpn, "*") || !strings.Contains(rpn, "+") {
		t.Errorf("rpn = %q", rpn)
	}
}

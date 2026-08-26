package validate

import (
	"fmt"
	"sort"

	"expr-eval/internal/parser"
)

type Issue struct {
	Pos      int
	Message  string
	Severity string
}

type Result struct {
	Issues []Issue
}

func (r *Result) OK() bool {
	for _, iss := range r.Issues {
		if iss.Severity == "error" {
			return false
		}
	}
	return true
}

func (r *Result) Errors() []Issue {
	var out []Issue
	for _, iss := range r.Issues {
		if iss.Severity == "error" {
			out = append(out, iss)
		}
	}
	return out
}

type FuncSpec struct {
	MinArgs int
	MaxArgs int
}

func DefaultFuncs() map[string]FuncSpec {
	return map[string]FuncSpec{
		"abs":   {1, 1},
		"sqrt":  {1, 1},
		"ceil":  {1, 1},
		"floor": {1, 1},
		"round": {1, 1},
		"log":   {1, 1},
		"log2":  {1, 1},
		"log10": {1, 1},
		"exp":   {1, 1},
		"sin":   {1, 1},
		"cos":   {1, 1},
		"tan":   {1, 1},
		"asin":  {1, 1},
		"acos":  {1, 1},
		"atan":  {1, 1},
		"atan2": {2, 2},
		"pow":   {2, 2},
		"min":   {1, -1},
		"max":   {1, -1},
		"clamp": {3, 3},
		"sign":  {1, 1},
		"hypot": {2, 2},
		"len":   {1, 1},
		"upper": {1, 1},
		"lower": {1, 1},
		"trim":  {1, 1},
	}
}

func Validate(node parser.Node, knownVars []string, funcs map[string]FuncSpec) *Result {
	r := &Result{}
	varSet := map[string]bool{}
	for _, v := range knownVars {
		varSet[v] = true
	}
	validate(node, varSet, funcs, r)
	return r
}

func validate(node parser.Node, vars map[string]bool, funcs map[string]FuncSpec, r *Result) {
	switch n := node.(type) {
	case *parser.NumberNode:
	case *parser.IdentNode:
		if !vars[n.Name] {
			r.Issues = append(r.Issues, Issue{
				Pos:      n.Pos(),
				Message:  fmt.Sprintf("undefined variable %q", n.Name),
				Severity: "error",
			})
		}
	case *parser.UnaryNode:
		validate(n.Expr, vars, funcs, r)
	case *parser.BinaryNode:
		validate(n.Left, vars, funcs, r)
		validate(n.Right, vars, funcs, r)
		if n.Op == "/" || n.Op == "%" {
			if num, ok := n.Right.(*parser.NumberNode); ok && num.Val == 0 {
				r.Issues = append(r.Issues, Issue{
					Pos:      n.Pos(),
					Message:  fmt.Sprintf("division/modulo by zero literal"),
					Severity: "warning",
				})
			}
		}
	case *parser.CallNode:
		spec, known := funcs[n.Name]
		if !known {
			r.Issues = append(r.Issues, Issue{
				Pos:      n.Pos(),
				Message:  fmt.Sprintf("unknown function %q", n.Name),
				Severity: "error",
			})
		} else {
			argc := len(n.Args)
			if argc < spec.MinArgs {
				r.Issues = append(r.Issues, Issue{
					Pos:      n.Pos(),
					Message:  fmt.Sprintf("%s requires at least %d argument(s), got %d", n.Name, spec.MinArgs, argc),
					Severity: "error",
				})
			}
			if spec.MaxArgs >= 0 && argc > spec.MaxArgs {
				r.Issues = append(r.Issues, Issue{
					Pos:      n.Pos(),
					Message:  fmt.Sprintf("%s accepts at most %d argument(s), got %d", n.Name, spec.MaxArgs, argc),
					Severity: "error",
				})
			}
		}
		for _, arg := range n.Args {
			validate(arg, vars, funcs, r)
		}
	}
}

func UndefinedVars(node parser.Node, knownVars []string) []string {
	varSet := map[string]bool{}
	for _, v := range knownVars {
		varSet[v] = true
	}
	used := map[string]bool{}
	collectVars(node, used)
	var undefined []string
	for v := range used {
		if !varSet[v] {
			undefined = append(undefined, v)
		}
	}
	sort.Strings(undefined)
	return undefined
}

func collectVars(node parser.Node, vars map[string]bool) {
	switch n := node.(type) {
	case *parser.IdentNode:
		vars[n.Name] = true
	case *parser.UnaryNode:
		collectVars(n.Expr, vars)
	case *parser.BinaryNode:
		collectVars(n.Left, vars)
		collectVars(n.Right, vars)
	case *parser.CallNode:
		for _, a := range n.Args {
			collectVars(a, vars)
		}
	}
}

func UsedVars(node parser.Node) []string {
	used := map[string]bool{}
	collectVars(node, used)
	out := make([]string, 0, len(used))
	for v := range used {
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

func Complexity(node parser.Node) int {
	switch n := node.(type) {
	case *parser.NumberNode, *parser.IdentNode:
		return 0
	case *parser.UnaryNode:
		return 1 + Complexity(n.Expr)
	case *parser.BinaryNode:
		return 1 + Complexity(n.Left) + Complexity(n.Right)
	case *parser.CallNode:
		c := 1
		for _, a := range n.Args {
			c += Complexity(a)
		}
		return c
	default:
		return 0
	}
}

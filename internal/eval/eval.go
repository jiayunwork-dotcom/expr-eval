package eval

import (
	"fmt"
	"math"

	"expr-eval/internal/parser"
)

func Eval(node parser.Node, vars map[string]float64) (float64, error) {
	v, err := evalAST(node, vars)
	if err != nil {
		return 0, err
	}
	return fillEval(v), nil
}

func evalAST(node parser.Node, vars map[string]float64) (float64, error) {
	switch n := node.(type) {
	case *parser.NumberNode:
		return n.Val, nil
	case *parser.IdentNode:
		v, ok := vars[n.Name]
		if !ok {
			return 0, fmt.Errorf("unknown variable %q", n.Name)
		}
		return v, nil
	case *parser.UnaryNode:
		val, err := evalAST(n.Expr, vars)
		if err != nil {
			return 0, err
		}
		if n.Op == '-' {
			return -val, nil
		}
		return val, nil
	case *parser.BinaryNode:
		return evalBinary(n, vars)
	case *parser.CallNode:
		return evalCall(n, vars)
	default:
		return 0, fmt.Errorf("unsupported node type %T", node)
	}
}

func evalBinary(n *parser.BinaryNode, vars map[string]float64) (float64, error) {
	left, err := evalAST(n.Left, vars)
	if err != nil {
		return 0, err
	}
	right, err := evalAST(n.Right, vars)
	if err != nil {
		return 0, err
	}
	switch n.Op {
	case "+":
		return left + right, nil
	case "-":
		return left - right, nil
	case "*":
		return left * right, nil
	case "/":
		if right == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return left / right, nil
	case "%":
		if right == 0 {
			return 0, fmt.Errorf("modulo by zero")
		}
		return math.Mod(left, right), nil
	case "^":
		return math.Pow(left, right), nil
	default:
		return 0, fmt.Errorf("unknown operator %q", n.Op)
	}
}

func evalCall(n *parser.CallNode, vars map[string]float64) (float64, error) {
	args := make([]float64, len(n.Args))
	for i, a := range n.Args {
		v, err := evalAST(a, vars)
		if err != nil {
			return 0, err
		}
		args[i] = v
	}
	switch n.Name {
	case "min":
		if len(args) == 0 {
			return 0, fmt.Errorf("min requires at least 1 argument")
		}
		m := args[0]
		for _, a := range args[1:] {
			if a < m {
				m = a
			}
		}
		return m, nil
	case "max":
		if len(args) == 0 {
			return 0, fmt.Errorf("max requires at least 1 argument")
		}
		m := args[0]
		for _, a := range args[1:] {
			if a > m {
				m = a
			}
		}
		return m, nil
	case "abs":
		if len(args) != 1 {
			return 0, fmt.Errorf("abs requires exactly 1 argument")
		}
		return math.Abs(args[0]), nil
	case "sqrt":
		if len(args) != 1 {
			return 0, fmt.Errorf("sqrt requires exactly 1 argument")
		}
		if args[0] < 0 {
			return 0, fmt.Errorf("square root of negative number")
		}
		return math.Sqrt(args[0]), nil
	default:
		return 0, fmt.Errorf("unknown function %q", n.Name)
	}
}

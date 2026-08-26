package format

import (
	"fmt"
	"strings"

	"expr-eval/internal/parser"
)

func Sprint(node parser.Node) string {
	return sprintNode(node, false)
}

func SprintParenthesized(node parser.Node) string {
	return sprintNode(node, true)
}

func sprintNode(node parser.Node, parens bool) string {
	switch n := node.(type) {
	case *parser.NumberNode:
		if n.Val == float64(int64(n.Val)) {
			return fmt.Sprintf("%d", int64(n.Val))
		}
		return fmt.Sprintf("%g", n.Val)
	case *parser.IdentNode:
		return n.Name
	case *parser.UnaryNode:
		inner := sprintNode(n.Expr, parens)
		if n.Op == '-' {
			return fmt.Sprintf("(-%s)", inner)
		}
		return inner
	case *parser.BinaryNode:
		left := sprintNode(n.Left, parens)
		right := sprintNode(n.Right, parens)
		if parens {
			return fmt.Sprintf("(%s %s %s)", left, n.Op, right)
		}
		return fmt.Sprintf("%s %s %s", left, n.Op, right)
	case *parser.CallNode:
		args := make([]string, len(n.Args))
		for i, a := range n.Args {
			args[i] = sprintNode(a, parens)
		}
		return fmt.Sprintf("%s(%s)", n.Name, strings.Join(args, ", "))
	default:
		return "<?>"
	}
}

func Indent(node parser.Node) string {
	var b strings.Builder
	indent(node, &b, 0)
	return b.String()
}

func indent(node parser.Node, b *strings.Builder, level int) {
	prefix := strings.Repeat("  ", level)
	switch n := node.(type) {
	case *parser.NumberNode:
		fmt.Fprintf(b, "%sNumber(%g)\n", prefix, n.Val)
	case *parser.IdentNode:
		fmt.Fprintf(b, "%sIdent(%s)\n", prefix, n.Name)
	case *parser.UnaryNode:
		fmt.Fprintf(b, "%sUnary(%c)\n", prefix, n.Op)
		indent(n.Expr, b, level+1)
	case *parser.BinaryNode:
		fmt.Fprintf(b, "%sBinary(%s)\n", prefix, n.Op)
		indent(n.Left, b, level+1)
		indent(n.Right, b, level+1)
	case *parser.CallNode:
		fmt.Fprintf(b, "%sCall(%s, %d args)\n", prefix, n.Name, len(n.Args))
		for _, a := range n.Args {
			indent(a, b, level+1)
		}
	default:
		fmt.Fprintf(b, "%s<?>\n", prefix)
	}
}

func RPNString(node parser.Node) string {
	var parts []string
	rpn(node, &parts)
	return strings.Join(parts, " ")
}

func rpn(node parser.Node, parts *[]string) {
	switch n := node.(type) {
	case *parser.NumberNode:
		if n.Val == float64(int64(n.Val)) {
			*parts = append(*parts, fmt.Sprintf("%d", int64(n.Val)))
		} else {
			*parts = append(*parts, fmt.Sprintf("%g", n.Val))
		}
	case *parser.IdentNode:
		*parts = append(*parts, n.Name)
	case *parser.UnaryNode:
		rpn(n.Expr, parts)
		*parts = append(*parts, "NEG")
	case *parser.BinaryNode:
		rpn(n.Left, parts)
		rpn(n.Right, parts)
		*parts = append(*parts, n.Op)
	case *parser.CallNode:
		for _, a := range n.Args {
			rpn(a, parts)
		}
		*parts = append(*parts, fmt.Sprintf("%s/%d", n.Name, len(n.Args)))
	}
}

package ast

// Visitor is called for each node during tree traversal.
type Visitor func(node Node) bool

// Walk traverses the AST in pre-order. If the visitor returns false,
// children of that node are not visited.
func Walk(node Node, visit Visitor) {
	if node == nil {
		return
	}
	if !visit(node) {
		return
	}
	switch n := node.(type) {
	case *UnaryNode:
		Walk(n.Expr, visit)
	case *BinaryNode:
		Walk(n.Left, visit)
		Walk(n.Right, visit)
	case *CallNode:
		for _, arg := range n.Args {
			Walk(arg, visit)
		}
	case *CondNode:
		Walk(n.Cond, visit)
		Walk(n.Then, visit)
		Walk(n.Else, visit)
	case *IndexNode:
		Walk(n.Expr, visit)
		Walk(n.Index, visit)
	}
}

// Collect returns all nodes in pre-order traversal.
func Collect(node Node) []Node {
	var nodes []Node
	Walk(node, func(n Node) bool {
		nodes = append(nodes, n)
		return true
	})
	return nodes
}

// Variables returns all unique variable names referenced in the AST.
func Variables(node Node) []string {
	seen := map[string]bool{}
	var result []string
	Walk(node, func(n Node) bool {
		if id, ok := n.(*IdentNode); ok {
			if !seen[id.Name] {
				seen[id.Name] = true
				result = append(result, id.Name)
			}
		}
		return true
	})
	return result
}

// FunctionCalls returns all unique function names called in the AST.
func FunctionCalls(node Node) []string {
	seen := map[string]bool{}
	var result []string
	Walk(node, func(n Node) bool {
		if call, ok := n.(*CallNode); ok {
			if !seen[call.Name] {
				seen[call.Name] = true
				result = append(result, call.Name)
			}
		}
		return true
	})
	return result
}

// Depth returns the maximum nesting depth of the AST.
func Depth(node Node) int {
	if node == nil {
		return 0
	}
	switch n := node.(type) {
	case *UnaryNode:
		return 1 + Depth(n.Expr)
	case *BinaryNode:
		ld := Depth(n.Left)
		rd := Depth(n.Right)
		if ld > rd {
			return 1 + ld
		}
		return 1 + rd
	case *CallNode:
		maxD := 0
		for _, arg := range n.Args {
			d := Depth(arg)
			if d > maxD {
				maxD = d
			}
		}
		return 1 + maxD
	case *CondNode:
		cd := Depth(n.Cond)
		td := Depth(n.Then)
		ed := Depth(n.Else)
		best := cd
		if td > best {
			best = td
		}
		if ed > best {
			best = ed
		}
		return 1 + best
	case *IndexNode:
		ed := Depth(n.Expr)
		id := Depth(n.Index)
		if ed > id {
			return 1 + ed
		}
		return 1 + id
	default:
		return 1
	}
}

// NodeCount returns the total number of nodes in the AST.
func NodeCount(node Node) int {
	count := 0
	Walk(node, func(_ Node) bool {
		count++
		return true
	})
	return count
}

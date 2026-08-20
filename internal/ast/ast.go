// Package ast defines the abstract syntax tree node types for the expression
// language. It supports numeric, boolean, and string literals, variables,
// unary/binary/comparison/logical operators, conditional expressions, and
// function calls.
package ast

// Node is the interface all AST nodes implement.
type Node interface {
	Pos() int
	nodeTag()
}

// NumberNode represents a floating-point literal.
type NumberNode struct {
	Val  float64
	Posi int
}

func (n *NumberNode) Pos() int { return n.Posi }
func (n *NumberNode) nodeTag() {}

// BoolNode represents a boolean literal (true/false).
type BoolNode struct {
	Val  bool
	Posi int
}

func (n *BoolNode) Pos() int { return n.Posi }
func (n *BoolNode) nodeTag() {}

// StringNode represents a string literal.
type StringNode struct {
	Val  string
	Posi int
}

func (n *StringNode) Pos() int { return n.Posi }
func (n *StringNode) nodeTag() {}

// IdentNode represents a variable reference.
type IdentNode struct {
	Name string
	Posi int
}

func (n *IdentNode) Pos() int { return n.Posi }
func (n *IdentNode) nodeTag() {}

// UnaryOp is a unary operator kind.
type UnaryOp int

const (
	OpNeg UnaryOp = iota // -
	OpNot                // !
)

// UnaryNode represents a unary operation.
type UnaryNode struct {
	Op   UnaryOp
	Expr Node
	Posi int
}

func (n *UnaryNode) Pos() int { return n.Posi }
func (n *UnaryNode) nodeTag() {}

// BinaryOp is a binary operator kind.
type BinaryOp int

const (
	OpAdd BinaryOp = iota
	OpSub
	OpMul
	OpDiv
	OpMod
	OpPow
	// Comparison
	OpEq
	OpNeq
	OpLt
	OpLte
	OpGt
	OpGte
	// Logical
	OpAnd
	OpOr
)

// BinaryNode represents a binary operation.
type BinaryNode struct {
	Op    BinaryOp
	Left  Node
	Right Node
	Posi  int
}

func (n *BinaryNode) Pos() int { return n.Posi }
func (n *BinaryNode) nodeTag() {}

// CallNode represents a function call.
type CallNode struct {
	Name string
	Args []Node
	Posi int
}

func (n *CallNode) Pos() int { return n.Posi }
func (n *CallNode) nodeTag() {}

// CondNode represents a conditional expression: cond ? then : else.
type CondNode struct {
	Cond Node
	Then Node
	Else Node
	Posi int
}

func (n *CondNode) Pos() int { return n.Posi }
func (n *CondNode) nodeTag() {}

// IndexNode represents array/string indexing: expr[index].
type IndexNode struct {
	Expr  Node
	Index Node
	Posi  int
}

func (n *IndexNode) Pos() int { return n.Posi }
func (n *IndexNode) nodeTag() {}

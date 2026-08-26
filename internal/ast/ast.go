package ast

type Node interface {
	Pos() int
	nodeTag()
}

type NumberNode struct {
	Val  float64
	Posi int
}

func (n *NumberNode) Pos() int { return n.Posi }
func (n *NumberNode) nodeTag() {}

type BoolNode struct {
	Val  bool
	Posi int
}

func (n *BoolNode) Pos() int { return n.Posi }
func (n *BoolNode) nodeTag() {}

type StringNode struct {
	Val  string
	Posi int
}

func (n *StringNode) Pos() int { return n.Posi }
func (n *StringNode) nodeTag() {}

type IdentNode struct {
	Name string
	Posi int
}

func (n *IdentNode) Pos() int { return n.Posi }
func (n *IdentNode) nodeTag() {}

type UnaryOp int

const (
	OpNeg UnaryOp = iota
	OpNot
)

type UnaryNode struct {
	Op   UnaryOp
	Expr Node
	Posi int
}

func (n *UnaryNode) Pos() int { return n.Posi }
func (n *UnaryNode) nodeTag() {}

type BinaryOp int

const (
	OpAdd BinaryOp = iota
	OpSub
	OpMul
	OpDiv
	OpMod
	OpPow
	OpEq
	OpNeq
	OpLt
	OpLte
	OpGt
	OpGte
	OpAnd
	OpOr
)

type BinaryNode struct {
	Op    BinaryOp
	Left  Node
	Right Node
	Posi  int
}

func (n *BinaryNode) Pos() int { return n.Posi }
func (n *BinaryNode) nodeTag() {}

type CallNode struct {
	Name string
	Args []Node
	Posi int
}

func (n *CallNode) Pos() int { return n.Posi }
func (n *CallNode) nodeTag() {}

type CondNode struct {
	Cond Node
	Then Node
	Else Node
	Posi int
}

func (n *CondNode) Pos() int { return n.Posi }
func (n *CondNode) nodeTag() {}

type IndexNode struct {
	Expr  Node
	Index Node
	Posi  int
}

func (n *IndexNode) Pos() int { return n.Posi }
func (n *IndexNode) nodeTag() {}

package ast

import "testing"

func TestWalkCollect(t *testing.T) {
	tree := &BinaryNode{
		Op:   OpAdd,
		Left: &NumberNode{Val: 1},
		Right: &BinaryNode{
			Op:    OpMul,
			Left:  &IdentNode{Name: "x"},
			Right: &NumberNode{Val: 2},
		},
	}
	nodes := Collect(tree)
	if len(nodes) != 5 {
		t.Errorf("nodes = %d, want 5", len(nodes))
	}
}

func TestVariables(t *testing.T) {
	tree := &BinaryNode{
		Op:    OpAdd,
		Left:  &IdentNode{Name: "x"},
		Right: &IdentNode{Name: "y"},
	}
	vars := Variables(tree)
	if len(vars) != 2 {
		t.Errorf("vars = %v", vars)
	}
}

func TestDepth(t *testing.T) {
	tree := &BinaryNode{
		Op:   OpAdd,
		Left: &NumberNode{Val: 1},
		Right: &BinaryNode{
			Op:    OpMul,
			Left:  &NumberNode{Val: 2},
			Right: &NumberNode{Val: 3},
		},
	}
	if d := Depth(tree); d != 3 {
		t.Errorf("depth = %d, want 3", d)
	}
}

func TestNodeCount(t *testing.T) {
	tree := &CallNode{
		Name: "min",
		Args: []Node{&NumberNode{Val: 1}, &NumberNode{Val: 2}},
	}
	if c := NodeCount(tree); c != 3 {
		t.Errorf("count = %d, want 3", c)
	}
}

func TestFunctionCalls(t *testing.T) {
	tree := &CallNode{
		Name: "max",
		Args: []Node{&CallNode{Name: "abs", Args: []Node{&IdentNode{Name: "x"}}}},
	}
	fns := FunctionCalls(tree)
	if len(fns) != 2 {
		t.Errorf("fns = %v", fns)
	}
}

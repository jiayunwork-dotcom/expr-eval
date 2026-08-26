package parser

import "testing"

func TestParsePrecedence(t *testing.T) {
	node, err := Parse("2 + 3 * 4")
	if err != nil {
		t.Fatal(err)
	}
	bin, ok := node.(*BinaryNode)
	if !ok {
		t.Fatalf("expected *BinaryNode, got %T", node)
	}
	if bin.Op != "+" {
		t.Fatalf("top op = %q", bin.Op)
	}
	rbin, ok := bin.Right.(*BinaryNode)
	if !ok {
		t.Fatalf("expected right *BinaryNode, got %T", bin.Right)
	}
	if rbin.Op != "*" {
		t.Fatalf("right op = %q", rbin.Op)
	}
}

func TestParseRightAssocPower(t *testing.T) {
	node, err := Parse("2 ^ 3 ^ 2")
	if err != nil {
		t.Fatal(err)
	}
	bin, ok := node.(*BinaryNode)
	if !ok {
		t.Fatalf("expected *BinaryNode, got %T", node)
	}
	if bin.Op != "^" {
		t.Fatalf("op = %q", bin.Op)
	}
	rbin, ok := bin.Right.(*BinaryNode)
	if !ok {
		t.Fatalf("right should be *BinaryNode for right-assoc, got %T", bin.Right)
	}
	if rbin.Op != "^" {
		t.Fatalf("right op = %q", rbin.Op)
	}
}

func TestParseParenMismatch(t *testing.T) {
	if _, err := Parse("(1 + 2"); err == nil {
		t.Fatal("expected error for mismatched paren")
	}
}

func TestParseEmpty(t *testing.T) {
	if _, err := Parse(""); err == nil {
		t.Fatal("expected error for empty expression")
	}
}

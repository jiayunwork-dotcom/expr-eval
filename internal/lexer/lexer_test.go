package lexer

import "testing"

func TestTokenizeNumber(t *testing.T) {
	toks, err := Tokenize("3.14")
	if err != nil {
		t.Fatal(err)
	}
	if len(toks) != 1 || toks[0].Type != TokenNumber || toks[0].Val != "3.14" {
		t.Fatalf("got %+v", toks)
	}
}

func TestTokenizeOperators(t *testing.T) {
	toks, err := Tokenize("+-*/%^(),")
	if err != nil {
		t.Fatal(err)
	}
	want := []TokenType{
		TokenPlus, TokenMinus, TokenStar, TokenSlash, TokenPercent,
		TokenCaret, TokenLParen, TokenRParen, TokenComma,
	}
	if len(toks) != len(want) {
		t.Fatalf("len=%d want %d", len(toks), len(want))
	}
	for i := range want {
		if toks[i].Type != want[i] {
			t.Fatalf("tok %d = %v want %v", i, toks[i].Type, want[i])
		}
	}
}

func TestTokenizeBadChar(t *testing.T) {
	if _, err := Tokenize("1 @ 2"); err == nil {
		t.Fatal("expected error for bad character")
	}
}

func TestTokenizeEmpty(t *testing.T) {
	toks, err := Tokenize("")
	if err != nil {
		t.Fatal(err)
	}
	if len(toks) != 0 {
		t.Fatalf("expected empty slice, got %d tokens", len(toks))
	}
}

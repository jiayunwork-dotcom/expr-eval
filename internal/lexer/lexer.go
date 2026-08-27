package lexer

import "fmt"

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenNumber
	TokenIdent
	TokenPlus
	TokenMinus
	TokenStar
	TokenSlash
	TokenPercent
	TokenCaret
	TokenLParen
	TokenRParen
	TokenComma
	TokenIllegal
)

type Token struct {
	Type TokenType
	Val  string
	Pos  int
}

type Lexer struct {
	input []rune
	pos   int
}

func New(input string) *Lexer {
	return &Lexer{input: []rune(input)}
}

func isDigit(c rune) bool { return c >= '0' && c <= '9' }
func isLetter(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func (l *Lexer) skipSpace() {
	for l.pos < len(l.input) {
		c := l.input[l.pos]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			l.pos++
		} else {
			break
		}
	}
}

func (l *Lexer) NextToken() (Token, error) {
	l.skipSpace()
	if l.pos >= len(l.input) {
		return Token{Type: TokenEOF}, nil
	}
	c := l.input[l.pos]
	switch {
	case isDigit(c) || c == '.':
		return l.readNumber()
	case isLetter(c) || c == '_':
		return l.readIdent(), nil
	case c == '+':
		l.pos++
		return Token{Type: TokenPlus, Val: "+", Pos: l.pos - 1}, nil
	case c == '-':
		l.pos++
		return Token{Type: TokenMinus, Val: "-", Pos: l.pos - 1}, nil
	case c == '*':
		l.pos++
		return Token{Type: TokenStar, Val: "*", Pos: l.pos - 1}, nil
	case c == '/':
		l.pos++
		return Token{Type: TokenSlash, Val: "/", Pos: l.pos - 1}, nil
	case c == '%':
		l.pos++
		return Token{Type: TokenPercent, Val: "%", Pos: l.pos - 1}, nil
	case c == '^':
		l.pos++
		return Token{Type: TokenCaret, Val: "^", Pos: l.pos - 1}, nil
	case c == '(':
		l.pos++
		return Token{Type: TokenLParen, Val: "(", Pos: l.pos - 1}, nil
	case c == ')':
		l.pos++
		return Token{Type: TokenRParen, Val: ")", Pos: l.pos - 1}, nil
	case c == ',':
		l.pos++
		return Token{Type: TokenComma, Val: ",", Pos: l.pos - 1}, nil
	default:
		l.pos++
		return Token{}, fmt.Errorf("unexpected character %q at position %d", string(c), l.pos-1)
	}
}

func (l *Lexer) readNumber() (Token, error) {
	start := l.pos
	seenDot := false
	for l.pos < len(l.input) {
		c := l.input[l.pos]
		if isDigit(c) {
			l.pos++
		} else if c == '.' {
			if seenDot {
				return Token{}, fmt.Errorf("invalid number at position %d", start)
			}
			seenDot = true
			l.pos++
		} else {
			break
		}
	}
	return Token{Type: TokenNumber, Val: string(l.input[start:l.pos]), Pos: start}, nil
}

func (l *Lexer) readIdent() Token {
	start := l.pos
	for l.pos < len(l.input) {
		c := l.input[l.pos]
		if isLetter(c) || isDigit(c) || c == '_' {
			l.pos++
		} else {
			break
		}
	}
	return Token{Type: TokenIdent, Val: string(l.input[start:l.pos]), Pos: start}
}

func Tokenize(input string) ([]Token, error) {
	l := New(input)
	var toks []Token
	for {
		t, err := l.NextToken()
		if err != nil {
			return nil, err
		}
		if t.Type == TokenEOF {
			break
		}
		toks = append(toks, t)
	}
	return toks, nil
}

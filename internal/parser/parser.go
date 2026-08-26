package parser

import (
	"fmt"
	"strconv"

	"expr-eval/internal/lexer"
)

type Node interface {
	Pos() int
}

type NumberNode struct {
	Val float64
	pos int
}

type IdentNode struct {
	Name string
	pos  int
}

type UnaryNode struct {
	Op   rune
	Expr Node
	pos  int
}

type BinaryNode struct {
	Op    string
	Left  Node
	Right Node
	pos   int
}

type CallNode struct {
	Name string
	Args []Node
	pos  int
}

func (n *NumberNode) Pos() int { return n.pos }
func (n *IdentNode) Pos() int  { return n.pos }
func (n *UnaryNode) Pos() int  { return n.pos }
func (n *BinaryNode) Pos() int { return n.pos }
func (n *CallNode) Pos() int   { return n.pos }

type Parser struct {
	toks []lexer.Token
	pos  int
}

func Parse(input string) (Node, error) {
	toks, err := lexer.Tokenize(input)
	if err != nil {
		return nil, err
	}
	if len(toks) == 0 {
		return nil, fmt.Errorf("empty expression")
	}
	p := &Parser{toks: toks}
	node, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.pos < len(p.toks) {
		t := p.toks[p.pos]
		return nil, fmt.Errorf("unexpected token %q at position %d", t.Val, t.Pos)
	}
	return node, nil
}

func (p *Parser) peek() lexer.Token {
	if p.pos < len(p.toks) {
		return p.toks[p.pos]
	}
	return lexer.Token{Type: lexer.TokenEOF}
}

func (p *Parser) advance() lexer.Token {
	t := p.peek()
	p.pos++
	return t
}

func (p *Parser) parseExpr() (Node, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	for {
		t := p.peek()
		if t.Type == lexer.TokenPlus || t.Type == lexer.TokenMinus {
			p.advance()
			right, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			op := "+"
			if t.Type == lexer.TokenMinus {
				op = "-"
			}
			left = &BinaryNode{Op: op, Left: left, Right: right, pos: t.Pos}
		} else {
			break
		}
	}
	return left, nil
}

func (p *Parser) parseTerm() (Node, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for {
		t := p.peek()
		if t.Type == lexer.TokenStar || t.Type == lexer.TokenSlash || t.Type == lexer.TokenPercent {
			p.advance()
			right, err := p.parseUnary()
			if err != nil {
				return nil, err
			}
			op := ""
			switch t.Type {
			case lexer.TokenStar:
				op = "*"
			case lexer.TokenSlash:
				op = "/"
			case lexer.TokenPercent:
				op = "%"
			}
			left = &BinaryNode{Op: op, Left: left, Right: right, pos: t.Pos}
		} else {
			break
		}
	}
	return left, nil
}

func (p *Parser) parseUnary() (Node, error) {
	t := p.peek()
	if t.Type == lexer.TokenMinus || t.Type == lexer.TokenPlus {
		p.advance()
		expr, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		if t.Type == lexer.TokenPlus {
			return expr, nil
		}
		return &UnaryNode{Op: '-', Expr: expr, pos: t.Pos}, nil
	}
	return p.parsePower()
}

func (p *Parser) parsePower() (Node, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	t := p.peek()
	if t.Type == lexer.TokenCaret {
		p.advance()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = &BinaryNode{Op: "^", Left: left, Right: right, pos: t.Pos}
	}
	return left, nil
}

func (p *Parser) parsePrimary() (Node, error) {
	t := p.peek()
	switch t.Type {
	case lexer.TokenNumber:
		p.advance()
		v, err := strconv.ParseFloat(t.Val, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number %q at position %d", t.Val, t.Pos)
		}
		return &NumberNode{Val: v, pos: t.Pos}, nil
	case lexer.TokenIdent:
		p.advance()
		if p.peek().Type == lexer.TokenLParen {
			p.advance()
			var args []Node
			if p.peek().Type != lexer.TokenRParen {
				for {
					a, err := p.parseExpr()
					if err != nil {
						return nil, err
					}
					args = append(args, a)
					if p.peek().Type == lexer.TokenComma {
						p.advance()
						continue
					}
					break
				}
			}
			rp := p.peek()
			if rp.Type != lexer.TokenRParen {
				return nil, fmt.Errorf("expected ')' at position %d", rp.Pos)
			}
			p.advance()
			return &CallNode{Name: t.Val, Args: args, pos: t.Pos}, nil
		}
		return &IdentNode{Name: t.Val, pos: t.Pos}, nil
	case lexer.TokenLParen:
		p.advance()
		node, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		rp := p.peek()
		if rp.Type != lexer.TokenRParen {
			return nil, fmt.Errorf("expected ')' at position %d", rp.Pos)
		}
		p.advance()
		return node, nil
	default:
		return nil, fmt.Errorf("unexpected token %q at position %d", t.Val, t.Pos)
	}
}

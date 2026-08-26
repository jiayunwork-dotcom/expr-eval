package lexer

import "fmt"

func TokenName(t TokenType) string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenNumber:
		return "NUMBER"
	case TokenIdent:
		return "IDENT"
	case TokenPlus:
		return "PLUS"
	case TokenMinus:
		return "MINUS"
	case TokenStar:
		return "STAR"
	case TokenSlash:
		return "SLASH"
	case TokenPercent:
		return "PERCENT"
	case TokenCaret:
		return "CARET"
	case TokenLParen:
		return "LPAREN"
	case TokenRParen:
		return "RPAREN"
	case TokenComma:
		return "COMMA"
	case TokenIllegal:
		return "ILLEGAL"
	default:
		return fmt.Sprintf("TOKEN_%d", int(t))
	}
}

func IsOperator(t TokenType) bool {
	switch t {
	case TokenPlus, TokenMinus, TokenStar, TokenSlash, TokenPercent, TokenCaret:
		return true
	default:
		return false
	}
}

func Precedence(t TokenType) int {
	switch t {
	case TokenPlus, TokenMinus:
		return 1
	case TokenStar, TokenSlash, TokenPercent:
		return 2
	case TokenCaret:
		return 3
	default:
		return 0
	}
}

func DescribeTokens(tokens []Token) string {
	var s string
	for i, tok := range tokens {
		if i > 0 {
			s += " "
		}
		s += fmt.Sprintf("%s(%q)", TokenName(tok.Type), tok.Val)
	}
	return s
}

func CountTokens(tokens []Token, typ TokenType) int {
	count := 0
	for _, t := range tokens {
		if t.Type == typ {
			count++
		}
	}
	return count
}

func FilterTokens(tokens []Token, typ TokenType) []Token {
	var out []Token
	for _, t := range tokens {
		if t.Type == typ {
			out = append(out, t)
		}
	}
	return out
}

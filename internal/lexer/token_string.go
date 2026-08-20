package lexer

import "fmt"

// TokenName returns the human-readable name for a token type.
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

// IsOperator returns true if the token is an arithmetic operator.
func IsOperator(t TokenType) bool {
	switch t {
	case TokenPlus, TokenMinus, TokenStar, TokenSlash, TokenPercent, TokenCaret:
		return true
	default:
		return false
	}
}

// Precedence returns the operator precedence for the token (higher = binds tighter).
// Returns 0 for non-operators.
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

// DescribeTokens returns a formatted string showing all tokens in a slice.
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

// CountTokens returns the number of tokens of a specific type.
func CountTokens(tokens []Token, typ TokenType) int {
	count := 0
	for _, t := range tokens {
		if t.Type == typ {
			count++
		}
	}
	return count
}

// FilterTokens returns only tokens matching the given type.
func FilterTokens(tokens []Token, typ TokenType) []Token {
	var out []Token
	for _, t := range tokens {
		if t.Type == typ {
			out = append(out, t)
		}
	}
	return out
}

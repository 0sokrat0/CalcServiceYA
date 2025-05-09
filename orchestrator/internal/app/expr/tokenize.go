package expr

import (
	"unicode"
)

type TokenType int

const (
	TokenNumber TokenType = iota
	TokenPlus
	TokenMinus
	TokenMul
	TokenDiv
	TokenLParen
	TokenRParen
	TokenEOF
)

type Token struct {
	Type    TokenType
	Literal string
}

func tokenize(input string) []Token {
	var tokens []Token
	i := 0
	for i < len(input) {
		ch := input[i]
		switch {
		case unicode.IsSpace(rune(ch)):
			i++
		case unicode.IsDigit(rune(ch)) || ch == '.':
			start := i
			for i < len(input) && (unicode.IsDigit(rune(input[i])) || input[i] == '.') {
				i++
			}
			tokens = append(tokens, Token{Type: TokenNumber, Literal: input[start:i]})
		case ch == '+':
			tokens = append(tokens, Token{Type: TokenPlus, Literal: string(ch)})
			i++
		case ch == '-':
			tokens = append(tokens, Token{Type: TokenMinus, Literal: string(ch)})
			i++
		case ch == '*':
			tokens = append(tokens, Token{Type: TokenMul, Literal: string(ch)})
			i++
		case ch == '/':
			tokens = append(tokens, Token{Type: TokenDiv, Literal: string(ch)})
			i++
		case ch == '(':
			tokens = append(tokens, Token{Type: TokenLParen, Literal: string(ch)})
			i++
		case ch == ')':
			tokens = append(tokens, Token{Type: TokenRParen, Literal: string(ch)})
			i++
		default:
			i++
		}
	}
	tokens = append(tokens, Token{Type: TokenEOF})
	return tokens
}

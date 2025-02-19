package expr

import (
	"fmt"
	"strconv"
)

type Expr interface{}

type Number struct {
	Value float64
}

type BinaryExpr struct {
	Left     Expr
	Operator TokenType
	Right    Expr
}

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func (p *Parser) currentToken() Token {
	return p.tokens[p.pos]
}

func (p *Parser) nextToken() {
	if p.pos < len(p.tokens)-1 {
		p.pos++
	}
}

func precedence(t TokenType) int {
	switch t {
	case TokenPlus, TokenMinus:
		return 1
	case TokenMul, TokenDiv:
		return 2
	default:
		return 0
	}
}

func (p *Parser) parseExpression(rbp int) Expr {
	left := p.parsePrimary()

	for {
		curr := p.currentToken()
		if curr.Type != TokenPlus && curr.Type != TokenMinus &&
			curr.Type != TokenMul && curr.Type != TokenDiv {
			break
		}
		prec := precedence(curr.Type)
		if prec < rbp {
			break
		}

		op := curr.Type
		p.nextToken()
		right := p.parseExpression(prec + 1)
		left = BinaryExpr{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}
	return left
}

func (p *Parser) parsePrimary() Expr {
	tok := p.currentToken()
	switch tok.Type {
	case TokenNumber:
		p.nextToken()
		val, err := strconv.ParseFloat(tok.Literal, 64)
		if err != nil {
			panic(fmt.Sprintf("неверное число: %s", tok.Literal))
		}
		return Number{Value: val}
	case TokenLParen:
		p.nextToken()
		expr := p.parseExpression(0)
		if p.currentToken().Type != TokenRParen {
			panic("ожидалась закрывающая скобка")
		}
		p.nextToken()
		return expr
	default:
		panic(fmt.Sprintf("неожиданный токен: %v", tok))
	}
}

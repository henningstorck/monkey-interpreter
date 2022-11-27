package parser

import (
	"github.com/henningstorck/monkey-interpreter/ast"
	"github.com/henningstorck/monkey-interpreter/token"
)

func (par *Parser) parseExpression(precedence int) ast.Expression {
	prefixParseFn := par.prefixParseFns[par.curToken.Type]

	if prefixParseFn == nil {
		par.noPrefixParseFnError(par.curToken.Type)
		return nil
	}

	leftExp := prefixParseFn()

	for !par.peekTokenIs(token.Semicolon) && precedence < par.peekPrecedence() {
		infix := par.infixParseFns[par.peekToken.Type]

		if infix == nil {
			return leftExp
		}

		par.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (par *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    par.curToken,
		Operator: par.curToken.Literal,
	}

	par.nextToken()
	exp.Right = par.parseExpression(Prefix)
	return exp
}

func (par *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token:    par.curToken,
		Operator: par.curToken.Literal,
		Left:     left,
	}

	precedence := par.curPrecedence()
	par.nextToken()
	exp.Right = par.parseExpression(precedence)
	return exp
}

func (par *Parser) parseGroupedExpression() ast.Expression {
	par.nextToken()
	exp := par.parseExpression(Lowest)

	if !par.expectPeek(token.RParen) {
		return nil
	}

	return exp
}

func (par *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: par.curToken}

	if !par.expectPeek(token.LParen) {
		return nil
	}

	par.nextToken()
	exp.Condition = par.parseExpression(Lowest)

	if !par.expectPeek(token.RParen) {
		return nil
	}

	if !par.expectPeek(token.LBrace) {
		return nil
	}

	exp.Consequence = par.parseBlockStatement()

	if par.peekTokenIs(token.Else) {
		par.nextToken()

		if !par.expectPeek(token.LBrace) {
			return nil
		}

		exp.Alternative = par.parseBlockStatement()
	}

	return exp
}

func (par *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:    par.curToken,
		Function: fn,
	}

	exp.Arguments = par.parseExpressionList(token.RParen)
	return exp
}

func (par *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if par.peekTokenIs(end) {
		par.nextToken()
		return list
	}

	par.nextToken()
	list = append(list, par.parseExpression(Lowest))

	for par.peekTokenIs(token.Comma) {
		par.nextToken()
		par.nextToken()
		list = append(list, par.parseExpression(Lowest))
	}

	if !par.expectPeek(end) {
		return nil
	}

	return list
}

package parser

import (
	"fmt"
	"strconv"

	"github.com/henningstorck/monkey-interpreter/ast"
	"github.com/henningstorck/monkey-interpreter/token"
)

func (par *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: par.curToken, Value: par.curToken.Literal}
}

func (par *Parser) parseIntegerLiteral() ast.Expression {
	intLiteral := &ast.IntegerLiteral{Token: par.curToken}
	value, err := strconv.ParseInt(par.curToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", par.curToken.Literal)
		par.errors = append(par.errors, msg)
		return nil
	}

	intLiteral.Value = value
	return intLiteral
}

func (par *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{
		Token: par.curToken,
		Value: par.curTokenIs(token.True),
	}
}

func (par *Parser) parseFunctionLiteral() ast.Expression {
	fnLiteral := &ast.FunctionLiteral{Token: par.curToken}

	if !par.expectPeek(token.LParen) {
		return nil
	}

	fnLiteral.Parameters = par.parseFunctionParameters()

	if !par.expectPeek(token.LBrace) {
		return nil
	}

	fnLiteral.Body = par.parseBlockStatement()
	return fnLiteral
}

func (par *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if par.peekTokenIs(token.RParen) {
		par.nextToken()
		return identifiers
	}

	par.nextToken()

	identifier := &ast.Identifier{
		Token: par.curToken,
		Value: par.curToken.Literal,
	}

	identifiers = append(identifiers, identifier)

	for par.peekTokenIs(token.Comma) {
		par.nextToken()
		par.nextToken()

		identifier := &ast.Identifier{
			Token: par.curToken,
			Value: par.curToken.Literal,
		}

		identifiers = append(identifiers, identifier)
	}

	if !par.expectPeek(token.RParen) {
		return nil
	}

	return identifiers
}

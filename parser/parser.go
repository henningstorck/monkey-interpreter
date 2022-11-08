package parser

import (
	"fmt"

	"github.com/henningstorck/monkey-interpreter/ast"
	"github.com/henningstorck/monkey-interpreter/lexer"
	"github.com/henningstorck/monkey-interpreter/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	Lowest
	Equals
	Sum
	Product
	Prefix
	Call
)

type Parser struct {
	lex    *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func NewParser(lex *lexer.Lexer) *Parser {
	par := &Parser{
		lex:    lex,
		errors: []string{},
	}

	par.populateCurAndPeekToken()

	par.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	par.registerPrefix(token.Ident, par.parseIdentifier)

	return par
}

func (par *Parser) populateCurAndPeekToken() {
	par.nextToken()
	par.nextToken()
}

func (par *Parser) nextToken() {
	par.curToken = par.peekToken
	par.peekToken = par.lex.NextToken()
}

func (par *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for par.curToken.Type != token.EOF {
		stmt := par.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		par.nextToken()
	}

	return program
}

func (par *Parser) parseStatement() ast.Statement {
	switch par.curToken.Type {
	case token.Let:
		return par.parseLetStatement()
	case token.Return:
		return par.parseReturnStatement()
	default:
		return par.parseExpressionStatement()
	}
}

func (par *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: par.curToken}

	if !par.expectPeek(token.Ident) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: par.curToken, Value: par.curToken.Literal}

	if !par.expectPeek(token.Assign) {
		return nil
	}

	// Skipping expression for now

	for !par.curTokenIs(token.Semicolon) {
		par.nextToken()
	}

	return stmt
}

func (par *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: par.curToken}
	par.nextToken()

	// Skipping expression for now

	for !par.curTokenIs(token.Semicolon) {
		par.nextToken()
	}

	return stmt
}

func (par *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: par.curToken}
	stmt.Expression = par.parseExpression(Lowest)

	if par.peekTokenIs(token.Semicolon) {
		par.nextToken()
	}

	return stmt
}

func (par *Parser) parseExpression(precedence int) ast.Expression {
	prefixParseFn := par.prefixParseFns[par.curToken.Type]

	if prefixParseFn == nil {
		return nil
	}

	leftExp := prefixParseFn()
	return leftExp
}

func (par *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: par.curToken, Value: par.curToken.Literal}
}

func (par *Parser) curTokenIs(tokenType token.TokenType) bool {
	return par.curToken.Type == tokenType
}

func (par *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return par.peekToken.Type == tokenType
}

func (par *Parser) expectPeek(tokenType token.TokenType) bool {
	if par.peekTokenIs(tokenType) {
		par.nextToken()
		return true
	} else {
		par.peekError(tokenType)
		return false
	}
}

func (par *Parser) Errors() []string {
	return par.errors
}

func (par *Parser) peekError(tokenType token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", tokenType, par.peekToken.Type)
	par.errors = append(par.errors, msg)
}

func (par *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	par.prefixParseFns[tokenType] = fn
}

func (par *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	par.infixParseFns[tokenType] = fn
}

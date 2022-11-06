package parser

import (
	"fmt"

	"github.com/henningstorck/monkey-interpreter/ast"
	"github.com/henningstorck/monkey-interpreter/lexer"
	"github.com/henningstorck/monkey-interpreter/token"
)

type Parser struct {
	lex *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []string
}

func NewParser(lex *lexer.Lexer) *Parser {
	par := &Parser{
		lex:    lex,
		errors: []string{},
	}

	par.populateCurAndPeekToken()
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
	default:
		return nil
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

package parser

import (
	"fmt"
	"strconv"

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
	LessGreater
	Sum
	Product
	Prefix
	Call
)

var precedences = map[token.TokenType]int{
	token.Eq:       Equals,
	token.NotEq:    Equals,
	token.LT:       LessGreater,
	token.GT:       LessGreater,
	token.Plus:     Sum,
	token.Minus:    Sum,
	token.Slash:    Product,
	token.Asterisk: Product,
}

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
	par.registerPrefix(token.Int, par.parseIntegerLiteral)
	par.registerPrefix(token.True, par.parseBooleanLiteral)
	par.registerPrefix(token.False, par.parseBooleanLiteral)
	par.registerPrefix(token.Bang, par.parsePrefixExpression)
	par.registerPrefix(token.Minus, par.parsePrefixExpression)
	par.registerPrefix(token.LParen, par.parseGroupedExpression)
	par.registerPrefix(token.If, par.parseIfExpression)

	par.infixParseFns = make(map[token.TokenType]infixParseFn)
	par.registerInfix(token.Plus, par.parseInfixExpression)
	par.registerInfix(token.Minus, par.parseInfixExpression)
	par.registerInfix(token.Slash, par.parseInfixExpression)
	par.registerInfix(token.Asterisk, par.parseInfixExpression)
	par.registerInfix(token.Eq, par.parseInfixExpression)
	par.registerInfix(token.NotEq, par.parseInfixExpression)
	par.registerInfix(token.LT, par.parseInfixExpression)
	par.registerInfix(token.GT, par.parseInfixExpression)

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
		program.Statements = append(program.Statements, stmt)
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

	par.nextToken()
	stmt.Value = par.parseExpression(Lowest)

	for !par.curTokenIs(token.Semicolon) {
		par.nextToken()
	}

	return stmt
}

func (par *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: par.curToken}
	par.nextToken()
	stmt.ReturnValue = par.parseExpression(Lowest)

	for !par.curTokenIs(token.Semicolon) {
		par.nextToken()
	}

	return stmt
}

func (par *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	defer untrace(trace("parseExpressionStatement"))
	stmt := &ast.ExpressionStatement{Token: par.curToken}
	stmt.Expression = par.parseExpression(Lowest)

	if par.peekTokenIs(token.Semicolon) {
		par.nextToken()
	}

	return stmt
}

func (par *Parser) parseBlockStatement() *ast.BlockStatement {
	blockStmt := &ast.BlockStatement{Token: par.curToken}
	blockStmt.Statements = []ast.Statement{}
	par.nextToken()

	for !par.curTokenIs(token.RBrace) && !par.curTokenIs(token.EOF) {
		stmt := par.parseStatement()
		blockStmt.Statements = append(blockStmt.Statements, stmt)
		par.nextToken()
	}

	return blockStmt
}

func (par *Parser) parseExpression(precedence int) ast.Expression {
	defer untrace(trace("parseExpression"))
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

func (par *Parser) parseIdentifier() ast.Expression {
	defer untrace(trace("parseIdentifier"))
	return &ast.Identifier{Token: par.curToken, Value: par.curToken.Literal}
}

func (par *Parser) parseIntegerLiteral() ast.Expression {
	defer untrace(trace("parseIntegerLiteral"))
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

func (par *Parser) parsePrefixExpression() ast.Expression {
	defer untrace(trace("parsePrefixExpression"))

	exp := &ast.PrefixExpression{
		Token:    par.curToken,
		Operator: par.curToken.Literal,
	}

	par.nextToken()
	exp.Right = par.parseExpression(Prefix)
	return exp
}

func (par *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	defer untrace(trace("parseInfixExpression"))

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
	defer untrace(trace("parseGroupedExpression"))
	par.nextToken()
	exp := par.parseExpression(Lowest)

	if !par.expectPeek(token.RParen) {
		return nil
	}

	return exp
}

func (par *Parser) parseIfExpression() ast.Expression {
	defer untrace(trace("parseIfExpression"))
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

func (par *Parser) noPrefixParseFnError(tokenType token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", tokenType)
	par.errors = append(par.errors, msg)
}

func (par *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	par.prefixParseFns[tokenType] = fn
}

func (par *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	par.infixParseFns[tokenType] = fn
}

func (par *Parser) peekPrecedence() int {
	if precedences, ok := precedences[par.peekToken.Type]; ok {
		return precedences
	}

	return Lowest
}

func (par *Parser) curPrecedence() int {
	if precedences, ok := precedences[par.curToken.Type]; ok {
		return precedences
	}

	return Lowest
}

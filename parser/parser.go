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
	LessGreater
	Sum
	Product
	Prefix
	Call
)

var precedences = map[token.TokenType]int{
	token.Eq:          Equals,
	token.NotEq:       Equals,
	token.LessThan:    LessGreater,
	token.GreaterThan: LessGreater,
	token.Plus:        Sum,
	token.Minus:       Sum,
	token.Slash:       Product,
	token.Asterisk:    Product,
	token.LParen:      Call,
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
	par.registerPrefix(token.Function, par.parseFunctionLiteral)
	par.registerPrefix(token.String, par.parseStringLiteral)
	par.registerPrefix(token.LBracket, par.parseArrayLiteral)

	par.infixParseFns = make(map[token.TokenType]infixParseFn)
	par.registerInfix(token.Plus, par.parseInfixExpression)
	par.registerInfix(token.Minus, par.parseInfixExpression)
	par.registerInfix(token.Slash, par.parseInfixExpression)
	par.registerInfix(token.Asterisk, par.parseInfixExpression)
	par.registerInfix(token.Eq, par.parseInfixExpression)
	par.registerInfix(token.NotEq, par.parseInfixExpression)
	par.registerInfix(token.LessThan, par.parseInfixExpression)
	par.registerInfix(token.GreaterThan, par.parseInfixExpression)
	par.registerInfix(token.LParen, par.parseCallExpression)

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

	for !par.curTokenIs(token.EOF) {
		stmt := par.parseStatement()
		program.Statements = append(program.Statements, stmt)
		par.nextToken()
	}

	return program
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

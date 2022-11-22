package parser

import (
	"github.com/henningstorck/monkey-interpreter/ast"
	"github.com/henningstorck/monkey-interpreter/token"
)

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

	if par.peekTokenIs(token.Semicolon) {
		par.nextToken()
	}

	return stmt
}

func (par *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: par.curToken}
	par.nextToken()
	stmt.ReturnValue = par.parseExpression(Lowest)

	if par.peekTokenIs(token.Semicolon) {
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

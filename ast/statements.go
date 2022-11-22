package ast

import (
	"bytes"

	"github.com/henningstorck/monkey-interpreter/token"
)

type Statement interface {
	Node
	statementNode()
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (letStmt *LetStatement) statementNode()       {}
func (letStmt *LetStatement) TokenLiteral() string { return letStmt.Token.Literal }

func (letStmt *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(letStmt.TokenLiteral() + " ")
	out.WriteString(letStmt.Name.String())
	out.WriteString(" = ")

	if letStmt.Value != nil {
		out.WriteString(letStmt.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (returnStmt *ReturnStatement) statementNode()       {}
func (returnStmt *ReturnStatement) TokenLiteral() string { return returnStmt.Token.Literal }

func (returnStmt *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(returnStmt.TokenLiteral() + " ")

	if returnStmt.ReturnValue != nil {
		out.WriteString(returnStmt.ReturnValue.String())
	}

	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (expressionStmt *ExpressionStatement) statementNode()       {}
func (expressionStmt *ExpressionStatement) TokenLiteral() string { return expressionStmt.Token.Literal }

func (expressionStmt *ExpressionStatement) String() string {
	if expressionStmt.Expression != nil {
		return expressionStmt.Expression.String()
	}

	return ""
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (blockStmt BlockStatement) statementNode()       {}
func (blockStmt BlockStatement) TokenLiteral() string { return blockStmt.Token.Literal }

func (blockStmt BlockStatement) String() string {
	var out bytes.Buffer

	for _, stmt := range blockStmt.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

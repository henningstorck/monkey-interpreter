package ast

import (
	"bytes"

	"github.com/henningstorck/monkey-interpreter/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (prog *Program) TokenLiteral() string {
	if len(prog.Statements) > 0 {
		return prog.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
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

type Identifier struct {
	Token token.Token
	Value string
}

func (ident *Identifier) expressionNode()      {}
func (ident *Identifier) TokenLiteral() string { return ident.Token.Literal }
func (ident *Identifier) String() string       { return ident.Value }

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

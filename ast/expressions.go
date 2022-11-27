package ast

import (
	"bytes"
	"strings"

	"github.com/henningstorck/monkey-interpreter/token"
)

type Expression interface {
	Node
	expressionNode()
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (prefixExp *PrefixExpression) expressionNode()      {}
func (prefixExp *PrefixExpression) TokenLiteral() string { return prefixExp.Token.Literal }

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (infixExp *InfixExpression) expressionNode()      {}
func (infixExp *InfixExpression) TokenLiteral() string { return infixExp.Token.Literal }

func (infixExp *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(infixExp.Left.String())
	out.WriteString(" " + infixExp.Operator + " ")
	out.WriteString(infixExp.Right.String())
	out.WriteString(")")
	return out.String()
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ifExp IfExpression) expressionNode()      {}
func (ifExp IfExpression) TokenLiteral() string { return ifExp.Token.Literal }

func (ifExp IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ifExp.Condition.String())
	out.WriteString(" ")
	out.WriteString(ifExp.Consequence.String())

	if ifExp.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ifExp.Alternative.String())
	}

	return out.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression // identifier or function literal
	Arguments []Expression
}

func (callExp *CallExpression) expressionNode()      {}
func (callExp *CallExpression) TokenLiteral() string { return callExp.Token.Literal }

func (callExp *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}

	for _, arg := range callExp.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString(callExp.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (indexExp IndexExpression) expressionNode()      {}
func (indexExp IndexExpression) TokenLiteral() string { return indexExp.Token.Literal }

func (indexExp IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(indexExp.Left.String())
	out.WriteString("[")
	out.WriteString(indexExp.Index.String())
	out.WriteString("])")
	return out.String()
}

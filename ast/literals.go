package ast

import (
	"bytes"
	"strings"

	"github.com/henningstorck/monkey-interpreter/token"
)

type Identifier struct {
	Token token.Token
	Value string
}

func (ident *Identifier) expressionNode()      {}
func (ident *Identifier) TokenLiteral() string { return ident.Token.Literal }
func (ident *Identifier) String() string       { return ident.Value }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (intLiteral *IntegerLiteral) expressionNode()      {}
func (intLiteral *IntegerLiteral) TokenLiteral() string { return intLiteral.Token.Literal }
func (intLiteral *IntegerLiteral) String() string       { return intLiteral.Token.Literal }

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (boolLiteral *BooleanLiteral) expressionNode()      {}
func (boolLiteral *BooleanLiteral) TokenLiteral() string { return boolLiteral.Token.Literal }
func (boolLiteral *BooleanLiteral) String() string       { return boolLiteral.Token.Literal }

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fnLiteral FunctionLiteral) expressionNode()      {}
func (fnLiteral FunctionLiteral) TokenLiteral() string { return fnLiteral.Token.Literal }

func (fnLiteral FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}

	for _, param := range fnLiteral.Parameters {
		params = append(params, param.String())
	}

	out.WriteString(fnLiteral.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fnLiteral.Body.String())
	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (stringLiteral *StringLiteral) expressionNode()      {}
func (stringLiteral *StringLiteral) TokenLiteral() string { return stringLiteral.Token.Literal }
func (stringLiteral *StringLiteral) String() string       { return stringLiteral.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (arrLiteral ArrayLiteral) expressionNode()      {}
func (arrLiteral ArrayLiteral) TokenLiteral() string { return arrLiteral.Token.Literal }

func (arrLiteral ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}

	for _, element := range arrLiteral.Elements {
		elements = append(elements, element.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

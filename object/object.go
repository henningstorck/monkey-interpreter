package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/henningstorck/monkey-interpreter/ast"
)

const (
	IntegerObj     = "INTEGER"
	BooleanObj     = "BOOLEAN"
	NullObj        = "NULL"
	ReturnValueObj = "RETURN_VALUE"
	ErrorObj       = "ERROR"
	FunctionObj    = "FUNCTION"
	StringObj      = "STRING"
	BuiltinObj     = "BUILTIN"
	ArrayObj       = "ARRAY"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (integer *Integer) Inspect() string  { return fmt.Sprintf("%d", integer.Value) }
func (integer *Integer) Type() ObjectType { return IntegerObj }

type Boolean struct {
	Value bool
}

func (boolean *Boolean) Inspect() string  { return fmt.Sprintf("%t", boolean.Value) }
func (boolean *Boolean) Type() ObjectType { return BooleanObj }

type Null struct{}

func (null *Null) Inspect() string  { return "null" }
func (null *Null) Type() ObjectType { return NullObj }

type ReturnValue struct {
	Value Object
}

func (returnValue *ReturnValue) Type() ObjectType { return ReturnValueObj }
func (returnValue *ReturnValue) Inspect() string  { return returnValue.Value.Inspect() }

type Error struct {
	Message string
}

func (err *Error) Type() ObjectType { return ErrorObj }
func (err *Error) Inspect() string  { return "ERROR: " + err.Message }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (fn *Function) Type() ObjectType { return FunctionObj }

func (fn *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}

	for _, param := range fn.Parameters {
		params = append(params, param.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(fn.Body.String())
	out.WriteString("\n}")
	return out.String()
}

type String struct {
	Value string
}

func (str *String) Type() ObjectType { return StringObj }
func (str *String) Inspect() string  { return str.Value }

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Function BuiltinFunction
}

func (builtin *Builtin) Type() ObjectType { return BuiltinObj }
func (builtin *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (arr *Array) Type() ObjectType { return ArrayObj }

func (arr *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}

	for _, element := range arr.Elements {
		elements = append(elements, element.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

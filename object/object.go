package object

import "fmt"

const (
	IntegerObj     = "INTEGER"
	BooleanObj     = "BOOLEAN"
	NullObj        = "NULL"
	ReturnValueObj = "RETURN_VALUE"
	ErrorObj       = "ERROR"
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

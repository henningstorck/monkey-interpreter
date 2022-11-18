package evaluator_test

import (
	"testing"

	"github.com/henningstorck/monkey-interpreter/evaluator"
	"github.com/henningstorck/monkey-interpreter/lexer"
	"github.com/henningstorck/monkey-interpreter/object"
	"github.com/henningstorck/monkey-interpreter/parser"
	"github.com/stretchr/testify/assert"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func testEval(input string) object.Object {
	lex := lexer.NewLexer(input)
	par := parser.NewParser(lex)
	program := par.ParseProgram()
	return evaluator.Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) {
	result, ok := obj.(*object.Integer)
	assert.True(t, ok)
	assert.Equal(t, expected, result.Value)
}

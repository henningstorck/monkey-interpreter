package evaluator_test

import (
	"testing"

	"github.com/henningstorck/monkey-interpreter/object"
	"github.com/stretchr/testify/assert"
)

func TestEvalBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got 2, but expected 1"},

		{"len([])", 0},
		{"len([1, 2, 3])", 3},

		{"first([])", nil},
		{"first([1])", 1},
		{"first([1, 2, 3])", 1},
		{"first(1)", "argument to `first` must be ARRAY, got INTEGER"},
		{"first([], [])", "wrong number of arguments. got 2, but expected 1"},

		{"last([])", nil},
		{"last([1])", 1},
		{"last([1, 2, 3])", 3},
		{"last(1)", "argument to `first` must be ARRAY, got INTEGER"},
		{"last([], [])", "wrong number of arguments. got 2, but expected 1"},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)

		switch expected := test.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			assert.True(t, ok)
			assert.Equal(t, expected, errObj.Message)
		case nil:
			_, ok := evaluated.(*object.Null)
			assert.True(t, ok)
		}
	}
}

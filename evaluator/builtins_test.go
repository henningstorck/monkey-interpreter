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
		}
	}
}

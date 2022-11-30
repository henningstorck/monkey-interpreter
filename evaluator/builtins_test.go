package evaluator_test

import (
	"testing"

	"github.com/henningstorck/monkey-interpreter/object"
	"github.com/stretchr/testify/assert"
)

func TestEvalBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "invalid argument. got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got 2, but expected 1"},

		{"len([])", 0},
		{"len([1, 2, 3])", 3},

		{"first([])", nil},
		{"first([1])", 1},
		{"first([1, 2, 3])", 1},
		{"first(1)", "invalid argument. got INTEGER, but expected ARRAY"},
		{"first([], [])", "wrong number of arguments. got 2, but expected 1"},

		{"last([])", nil},
		{"last([1])", 1},
		{"last([1, 2, 3])", 3},
		{"last(1)", "invalid argument. got INTEGER, but expected ARRAY"},
		{"last([], [])", "wrong number of arguments. got 2, but expected 1"},

		{"rest([1]);", []int{}},
		{"rest([1, 2, 3, 4]);", []int{2, 3, 4}},
		{"rest(1)", "invalid argument. got INTEGER, but expected ARRAY"},
		{"rest([], [])", "wrong number of arguments. got 2, but expected 1"},

		{"push([], 1);", []int{1}},
		{"push([1, 2, 3], 4);", []int{1, 2, 3, 4}},
		{"push([1, 2, 3], 8 / 2);", []int{1, 2, 3, 4}},
		{"push(1, 2)", "invalid argument. got INTEGER, but expected ARRAY"},
		{"push([], 1, 2)", "wrong number of arguments. got 3, but expected 2"},
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
		case []int:
			arrObj, ok := evaluated.(*object.Array)
			assert.True(t, ok)
			assert.Equal(t, len(expected), len(arrObj.Elements))

			for i, expectedItem := range expected {
				testIntegerObject(t, arrObj.Elements[i], int64(expectedItem))
			}
		case nil:
			_, ok := evaluated.(*object.Null)
			assert.True(t, ok)
		}
	}
}

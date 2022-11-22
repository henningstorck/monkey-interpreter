package parser_test

import (
	"fmt"
	"testing"

	"github.com/henningstorck/monkey-interpreter/ast"
	"github.com/stretchr/testify/assert"
)

func TestParseIdentifier(t *testing.T) {
	input := "myVar;"
	program := testParse(t, input)
	assert.Equal(t, 1, len(program.Statements))
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	testLiteral(t, stmt.Expression, "myVar")
}

func TestParseIntegerLiteral(t *testing.T) {
	input := "5;"
	program := testParse(t, input)
	assert.Equal(t, 1, len(program.Statements))
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	testLiteral(t, stmt.Expression, 5)
}

func TestParseBooleanLiteral(t *testing.T) {
	input := "true;"
	program := testParse(t, input)
	assert.Equal(t, 1, len(program.Statements))
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	testLiteral(t, stmt.Expression, true)
}

func TestParseFunctionLiteral(t *testing.T) {
	input := "fn(x, y) { x + y; }"
	program := testParse(t, input)
	assert.Equal(t, 1, len(program.Statements))
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	fnLiteral, ok := stmt.Expression.(*ast.FunctionLiteral)
	assert.True(t, ok)
	assert.Equal(t, 2, len(fnLiteral.Parameters))
	testLiteral(t, fnLiteral.Parameters[0], "x")
	testLiteral(t, fnLiteral.Parameters[1], "y")
	assert.Equal(t, 1, len(fnLiteral.Body.Statements))
	bodyStmt, ok := fnLiteral.Body.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestParseFunctionParameters(t *testing.T) {
	tests := []struct {
		input  string
		params []string
	}{
		{input: "fn() {};", params: []string{}},
		{input: "fn(x) {};", params: []string{"x"}},
		{input: "fn(x, y, z) {};", params: []string{"x", "y", "z"}},
	}

	for _, test := range tests {
		program := testParse(t, test.input)
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		fnLiteral := stmt.Expression.(*ast.FunctionLiteral)
		assert.Equal(t, len(test.params), len(fnLiteral.Parameters))

		for i, ident := range test.params {
			testLiteral(t, fnLiteral.Parameters[i], ident)
		}
	}
}

func testLiteral(t *testing.T, exp ast.Expression, expected any) {
	switch value := expected.(type) {
	case int:
		testIntegerLiteral(t, exp, int64(value))
	case int64:
		testIntegerLiteral(t, exp, value)
	case bool:
		testBooleanLiteral(t, exp, value)
	case string:
		testIdentifier(t, exp, value)
	default:
		t.Fail()
	}
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) {
	intLiteral, ok := exp.(*ast.IntegerLiteral)
	assert.True(t, ok)
	assert.Equal(t, value, intLiteral.Value)
	assert.Equal(t, fmt.Sprintf("%d", value), intLiteral.TokenLiteral())
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) {
	boolLiteral, ok := exp.(*ast.BooleanLiteral)
	assert.True(t, ok)
	assert.Equal(t, value, boolLiteral.Value)
	assert.Equal(t, fmt.Sprintf("%t", value), boolLiteral.TokenLiteral())
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) {
	ident, ok := exp.(*ast.Identifier)
	assert.True(t, ok)
	assert.Equal(t, value, ident.Value)
	assert.Equal(t, value, ident.TokenLiteral())
}

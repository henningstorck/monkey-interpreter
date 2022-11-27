package parser_test

import (
	"testing"

	"github.com/henningstorck/monkey-interpreter/ast"
	"github.com/stretchr/testify/assert"
)

func TestParsePrefixExpressions(t *testing.T) {
	tests := []struct {
		input        string
		operator     string
		integerValue any
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"-alice", "-", "alice"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, test := range tests {
		program := testParse(t, test.input)
		assert.Equal(t, 1, len(program.Statements))
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		assert.True(t, ok)
		assert.Equal(t, test.operator, exp.Operator)
		testLiteral(t, exp.Right, test.integerValue)
	}
}

func TestParseInfixExpressions(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"alice * bob;", "alice", "*", "bob"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, test := range tests {
		program := testParse(t, test.input)
		assert.Equal(t, 1, len(program.Statements))
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)
		testInfixExpression(t, stmt.Expression, test.leftValue, test.operator, test.rightValue)
	}
}

func TestParseIfExpression(t *testing.T) {
	input := "if (x < y) { x }"
	program := testParse(t, input)

	assert.Equal(t, 1, len(program.Statements))
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	exp, ok := stmt.Expression.(*ast.IfExpression)
	assert.True(t, ok)
	testInfixExpression(t, exp.Condition, "x", "<", "y")

	assert.Equal(t, 1, len(exp.Consequence.Statements))
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	testLiteral(t, consequence.Expression, "x")

	assert.Nil(t, exp.Alternative)
}

func TestParseIfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else { y }"
	program := testParse(t, input)

	assert.Equal(t, 1, len(program.Statements))
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	exp, ok := stmt.Expression.(*ast.IfExpression)
	assert.True(t, ok)
	testInfixExpression(t, exp.Condition, "x", "<", "y")

	assert.Equal(t, 1, len(exp.Consequence.Statements))
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	testLiteral(t, consequence.Expression, "x")

	assert.Equal(t, 1, len(exp.Alternative.Statements))
	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	testLiteral(t, alternative.Expression, "y")
}

func TestParseCallExpression(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	program := testParse(t, input)
	assert.Equal(t, 1, len(program.Statements))
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	exp, ok := stmt.Expression.(*ast.CallExpression)
	assert.True(t, ok)
	testIdentifier(t, exp.Function, "add")
	assert.Equal(t, 3, len(exp.Arguments))
	testLiteral(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestParseCallArguments(t *testing.T) {
	tests := []struct {
		input string
		args  []any
	}{
		{input: "add();", args: []any{}},
		{input: "add(1);", args: []any{1}},
		{input: "add(1, 2, 3);", args: []any{1, 2, 3}},
	}

	for _, test := range tests {
		program := testParse(t, test.input)
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)
		exp, ok := stmt.Expression.(*ast.CallExpression)
		assert.True(t, ok)
		assert.Equal(t, len(test.args), len(exp.Arguments))

		for i, arg := range test.args {
			testLiteral(t, exp.Arguments[i], arg)
		}
	}
}

func TestParseIndexExpression(t *testing.T) {
	input := "myArray[1 + 1];"
	program := testParse(t, input)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	exp, ok := stmt.Expression.(*ast.IndexExpression)
	assert.True(t, ok)
	testIdentifier(t, exp.Left, "myArray")
	testInfixExpression(t, exp.Index, 1, "+", 1)
}

func testInfixExpression(t *testing.T, exp ast.Expression, leftValue any, operator string, rightValue any) {
	infixExp, ok := exp.(*ast.InfixExpression)
	assert.True(t, ok)
	testLiteral(t, infixExp.Left, leftValue)
	assert.Equal(t, operator, infixExp.Operator)
	testLiteral(t, infixExp.Right, rightValue)
}

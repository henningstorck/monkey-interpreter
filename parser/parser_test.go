package parser_test

import (
	"fmt"
	"testing"

	"github.com/henningstorck/monkey-interpreter/ast"
	"github.com/henningstorck/monkey-interpreter/lexer"
	"github.com/henningstorck/monkey-interpreter/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseLetStatements(t *testing.T) {
	tests := []struct {
		input string
		ident string
		value any
	}{
		{"let x=5;", "x", 5},
		{"let y = 10;", "y", 10},
		{"let z = 838383;", "z", 838383},
	}

	for _, test := range tests {
		program := testParse(t, test.input)
		assert.Equal(t, 1, len(program.Statements))
		stmt, ok := program.Statements[0].(*ast.LetStatement)
		assert.True(t, ok)
		testLetStatememt(t, stmt, test.ident)
		testLiteral(t, stmt.Value, test.value)
	}
}

func TestParseReturnStatements(t *testing.T) {
	tests := []struct {
		input string
		value any
	}{
		{"return 5;", 5},
		{"return 10;", 10},
		{"return 993322;", 993322},
	}

	for _, test := range tests {
		program := testParse(t, test.input)
		assert.Equal(t, 1, len(program.Statements))
		stmt, ok := program.Statements[0].(*ast.ReturnStatement)
		assert.True(t, ok)
		assert.Equal(t, "return", stmt.TokenLiteral())
		testLiteral(t, stmt.ReturnValue, test.value)
	}
}

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
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp := stmt.Expression.(*ast.CallExpression)
		assert.Equal(t, len(test.args), len(exp.Arguments))

		for i, arg := range test.args {
			testLiteral(t, exp.Arguments[i], arg)
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for _, test := range tests {
		program := testParse(t, test.input)
		assert.Equal(t, test.expected, program.String())
	}
}

func TestParseMissingSemicolon(t *testing.T) {
	input := "let x = 1 * 2 * 3 * 4 * 5"
	program := testParse(t, input)
	assert.Equal(t, "let x = ((((1 * 2) * 3) * 4) * 5);", program.String())
}

func testParse(t *testing.T, input string) *ast.Program {
	lex := lexer.NewLexer(input)
	par := parser.NewParser(lex)
	program := par.ParseProgram()
	checkParserErrors(t, par)
	return program
}

func testLetStatememt(t *testing.T, stmt ast.Statement, name string) {
	assert.Equal(t, "let", stmt.TokenLiteral())
	letStmt, ok := stmt.(*ast.LetStatement)
	assert.True(t, ok)
	assert.Equal(t, name, letStmt.Name.Value)
	assert.Equal(t, name, letStmt.Name.TokenLiteral())
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

func testInfixExpression(t *testing.T, exp ast.Expression, leftValue any, operator string, rightValue any) {
	infixExp, ok := exp.(*ast.InfixExpression)
	assert.True(t, ok)
	testLiteral(t, infixExp.Left, leftValue)
	assert.Equal(t, operator, infixExp.Operator)
	testLiteral(t, infixExp.Right, rightValue)
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

func checkParserErrors(t *testing.T, par *parser.Parser) {
	errors := par.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

package parser_test

import (
	"fmt"
	"testing"

	"github.com/henningstorck/monkey-interpreter/ast"
	"github.com/henningstorck/monkey-interpreter/lexer"
	"github.com/henningstorck/monkey-interpreter/parser"
	"github.com/stretchr/testify/assert"
)

func TestLetStatements(t *testing.T) {
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
		lex := lexer.NewLexer(test.input)
		par := parser.NewParser(lex)
		program := par.ParseProgram()
		checkParserErrors(t, par)

		assert.Equal(t, 1, len(program.Statements))
		stmt, ok := program.Statements[0].(*ast.LetStatement)
		assert.True(t, ok)
		testLetStatememt(t, stmt, test.ident)
		testLiteralExpression(t, stmt.Value, test.value)
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input string
		value any
	}{
		{"return 5;", 5},
		{"return 10;", 10},
		{"return 993322;", 993322},
	}

	for _, test := range tests {
		lex := lexer.NewLexer(test.input)
		par := parser.NewParser(lex)
		program := par.ParseProgram()
		checkParserErrors(t, par)

		assert.Equal(t, 1, len(program.Statements))
		stmt, ok := program.Statements[0].(*ast.ReturnStatement)
		assert.True(t, ok)
		assert.Equal(t, "return", stmt.TokenLiteral())
		testLiteralExpression(t, stmt.ReturnValue, test.value)
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "myVar;"

	lex := lexer.NewLexer(input)
	par := parser.NewParser(lex)
	program := par.ParseProgram()
	checkParserErrors(t, par)

	assert.Equal(t, 1, len(program.Statements))
	expressionStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	ident, ok := expressionStmt.Expression.(*ast.Identifier)
	assert.True(t, ok)
	assert.Equal(t, "myVar", ident.Value)
	assert.Equal(t, "myVar", ident.TokenLiteral())
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	lex := lexer.NewLexer(input)
	par := parser.NewParser(lex)
	program := par.ParseProgram()
	checkParserErrors(t, par)

	assert.Equal(t, 1, len(program.Statements))
	expressionStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	intLiteral, ok := expressionStmt.Expression.(*ast.IntegerLiteral)
	assert.True(t, ok)
	assert.Equal(t, int64(5), intLiteral.Value)
	assert.Equal(t, "5", intLiteral.TokenLiteral())
}

func TestBooleanLiteralExpression(t *testing.T) {
	input := "true;"

	lex := lexer.NewLexer(input)
	par := parser.NewParser(lex)
	program := par.ParseProgram()
	checkParserErrors(t, par)

	assert.Equal(t, 1, len(program.Statements))
	expressionStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	boolLiteral, ok := expressionStmt.Expression.(*ast.BooleanLiteral)
	assert.True(t, ok)
	assert.True(t, boolLiteral.Value)
	assert.Equal(t, "true", boolLiteral.TokenLiteral())
}

func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		input        string
		operator     string
		integerValue any
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"-alice", "-", "alice"},
	}

	for _, test := range tests {
		lex := lexer.NewLexer(test.input)
		par := parser.NewParser(lex)
		program := par.ParseProgram()
		checkParserErrors(t, par)

		assert.Equal(t, 1, len(program.Statements))
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		assert.True(t, ok)
		assert.Equal(t, test.operator, exp.Operator)
		testLiteralExpression(t, exp.Right, test.integerValue)
	}
}

func TestInfixExpressions(t *testing.T) {
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
	}

	for _, test := range tests {
		lex := lexer.NewLexer(test.input)
		par := parser.NewParser(lex)
		program := par.ParseProgram()
		checkParserErrors(t, par)

		assert.Equal(t, 1, len(program.Statements))
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		assert.True(t, ok)
		testLiteralExpression(t, exp.Left, test.leftValue)
		assert.Equal(t, test.operator, exp.Operator)
		testLiteralExpression(t, exp.Right, test.rightValue)
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
	}

	for _, test := range tests {
		lex := lexer.NewLexer(test.input)
		par := parser.NewParser(lex)
		program := par.ParseProgram()
		checkParserErrors(t, par)
		assert.Equal(t, test.expected, program.String())
	}
}

func testLetStatememt(t *testing.T, stmt ast.Statement, name string) {
	assert.Equal(t, "let", stmt.TokenLiteral())
	letStmt, ok := stmt.(*ast.LetStatement)
	assert.True(t, ok)
	assert.Equal(t, name, letStmt.Name.Value)
	assert.Equal(t, name, letStmt.Name.TokenLiteral())
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) {
	switch value := expected.(type) {
	case int:
		testIntegerLiteral(t, exp, int64(value))
	case int64:
		testIntegerLiteral(t, exp, value)
	case string:
		testIdentifier(t, exp, value)
	}
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) {
	intLiteral, ok := exp.(*ast.IntegerLiteral)
	assert.True(t, ok)
	assert.Equal(t, value, intLiteral.Value)
	assert.Equal(t, fmt.Sprintf("%d", value), intLiteral.TokenLiteral())
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

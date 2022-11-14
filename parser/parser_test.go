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
	input := `let x=5;
let y = 10;
let z = 838383;`

	lex := lexer.NewLexer(input)
	par := parser.NewParser(lex)
	program := par.ParseProgram()
	checkParserErrors(t, par)

	assert.NotNil(t, program)
	assert.Equal(t, 3, len(program.Statements))

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"z"},
	}

	for i, test := range tests {
		stmt := program.Statements[i]
		testLetStatememt(t, stmt, test.expectedIdentifier)
	}
}

func TestReturnStatements(t *testing.T) {
	input := `return 5;
return 10;
return 993322;`

	lex := lexer.NewLexer(input)
	par := parser.NewParser(lex)
	program := par.ParseProgram()
	checkParserErrors(t, par)

	assert.NotNil(t, program)
	assert.Equal(t, 3, len(program.Statements))

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		assert.True(t, ok)
		assert.Equal(t, "return", returnStmt.TokenLiteral())
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

func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
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
		testIntegerLiteral(t, exp.Right, test.integerValue)
	}
}

func testLetStatememt(t *testing.T, stmt ast.Statement, name string) {
	assert.Equal(t, "let", stmt.TokenLiteral())
	letStmt, ok := stmt.(*ast.LetStatement)
	assert.True(t, ok)
	assert.Equal(t, name, letStmt.Name.Value)
	assert.Equal(t, name, letStmt.Name.TokenLiteral())
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) {
	intLiteral, ok := exp.(*ast.IntegerLiteral)
	assert.True(t, ok)
	assert.Equal(t, value, intLiteral.Value)
	assert.Equal(t, fmt.Sprintf("%d", value), intLiteral.TokenLiteral())
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

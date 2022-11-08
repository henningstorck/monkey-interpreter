package parser_test

import (
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

func testLetStatememt(t *testing.T, stmt ast.Statement, name string) {
	assert.Equal(t, "let", stmt.TokenLiteral())
	letStmt, ok := stmt.(*ast.LetStatement)
	assert.True(t, ok)
	assert.Equal(t, name, letStmt.Name.Value)
	assert.Equal(t, name, letStmt.Name.TokenLiteral())
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

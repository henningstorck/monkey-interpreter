package parser_test

import (
	"testing"

	"github.com/henningstorck/monkey-interpreter/ast"
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

func testLetStatememt(t *testing.T, stmt ast.Statement, name string) {
	assert.Equal(t, "let", stmt.TokenLiteral())
	letStmt, ok := stmt.(*ast.LetStatement)
	assert.True(t, ok)
	assert.Equal(t, name, letStmt.Name.Value)
	assert.Equal(t, name, letStmt.Name.TokenLiteral())
}

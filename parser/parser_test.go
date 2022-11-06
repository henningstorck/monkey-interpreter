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

func testLetStatememt(t *testing.T, stmt ast.Statement, name string) {
	assert.Equal(t, "let", stmt.TokenLiteral())
	letStmt, ok := stmt.(*ast.LetStatement)
	assert.True(t, ok)
	assert.Equal(t, name, letStmt.Name.Value)
	assert.Equal(t, name, letStmt.Name.TokenLiteral())
}

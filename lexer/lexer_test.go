package lexer_test

import (
	"testing"

	"github.com/henningstorck/monkey-interpreter/lexer"
	"github.com/henningstorck/monkey-interpreter/token"
	"github.com/stretchr/testify/assert"
)

func TestNextTokenSimpleInput(t *testing.T) {
	input := "=+(){},;"

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.Assign, "="},
		{token.Plus, "+"},
		{token.LParen, "("},
		{token.RParen, ")"},
		{token.LBrace, "{"},
		{token.RBrace, "}"},
		{token.Comma, ","},
		{token.Semicolon, ";"},
		{token.EOF, ""},
	}

	lex := lexer.NewLexer(input)

	for _, test := range tests {
		tok := lex.NextToken()
		assert.Equal(t, test.expectedType, tok.Type)
		assert.Equal(t, test.expectedLiteral, tok.Literal)
	}
}

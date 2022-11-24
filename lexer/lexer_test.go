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

func TestNextTokenSimpleProgram(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
	x + y;
};

let result = add(five, ten);
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.Let, "let"},
		{token.Ident, "five"},
		{token.Assign, "="},
		{token.Int, "5"},
		{token.Semicolon, ";"},

		{token.Let, "let"},
		{token.Ident, "ten"},
		{token.Assign, "="},
		{token.Int, "10"},
		{token.Semicolon, ";"},

		{token.Let, "let"},
		{token.Ident, "add"},
		{token.Assign, "="},
		{token.Function, "fn"},
		{token.LParen, "("},
		{token.Ident, "x"},
		{token.Comma, ","},
		{token.Ident, "y"},
		{token.RParen, ")"},
		{token.LBrace, "{"},
		{token.Ident, "x"},
		{token.Plus, "+"},
		{token.Ident, "y"},
		{token.Semicolon, ";"},
		{token.RBrace, "}"},
		{token.Semicolon, ";"},

		{token.Let, "let"},
		{token.Ident, "result"},
		{token.Assign, "="},
		{token.Ident, "add"},
		{token.LParen, "("},
		{token.Ident, "five"},
		{token.Comma, ","},
		{token.Ident, "ten"},
		{token.RParen, ")"},
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

func TestNextTokenMoreOperators(t *testing.T) {
	input := `!-/*5;
5 < 10 > 5;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.Bang, "!"},
		{token.Minus, "-"},
		{token.Slash, "/"},
		{token.Asterisk, "*"},
		{token.Int, "5"},
		{token.Semicolon, ";"},
		{token.Int, "5"},
		{token.LessThan, "<"},
		{token.Int, "10"},
		{token.GreaterThan, ">"},
		{token.Int, "5"},
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

func TestNextTokenMoreKeywords(t *testing.T) {
	input := `if (5 < 10) {
	return true;
} else {
	return false;
}`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.If, "if"},
		{token.LParen, "("},
		{token.Int, "5"},
		{token.LessThan, "<"},
		{token.Int, "10"},
		{token.RParen, ")"},
		{token.LBrace, "{"},

		{token.Return, "return"},
		{token.True, "true"},
		{token.Semicolon, ";"},

		{token.RBrace, "}"},
		{token.Else, "else"},
		{token.LBrace, "{"},

		{token.Return, "return"},
		{token.False, "false"},
		{token.Semicolon, ";"},

		{token.RBrace, "}"},
		{token.EOF, ""},
	}

	lex := lexer.NewLexer(input)

	for _, test := range tests {
		tok := lex.NextToken()
		assert.Equal(t, test.expectedType, tok.Type)
		assert.Equal(t, test.expectedLiteral, tok.Literal)
	}
}

func TestNextTokenComposedOperators(t *testing.T) {
	input := `10 == 10;
10 != 9;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.Int, "10"},
		{token.Eq, "=="},
		{token.Int, "10"},
		{token.Semicolon, ";"},

		{token.Int, "10"},
		{token.NotEq, "!="},
		{token.Int, "9"},
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

func TestNextTokenStrings(t *testing.T) {
	input := `"meow";
"hello world";
"";`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.String, "meow"},
		{token.Semicolon, ";"},

		{token.String, "hello world"},
		{token.Semicolon, ";"},

		{token.String, ""},
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

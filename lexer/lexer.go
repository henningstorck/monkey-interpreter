package lexer

import "github.com/henningstorck/monkey-interpreter/token"

type Lexer struct {
	input        string
	position     int  // current char
	readPosition int  // current char + 1
	char         byte // current char under examination
}

func NewLexer(input string) *Lexer {
	lex := &Lexer{input: input}
	lex.readChar()
	return lex
}

func (lex *Lexer) readChar() {
	if lex.readPosition >= len(lex.input) {
		lex.char = 0
	} else {
		lex.char = lex.input[lex.readPosition]
	}

	lex.position = lex.readPosition
	lex.readPosition++
}

func (lex *Lexer) NextToken() token.Token {
	var tok token.Token

	switch lex.char {
	case '=':
		tok = token.NewToken(token.Assign, lex.char)
	case '+':
		tok = token.NewToken(token.Plus, lex.char)
	case ';':
		tok = token.NewToken(token.Semicolon, lex.char)
	case ',':
		tok = token.NewToken(token.Comma, lex.char)
	case '(':
		tok = token.NewToken(token.LParen, lex.char)
	case ')':
		tok = token.NewToken(token.RParen, lex.char)
	case '{':
		tok = token.NewToken(token.LBrace, lex.char)
	case '}':
		tok = token.NewToken(token.RBrace, lex.char)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}

	lex.readChar()
	return tok
}

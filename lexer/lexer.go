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

	lex.skipWhitespace()

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
	default:
		if isLetter(lex.char) {
			tok.Literal = lex.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(lex.char) {
			tok.Literal = lex.readNumber()
			tok.Type = token.Int
			return tok
		} else {
			tok = token.NewToken(token.Illegal, lex.char)
		}
	}

	lex.readChar()
	return tok
}

func (lex *Lexer) skipWhitespace() {
	for lex.char == ' ' || lex.char == '\t' || lex.char == '\n' || lex.char == '\r' {
		lex.readChar()
	}
}

func (lex *Lexer) readIdentifier() string {
	position := lex.position

	for isLetter(lex.char) {
		lex.readChar()
	}

	return lex.input[position:lex.position]
}

func (lex *Lexer) readNumber() string {
	position := lex.position

	for isDigit(lex.char) {
		lex.readChar()
	}

	return lex.input[position:lex.position]
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

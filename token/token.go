package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	Illegal = "ILLEGAL"
	EOF     = "EOF"

	Ident = "IDENT"
	Int   = "INT"

	Assign = "="
	Plus   = "+"

	Comma     = ","
	Semicolon = ";"

	LParen = "("
	RParen = ")"
	LBrace = "{"
	RBrace = "{"

	Function = "FUNCTION"
	Let      = "LET"
)

func NewToken(tokenType TokenType, char byte) Token {
	return Token{Type: tokenType, Literal: string(char)}
}
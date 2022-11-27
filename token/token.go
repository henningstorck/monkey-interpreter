package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	Illegal = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and literals
	Ident  = "IDENT"
	Int    = "INT"
	String = "STRING"

	// Operators
	Assign   = "="
	Bang     = "!"
	Plus     = "+"
	Minus    = "-"
	Asterisk = "*"
	Slash    = "/"

	LessThan    = "<"
	GreaterThan = ">"

	Eq    = "=="
	NotEq = "!="

	// Delimeters
	Comma     = ","
	Semicolon = ";"

	LParen   = "("
	RParen   = ")"
	LBrace   = "{"
	RBrace   = "}"
	LBracket = "["
	RBracket = "]"

	// Keywords
	Function = "FUNCTION"
	Let      = "LET"
	True     = "TRUE"
	False    = "FALSE"
	If       = "IF"
	Else     = "ELSE"
	Return   = "RETURN"
)

func NewToken(tokenType TokenType, char byte) Token {
	return Token{Type: tokenType, Literal: string(char)}
}

var keywords = map[string]TokenType{
	"fn":     Function,
	"let":    Let,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return Ident
}

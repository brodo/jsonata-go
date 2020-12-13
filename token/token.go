package token

type TokType string

type Token struct {
	Type     TokType
	Literal  string
	Position int
}

const (
	// Identifiers + literals
	IDENT  = "IDENT" // variable names, function names
	NUMBER = "NUMBER"
	STRING = "STRING"

	// Keywords
	TRUE  = "TRUE"
	FALSE = "FALSE"
	NULL  = "NULL"

	// Operators
	DOT         = "."
	LSBRACKET   = "["
	RSBRACKET   = "]"
	LBRACE      = "{"
	RBRACE      = "}"
	LPAREN      = "("
	RPAREN      = ")"
	COMMA       = ","
	AT          = "@"
	HASH        = "#"
	SEMICOLON   = ";"
	COLON       = ":"
	QUERY       = "?"
	PLUS        = "+"
	MINUS       = "-"
	ASTERISK    = "*"
	SLASH       = "/"
	PERCENT     = "%"
	PIPE        = "|"
	EQUALS      = "="
	LT          = "<"
	GT          = ">"
	CARET       = "^"
	DESCENDANTS = "**"
	RANGE       = ".."
	BIND        = ":="
	NQE         = "!="
	LTE         = "<="
	GTE         = ">="
	CHAIN       = "~>"
	AND         = "and"
	OR          = "OR"
	IN          = "IN"
	CONCAT      = "&"
	BANG        = "!"
	TILDE       = "~"

	// Other
	EOF     = "EOF"
	INVALID = "INVALID"
	REGEX   = "REGEX"
)

var keywords = map[string]TokType{
	"and":   AND,
	"or":    OR,
	"in":    IN,
	"true":  TRUE,
	"false": FALSE,
	"null":  NULL,
}

func LookupIdent(ident string) TokType {
	if tok, exists := keywords[ident]; exists {
		return tok
	}
	return IDENT
}

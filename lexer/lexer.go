package lexer

import (
	"github.com/brodo/jsonata-go/token"
	"unicode"
)

type Lexer struct {
	input        []rune
	position     int  // current position in input (points to current rune)
	readPosition int  // current reading position in input (after current rune)
	ch           rune // current rune under examination
}

func (l *Lexer) readRune() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekRune() rune {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func NewLexer(input string) *Lexer {
	runes := []rune(input)
	l := &Lexer{input: runes, position: 0, ch: runes[0]}
	l.readRune()
	return l

}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()
	switch l.ch {

	case '.':
		tok = l.makeTwoCharToken('.', token.DOT, token.RANGE)
	case '[':
		tok = l.newToken(token.LSBRACKET, l.ch)
	case ']':
		tok = l.newToken(token.RSBRACKET, l.ch)
	case '{':
		tok = l.newToken(token.LBRACE, l.ch)
	case '}':
		tok = l.newToken(token.RBRACE, l.ch)
	case '(':
		tok = l.newToken(token.LPAREN, l.ch)
	case ')':
		tok = l.newToken(token.RPAREN, l.ch)
	case ',':
		tok = l.newToken(token.COMMA, l.ch)
	case '@':
		tok = l.newToken(token.AT, l.ch)
	case '#':
		tok = l.newToken(token.HASH, l.ch)
	case ';':
		tok = l.newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = l.makeTwoCharToken('=', token.COLON, token.BIND)
	case '?':
		tok = l.newToken(token.QUERY, l.ch)
	case '+':
		tok = l.newToken(token.PLUS, l.ch)
	case '-':
		tok = l.newToken(token.MINUS, l.ch)
	case '*':
		tok = l.makeTwoCharToken('*', token.ASTERISK, token.DESCENDANTS)
	case '/':
		tok = l.newToken(token.SLASH, l.ch)
	case '%':
		tok = l.newToken(token.PERCENT, l.ch)
	case '|':
		tok = l.newToken(token.PIPE, l.ch)
	case '=':
		tok = l.newToken(token.EQUALS, l.ch)
	case '<':
		tok = l.makeTwoCharToken('=', token.LT, token.LTE)
	case '>':
		tok = l.makeTwoCharToken('=', token.GT, token.GTE)
	case '^':
		tok = l.newToken(token.CARET, l.ch)
	case '&':
		tok = l.newToken(token.CONCAT, l.ch)
	case '!':
		tok = l.makeTwoCharToken('=', token.BANG, token.NQE)
	case '~':
		tok = l.makeTwoCharToken('>', token.TILDE, token.CHAIN)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	default:
		if unicode.IsLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if unicode.IsDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = l.newToken(token.INVALID, l.ch)
		}
	}
	l.readRune()
	return tok
}

func (l *Lexer) newToken(tokenType token.TokType, literal rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(literal), Position: l.position}
}

func (l *Lexer) makeTwoCharToken(nextToken rune, oneCharType token.TokType, twoCharType token.TokType) token.Token {
	var tok token.Token
	if l.peekRune() == nextToken {
		ch := l.ch
		l.readRune()
		literal := string(ch) + string(l.ch)
		tok = token.Token{Type: twoCharType, Literal: literal}
	} else {
		tok = l.newToken(oneCharType, l.ch)
	}
	return tok
}

func (l *Lexer) skipWhitespace() {
	for unicode.Is(unicode.White_Space, l.ch) {
		l.readRune()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for unicode.IsLetter(l.ch) {
		l.readRune()
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) readNumber() string {
	position := l.position
	for unicode.IsDigit(l.ch) {
		l.readRune()
	}
	return string(l.input[position:l.position])
}

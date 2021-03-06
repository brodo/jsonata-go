package lexer

import (
	"github.com/brodo/jsonata-go/token"
	"regexp"
	"unicode"
)

var numberRegex = regexp.MustCompile("-?(0|([1-9][0-9]*))(\\.[0-9]+)?([Ee][-+]?[0-9]+)?")

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
		tok = l.makeTwoRuneToken('.', token.DOT, token.RANGE)
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
		tok = l.makeTwoRuneToken('=', token.COLON, token.BIND)
	case '?':
		tok = l.newToken(token.QUERY, l.ch)
	case '+':
		tok = l.newToken(token.PLUS, l.ch)
	case '-':
		if numberRegex.MatchString(string(l.input[l.position:])) {
			return l.readNumber()
		}
		tok = l.newToken(token.MINUS, l.ch)
	case '*':
		tok = l.makeTwoRuneToken('*', token.ASTERISK, token.DESCENDANTS)
	case '/':
		l.readRune()
		if l.ch == '*' {
			tok = l.readComment()
		} else {
			tok = l.newToken(token.SLASH, '/')
		}
	case '%':
		tok = l.newToken(token.PERCENT, l.ch)
	case '|':
		tok = l.newToken(token.PIPE, l.ch)
	case '=':
		tok = l.newToken(token.EQUALS, l.ch)
	case '<':
		tok = l.makeTwoRuneToken('=', token.LT, token.LTE)
	case '>':
		tok = l.makeTwoRuneToken('=', token.GT, token.GTE)
	case '^':
		tok = l.newToken(token.CARET, l.ch)
	case '&':
		tok = l.newToken(token.CONCAT, l.ch)
	case '!':
		tok = l.makeTwoRuneToken('=', token.BANG, token.NQE)
	case '~':
		tok = l.makeTwoRuneToken('>', token.TILDE, token.CHAIN)
	case '`':
		tok = l.readUntilRune('`', token.IDENT)
	case '\'':
		tok = l.readUntilRune('\'', token.STRING)
	case '"':
		tok = l.readUntilRune('"', token.STRING)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	default:
		if !unicode.Is(unicode.White_Space, l.ch) && !isReservedCharacter(l.peekRune()) && l.peekRune() != 0 && !unicode.IsDigit(l.ch) {
			tok.Start = l.position
			tok.Literal = l.readIdentifier()
			tok.End = l.position
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if numberRegex.MatchString(string(l.input[l.position:])) {
			return l.readNumber()
		} else {
			tok = l.newToken(token.INVALID, l.ch)
		}
	}
	l.readRune()
	return tok
}

func (l *Lexer) readUntilRune(enclosingRune rune, tokenType token.TokType) token.Token {
	var tok token.Token
	tok.Start = l.position
	l.readRune()
	tok.Type = tokenType
	position := l.position
	numBackslashes := 0
	isEscaped := false
	for l.ch != enclosingRune || isEscaped {
		if l.ch == '\\' {
			numBackslashes++
		} else {
			numBackslashes = max(numBackslashes, 0)
		}
		if numBackslashes%2 != 0 && l.peekRune() == enclosingRune {
			isEscaped = true
		} else {
			isEscaped = false
		}
		l.readRune()
		if l.ch == 0 {
			tok.Type = token.INVALID
			break
		}
	}
	tok.Literal = string(l.input[position:l.position])
	l.readRune()
	tok.End = l.position
	return tok
}

func (l *Lexer) readNumber() token.Token {
	var tok token.Token
	match := numberRegex.FindStringIndex(string(l.input[l.position:]))
	tok.Type = token.NUMBER
	tok.Literal = string(l.input[l.position : l.position+match[1]])
	tok.Start = l.position
	l.position += match[1]
	tok.End = l.position
	l.readPosition = l.position + 1
	if l.position < len(l.input) {
		l.ch = l.input[l.position]
	} else {
		l.ch = 0
	}
	return tok
}

//func (l *Lexer) readRegex() token.Token  {
//	var tok token.Token
//
//
//	tok.Type = token.REGEX
//
//}

func (l *Lexer) newToken(tokenType token.TokType, literal rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(literal), Start: l.position, End: l.position + 1}
}

func (l *Lexer) makeTwoRuneToken(nextToken rune, oneCharType token.TokType, twoCharType token.TokType) token.Token {
	var tok token.Token
	tok.Start = l.position
	if l.peekRune() == nextToken {
		ch := l.ch
		l.readRune()
		literal := string(ch) + string(l.ch)
		tok = token.Token{Type: twoCharType, Literal: literal}
	} else {
		tok = l.newToken(oneCharType, l.ch)
	}
	tok.End = tok.Start + 2
	return tok
}

func (l *Lexer) skipWhitespace() {
	for unicode.Is(unicode.White_Space, l.ch) {
		l.readRune()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for !unicode.Is(unicode.White_Space, l.ch) && !isReservedCharacter(l.ch) && l.ch != 0 {
		l.readRune()
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) readComment() token.Token {
	var tok token.Token
	tok.Start = l.position - 1 // the '/' was already read at this point
	for {
		if l.ch == '*' && l.peekRune() == '/' {
			tok.Type = token.COMMENT
			l.readRune() // closing '*'
			l.readRune() // closing '/'
			break
		}

		if l.ch == 0 {
			tok.Type = token.INVALID
			break
		}
		l.readRune()
	}
	tok.End = l.position
	tok.Literal = string(l.input[tok.Start:tok.End])
	return tok
}

func isReservedCharacter(r rune) bool {
	return r == '.' ||
		r == '[' ||
		r == ']' ||
		r == '{' ||
		r == '}' ||
		r == '(' ||
		r == ')' ||
		r == ',' ||
		r == '@' ||
		r == '#' ||
		r == ';' ||
		r == ':' ||
		r == '?' ||
		r == '+' ||
		r == '-' ||
		r == '*' ||
		r == '/' ||
		r == '%' ||
		r == '|' ||
		r == '=' ||
		r == '<' ||
		r == '>' ||
		r == '^' ||
		r == '&' ||
		r == '!' ||
		r == '~' ||
		r == '\'' ||
		r == '"'
}

func max(a, b int) int {
	if a < b {
		return a
	}
	return b
}

package lexer

import (
	"github.com/brodo/jsonata-go/token"
	"testing"
)

func TestNextTokenNonAscii(t *testing.T) {
	input := `  + - 😃 +`
	tests := []struct {
		expectedType    token.TokType
		expectedLiteral string
	}{
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.IDENT, "😃"},
		{token.PLUS, "+"},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] - tokentype wrong. expected=%q, got=%q.", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextTokenOperators(t *testing.T) {
	input := `+-/* := ~> ! !=`
	tests := []struct {
		expectedType    token.TokType
		expectedLiteral string
	}{
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.BIND, ":="},
		{token.CHAIN, "~>"},
		{token.BANG, "!"},
		{token.NQE, "!="},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] - tokentype wrong. expected=%q, got=%q.", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextTokenThreeLetterOperators(t *testing.T) {
	input := `and andy or in indian in true false null`
	tests := []struct {
		expectedType    token.TokType
		expectedLiteral string
	}{
		{token.AND, "and"},
		{token.IDENT, "andy"},
		{token.OR, "or"},
		{token.IN, "in"},
		{token.IDENT, "indian"},
		{token.IN, "in"},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.NULL, "null"},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] - tokentype wrong. expected=%q, got=%q.", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextTokenNumbers(t *testing.T) {
	input := `100 1.4 2455.243 -100 -1234.33 10e-12 0`
	tests := []struct {
		expectedType    token.TokType
		expectedLiteral string
	}{
		{token.NUMBER, "100"},
		{token.NUMBER, "1.4"},
		{token.NUMBER, "2455.243"},
		{token.NUMBER, "-100"},
		{token.NUMBER, "-1234.33"},
		{token.NUMBER, "10e-12"},
		{token.NUMBER, "0"},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] - tokentype wrong. expected=%q, got=%q.", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

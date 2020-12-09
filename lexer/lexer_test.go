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
		{token.INVALID, "😃"},
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
	input := `and andy or in indian in`
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

package lexer

import (
	"fmt"
	"github.com/brodo/jsonata-go/token"
	"testing"
)

type LexTest struct {
	expectedType    token.TokType
	expectedLiteral string
}

func runTests(tests []LexTest, input string, t *testing.T) {
	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()
		fmt.Printf("%+v\n", tok)
		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] - tokentype wrong. expected=%q, got=%q.", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
		if tok.Start == tok.End {
			t.Fatalf("test[%d] - start is eqal to end: %d, literal: %q ", i, tok.End, tok.Literal)
		}
	}
}

func TestNextTokenNonAscii(t *testing.T) {
	input := `  + - 😃 +`
	tests := []LexTest{
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.IDENT, "😃"},
		{token.PLUS, "+"},
	}

	runTests(tests, input, t)
}

func TestNextTokenOperators(t *testing.T) {
	input := `+-/ * := ~> ! !=`
	tests := []LexTest{
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.BIND, ":="},
		{token.CHAIN, "~>"},
		{token.BANG, "!"},
		{token.NQE, "!="},
	}
	runTests(tests, input, t)
}

func TestNextTokenThreeLetterOperators(t *testing.T) {
	input := `and andy or in indian in true false null`
	tests := []LexTest{
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
	runTests(tests, input, t)
}

func TestNextTokenNumbers(t *testing.T) {
	input := `100 1.4 2455.243 -100 -1234.33 10e-12 0`
	tests := []LexTest{
		{token.NUMBER, "100"},
		{token.NUMBER, "1.4"},
		{token.NUMBER, "2455.243"},
		{token.NUMBER, "-100"},
		{token.NUMBER, "-1234.33"},
		{token.NUMBER, "10e-12"},
		{token.NUMBER, "0"},
	}

	runTests(tests, input, t)
}

func TestNextTokenWhitespaceNames(t *testing.T) {
	input := "Other.`Over 18 ?` lala.`this \\` is a test` `another \\\\` lala"
	tests := []LexTest{
		{token.IDENT, "Other"},
		{token.DOT, "."},
		{token.IDENT, "Over 18 ?"},
		{token.IDENT, "lala"},
		{token.DOT, "."},
		{token.IDENT, "this \\` is a test"},
		{token.IDENT, "another \\\\"},
	}

	runTests(tests, input, t)
}

func TestNextTokenStrings(t *testing.T) {
	input := `"this is a string test" 'this is another string test' "this \" is an escape test" 
'this is "a string"' 'tripple \\\' escape!'`
	tests := []LexTest{
		{token.STRING, "this is a string test"},
		{token.STRING, "this is another string test"},
		{token.STRING, "this \\\" is an escape test"},
		{token.STRING, "this is \"a string\""},
		{token.STRING, "tripple \\\\\\' escape!"},
	}

	runTests(tests, input, t)
}

func TestNextTokenInvalid(t *testing.T) {
	input := `"this is a string test`
	tests := []LexTest{
		{token.INVALID, "this is a string test"},
	}
	runTests(tests, input, t)
}

func TestNextComment(t *testing.T) {
	input := `1 /* Comment! */ /* non-colosed comment!`
	tests := []LexTest{
		{token.NUMBER, "1"},
		{token.COMMENT, "/* Comment! */"},
		{token.INVALID, "/* non-colosed comment!"},
	}
	runTests(tests, input, t)
}

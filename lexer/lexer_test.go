package lexer

import (
	"simpsel/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `load $0 #10
aold
mul $0 $1 $2`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	} {
		{token.LOAD, "load"},
		{token.REGISTER, "0"},
		{token.INT, "10"},
		{token.ILLEGAL, "aold"},
		{token.MUL, "mul"},
		{token.REGISTER, "0"},
		{token.REGISTER, "1"},
		{token.REGISTER, "2"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

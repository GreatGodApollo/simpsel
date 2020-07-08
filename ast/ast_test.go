package ast

import (
	"simpsel/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Instructions: []Instruction{
			&AssemblerInstruction{
				token.Token{Type: token.LOAD, Literal: "load"},
				&RegisterLiteral{
					Token: token.Token{Type: token.REGISTER, Literal: "0"},
					Value: 0,
				},
				&IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "10"},
					Value: 10,
				},
				nil,
			},
		},
	}

	if program.String() != "load $0 #10;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

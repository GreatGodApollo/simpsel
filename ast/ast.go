package ast

import (
	"bytes"
	"simpsel/token"
)

// Base node
type Node interface {
	TokenLiteral() string
	String() string
}

// A singular instruction
type Instruction interface {
	Node
	instructionNode()
}

// An expression
type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Instructions []Instruction
}

func (p *Program) TokenLiteral() string {
	if len(p.Instructions) > 0 {
		return p.Instructions[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Instructions {
		out.WriteString(s.String())
	}

	return out.String()
}

type AssemblerInstruction struct {
	Opcode token.Token   // The token, ie `token.LOAD`
	Operand1 Expression
	Operand2 Expression
	Operand3 Expression
}

func (ai *AssemblerInstruction) instructionNode() {}
func (ai *AssemblerInstruction) TokenLiteral() string { return ai.Opcode.Literal }
func (ai *AssemblerInstruction) String() string {
	var out bytes.Buffer

	out.WriteString(ai.Opcode.Literal)

	if ai.Operand1 != nil {
		out.WriteString(" " + ai.Operand1.String())
	}
	if ai.Operand2 != nil {
		out.WriteString(" " + ai.Operand2.String())
	}
	if ai.Operand3 != nil {
		out.WriteString(" " + ai.Operand3.String())
	}

	out.WriteString(";")

	return out.String()
}

type RegisterLiteral struct {
	Token token.Token
	Value uint8
}

func (rl *RegisterLiteral) expressionNode()      {}
func (rl *RegisterLiteral) TokenLiteral() string { return rl.Token.Literal }
func (rl *RegisterLiteral) String() string       { return "$" + rl.Token.Literal }

type IntegerLiteral struct {
	Token token.Token
	Value uint16
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return "#" + il.Token.Literal }
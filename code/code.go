package code

import (
	"bytes"
	"encoding/hex"
	"simpsel/token"
)

type Instructions []byte

type Opcode byte

const (
	OpLoad Opcode = iota // 00
	OpAdd // 01
	OpSub // 02
	OpMul // 03
	OpDiv // 04
	OpHlt // 05
	OpIgl // 06
	OpJmp // 07
	OpJmpf // 08
	OpJmpb // 09
	OpEq // 0A
	OpNeq // 0B
	OpGt // 0C
	OpLt // 0D
	OpGte // 0E
	OpLte // 0F
	OpJmpe // 10
	OpNop // 11
)

func (ins Instructions) String() string {
	str := hex.EncodeToString(ins)
	return splitNth(str, 2)
}

func FromToken(tok token.Token) Opcode {
	switch tok.Type {
	case token.LOAD:
		return OpLoad
	case token.ADD:
		return OpAdd
	case token.SUB:
		return OpSub
	case token.MUL:
		return OpMul
	case token.DIV:
		return OpDiv
	case token.HLT:
		return OpHlt
	case token.ILLEGAL:
		return OpIgl
	case token.JMP:
		return OpJmp
	case token.JMPF:
		return OpJmpf
	case token.JMPB:
		return OpJmpb
	case token.EQ:
		return OpEq
	case token.NEQ:
		return OpNeq
	case token.GT:
		return OpGt
	case token.LT:
		return OpLt
	case token.GTE:
		return OpGte
	case token.LTE:
		return OpLte
	case token.JMPE:
		return OpJmpe
	case token.NOP:
		return OpNop
	default:
		return OpIgl
	}
}

func splitNth(s string, n int) string {
	var buffer bytes.Buffer
	var x = n - 1
	var l = len(s) - 1
	for i, char := range s {
		buffer.WriteRune(char)
		if i%n == x && i != l {
			buffer.WriteRune(' ')
		}
	}
	return buffer.String()
}
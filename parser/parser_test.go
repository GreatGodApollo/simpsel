package parser

import (
	"fmt"
	"simpsel/ast"
	"simpsel/lexer"
	"simpsel/token"
	"testing"
)

func TestLoadOpcode(t *testing.T) {
	input := "load $1 #10"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Instructions) != 1 {
		t.Fatalf("program.Instructions does not contain %d statements. got=%d\n",
			1, len(program.Instructions))
	}

	inst, ok := program.Instructions[0].(*ast.AssemblerInstruction)
	if !ok {
		t.Fatalf("inst is not ast.AssemblerInstruction. got=%T",
			program.Instructions[0])
	}

	if inst.Opcode.Type != token.LOAD {
		t.Fatalf("inst.Opcode.Type is not token.LOAD. got=%q",
			inst.Opcode.Type)
	}

	if !testRegister(t, inst.Operand1, 1) { return }

	if !testInteger(t, inst.Operand2, 10) { return }

	if inst.Operand3 != nil {
		t.Fatalf("inst.Operand is not nil. got=%q",
			inst.Operand3)
	}

}

func TestAddOpcode(t *testing.T) {
	input := "add $0 $1 $2"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Instructions) != 1 {
		t.Fatalf("program.Instructions does not contain %d statements. got=%d\n",
			1, len(program.Instructions))
	}

	inst, ok := program.Instructions[0].(*ast.AssemblerInstruction)
	if !ok {
		t.Fatalf("inst is not ast.AssemblerInstruction. got=%T",
			program.Instructions[0])
	}

	if inst.Opcode.Type != token.ADD {
		t.Fatalf("inst.Opcode.Type is not token.ADD. got=%q",
			inst.Opcode.Type)
	}

	if !testRegister(t, inst.Operand1, 0) { return }

	if !testRegister(t, inst.Operand2, 1) { return }

	if !testRegister(t, inst.Operand3, 2) { return }
}

func TestHltOpcode(t *testing.T) {
	input := "hlt"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Instructions) != 1 {
		t.Fatalf("program.Instructions does not contain %d statements. got=%d\n",
			1, len(program.Instructions))
	}

	inst, ok := program.Instructions[0].(*ast.AssemblerInstruction)
	if !ok {
		t.Fatalf("inst is not ast.AssemblerInstruction. got=%T",
			program.Instructions[0])
	}

	if inst.Opcode.Type != token.HLT {
		t.Fatalf("inst.Opcode.Type is not token.HLT. got=%q",
			inst.Opcode.Type)
	}

	if !testNil(t, inst.Operand1, true) { return }

	if !testNil(t, inst.Operand2, true) { return }

	if !testNil(t, inst.Operand3, true) { return }
}

func TestJmpOpcode(t *testing.T) {
	input := "jmp $1"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Instructions) != 1 {
		t.Fatalf("program.Instructions does not contain %d statements. got=%d\n",
			1, len(program.Instructions))
	}

	inst, ok := program.Instructions[0].(*ast.AssemblerInstruction)
	if !ok {
		t.Fatalf("inst is not ast.AssemblerInstruction. got=%T",
			program.Instructions[0])
	}

	if inst.Opcode.Type != token.JMP {
		t.Fatalf("inst.Opcode.Type is not token.JMP. got=%q",
			inst.Opcode.Type)
	}

	if !testRegister(t, inst.Operand1, 1) { return }

	if !testNil(t, inst.Operand2, true) { return }

	if !testNil(t, inst.Operand3, true) { return }
}

func testRegister(t *testing.T, exp ast.Expression, value uint8) bool {
	reg, ok := exp.(*ast.RegisterLiteral)
	if !ok {
		t.Errorf("exp not *ast.RegisterLiteral. got=%T", exp)
		return false
	}

	if reg.Value != value {
		t.Errorf("reg.Value not %d. got=%d", value, reg.Value)
		return false
	}

	if reg.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("reg.TokenLiteral not %d. got=%s", value,
			reg.TokenLiteral())
		return false
	}

	return true
}

func testInteger(t *testing.T, exp ast.Expression, value uint16) bool {
	int, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("exp not *ast.IntegerLiteral. got=%T", exp)
		return false
	}

	if int.Value != value {
		t.Errorf("int.Value not %d. got=%d", value, int.Value)
		return false
	}

	if int.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("int.TokenLiteral not %d. got=%s", value,
			int.TokenLiteral())
		return false
	}

	return true
}

func testNil(t *testing.T, exp ast.Expression, expect bool) bool {
	if (exp == nil) && !expect {
		t.Errorf("exp not nil. got=%q", exp)
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

package parser

import (
	"fmt"
	"simpsel/ast"
	"simpsel/lexer"
	"simpsel/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	OPCODE     // load, add, etc.
	DIRECTIVES // .code / .data
)

var precedences = map[token.TokenType]int {
	token.LOAD: OPCODE,
	token.ADD: OPCODE,
	token.SUB: OPCODE,
	token.MUL: OPCODE,
	token.DIV: OPCODE,
	token.HLT: OPCODE,
	token.JMP: OPCODE,
	token.JMPF: OPCODE,
	token.JMPB: OPCODE,
	token.EQ: OPCODE,
	token.NEQ: OPCODE,
	token.GT: OPCODE,
	token.LT: OPCODE,
	token.GTE: OPCODE,
	token.LTE: OPCODE,
	token.JMPE: OPCODE,
	token.NOP: OPCODE,
}

type (
	opCodeParseFn func() ast.Instruction
)

type Parser struct {
	l *lexer.Lexer
	errors []string

	curToken token.Token
	peekToken token.Token

	opCodeParseFns map[token.TokenType]opCodeParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
		errors: []string{},
		opCodeParseFns: make(map[token.TokenType]opCodeParseFn),
	}

	// op
	p.registerOpcode(token.HLT, p.parseBlank)
	p.registerOpcode(token.ILLEGAL, p.parseBlank)
	p.registerOpcode(token.NOP, p.parseBlank)

	// op $Reg
	p.registerOpcode(token.JMP, p.parseRegister)
	p.registerOpcode(token.JMPF, p.parseRegister)
	p.registerOpcode(token.JMPB, p.parseRegister)
	p.registerOpcode(token.JMPE, p.parseRegister)

	// op $Reg #Int
	p.registerOpcode(token.LOAD, p.parseRegisterInt)

	// op $Reg $Reg
	p.registerOpcode(token.EQ, p.parseRegisterRegister)
	p.registerOpcode(token.NEQ, p.parseRegisterRegister)
	p.registerOpcode(token.GT, p.parseRegisterRegister)
	p.registerOpcode(token.LT, p.parseRegisterRegister)
	p.registerOpcode(token.GTE, p.parseRegisterRegister)
	p.registerOpcode(token.LTE, p.parseRegisterRegister)

	// op $Reg $Reg $Reg
	p.registerOpcode(token.ADD, p.parseRegisterRegisterRegister)
	p.registerOpcode(token.SUB, p.parseRegisterRegisterRegister)
	p.registerOpcode(token.MUL, p.parseRegisterRegisterRegister)
	p.registerOpcode(token.DIV, p.parseRegisterRegisterRegister)

	// Read two tokens, so both curToken and peekToken are set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noOpcodeParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no opcode parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerOpcode(tokenType token.TokenType, fn opCodeParseFn) {
	p.opCodeParseFns[tokenType] = fn
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Instructions = []ast.Instruction{}

	for !p.curTokenIs(token.EOF) {
		inst := p.parseInstruction()
		if inst != nil {
			program.Instructions = append(program.Instructions, inst)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseInstruction() ast.Instruction {
	fn := p.opCodeParseFns[p.curToken.Type]
	if fn == nil {
		p.noOpcodeParseFnError(p.curToken.Type)
		return nil
	}
	inst := fn()

	return inst
}

func (p *Parser) parseRegisterInt() ast.Instruction {
	inst := &ast.AssemblerInstruction{Opcode: p.curToken}

	if !p.expectPeek(token.REGISTER) {
		return nil
	}

	reg, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		return nil
	}
	inst.Operand1 = &ast.RegisterLiteral{
		Token: p.curToken,
		Value: uint8(reg),
	}

	if !p.expectPeek(token.INT) {
		return nil
	}

	val, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		return nil
	}
	inst.Operand2 = &ast.IntegerLiteral{
		Token: p.curToken,
		Value: uint16(val),
	}

	inst.Operand3 = nil
	return inst
}

func (p *Parser) parseRegisterRegisterRegister() ast.Instruction {
	inst := &ast.AssemblerInstruction{Opcode: p.curToken}

	if !p.expectPeek(token.REGISTER) {
		return nil
	}

	reg1, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		return nil
	}
	inst.Operand1 = &ast.RegisterLiteral{
		Token: p.curToken,
		Value: byte(reg1),
	}

	if !p.expectPeek(token.REGISTER) {
		return nil
	}

	reg2, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		return nil
	}
	inst.Operand2 = &ast.RegisterLiteral{
		Token: p.curToken,
		Value: byte(reg2),
	}

	if !p.expectPeek(token.REGISTER) {
		return nil
	}

	reg3, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		return nil
	}
	inst.Operand3 = &ast.RegisterLiteral{
		Token: p.curToken,
		Value: byte(reg3),
	}

	return inst
}

func (p *Parser) parseBlank() ast.Instruction {
	inst := &ast.AssemblerInstruction{Opcode: p.curToken}

	inst.Operand1 = nil
	inst.Operand2 = nil
	inst.Operand3 = nil

	return inst
}

func (p *Parser) parseRegister() ast.Instruction {
	inst := &ast.AssemblerInstruction{Opcode: p.curToken}

	if !p.expectPeek(token.REGISTER) {
		return nil
	}

	reg, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		return nil
	}
	inst.Operand1 = &ast.RegisterLiteral{
		Token: p.curToken,
		Value: uint8(reg),
	}

	inst.Operand2 = nil
	inst.Operand3 = nil

	return inst
}

func (p *Parser) parseRegisterRegister() ast.Instruction {
	inst := &ast.AssemblerInstruction{Opcode: p.curToken}

	if !p.expectPeek(token.REGISTER) {
		return nil
	}

	reg1, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		return nil
	}
	inst.Operand1 = &ast.RegisterLiteral{
		Token: p.curToken,
		Value: uint8(reg1),
	}

	if !p.expectPeek(token.REGISTER) {
		return nil
	}

	reg2, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		return nil
	}
	inst.Operand2 = &ast.RegisterLiteral{
		Token: p.curToken,
		Value: uint8(reg2),
	}

	inst.Operand3 = nil

	return inst
}


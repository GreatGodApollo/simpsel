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

	// Ignore comments bb
	p.registerParseFn(token.COMMENT, p.parseIgnore)

	// directive

	// op
	p.registerParseFn(token.HLT, p.parseBlank)
	p.registerParseFn(token.ILLEGAL, p.parseBlank)
	p.registerParseFn(token.NOP, p.parseBlank)

	// op $Reg
	p.registerParseFn(token.JMP, p.parseRegister)
	p.registerParseFn(token.JMPF, p.parseRegister)
	p.registerParseFn(token.JMPB, p.parseRegister)
	p.registerParseFn(token.JMPE, p.parseRegister)

	// op $Reg #Int
	p.registerParseFn(token.LOAD, p.parseRegisterInt)

	// op $Reg $Reg
	p.registerParseFn(token.EQ, p.parseRegisterRegister)
	p.registerParseFn(token.NEQ, p.parseRegisterRegister)
	p.registerParseFn(token.GT, p.parseRegisterRegister)
	p.registerParseFn(token.LT, p.parseRegisterRegister)
	p.registerParseFn(token.GTE, p.parseRegisterRegister)
	p.registerParseFn(token.LTE, p.parseRegisterRegister)

	// op $Reg $Reg $Reg
	p.registerParseFn(token.ADD, p.parseRegisterRegisterRegister)
	p.registerParseFn(token.SUB, p.parseRegisterRegisterRegister)
	p.registerParseFn(token.MUL, p.parseRegisterRegisterRegister)
	p.registerParseFn(token.DIV, p.parseRegisterRegisterRegister)

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

func (p *Parser) noParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerTooBigError(regNum uint8) bool {
	if regNum > uint8(31) {
		msg := fmt.Sprintf("register number too big, must be less than 32. got=%d", regNum)
		p.errors = append(p.errors, msg)
		return true
	}
	return false
}

func (p *Parser) registerParseFn(tokenType token.TokenType, fn opCodeParseFn) {
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
		p.noParseFnError(p.curToken.Type)
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

	if p.registerTooBigError(byte(reg)) {
		return nil
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

	if p.registerTooBigError(byte(reg1)) {
		return nil
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

	if p.registerTooBigError(byte(reg2)) {
		return nil
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

	if p.registerTooBigError(byte(reg3)) {
		return nil
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

	if p.registerTooBigError(byte(reg)) {
		return nil
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

	if p.registerTooBigError(byte(reg1)) {
		return nil
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

	if p.registerTooBigError(byte(reg2)) {
		return nil
	}

	inst.Operand3 = nil

	return inst
}

func (p *Parser) parseIgnore() ast.Instruction {
	return nil
}
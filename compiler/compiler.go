package compiler

import (
	"encoding/binary"
	"simpsel/ast"
	"simpsel/code"
)

type Compiler struct {
	instructions code.Instructions
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, i := range node.Instructions {
			err := c.Compile(i)
			if err != nil {
				return err
			}
		}

	case *ast.AssemblerInstruction:
		c.emit(code.FromToken(node.Opcode), node.Operand1, node.Operand2, node.Operand3)
	}

	return nil
}

func (c *Compiler) emit(op code.Opcode, operands ...ast.Expression) int {
	ins := make([]byte, 4)

	i := 0
	for len(ins) < 4 {
		ins[i] = 0
		i++
	}

	ins[0] = byte(op)
	p := 1
	for _, operand := range operands {
		switch operand := operand.(type) {
		case *ast.IntegerLiteral:
			if len(ins)-p > 1 {
				binary.LittleEndian.PutUint16(ins[p:], operand.Value)
				p += 2
			}
		case *ast.RegisterLiteral:
			if  len(ins) - p > 0 {
				ins[p] = operand.Value
				p += 1
			}
		}
	}

	pos := c.addInstruction(ins)
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
	}
}

type Bytecode struct {
	Instructions code.Instructions
}

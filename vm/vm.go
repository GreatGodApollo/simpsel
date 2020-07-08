package vm

import (
	"encoding/binary"
	"simpsel/code"
	"simpsel/compiler"
)

type VM struct {
	Registers []int32
	Program   code.Instructions
	Counter   int
	Remainder int32
	EqualFlag bool
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		Registers: make([]int32, 32),
		Program:   bytecode.Instructions,
		Counter:   0,
		Remainder: 0,
		EqualFlag: false,
	}
}

func (vm *VM) Run() {
	isDone := false
	for !isDone {
		isDone = vm.executeInstruction()
	}
}

func (vm *VM) RunOnce() {
	vm.executeInstruction()
}

func (vm *VM) executeInstruction() bool {
	if vm.Counter >= len(vm.Program) {
		return true
	}
	switch vm.decodeOpcode() {
	case code.OpLoad:
		register := vm.nextByte()
		num := int32(vm.next2Bytes())
		vm.Registers[register] = num
	case code.OpAdd:
		register1 := vm.Registers[vm.nextByte()]
		register2 := vm.Registers[vm.nextByte()]
		vm.Registers[vm.nextByte()] = register1 + register2
	case code.OpSub:
		register1 := vm.Registers[vm.nextByte()]
		register2 := vm.Registers[vm.nextByte()]
		vm.Registers[vm.nextByte()] = register1 - register2
	case code.OpMul:
		register1 := vm.Registers[vm.nextByte()]
		register2 := vm.Registers[vm.nextByte()]
		vm.Registers[vm.nextByte()] = register1 * register2
	case code.OpDiv:
		register1 := vm.Registers[vm.nextByte()]
		register2 := vm.Registers[vm.nextByte()]
		vm.Registers[vm.nextByte()] = register1 / register2
		vm.Remainder = register1 % register2
	case code.OpHlt:
		println("HLT Encountered")
		return true
	case code.OpIgl:
		println("Illegal opcode @ " + string(vm.Counter))
		return true
	case code.OpJmp:
		target := vm.Registers[vm.nextByte()]
		vm.Counter = int(target)
	case code.OpJmpf:
		value := vm.Registers[vm.nextByte()]
		vm.Counter += int(value)
	case code.OpJmpb:
		value := vm.Registers[vm.nextByte()]
		vm.Counter -= int(value)
	case code.OpEq:
		register1 := vm.Registers[vm.nextByte()]
		register2 := vm.Registers[vm.nextByte()]
		vm.EqualFlag = register1 == register2
		vm.nextByte()
	case code.OpNeq:
		register1 := vm.Registers[vm.nextByte()]
		register2 := vm.Registers[vm.nextByte()]
		vm.EqualFlag = register1 != register2
		vm.nextByte()
	case code.OpGt:
		register1 := vm.Registers[vm.nextByte()]
		register2 := vm.Registers[vm.nextByte()]
		vm.EqualFlag = register1 > register2
		vm.nextByte()
	case code.OpLt:
		register1 := vm.Registers[vm.nextByte()]
		register2 := vm.Registers[vm.nextByte()]
		vm.EqualFlag = register1 < register2
		vm.nextByte()
	case code.OpGte:
		register1 := vm.Registers[vm.nextByte()]
		register2 := vm.Registers[vm.nextByte()]
		vm.EqualFlag = register1 >= register2
		vm.nextByte()
	case code.OpLte:
		register1 := vm.Registers[vm.nextByte()]
		register2 := vm.Registers[vm.nextByte()]
		vm.EqualFlag = register1 <= register2
		vm.nextByte()
	case code.OpJmpe:
		if vm.EqualFlag {
			target := vm.Registers[vm.nextByte()]
			vm.Counter = int(target)
		} else {
			vm.nextByte()
			vm.nextByte()
			vm.nextByte()
		}
	case code.OpNop:
		vm.nextByte()
		vm.nextByte()
		vm.nextByte()
	}

	return false
}

func (vm *VM) decodeOpcode() code.Opcode {
	opcode := code.Opcode(vm.Program[vm.Counter])
	vm.Counter++
	return opcode
}

func (vm *VM) nextByte() byte {
	result := vm.Program[vm.Counter]
	vm.Counter++
	return result
}

func (vm *VM) next2Bytes() uint16 {
	result := binary.LittleEndian.Uint16(vm.Program[vm.Counter:])
	vm.Counter += 2
	return result
}
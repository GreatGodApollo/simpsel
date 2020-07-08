package vm

import (
	"bytes"
	"fmt"
	"simpsel/ast"
	"simpsel/compiler"
	"simpsel/lexer"
	"simpsel/parser"
	"testing"
)

type vmTestCase struct {
	input    string
	numInst  int
	expected interface{}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testInteger(expected, actual int) error {
	if actual != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d",
			actual, expected)
	}

	return nil
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(comp.Bytecode())

		for i := 0; i < tt.numInst; i ++ {
			vm.executeInstruction(bytes.NewBuffer([]byte{}))
		}

		reg31 := vm.Registers[31]

		testExpectedObject(t, tt.expected, int(reg31))
	}
}

func testExpectedObject(t *testing.T, expected interface{}, actual interface{}) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testInteger(expected, actual.(int))
		if err != nil {
			t.Errorf("testInteger failed: %s", err)
		}
	}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"load $31 #10", 1, 10},
		{"load $0 #1\nload $1 #2\nadd $0 $1 $31", 3, 3},
	}

	runVmTests(t, tests)
}
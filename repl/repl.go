package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"simpsel/compiler"
	"simpsel/lexer"
	"simpsel/parser"
	"simpsel/vm"
)

const PROMPT = ">>> "

func Start(in io.Reader, out io.Writer) {
	run := true
	scanner := bufio.NewScanner(in)
	machine := vm.New(&compiler.Bytecode{Instructions: []byte{}})

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		switch line {
		case ".registers":
			fmt.Fprintf(out, "%v\n", machine.Registers)
		case ".program":
			fmt.Fprintf(out, "%v\n", machine.Program.String())
		case ".pc":
			fmt.Fprintf(out, "%d\n", machine.Counter)
		case ".rpc":
			fmt.Fprint(out, "Counter reset\n")
			machine.Counter = 0
		case ".run":
			machine.Run()
		case ".runonce":
			machine.RunOnce()
		case ".tr":
			run = !run
			fmt.Fprintf(out, "Running? %t\n", run)
		case ".quit":
			os.Exit(0)
		default:
			l := lexer.New(line)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				continue
			}

			comp := compiler.New()
			err := comp.Compile(program)
			if err != nil {
				fmt.Fprintf(out, "Woophs! Compilation failed:\n %s\n", err)
				continue
			}

			machine.Program = append(machine.Program, comp.Bytecode().Instructions...)
			if run {
				machine.RunOnce()
			}
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
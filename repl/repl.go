package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"io/ioutil"
	"os"
	"simpsel/code"
	"simpsel/compiler"
	"simpsel/lexer"
	"simpsel/parser"
	"simpsel/vm"
	"strings"
)

const PROMPT = ">>> "

func Start(in io.Reader, out io.Writer) {
	run := true
	closed := false
	scanner := bufio.NewScanner(in)
	machine := vm.New(&compiler.Bytecode{Instructions: []byte{}})
	fmt.Fprint(out, "Welcome to simpsel. Let's be productive!\n\n")

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		run, closed = handleInput(out, line, machine, run)
		if closed {
			os.Exit(1)
		}
	}
}

func StartTerminal(s ssh.Session, term *terminal.Terminal) {
	run := true
	closed := false
	machine := vm.New(&compiler.Bytecode{Instructions: []byte{}})
	term.Write([]byte("Welcome to simpsel. Let's be productive!\n\n"))

	for closed != true {
		term.Write([]byte(PROMPT))
		line, err := term.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				closed = true
				s.Close()
				return
			}
			fmt.Fprintf(os.Stdout, "An error occured: %s", err)
			s.Close()
			return
		}
		out := bytes.NewBuffer([]byte{})
		run, closed = handleInput(out, line, machine, run)
		term.Write(out.Bytes())
		if closed {
			s.Close()
			return
		}
	}
}

func handleInput(out io.Writer, input string, machine *vm.VM, run bool) (runO, close bool){
	switch input {
	case ".clear_registers":
		machine.Registers = make([]int32, 32)
		fmt.Fprint(out, "Registers cleared\n")
	case ".registers":
		fmt.Fprintf(out, "%v\n", machine.Registers)
	case ".clear_program":
		machine.Program = code.Instructions{}
		fmt.Fprint(out, "Program cleared\n")
	case ".program":
		fmt.Fprintf(out, "BEGIN PROGRAM LISTING\n%v\nEND PROGRAM LISTING\n", machine.Program)
	case ".pc":
		fmt.Fprintf(out, "%d\n", machine.Counter)
	case ".rpc":
		fmt.Fprint(out, "Counter reset\n")
		machine.Counter = 0
	case ".run":
		machine.Run(out)
	case ".run_once":
		machine.RunOnce(out)
	case ".tr":
		runO = !run
		fmt.Fprintf(out, "Running? %t\n", runO)
		return runO, false
	case ".quit":
		fmt.Fprint(out, "Goodbye!\n")
		return run, true

	default:
		if strings.HasPrefix(input, ".load_file") {
			inArr := strings.Split(input, " ")
			if len(inArr) < 2 {
				fmt.Fprintf(out, "You must provide a file to load!")
				return run, false
			}
			fi, err := os.Stat(inArr[1])
			if err != nil {
				fmt.Fprintf(out, "Invalid file!")
				return run, false
			}
			if fi.IsDir() {
				fmt.Fprintf(out, "You can't load a directory!")
				return run, false
			}
			file, _ := os.Open(inArr[1])
			inputb, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Fprintf(out, "An error occured trying to load the file! %s", err)
				return run, false
			}
			input = string(inputb)
		} else if strings.HasPrefix(input, ".") {
			fmt.Fprintf(out, "Unknown command %s\n", input)
			return run, false
		}
		l := lexer.New(input)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			PrintParserErrors(out, p.Errors())
			return run, false
		}

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woophs! Compilation failed:\n %s\n", err)
			return run, false
		}

		machine.Program = append(machine.Program, comp.Bytecode().Instructions...)
		if run {
			machine.RunOnce(out)
		}
	}
	return run, false
}

func PrintParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
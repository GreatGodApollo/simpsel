package main

import (
	"flag"
	"fmt"
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"simpsel/compiler"
	"simpsel/lexer"
	"simpsel/parser"
	"simpsel/repl"
	"simpsel/vm"
)

func main() {
	addr := flag.String("addr", ":2222", "Address to listen on")
	runSsh := flag.Bool("ssh", false, "Run the ssh server?")
	file := flag.String("file", "", "File to run")

	flag.Parse()

	if *file != "" {
		fi, err := os.Stat(*file)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Invalid file!")
			return
		}
		if fi.IsDir() {
			fmt.Fprintf(os.Stdout, "You can't load a directory!")
			return
		}
		f, _ := os.Open(*file)
		input, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Fprintf(os.Stdout, "An error occured trying to load the file! %s", err)
			return
		}

		l := lexer.New(string(input))
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			repl.PrintParserErrors(os.Stdout, p.Errors())
			return
		}

		comp := compiler.New()
		err = comp.Compile(program)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Woophs! Compilation failed:\n %s\n", err)
			return
		}

		machine := vm.New(comp.Bytecode())
		machine.Run(os.Stdout)
		fmt.Fprintf(os.Stdout, "------\nOutput:\nCounter: %d\nRegisters: %v\n", machine.Counter, machine.Registers)
		return
	}

	if *runSsh {
		startSshServer(*addr)
	} else {
		repl.Start(os.Stdin, os.Stdout)
	}
}

func startSshServer(addr string) {
	ssh.Handle(func(s ssh.Session) {
		fmt.Fprintf(os.Stdout, "New connection from: %s\n", s.RemoteAddr())
		term := terminal.NewTerminal(s, "")
		repl.StartTerminal(s, term)
	})

	fmt.Fprintf(os.Stdout, "Starting SSH server @ %s\n", addr)
	err := ssh.ListenAndServe(addr, nil, ssh.HostKeyFile("./host.key"))
	if err != nil {
		fmt.Fprintf(os.Stdout,"Failed to start SSH server: %s", err)
	}
}
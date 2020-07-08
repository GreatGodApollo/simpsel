package main

import (
	"os"
	"simpsel/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}

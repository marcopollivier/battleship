package main

import (
	"os"

	"github.com/marcopollivier/battleship/cmd"
)

func main() {
	os.Exit(cmd.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

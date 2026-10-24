package main

import (
	"fmt"
	"os"

	"github.com/marcopollivier/battleship/cmd"
)

func main() {
	if err := cmd.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

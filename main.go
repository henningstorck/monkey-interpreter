package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/henningstorck/monkey-interpreter/repl"
)

func main() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Hey %s! This is the Monkey programming language.\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}

package main

import (
	"cobble/help"
	"cobble/new"
	"fmt"
	"os"
)

func main() {

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("No arguments provided")
		help.Help()
		os.Exit(0)
	}

	switch args[0] {
	case "new":
		new.Run(args[1:])
	default:
		help.Help()
	}

}

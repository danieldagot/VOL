package main

import (
	"fmt"
	"os"

	"vol/internal/lang"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: vol <file.vol>")
		os.Exit(2)
	}

	path := os.Args[1]
	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "vol: %v\n", err)
		os.Exit(1)
	}

	program, diagnostic := lang.Parse(path, string(source))
	if diagnostic != nil {
		fmt.Fprintln(os.Stderr, diagnostic.Human(string(source)))
		os.Exit(1)
	}

	if diagnostic = lang.Execute(program, os.Stdout); diagnostic != nil {
		fmt.Fprintln(os.Stderr, diagnostic.Human(string(source)))
		os.Exit(1)
	}
}

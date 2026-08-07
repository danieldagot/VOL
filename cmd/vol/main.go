package main

import (
	"fmt"
	"io"
	"os"

	"vol/internal/lang"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	var path string
	switch {
	case len(args) >= 2 && args[0] == "run":
		path = args[1]
	case len(args) == 1 && args[0] != "run":
		// Keep the original command form working for existing users.
		path = args[0]
	default:
		fmt.Fprintln(stderr, "usage: vol run <file.vol>")
		return 2
	}

	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(stderr, "vol: %v\n", err)
		return 1
	}

	program, diagnostic := lang.Parse(path, string(source))
	if diagnostic != nil {
		fmt.Fprintln(stderr, diagnostic.Human(string(source)))
		return 1
	}

	programArgs := []string{}
	if len(args) > 2 {
		programArgs = args[2:]
		if programArgs[0] == "--" {
			programArgs = programArgs[1:]
		}
	}
	if diagnostic = lang.ExecuteWithOptions(program, stdout, lang.ExecuteOptions{Input: os.Stdin, Args: programArgs}); diagnostic != nil {
		fmt.Fprintln(stderr, diagnostic.Human(string(source)))
		return 1
	}
	return 0
}

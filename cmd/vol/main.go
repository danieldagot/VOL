package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"vol/internal/lang"
)

const usage = "usage: vol [--json] run <file.vol> | vol [--json] fmt [--check] <path>"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	// Strip --json flag from any position.
	jsonMode := false
	filtered := make([]string, 0, len(args))
	for _, a := range args {
		if a == "--json" {
			jsonMode = true
		} else {
			filtered = append(filtered, a)
		}
	}
	args = filtered

	reportDiagnostic := func(d *lang.Diagnostic) {
		if jsonMode {
			encoded, _ := json.Marshal(d)
			fmt.Fprintln(stderr, string(encoded))
			return
		}
		source := ""
		if d.File != "" {
			if data, err := os.ReadFile(d.File); err == nil {
				source = string(data)
			}
		}
		fmt.Fprintln(stderr, d.Human(source))
	}

	if len(args) >= 1 && args[0] == "fmt" {
		return runFmt(args[1:], stderr, reportDiagnostic)
	}

	var path string
	switch {
	case len(args) >= 2 && args[0] == "run":
		path = args[1]
	case len(args) == 1 && args[0] != "run" && args[0] != "fmt":
		// Keep the original command form working for existing users.
		path = args[0]
	default:
		fmt.Fprintln(stderr, usage)
		return 2
	}

	if _, err := os.Stat(path); err != nil {
		fmt.Fprintf(stderr, "vol: %v\n", err)
		return 1
	}

	programArgs := []string{}
	if len(args) >= 2 && args[0] == "run" && len(args) > 2 {
		programArgs = args[2:]
		if programArgs[0] == "--" {
			programArgs = programArgs[1:]
		}
	}

	if diagnostic := lang.RunFile(path, stdout, lang.ExecuteOptions{Input: os.Stdin, Args: programArgs}); diagnostic != nil {
		reportDiagnostic(diagnostic)
		return 1
	}
	return 0
}

func runFmt(args []string, stderr io.Writer, reportDiagnostic func(*lang.Diagnostic)) int {
	check := false
	paths := make([]string, 0, len(args))
	for _, a := range args {
		if a == "--check" {
			check = true
			continue
		}
		paths = append(paths, a)
	}
	if len(paths) != 1 {
		fmt.Fprintln(stderr, usage)
		return 2
	}
	path := paths[0]
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(stderr, "vol: %v\n", err)
		return 1
	}
	if _, diagnostic := lang.Parse(path, string(data)); diagnostic != nil {
		reportDiagnostic(diagnostic)
		return 1
	}
	mode := "rewrite"
	if check {
		mode = "format --check"
	}
	fmt.Fprintf(stderr, "vol fmt: %s not implemented yet (see IDEAS.md Formatter and SPEC.md §12)\n", mode)
	return 2
}

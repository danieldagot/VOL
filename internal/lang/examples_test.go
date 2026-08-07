package lang

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExamplesRemainExecutable(t *testing.T) {
	tests := []struct {
		file   string
		input  string
		args   []string
		output string
	}{
		{file: "first.vol", output: "16\n"},
		{file: "hello.vol", output: "Hello from VOL\nPrototype version 1\n"},
		{file: "conditions.vol", output: "Good weather for a walk.\n"},
		{file: "loops.vol", output: "Countdown\n3\n2\n1\nGo!\nGo!\n"},
		{file: "arrays.vol", output: "[72, 95, 81, 64]\nStudents: 4\nHigh score: 95\nHigh score: 81\n"},
		{file: "collections.vol", output: "[7, 9, 12]\nSum: 28\n"},
		{file: "functions.vol", output: "Hello, friend\nSix squared is 36\n"},
		{file: "scope.vol", output: "inside\noutside\n"},
		{file: "arguments.vol", args: []string{"apple", "banana"}, output: "Argument count: 2\nReceived: apple\nReceived: banana\n"},
		{file: "interaction.vol", input: "Ada\n", args: []string{"one"}, output: "What is your name? Hello, Ada\nYou passed 1 command arguments.\n"},
	}
	for _, test := range tests {
		t.Run(test.file, func(t *testing.T) {
			path := filepath.Join("..", "..", "examples", test.file)
			source, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			program, diagnostic := Parse(path, string(source))
			if diagnostic != nil {
				t.Fatalf("parse diagnostic: %s", diagnostic.Human(string(source)))
			}
			var output bytes.Buffer
			diagnostic = ExecuteWithOptions(program, &output, ExecuteOptions{Input: strings.NewReader(test.input), Args: test.args})
			if diagnostic != nil {
				t.Fatalf("runtime diagnostic: %s", diagnostic.Human(string(source)))
			}
			if output.String() != test.output {
				t.Fatalf("output\n got: %q\nwant: %q", output.String(), test.output)
			}
		})
	}
}

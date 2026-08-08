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
		{file: filepath.Join("basics", "first.vol"), output: "16\n"},
		{file: filepath.Join("basics", "hello.vol"), output: "Hello from VOL\nPrototype version 1\n"},
		{file: filepath.Join("basics", "conditions.vol"), output: "Good weather for a walk.\n"},
		{file: filepath.Join("basics", "loops.vol"), output: "Countdown\n3\n2\n1\nGo!\nGo!\n"},
		{file: filepath.Join("basics", "arrays.vol"), output: "[72, 95, 81, 64]\nStudents: 4\nHigh score: 95\nHigh score: 81\n[9, 2, 3]\n[1, 2, 3]\n[10, 20]\n"},
		{file: filepath.Join("basics", "collections.vol"), output: "[7, 9, 12]\nSum: 28\n"},
		{file: filepath.Join("basics", "functions.vol"), output: "Hello, friend\nSix squared is 36\n"},
		{file: filepath.Join("basics", "scope.vol"), output: "inside\noutside\n"},
		{file: filepath.Join("basics", "arguments.vol"), args: []string{"apple", "banana"}, output: "Argument count: 2\nReceived: apple\nReceived: banana\n"},
		{file: filepath.Join("basics", "interaction.vol"), input: "Ada\n", args: []string{"one"}, output: "What is your name? Hello, Ada\nYou passed 1 command arguments.\n"},
		{file: filepath.Join("features", "anonymous.vol"), output: "42\n7\n36\n"},
		{file: filepath.Join("features", "option.vol"), output: "some(VOL)\nHello, VOL\nVOL\nnone\nnot found\nguest\n"},
		{file: filepath.Join("features", "result.vol"), output: "7\nmissing\n"},
		{file: filepath.Join("features", "result_helpers.vol"), output: "5\ndivide by zero\n"},
		{file: filepath.Join("features", "option_result.vol"), output: "VOL\nempty name\n3\nnegative\n"},
		{file: filepath.Join("features", "struct.vol"), output: "Ada\n37\nUser { name: Ada, age: 37 }\nGrace\n"},
		{file: filepath.Join("features", "struct_nested.vol"), output: "A\n1\n9\n9\n"},
		{file: filepath.Join("features", "const_struct.vol"), output: "1\n"},
		{file: filepath.Join("features", "print_multi.vol"), output: "A grades: 1\nSum: 176\nCount: 3\n"},
		{file: filepath.Join("features", "modules", "main.vol"), output: "42\n"},
		{file: filepath.Join("features", "modules", "import_struct.vol"), output: "Hello, Grace\n40\n"},
		{file: filepath.Join("features", "modules", "aliases", "main.vol"), output: "42\n"},
		{file: filepath.Join("projects", "gradebook", "main.vol"), output: "Roster: [Ada, Grace, Alan, Katherine]\nAverage: 82\nPassing: 3\nHonor roll: [Ada, Katherine]\n"},
		{file: filepath.Join("projects", "contacts", "main.vol"), output: "Ada <ada@example.com>\nno contact named Alan\n"},
		{file: filepath.Join("projects", "shop", "main.vol"), output: "Subtotal: 39\ncharged 39\nempty cart\n"},
		{file: filepath.Join("projects", "fibonacci", "main.vol"), output: "0\n1\n1\n2\n3\n5\n8\n13\n"},
	}
	for _, test := range tests {
		t.Run(test.file, func(t *testing.T) {
			path := filepath.Join("..", "..", "examples", test.file)
			var output bytes.Buffer
			diagnostic := RunFile(path, &output, ExecuteOptions{Input: strings.NewReader(test.input), Args: test.args})
			if diagnostic != nil {
				source, _ := os.ReadFile(path)
				t.Fatalf("diagnostic: %s", diagnostic.Human(string(source)))
			}
			if output.String() != test.output {
				t.Fatalf("output\n got: %q\nwant: %q", output.String(), test.output)
			}
		})
	}
}

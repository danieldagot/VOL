package lang

import (
	"bytes"
	"strings"
	"testing"
)

func run(t *testing.T, source string) (string, *Diagnostic) {
	t.Helper()
	program, diagnostic := Parse("test.vol", source)
	if diagnostic != nil {
		return "", diagnostic
	}
	var output bytes.Buffer
	diagnostic = Execute(program, &output)
	return output.String(), diagnostic
}

func TestFirstProgram(t *testing.T) {
	source := `numbers := [4, 7, 2, 9]
total := 0
numbers.each number {
    if number > 5 {
        total = total + number
    }
}
print total`
	output, diagnostic := run(t, source)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	if output != "16\n" {
		t.Fatalf("got %q, want %q", output, "16\\n")
	}
}

func TestRepeatWhileAndScope(t *testing.T) {
	source := `count := 0
repeat 2 { count = count + 1 }
while count < 4 { count = count + 1 }
if count == 4 { print "done" } else { print "bad" }`
	output, diagnostic := run(t, source)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	if output != "done\n" {
		t.Fatalf("got %q", output)
	}
}

func TestArrayIndexAndLength(t *testing.T) {
	output, diagnostic := run(t, `items := [1, 2, 3]
items[1] = 8
print items
print items.length`)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	if output != "[1, 8, 3]\n3\n" {
		t.Fatalf("got %q", output)
	}
}

func TestDiagnosticHasSourceLocation(t *testing.T) {
	_, diagnostic := Parse("broken.vol", "value := 1 + }")
	if diagnostic == nil {
		t.Fatal("expected diagnostic")
	}
	human := diagnostic.Human("value := 1 + }")
	if !strings.Contains(human, "broken.vol:1:14") || !strings.Contains(human, "^") {
		t.Fatalf("unexpected diagnostic:\n%s", human)
	}
}

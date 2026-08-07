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

func TestFunctionsCallsAndVisibility(t *testing.T) {
	source := `fn Add(a, b) {
    return a + b
}
fn double(value) {
    return value * 2
}
export double
print double(Add(2, 3))`
	program, diagnostic := Parse("functions.vol", source)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	public := program.Statements[0].(*FunctionDeclaration)
	private := program.Statements[1].(*FunctionDeclaration)
	if public.Public || !private.Public {
		t.Fatal("visibility was not derived from the export declaration")
	}
	var output bytes.Buffer
	if diagnostic = Execute(program, &output); diagnostic != nil {
		t.Fatal(diagnostic)
	}
	if output.String() != "10\n" {
		t.Fatalf("got %q", output.String())
	}
}

func TestExportMayAppearBeforeDefinitions(t *testing.T) {
	program, diagnostic := Parse("exports.vol", `export start, Config
fn start() { return 1 }
fn Config() { return 2 }`)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	if !program.Statements[1].(*FunctionDeclaration).Public || !program.Statements[2].(*FunctionDeclaration).Public {
		t.Fatal("forward exports were not resolved")
	}
}

func TestUnknownExportIsRejected(t *testing.T) {
	_, diagnostic := Parse("exports.vol", "export missing")
	if diagnostic == nil || diagnostic.Code != "E118" {
		t.Fatalf("got %#v", diagnostic)
	}
}

func TestReturnOutsideFunctionIsRejected(t *testing.T) {
	_, diagnostic := Parse("return.vol", "return 1")
	if diagnostic == nil || diagnostic.Code != "E115" {
		t.Fatalf("got %#v", diagnostic)
	}
}

func TestWhereAndSumSemanticForm(t *testing.T) {
	output, diagnostic := run(t, `numbers := [4, 7, 2, 9]
total := numbers.where(_ > 5).sum
print total`)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	if output != "16\n" {
		t.Fatalf("got %q", output)
	}
}

func TestWhereRequiresBooleanCondition(t *testing.T) {
	_, diagnostic := run(t, `numbers := [1, 2]
result := numbers.where(_ + 1)`)
	if diagnostic == nil || diagnostic.Code != "R022" {
		t.Fatalf("got %#v", diagnostic)
	}
}

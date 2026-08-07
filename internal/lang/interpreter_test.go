package lang

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestInterpreterSupportedValuesAndOperators(t *testing.T) {
	source := `print 1 + 2 * 3
print (1 + 2) * 3
print 7 / 2
print 7.0 / 2
print -2 + 5.5
print 5.5 - 2
print 2.5 * 2
print 1 < 1.5
print 1.5 <= 1.5
print 2.5 > 2
print 2.5 >= 2.5
print 2 <= 2
print 3 > 2
print 3 >= 3
print 1 == 1.0
print 1 != 1.0
print [1, [2, 3]] == [1, [2, 3]]
print "a" + "b"
print not false
print true and true
print false or true`
	want := "7\n9\n3\n3.5\n3.5\n3.5\n5\ntrue\ntrue\ntrue\ntrue\ntrue\ntrue\ntrue\ntrue\nfalse\ntrue\nab\ntrue\ntrue\ntrue\n"
	output, diagnostic := run(t, source)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	if output != want {
		t.Fatalf("output\n got: %q\nwant: %q", output, want)
	}
}

func TestMixedNumericEqualityDoesNotLoseIntegerPrecision(t *testing.T) {
	source := `print 9007199254740992 == 9007199254740992.0
print 9007199254740993 == 9007199254740992.0
print 9007199254740993 != 9007199254740992.0
print -9007199254740993 == -9007199254740992.0
print 9007199254740992.0 == 9007199254740993`
	output, diagnostic := run(t, source)
	if diagnostic != nil || output != "true\nfalse\ntrue\nfalse\nfalse\n" {
		t.Fatalf("output = %q, diagnostic = %#v", output, diagnostic)
	}
}

func TestBooleanOperatorsShortCircuit(t *testing.T) {
	output, diagnostic := run(t, `print false and (1 / 0 == 0)
print true or missingCall()`)
	if diagnostic == nil || diagnostic.Code != "S002" {
		// Name resolution intentionally validates unreachable names before execution.
		t.Fatalf("diagnostic = %#v", diagnostic)
	}
	output, diagnostic = run(t, `fn fail() { return 1 / 0 }
print false and fail()
print true or fail()`)
	if diagnostic != nil || output != "false\ntrue\n" {
		t.Fatalf("output = %q, diagnostic = %#v", output, diagnostic)
	}
}

func TestStandaloneBlocksHaveLexicalScope(t *testing.T) {
	output, diagnostic := run(t, `value := "outer"
{
    value := "inner"
    value = value + "!"
    print value
}
print value`)
	if diagnostic != nil || output != "inner!\nouter\n" {
		t.Fatalf("output = %q, diagnostic = %#v", output, diagnostic)
	}
}

func TestControlFlowEdgeCases(t *testing.T) {
	source := `count := 0
if false { print "bad" } else { print "else" }
repeat 0 { print "bad" }
repeat 3 { count = count + 1
local := count }
while count < 5 { count = count + 1 }
while false { print "bad" }
print count`
	output, diagnostic := run(t, source)
	if diagnostic != nil || output != "else\n5\n" {
		t.Fatalf("output = %q, diagnostic = %#v", output, diagnostic)
	}
}

func TestCollectionsCoverEmptyNestedUnicodeAndMutation(t *testing.T) {
	source := `empty := []
print empty
print empty.length
print empty.sum
print empty.where(_ > 0)
items := [[1], [2, 3]]
items[0][0] = 9
print items
alias := items[1]
alias[0] = 8
print items
print "🙂é".length
total := [1, 2.5, 3].sum
print total
items.each item { print item.length }`
	want := "[]\n0\n0\n[]\n[[9], [2, 3]]\n[[9], [8, 3]]\n2\n6.5\n1\n2\n"
	output, diagnostic := run(t, source)
	if diagnostic != nil || output != want {
		t.Fatalf("output = %q, diagnostic = %#v", output, diagnostic)
	}
}

func TestWhereUsesCurrentItemAndOuterScope(t *testing.T) {
	output, diagnostic := run(t, `minimum := 2
values := [1, 2, 3, 4]
print values.where(_ > minimum).where(_ < 4)`)
	if diagnostic != nil || output != "[3]\n" {
		t.Fatalf("output = %q, diagnostic = %#v", output, diagnostic)
	}
}

func TestFunctionsSupportClosuresRecursionAndEarlyReturn(t *testing.T) {
	source := `fn factorial(n) {
    if n <= 1 { return 1 }
    return n * factorial(n - 1)
}
fn make(value) {
    fn get() { return value }
    return get()
}
fn choose(value) {
    if value { return "yes" }
    return "no"
}
fn no_result() { print "side effect" }
print factorial(5)
print make(7)
print choose(true)
print choose(false)
no_result()`
	want := "120\n7\nyes\nno\nside effect\n"
	output, diagnostic := run(t, source)
	if diagnostic != nil || output != want {
		t.Fatalf("output = %q, diagnostic = %#v", output, diagnostic)
	}
}

func TestBuiltinsHandleInputArgumentsAssertionsAndConversion(t *testing.T) {
	source := `first := input()
second := input("next: ")
third := input()
assert(first == "one")
assert(second == "two", "bad second")
assert(third == "")
print args
print string([1, true, ["x"]])`
	output, diagnostic := runWith(t, source, "one\r\ntwo\n", "a", "b")
	want := "next: [a, b]\n[1, true, [x]]\n"
	if diagnostic != nil || output != want {
		t.Fatalf("output = %q, diagnostic = %#v", output, diagnostic)
	}
}

func TestRuntimeDiagnosticsFromSource(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		code    string
		message string
	}{
		{name: "index non-array read", source: `print "x"[0]`, code: "R003", message: "Only arrays"},
		{name: "index non-array write", source: "text := \"x\"\ntext[0] = \"y\"", code: "R003", message: "Only arrays"},
		{name: "if condition", source: "if 1 {}", code: "R004", message: "Boolean"},
		{name: "while condition", source: "while 1 {}", code: "R004", message: "Boolean"},
		{name: "negative repeat", source: "repeat -1 {}", code: "R005", message: "non-negative integer"},
		{name: "float repeat", source: "repeat 1.0 {}", code: "R005", message: "non-negative integer"},
		{name: "each non-array", source: `"x".each item {}`, code: "R006", message: "requires an array"},
		{name: "unknown property", source: "print [1].missing", code: "R007", message: "Unknown property"},
		{name: "length wrong type", source: "print 1.length", code: "R007", message: "Unknown property"},
		{name: "not wrong type", source: "print not 1", code: "R008", message: "Boolean"},
		{name: "negative wrong type", source: `print -"x"`, code: "R009", message: "number"},
		{name: "and left wrong type", source: "print 1 and true", code: "R010", message: "Boolean"},
		{name: "or left wrong type", source: "print 1 or false", code: "R011", message: "Boolean"},
		{name: "and right wrong type", source: "print true and 1", code: "R012", message: "Boolean"},
		{name: "or right wrong type", source: "print false or 1", code: "R012", message: "Boolean"},
		{name: "operator types", source: `print "x" + 1`, code: "R013", message: "string and integer"},
		{name: "sum item type", source: `print [1, "x"].sum`, code: "R013", message: "integer and string"},
		{name: "integer division zero", source: "print 1 / 0", code: "R014", message: "Division by zero"},
		{name: "float division zero", source: "print 1.0 / 0.0", code: "R014", message: "Division by zero"},
		{name: "index float", source: "print [1][0.0]", code: "R015", message: "integer"},
		{name: "index negative", source: "print [1][-1]", code: "R016", message: "outside length 1"},
		{name: "index at length", source: "print [1][1]", code: "R016", message: "outside length 1"},
		{name: "call non-function", source: "value := 1\nvalue()", code: "R017", message: "Only functions"},
		{name: "indirect arity", source: "fn one(value) { return value }\ncall := one\ncall()", code: "R018", message: "expects 1 arguments, got 0"},
		{name: "sum non-array", source: `print "x".sum`, code: "R019", message: "requires an array"},
		{name: "where non-array", source: `print "x".where(_ == "x")`, code: "R021", message: "requires an array"},
		{name: "where condition", source: "print [1].where(_ + 1)", code: "R022", message: "must be Boolean"},
		{name: "input prompt", source: "input(1)", code: "R023", message: "prompt must be a string"},
		{name: "assert condition", source: "assert(1)", code: "R025", message: "condition must be Boolean"},
		{name: "assert message", source: "assert(false, 1)", code: "R026", message: "message must be a string"},
		{name: "assert default", source: "assert(false)", code: "R027", message: "Assertion failed"},
		{name: "assert custom", source: `assert(false, "custom")`, code: "R027", message: "custom"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, diagnostic := run(t, test.source)
			if output != "" {
				t.Fatalf("unexpected output %q", output)
			}
			if diagnostic == nil || diagnostic.Code != test.code || !strings.Contains(diagnostic.Message, test.message) {
				t.Fatalf("diagnostic = %#v, want %s containing %q", diagnostic, test.code, test.message)
			}
			if diagnostic.File != "test.vol" || diagnostic.Pos.Line < 1 || diagnostic.Pos.Column < 1 {
				t.Fatalf("invalid diagnostic location: %#v", diagnostic)
			}
		})
	}
}

type errorReader struct{ err error }

func (reader errorReader) Read([]byte) (int, error) { return 0, reader.err }

func TestInputReadFailureIsDiagnostic(t *testing.T) {
	program, diagnostic := Parse("input.vol", "input()")
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	var output bytes.Buffer
	diagnostic = ExecuteWithOptions(program, &output, ExecuteOptions{Input: errorReader{err: errors.New("broken reader")}})
	if diagnostic == nil || diagnostic.Code != "R024" || !strings.Contains(diagnostic.Message, "broken reader") {
		t.Fatalf("diagnostic = %#v", diagnostic)
	}
}

func TestInputReturnsPartialLineBeforeReaderError(t *testing.T) {
	program, diagnostic := Parse("input.vol", "print input()")
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	var output bytes.Buffer
	diagnostic = ExecuteWithOptions(program, &output, ExecuteOptions{Input: io.MultiReader(strings.NewReader("partial"), errorReader{err: errors.New("late failure")})})
	if diagnostic != nil || output.String() != "partial\n" {
		t.Fatalf("output = %q, diagnostic = %#v", output.String(), diagnostic)
	}
}

type unknownExpression struct{ pos Position }

func (unknownExpression) expression()                   {}
func (expression unknownExpression) Position() Position { return expression.pos }

func TestInterpreterDefensiveRuntimeDiagnostics(t *testing.T) {
	i := &interpreter{output: io.Discard, input: bufio.NewReader(strings.NewReader("")), env: newEnvironment(nil), file: "direct.vol"}
	position := Position{Line: 2, Column: 3}

	if diagnostic := i.execute(&Declaration{Name: Token{Lexeme: "value", Pos: position}, Value: &Literal{Value: int64(1)}}); diagnostic != nil {
		t.Fatal(diagnostic)
	}
	if diagnostic := i.execute(&Declaration{Name: Token{Lexeme: "value", Pos: position}, Value: &Literal{Value: int64(2)}}); diagnostic == nil || diagnostic.Code != "R001" {
		t.Fatalf("duplicate declaration diagnostic = %#v", diagnostic)
	}
	if _, diagnostic := i.evaluate(&Variable{Name: Token{Lexeme: "missing", Pos: position}}); diagnostic == nil || diagnostic.Code != "R002" {
		t.Fatalf("unknown variable diagnostic = %#v", diagnostic)
	}
	assignment := &Assignment{Target: &Variable{Name: Token{Lexeme: "missing", Pos: position}}, Equal: Token{Pos: position}, Value: &Literal{Value: int64(1)}}
	if diagnostic := i.execute(assignment); diagnostic == nil || diagnostic.Code != "R002" {
		t.Fatalf("unknown assignment diagnostic = %#v", diagnostic)
	}
	property := &Property{Object: &ArrayLiteral{}, Name: Token{Lexeme: "where", Pos: position}}
	if _, diagnostic := i.evaluateWhere(property, nil); diagnostic == nil || diagnostic.Code != "R020" {
		t.Fatalf("where arity diagnostic = %#v", diagnostic)
	}
	if _, diagnostic := i.evaluate(unknownExpression{pos: position}); diagnostic == nil || diagnostic.Code != "R999" {
		t.Fatalf("unknown expression diagnostic = %#v", diagnostic)
	}
}

func TestExecuteWithNilInputUsesEOF(t *testing.T) {
	program, diagnostic := Parse("nil-input.vol", "print input()")
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	var output bytes.Buffer
	if diagnostic = ExecuteWithOptions(program, &output, ExecuteOptions{}); diagnostic != nil || output.String() != "\n" {
		t.Fatalf("output = %q, diagnostic = %#v", output.String(), diagnostic)
	}
}

func TestTypeNamesCoverRuntimeValues(t *testing.T) {
	i := &interpreter{}
	values := []struct {
		value any
		want  string
	}{
		{value: nil, want: "nothing"},
		{value: int64(1), want: "integer"},
		{value: 1.0, want: "float"},
		{value: true, want: "Boolean"},
		{value: "text", want: "string"},
		{value: []any{}, want: "array"},
		{value: &function{}, want: "function"},
		{value: &builtinFunction{}, want: "function"},
		{value: i, want: "value"},
	}
	for _, test := range values {
		if got := typeName(test.value); got != test.want {
			t.Errorf("typeName(%T) = %q, want %q", test.value, got, test.want)
		}
	}
}

func TestEnvironmentAssignmentWalksParents(t *testing.T) {
	parent := newEnvironment(nil)
	parent.define("value", int64(1))
	child := newEnvironment(parent)
	if !child.assign("value", int64(2)) {
		t.Fatal("assignment did not find parent binding")
	}
	if value, ok := parent.get("value"); !ok || value != int64(2) {
		t.Fatalf("parent value = %#v, found = %v", value, ok)
	}
	if child.assign("missing", int64(3)) {
		t.Fatal("assignment unexpectedly created a missing binding")
	}
}

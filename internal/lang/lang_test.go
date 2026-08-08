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

func runWith(t *testing.T, source, input string, args ...string) (string, *Diagnostic) {
	t.Helper()
	program, diagnostic := Parse("test.vol", source)
	if diagnostic != nil {
		return "", diagnostic
	}
	var output bytes.Buffer
	diagnostic = ExecuteWithOptions(program, &output, ExecuteOptions{Input: strings.NewReader(input), Args: args})
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
print items.len`)
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
total := numbers.where(_ > 5).sum()
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

func TestResolverRejectsUndefinedNameBeforeExecution(t *testing.T) {
	_, diagnostic := run(t, "print missing")
	if diagnostic == nil || diagnostic.Code != "S002" || diagnostic.Fix == "" {
		t.Fatalf("got %#v", diagnostic)
	}
}

func TestResolverRejectsUndefinedAssignment(t *testing.T) {
	_, diagnostic := run(t, "missing = 1")
	if diagnostic == nil || diagnostic.Code != "S002" {
		t.Fatalf("got %#v", diagnostic)
	}
}

func TestResolverRejectsDuplicateNames(t *testing.T) {
	cases := []string{
		"value := 1\nvalue := 2",
		"fn work() { return 1 }\nfn work() { return 2 }",
		"fn work(item, item) { return item }",
		"input := 1",
	}
	for _, source := range cases {
		_, diagnostic := run(t, source)
		if diagnostic == nil || diagnostic.Code != "S001" {
			t.Errorf("source %q: got %#v", source, diagnostic)
		}
	}
}

func TestResolverAllowsNestedShadowing(t *testing.T) {
	output, diagnostic := run(t, "value := 1\nif true { value := 2\nprint value }\nprint value")
	if diagnostic != nil || output != "2\n1\n" {
		t.Fatalf("output %q, diagnostic %#v", output, diagnostic)
	}
}

func TestResolverRejectsWrongFunctionArgumentCount(t *testing.T) {
	_, diagnostic := run(t, "fn add(a, b) { return a + b }\nprint add(1)")
	if diagnostic == nil || diagnostic.Code != "S003" {
		t.Fatalf("got %#v", diagnostic)
	}
}

func TestForwardFunctionCall(t *testing.T) {
	output, diagnostic := run(t, "print answer()\nfn answer() { return 42 }")
	if diagnostic != nil || output != "42\n" {
		t.Fatalf("output %q, diagnostic %#v", output, diagnostic)
	}
}

func TestNestedFunctionDeclaration(t *testing.T) {
	output, diagnostic := run(t, "fn outer() { fn inner() { return 9 }\nreturn inner() }\nprint outer()")
	if diagnostic != nil || output != "9\n" {
		t.Fatalf("output %q, diagnostic %#v", output, diagnostic)
	}
}

func TestInputStringAssertAndArgs(t *testing.T) {
	source := `name := input("Name: ")
assert(name == "Ada", "unexpected name")
print "Hello, " + name
print string(42)
print args.len
print args[0]`
	output, diagnostic := runWith(t, source, "Ada\n", "first", "second")
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	if output != "Name: Hello, Ada\n42\n2\nfirst\n" {
		t.Fatalf("got %q", output)
	}
}

func TestAssertFailureUsesMessage(t *testing.T) {
	_, diagnostic := run(t, `assert(false, "not ready")`)
	if diagnostic == nil || diagnostic.Code != "R027" || diagnostic.Message != "not ready" {
		t.Fatalf("got %#v", diagnostic)
	}
}

func TestImprovedRuntimeTypeMessage(t *testing.T) {
	_, diagnostic := run(t, `print 2 + "count"`)
	if diagnostic == nil || diagnostic.Code != "R013" || !strings.Contains(diagnostic.Message, "integer and string") {
		t.Fatalf("got %#v", diagnostic)
	}
}

func TestPrintDisplayForms(t *testing.T) {
	output, diagnostic := run(t, `
print 42
print true
print "hi"
print [1, 2]
print "n=" + 7
print "ok" + "ay"
print "A grades:", 3
print "flag:", true
`)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	want := "42\ntrue\nhi\n[1, 2]\nn=7\nokay\nA grades: 3\nflag: true\n"
	if output != want {
		t.Fatalf("got %q want %q", output, want)
	}
}

func TestPrintRejectsNothing(t *testing.T) {
	_, diagnostic := run(t, `
fn sink() {}
print sink()
`)
	if diagnostic == nil || diagnostic.Code != "R029" {
		t.Fatalf("got %#v", diagnostic)
	}
	_, diagnostic = run(t, `
fn sink() {}
print "x", sink()
`)
	if diagnostic == nil || diagnostic.Code != "R029" {
		t.Fatalf("multi-arg nothing: got %#v", diagnostic)
	}
}

func TestStringConcatCoercion(t *testing.T) {
	output, diagnostic := run(t, `
print "count: " + 2
print "flag: " + true
print "xs: " + [1, 2]
`)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	want := "count: 2\nflag: true\nxs: [1, 2]\n"
	if output != want {
		t.Fatalf("got %q want %q", output, want)
	}
}

func TestStringConcatTypeMismatches(t *testing.T) {
	_, diagnostic := run(t, `print 2 + "count"`)
	if diagnostic == nil || diagnostic.Code != "R013" || !strings.Contains(diagnostic.Message, "integer and string") {
		t.Fatalf("got %#v", diagnostic)
	}
}

func TestCountZeroArgIsLen(t *testing.T) {
	output, diagnostic := run(t, `
nums := [3, 8, 1]
print nums.count()
print nums.where(_ > 5).count()
print "hi".count()
print nums.count(_ > 5)
`)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	want := "3\n1\n2\n1\n"
	if output != want {
		t.Fatalf("got %q want %q", output, want)
	}
}

func TestBuiltinArityIsResolved(t *testing.T) {
	for _, source := range []string{"input(1, 2)", "assert()", "string()"} {
		_, diagnostic := run(t, source)
		if diagnostic == nil || diagnostic.Code != "S003" {
			t.Errorf("source %q: got %#v", source, diagnostic)
		}
	}
}

func TestOptionSomeNoneIfLetAndCoalesce(t *testing.T) {
	output, diagnostic := run(t, `
maybe := some(7)
print maybe
if some n := maybe {
    print n
} else {
    print "missing"
}
print maybe ?? 0
empty := none
print empty
if some n := empty {
    print n
} else {
    print "missing"
}
print empty ?? "guest"
`)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	if output != "some(7)\n7\n7\nnone\nmissing\nguest\n" {
		t.Fatalf("got %q", output)
	}
}

func TestAnonymousFunctionExpression(t *testing.T) {
	output, diagnostic := run(t, `
double := fn(x) {
    return x * 2
}
print double(21)
print fn(a, b) {
    return a + b
}(3, 4)
triple := fn(x) x * 3
print triple(7)
`)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	if output != "42\n7\n21\n" {
		t.Fatalf("got %q", output)
	}
}

func TestMultiAssignAndDestructuring(t *testing.T) {
	output, diagnostic := run(t, `
a, b := 0, 1
print a
print b
a, b = b, a + b
print a
print b
`)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	if output != "0\n1\n1\n1\n" {
		t.Fatalf("got %q", output)
	}
}

func TestResultPropagate(t *testing.T) {
	output, diagnostic := run(t, `
fn divide(a, b) {
    if b == 0 {
        return err("zero")
    }
    return ok(a / b)
}
fn twice(a, b) {
    n := divide(a, b)?
    return ok(n * 2)
}
if ok v := twice(10, 2) {
    print v
} else err e {
    print e
}
if ok v := twice(10, 0) {
    print v
} else err e {
    print e
}
`)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	if output != "10\nzero\n" {
		t.Fatalf("got %q", output)
	}
}

func TestMapAndCount(t *testing.T) {
	output, diagnostic := run(t, `
nums := [1, 2, 3, 4]
print nums.map(_ * 2)
print nums.count(_ > 2)
`)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	if output != "[2, 4, 6, 8]\n2\n" {
		t.Fatalf("got %q", output)
	}
}

func TestPositionalStructLiteral(t *testing.T) {
	output, diagnostic := run(t, `
struct User {
    name
    age
}
u := User { "Ada", 36 }
print u.name
print u.age
`)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	if output != "Ada\n36\n" {
		t.Fatalf("got %q", output)
	}
}

func TestIfLetRequiresOption(t *testing.T) {
	_, diagnostic := run(t, `if some n := 1 {
    print n
} else {
    print "x"
}`)
	if diagnostic == nil || diagnostic.Code != "R034" {
		t.Fatalf("got %#v", diagnostic)
	}
}

func TestResultOkErrAndIfLet(t *testing.T) {
	output, diagnostic := run(t, `
r := ok(7)
print r
if ok n := r {
    print n
} else err msg {
    print msg
}
e := err("nope")
print e
if ok n := e {
    print n
} else err msg {
    print msg
}
`)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	if output != "ok(7)\n7\nerr(nope)\nnope\n" {
		t.Fatalf("got %q", output)
	}
}

func TestIfLetRequiresResult(t *testing.T) {
	_, diagnostic := run(t, `if ok n := 1 {
    print n
} else err e {
    print e
}`)
	if diagnostic == nil || diagnostic.Code != "R037" {
		t.Fatalf("got %#v", diagnostic)
	}
}

func TestMatchRejected(t *testing.T) {
	_, diagnostic := run(t, `match none {
    some n { print n }
    none { print "x" }
}`)
	if diagnostic == nil || diagnostic.Code != "E153" {
		t.Fatalf("got %#v", diagnostic)
	}
	if diagnostic.Fix == "" {
		t.Fatal("expected fix suggestion")
	}
}

func TestCoalesceRejectsResult(t *testing.T) {
	_, diagnostic := run(t, `print ok(1) ?? 0`)
	if diagnostic == nil || diagnostic.Code != "R042" {
		t.Fatalf("got %#v", diagnostic)
	}
}

func TestCoalesceShortCircuit(t *testing.T) {
	output, diagnostic := run(t, `
side := 0
fn bump() {
    side = side + 1
    return 9
}
print some(3) ?? bump()
print side
print none ?? bump()
print side
`)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	if output != "3\n0\n9\n1\n" {
		t.Fatalf("got %q", output)
	}
}

func TestStructConstructAccessAndAssign(t *testing.T) {
	output, diagnostic := run(t, `
struct User {
    name
    age
}
u := User { name: "Ada", age: 36 }
print u.name
u.age = 37
print u.age
print u
`)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	if output != "Ada\n37\nUser { name: Ada, age: 37 }\n" {
		t.Fatalf("got %q", output)
	}
}

func TestStructLiteralRequiresAllFields(t *testing.T) {
	_, diagnostic := run(t, `
struct Point {
    x
    y
}
p := Point { x: 1 }
`)
	if diagnostic == nil || diagnostic.Code != "R039" {
		t.Fatalf("got %#v", diagnostic)
	}
}

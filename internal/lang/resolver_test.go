package lang

import (
	"strings"
	"testing"
)

func TestResolverValidatesEveryExpressionShape(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   string
	}{
		{name: "unary", source: "print -missing", want: "missing"},
		{name: "binary left", source: "print missing + 1", want: "missing"},
		{name: "binary right", source: "print 1 + missing", want: "missing"},
		{name: "array element", source: "print [missing]", want: "missing"},
		{name: "index collection", source: "print missing[0]", want: "missing"},
		{name: "index expression", source: "items := [1]\nprint items[missing]", want: "missing"},
		{name: "property object", source: "print missing.len", want: "missing"},
		{name: "call argument", source: "print string(missing)", want: "missing"},
		{name: "where object", source: "print missing.where(_ > 0)", want: "missing"},
		{name: "where condition", source: "items := [1]\nprint items.where(_ > missing)", want: "missing"},
		{name: "assignment value", source: "value := 1\nvalue = missing", want: "missing"},
		{name: "condition", source: "if missing {}", want: "missing"},
		{name: "repeat count", source: "repeat missing {}", want: "missing"},
		{name: "while condition", source: "while missing {}", want: "missing"},
		{name: "each collection", source: "missing.each item {}", want: "missing"},
		{name: "return value", source: "fn work() { return missing }", want: "missing"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, diagnostic := run(t, test.source)
			if diagnostic == nil || diagnostic.Code != "S002" || !strings.Contains(diagnostic.Message, "`"+test.want+"`") {
				t.Fatalf("diagnostic = %#v", diagnostic)
			}
			if diagnostic.Fix == "" {
				t.Fatal("resolver diagnostic has no suggested fix")
			}
		})
	}
}

func TestUndefinedNameHintsStdAliases(t *testing.T) {
	tests := []struct {
		name   string
		source string
		fix    string
	}{
		{name: "contains", source: `print contains("a", "a")`, fix: "strings.has"},
		{name: "json module path", source: `print json.parse("{}")`, fix: "Import `@std/json`"},
		{name: "stringify", source: `print stringify(1)`, fix: "json.dump"},
		{name: "flat trim", source: `print trim(" x ")`, fix: "strings.trim"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, diagnostic := run(t, test.source)
			if diagnostic == nil || diagnostic.Code != "S002" {
				t.Fatalf("diagnostic = %#v", diagnostic)
			}
			if !strings.Contains(diagnostic.Fix, test.fix) {
				t.Fatalf("Fix = %q want substring %q", diagnostic.Fix, test.fix)
			}
		})
	}
}

func TestResolverRejectsDuplicatesInEveryScope(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{name: "module variables", source: "value := 1\nvalue := 2"},
		{name: "module functions", source: "fn work() { return 1 }\nfn work() { return 2 }"},
		{name: "function and value", source: "fn work() { return 1 }\nwork := 2"},
		{name: "builtin", source: "args := []"},
		{name: "parameters", source: "fn work(item, item) { return item }"},
		{name: "local declarations", source: "fn work() { item := 1\nitem := 2\nreturn item }"},
		{name: "nested function and local", source: "fn work() { fn item() { return 1 }\nitem := 2\nreturn item }"},
		{name: "each item and local", source: "[1].each item { item := 2 }"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, diagnostic := run(t, test.source)
			if diagnostic == nil || diagnostic.Code != "S001" || diagnostic.Fix == "" {
				t.Fatalf("diagnostic = %#v", diagnostic)
			}
		})
	}
}

func TestResolverChecksKnownCallArities(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		message string
	}{
		{name: "function too few", source: "fn add(a, b) { return a + b }\nadd(1)", message: "expects 2 arguments, got 1"},
		{name: "function too many", source: "fn one(a) { return a }\none(1, 2)", message: "expects 1 arguments, got 2"},
		{name: "string", source: "string()", message: "expects 1 arguments, got 0"},
		{name: "input", source: "input(1, 2)", message: "expects zero or one arguments, got 2"},
		{name: "assert none", source: "assert()", message: "expects one or two arguments, got 0"},
		{name: "assert too many", source: "assert(true, \"x\", \"y\")", message: "expects one or two arguments, got 3"},
		{name: "where none", source: "[1].where()", message: "expects 1 arguments, got 0"},
		{name: "where many", source: "[1].where(true, false)", message: "expects 1 arguments, got 2"},
		{name: "count many", source: "[1].count(true, false)", message: "expects 1 arguments, got 2"},
		{name: "count none", source: "[1].count()", message: "expects 1 arguments, got 0"},
		{name: "sum with filter", source: "print [1].sum(_ > 0)", message: "expects 0 arguments, got 1"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, diagnostic := run(t, test.source)
			if diagnostic == nil || diagnostic.Code != "S003" || !strings.Contains(diagnostic.Message, test.message) || diagnostic.Fix == "" {
				t.Fatalf("diagnostic = %#v", diagnostic)
			}
			if test.name == "sum with filter" && !strings.Contains(diagnostic.Fix, ".where(condition).sum()") {
				t.Fatalf("sum arity Fix = %q", diagnostic.Fix)
			}
		})
	}
}

func TestResolverScopesAndForwardReferences(t *testing.T) {
	source := `global := 3
print later()
fn later() {
    local := global
    [1].each item {
        { local := item
          print local }
    }
    return local
}`
	output, diagnostic := run(t, source)
	if diagnostic != nil {
		t.Fatalf("diagnostic = %#v", diagnostic)
	}
	if output != "1\n3\n" {
		t.Fatalf("output = %q", output)
	}
}

func TestResolverKeepsLocalDeclarationsOrderSensitive(t *testing.T) {
	_, diagnostic := run(t, "fn work() { print later\nlater := 1\nreturn later }")
	if diagnostic == nil || diagnostic.Code != "S002" {
		t.Fatalf("diagnostic = %#v", diagnostic)
	}
}

func TestResolverRejectsConstReassignment(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{name: "module const rebind", source: "const limit := 10\nlimit = 11"},
		{name: "local const rebind", source: "fn work() { const n := 1\nn = 2\nreturn n }"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, diagnostic := run(t, test.source)
			if diagnostic == nil || diagnostic.Code != "S030" || diagnostic.Fix == "" {
				t.Fatalf("diagnostic = %#v", diagnostic)
			}
		})
	}
}

func TestResolverAllowsRecursiveAndMutuallyRecursiveModuleFunctions(t *testing.T) {
	source := `fn even(n) { if n == 0 { return true } return odd(n - 1) }
fn odd(n) { if n == 0 { return false } return even(n - 1) }
print even(6)
print odd(7)`
	output, diagnostic := run(t, source)
	if diagnostic != nil || output != "true\ntrue\n" {
		t.Fatalf("output = %q, diagnostic = %#v", output, diagnostic)
	}
}

func TestResolverUsesResolvedSymbolForShadowedBuiltinNames(t *testing.T) {
	source := `fn use() {
    fn input(a, b) { return a + b }
    fn assert(a, b, c) { return a + b + c }
    return input(1, 2) + assert(1, 2, 3)
}
print use()`
	output, diagnostic := run(t, source)
	if diagnostic != nil || output != "9\n" {
		t.Fatalf("output = %q, diagnostic = %#v", output, diagnostic)
	}
}

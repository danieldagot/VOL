package lang

import "testing"

func TestParserBuildsStandaloneBlocksAndPostfixChains(t *testing.T) {
	program, diagnostic := Parse("shape.vol", `{ value := [[1, 2]]
print value[0][1]
print value.where(_ == value[0]).len }`)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	if len(program.Statements) != 1 {
		t.Fatalf("got %d statements", len(program.Statements))
	}
	block, ok := program.Statements[0].(*BlockStatement)
	if !ok || len(block.Body.Statements) != 3 {
		t.Fatalf("statement = %#v", program.Statements[0])
	}
	printStatement := block.Body.Statements[1].(*PrintStatement)
	outerIndex, ok := printStatement.Value.(*Index)
	if !ok {
		t.Fatalf("print expression = %T", printStatement.Value)
	}
	if _, ok := outerIndex.Collection.(*Index); !ok {
		t.Fatalf("nested index collection = %T", outerIndex.Collection)
	}
	property := block.Body.Statements[2].(*PrintStatement).Value.(*Property)
	if property.Name.Lexeme != "len" {
		t.Fatalf("property name = %q", property.Name.Lexeme)
	}
	if _, ok := property.Object.(*Call); !ok {
		t.Fatalf("property object = %T", property.Object)
	}
}

func TestParserOperatorPrecedence(t *testing.T) {
	program, diagnostic := Parse("precedence.vol", "print true or false and 1 + 2 * 3 == 7")
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	expression := program.Statements[0].(*PrintStatement).Value
	or, ok := expression.(*Binary)
	if !ok || or.Operator.Kind != TokenOr {
		t.Fatalf("root = %#v", expression)
	}
	and, ok := or.Right.(*Binary)
	if !ok || and.Operator.Kind != TokenAnd {
		t.Fatalf("or right = %#v", or.Right)
	}
	equality, ok := and.Right.(*Binary)
	if !ok || equality.Operator.Kind != TokenEqualEqual {
		t.Fatalf("and right = %#v", and.Right)
	}
	addition := equality.Left.(*Binary)
	if addition.Operator.Kind != TokenPlus || addition.Right.(*Binary).Operator.Kind != TokenStar {
		t.Fatalf("arithmetic tree = %#v", addition)
	}
}

func TestParserTracksExportsAndVisibility(t *testing.T) {
	program, diagnostic := Parse("exports.vol", `export value, work
value := 1
fn work() { return value }`)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	if len(program.Exports) != 2 || program.Exports[0].Lexeme != "value" || program.Exports[1].Lexeme != "work" {
		t.Fatalf("exports = %#v", program.Exports)
	}
	if !program.Statements[2].(*FunctionDeclaration).Public {
		t.Fatal("exported function was not marked public")
	}
}

func TestParserDiagnostics(t *testing.T) {
	tests := []struct {
		name   string
		source string
		code   string
	}{
		{name: "expression", source: "print }", code: "E101"},
		{name: "assignment target", source: "1 = 2", code: "E102"},
		{name: "each item", source: "[1].each {}", code: "E103"},
		{name: "block open", source: "if true print 1", code: "E104"},
		{name: "block close", source: "if true { print 1", code: "E105"},
		{name: "index close", source: "print [1][0", code: "E106"},
		{name: "property name", source: "print [1].", code: "E107"},
		{name: "group close", source: "print (1", code: "E108"},
		{name: "array close", source: "print [1", code: "E109"},
		{name: "function name", source: "fn () {}", code: "E110"},
		{name: "function parameters open", source: "fn work {}", code: "E111"},
		{name: "parameter name", source: "fn work(a, ) {}", code: "E112"},
		{name: "function parameters close", source: "fn work(a {}", code: "E113"},
		{name: "call close", source: "work(1", code: "E114"},
		{name: "return scope", source: "return 1", code: "E115"},
		{name: "export name", source: "export", code: "E116"},
		{name: "duplicate export", source: "value := 1\nexport value, value", code: "E117"},
		{name: "unknown export", source: "export missing", code: "E118"},
		{name: "top-level statement separator", source: "print 1 print 2", code: "E119"},
		{name: "block statement separator", source: "fn work() { print 1 print 2 }", code: "E119"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, diagnostic := Parse("bad.vol", test.source)
			if diagnostic == nil || diagnostic.Code != test.code {
				t.Fatalf("diagnostic = %#v, want code %s", diagnostic, test.code)
			}
			if diagnostic.File != "bad.vol" || diagnostic.Pos.Line < 1 || diagnostic.Pos.Column < 1 {
				t.Fatalf("invalid diagnostic location: %#v", diagnostic)
			}
		})
	}
}

func TestParserAcceptsReturnInsideNestedFunctionBlocks(t *testing.T) {
	_, diagnostic := Parse("return.vol", "fn choose() { if true { { return 1 } } return 2 }")
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
}

func TestNewlinesSeparateStatements(t *testing.T) {
	output, diagnostic := run(t, `value := 1
[2].each item {
    print item
}

// The alternative may begin on the following line.
if false {
    print "bad"
}
else {
    print value
}`)
	if diagnostic != nil || output != "2\n1\n" {
		t.Fatalf("output = %q, diagnostic = %#v", output, diagnostic)
	}
}

func TestBlockBodiedStatementsAreSelfDelimiting(t *testing.T) {
	output, diagnostic := run(t, `if true { print 1 } print 2
repeat 1 { print 3 } print 4`)
	if diagnostic != nil || output != "1\n2\n3\n4\n" {
		t.Fatalf("output = %q, diagnostic = %#v", output, diagnostic)
	}
}

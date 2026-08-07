package lang

import (
	"reflect"
	"strings"
	"testing"
)

func TestLexerRecognizesCompleteVocabulary(t *testing.T) {
	source := `name _value true false if else repeat while print fn return export and or not 42 3.5 "text" := = == != < <= > >= + - * / { } [ ] ( ) , .`
	tokens, diagnostic := lex("tokens.vol", source)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}

	want := []TokenKind{
		TokenIdentifier, TokenIdentifier,
		TokenTrue, TokenFalse, TokenIf, TokenElse, TokenRepeat, TokenWhile,
		TokenPrint, TokenFn, TokenReturn, TokenExport, TokenAnd, TokenOr, TokenNot,
		TokenInteger, TokenFloat, TokenString,
		TokenColonEqual, TokenEqual, TokenEqualEqual, TokenBangEqual,
		TokenLess, TokenLessEqual, TokenGreater, TokenGreaterEqual,
		TokenPlus, TokenMinus, TokenStar, TokenSlash,
		TokenLeftBrace, TokenRightBrace, TokenLeftBracket, TokenRightBracket,
		TokenLeftParen, TokenRightParen, TokenComma, TokenDot, TokenEOF,
	}
	got := make([]TokenKind, len(tokens))
	for index, token := range tokens {
		got[index] = token.Kind
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("token kinds\n got: %#v\nwant: %#v", got, want)
	}
	if tokens[0].Lexeme != "name" || tokens[15].Lexeme != "42" || tokens[17].Lexeme != "text" {
		t.Fatalf("unexpected lexemes: %#v", tokens)
	}
}

func TestLexerStringEscapesCommentsUnicodeAndPositions(t *testing.T) {
	source := "// ignored α\r\nβ := \"line\\n\\r\\t\\\\\\\"\" // tail"
	tokens, diagnostic := lex("unicode.vol", source)
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	if len(tokens) != 5 {
		t.Fatalf("got %d tokens: %#v", len(tokens), tokens)
	}
	if got := tokens[0]; got.Kind != TokenNewline || got.Pos.Line != 1 {
		t.Fatalf("newline token = %#v", got)
	}
	if got := tokens[1]; got.Lexeme != "β" || got.Pos.Line != 2 || got.Pos.Column != 1 || got.Pos.Offset != 15 {
		t.Fatalf("identifier token = %#v", got)
	}
	if got, want := tokens[3].Lexeme, "line\n\r\t\\\""; got != want {
		t.Fatalf("decoded string = %q, want %q", got, want)
	}
	if got := tokens[4].Pos; got.Line != 2 || got.Column != 30 {
		t.Fatalf("EOF position = %#v", got)
	}
}

func TestLexerKeepsDotSeparateWithoutFractionDigits(t *testing.T) {
	tokens, diagnostic := lex("numbers.vol", "1. 2.0")
	if diagnostic != nil {
		t.Fatal(diagnostic)
	}
	want := []TokenKind{TokenInteger, TokenDot, TokenFloat, TokenEOF}
	for index, kind := range want {
		if tokens[index].Kind != kind {
			t.Fatalf("token %d kind = %v, want %v", index, tokens[index].Kind, kind)
		}
	}
}

func TestLexerDiagnostics(t *testing.T) {
	tests := []struct {
		name   string
		source string
		code   string
		line   int
		column int
	}{
		{name: "colon without equals", source: ":", code: "E001", line: 1, column: 1},
		{name: "bang without equals", source: "!", code: "E002", line: 1, column: 1},
		{name: "unexpected character", source: "\n  @", code: "E003", line: 2, column: 3},
		{name: "string ends at newline", source: "\"open\n", code: "E004", line: 1, column: 1},
		{name: "string ends at eof", source: "\"open", code: "E004", line: 1, column: 1},
		{name: "dangling escape", source: "\"open\\", code: "E004", line: 1, column: 1},
		{name: "unknown escape", source: "\"\\q\"", code: "E005", line: 1, column: 1},
		{name: "integer overflow", source: "9223372036854775808", code: "E006", line: 1, column: 1},
		{name: "float overflow", source: strings.Repeat("9", 400) + ".0", code: "E006", line: 1, column: 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, diagnostic := lex("bad.vol", test.source)
			if diagnostic == nil {
				t.Fatal("expected diagnostic")
			}
			if diagnostic.Code != test.code || diagnostic.File != "bad.vol" || diagnostic.Pos.Line != test.line || diagnostic.Pos.Column != test.column {
				t.Fatalf("diagnostic = %#v", diagnostic)
			}
		})
	}
}

func TestTokenString(t *testing.T) {
	if got := (Token{Kind: TokenEOF}).String(); got != "end of file" {
		t.Fatalf("EOF string = %q", got)
	}
	if got := (Token{Kind: TokenIdentifier, Lexeme: "name"}).String(); got != `"name"` {
		t.Fatalf("identifier string = %q", got)
	}
}

func FuzzParseNeverPanics(f *testing.F) {
	for _, source := range []string{
		"",
		"print 1",
		"fn add(a, b) { return a + b }",
		"numbers := [1, 2, 3]\nnumbers.each number { print number }",
		"\x00\xff{[(\"",
	} {
		f.Add(source)
	}
	f.Fuzz(func(t *testing.T, source string) {
		Parse("fuzz.vol", source)
	})
}

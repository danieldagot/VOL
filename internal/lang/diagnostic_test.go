package lang

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDiagnosticHumanRendering(t *testing.T) {
	diagnostic := &Diagnostic{
		Code:    "S002",
		Message: "Undefined name `missing`.",
		File:    "source.vol",
		Pos:     Position{Offset: 16, Line: 2, Column: 7},
		Fix:     "Declare the name before using it.",
	}
	source := "value := 1\r\nprint missing\r\n"
	want := "error[S002] source.vol:2:7\n\nUndefined name `missing`.\n\n   2 | print missing\n     |       ^\n\nSuggestion: Declare the name before using it."
	if got := diagnostic.Human(source); got != want {
		t.Fatalf("human diagnostic\n got: %q\nwant: %q", got, want)
	}
	if diagnostic.Error() != diagnostic.Message {
		t.Fatalf("Error() = %q", diagnostic.Error())
	}
}

func TestDiagnosticHumanRenderingWithoutSourceOrFix(t *testing.T) {
	diagnostic := &Diagnostic{Code: "E101", Message: "Expected an expression.", File: "empty.vol", Pos: Position{Line: 3, Column: 1}}
	got := diagnostic.Human("")
	if !strings.Contains(got, "empty.vol:3:1") || !strings.Contains(got, "   3 | \n     | ^") || strings.Contains(got, "Suggestion:") {
		t.Fatalf("human diagnostic = %q", got)
	}
}

func TestDiagnosticJSONIsStableAndMachineReadable(t *testing.T) {
	diagnostic := &Diagnostic{
		Code: "E003", Message: "Unexpected character.", File: "bad.vol",
		Pos: Position{Offset: 4, Line: 2, Column: 3}, Fix: "Remove it.",
		Expected: "valid character", Actual: "unexpected", Operation: "lex",
		Repairs: []Repair{{Description: "Remove it."}},
	}
	encoded, err := json.Marshal(diagnostic)
	if err != nil {
		t.Fatal(err)
	}
	var value map[string]any
	if err := json.Unmarshal(encoded, &value); err != nil {
		t.Fatal(err)
	}
	if value["code"] != "E003" || value["message"] != "Unexpected character." || value["file"] != "bad.vol" || value["fix"] != "Remove it." {
		t.Fatalf("JSON = %s", encoded)
	}
	if value["expected"] != "valid character" || value["actual"] != "unexpected" || value["operation"] != "lex" {
		t.Fatalf("agent fields JSON = %s", encoded)
	}
	repairs, ok := value["repairs"].([]any)
	if !ok || len(repairs) != 1 {
		t.Fatalf("repairs JSON = %#v", value["repairs"])
	}
	position, ok := value["position"].(map[string]any)
	if !ok || position["Line"] != float64(2) || position["Column"] != float64(3) || position["Offset"] != float64(4) {
		t.Fatalf("JSON position = %#v", value["position"])
	}
}

func TestSourceLineBoundaries(t *testing.T) {
	if got := sourceLine("one\ntwo", 1); got != "one" {
		t.Fatalf("line 1 = %q", got)
	}
	if got := sourceLine("one\ntwo", 2); got != "two" {
		t.Fatalf("line 2 = %q", got)
	}
	for _, line := range []int{-1, 0, 3} {
		if got := sourceLine("one\ntwo", line); got != "" {
			t.Fatalf("line %d = %q", line, got)
		}
	}
}

package lang

import (
	"fmt"
	"strings"
)

// Repair is a structured agent-oriented fix hint (SF-3.1).
type Repair struct {
	Description string `json:"description"`
	Replacement string `json:"replacement,omitempty"`
}

type Diagnostic struct {
	Code      string   `json:"code"`
	Message   string   `json:"message"`
	File      string   `json:"file"`
	Pos       Position `json:"position"`
	Fix       string   `json:"fix,omitempty"`
	Expected  string   `json:"expected,omitempty"`
	Actual    string   `json:"actual,omitempty"`
	Operation string   `json:"operation,omitempty"`
	Repairs   []Repair `json:"repairs,omitempty"`
}

func (d *Diagnostic) Error() string { return d.Message }

func (d *Diagnostic) Human(source string) string {
	line := sourceLine(source, d.Pos.Line)
	pointer := strings.Repeat(" ", max(d.Pos.Column-1, 0)) + "^"
	location := fmt.Sprintf("%s:%d:%d", d.File, d.Pos.Line, d.Pos.Column)
	result := fmt.Sprintf("error[%s] %s\n\n%s\n\n%4d | %s\n     | %s", d.Code, location, d.Message, d.Pos.Line, line, pointer)
	if d.Fix != "" {
		result += "\n\nSuggestion: " + d.Fix
	}
	return result
}

func sourceLine(source string, wanted int) string {
	lines := strings.Split(source, "\n")
	if wanted < 1 || wanted > len(lines) {
		return ""
	}
	return strings.TrimSuffix(lines[wanted-1], "\r")
}

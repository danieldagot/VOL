package lang

import (
	"strings"
)

func stdStringsModule() *stdModule {
	return &stdModule{exports: map[string]any{
		"trim":    stdBuiltin("trim", 1, 1, stdTrim),
		"split":   stdBuiltin("split", 2, 2, stdSplit),
		"join":    stdBuiltin("join", 2, 2, stdJoin),
		"has":     stdBuiltin("has", 2, 2, stdHas),
		"prefix":  stdBuiltin("prefix", 2, 2, stdPrefix),
		"suffix":  stdBuiltin("suffix", 2, 2, stdSuffix),
		"replace": stdBuiltin("replace", 3, 3, stdReplace),
	}}
}

func stdTrim(values []any, pos Position) (any, *Diagnostic) {
	s, d := requireStringArg("trim", values, 0, pos)
	if d != nil {
		return nil, d
	}
	return strings.TrimSpace(s), nil
}

func stdSplit(values []any, pos Position) (any, *Diagnostic) {
	s, d := requireStringArg("split", values, 0, pos)
	if d != nil {
		return nil, d
	}
	sep, d := requireStringArg("split", values, 1, pos)
	if d != nil {
		return nil, d
	}
	parts := strings.Split(s, sep)
	out := make([]any, len(parts))
	for i, part := range parts {
		out[i] = part
	}
	return out, nil
}

func stdJoin(values []any, pos Position) (any, *Diagnostic) {
	parts, ok := values[0].([]any)
	if !ok {
		return nil, &Diagnostic{Code: "R047", Message: "`join` argument 1 must be an array, got " + typeName(values[0]) + ".", Pos: pos}
	}
	sep, d := requireStringArg("join", values, 1, pos)
	if d != nil {
		return nil, d
	}
	texts := make([]string, len(parts))
	for i, part := range parts {
		text, ok := part.(string)
		if !ok {
			return nil, &Diagnostic{Code: "R047", Message: "`join` array elements must be strings, got " + typeName(part) + ".", Pos: pos}
		}
		texts[i] = text
	}
	return strings.Join(texts, sep), nil
}

func stdHas(values []any, pos Position) (any, *Diagnostic) {
	s, d := requireStringArg("has", values, 0, pos)
	if d != nil {
		return nil, d
	}
	sub, d := requireStringArg("has", values, 1, pos)
	if d != nil {
		return nil, d
	}
	return strings.Contains(s, sub), nil
}

func stdPrefix(values []any, pos Position) (any, *Diagnostic) {
	s, d := requireStringArg("prefix", values, 0, pos)
	if d != nil {
		return nil, d
	}
	prefix, d := requireStringArg("prefix", values, 1, pos)
	if d != nil {
		return nil, d
	}
	return strings.HasPrefix(s, prefix), nil
}

func stdSuffix(values []any, pos Position) (any, *Diagnostic) {
	s, d := requireStringArg("suffix", values, 0, pos)
	if d != nil {
		return nil, d
	}
	suffix, d := requireStringArg("suffix", values, 1, pos)
	if d != nil {
		return nil, d
	}
	return strings.HasSuffix(s, suffix), nil
}

func stdReplace(values []any, pos Position) (any, *Diagnostic) {
	s, d := requireStringArg("replace", values, 0, pos)
	if d != nil {
		return nil, d
	}
	old, d := requireStringArg("replace", values, 1, pos)
	if d != nil {
		return nil, d
	}
	neu, d := requireStringArg("replace", values, 2, pos)
	if d != nil {
		return nil, d
	}
	return strings.ReplaceAll(s, old, neu), nil
}

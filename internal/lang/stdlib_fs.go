package lang

import (
	"os"
	"path/filepath"
)

func stdFSModule() *stdModule {
	return &stdModule{exports: map[string]any{
		"exists": stdBuiltin("exists", 1, 1, stdExists),
		"read":   stdBuiltin("read", 1, 1, stdRead),
		"write":  stdBuiltin("write", 2, 2, stdWrite),
		"list":   stdBuiltin("list", 1, 1, stdList),
	}}
}

func stdExists(values []any, pos Position) (any, *Diagnostic) {
	path, d := requireStringArg("exists", values, 0, pos)
	if d != nil {
		return nil, d
	}
	_, err := os.Stat(path)
	return err == nil, nil
}

func stdRead(values []any, pos Position) (any, *Diagnostic) {
	path, d := requireStringArg("read", values, 0, pos)
	if d != nil {
		return nil, d
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return errResult(err.Error()), nil
	}
	return okResult(string(data)), nil
}

func stdWrite(values []any, pos Position) (any, *Diagnostic) {
	path, d := requireStringArg("write", values, 0, pos)
	if d != nil {
		return nil, d
	}
	text, d := requireStringArg("write", values, 1, pos)
	if d != nil {
		return nil, d
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && filepath.Dir(path) != "." {
		return errResult(err.Error()), nil
	}
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		return errResult(err.Error()), nil
	}
	return okResult(true), nil
}

func stdList(values []any, pos Position) (any, *Diagnostic) {
	path, d := requireStringArg("list", values, 0, pos)
	if d != nil {
		return nil, d
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return errResult(err.Error()), nil
	}
	names := make([]any, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return okResult(names), nil
}

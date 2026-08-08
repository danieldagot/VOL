package lang

import (
	"path/filepath"
)

func stdPathModule() *stdModule {
	return &stdModule{exports: map[string]any{
		"join": stdBuiltin("join", 1, -1, stdPathJoin),
		"base": stdBuiltin("base", 1, 1, stdPathBase),
		"dir":  stdBuiltin("dir", 1, 1, stdPathDir),
		"ext":  stdBuiltin("ext", 1, 1, stdPathExt),
	}}
}

func stdPathJoin(values []any, pos Position) (any, *Diagnostic) {
	parts := make([]string, len(values))
	for i, value := range values {
		s, ok := value.(string)
		if !ok {
			return nil, &Diagnostic{Code: "R047", Message: "`join` arguments must be strings, got " + typeName(value) + ".", Pos: pos}
		}
		parts[i] = s
	}
	return filepath.ToSlash(filepath.Join(parts...)), nil
}

func stdPathBase(values []any, pos Position) (any, *Diagnostic) {
	path, d := requireStringArg("base", values, 0, pos)
	if d != nil {
		return nil, d
	}
	return filepath.Base(filepath.FromSlash(path)), nil
}

func stdPathDir(values []any, pos Position) (any, *Diagnostic) {
	path, d := requireStringArg("dir", values, 0, pos)
	if d != nil {
		return nil, d
	}
	return filepath.ToSlash(filepath.Dir(filepath.FromSlash(path))), nil
}

func stdPathExt(values []any, pos Position) (any, *Diagnostic) {
	path, d := requireStringArg("ext", values, 0, pos)
	if d != nil {
		return nil, d
	}
	return filepath.Ext(filepath.FromSlash(path)), nil
}

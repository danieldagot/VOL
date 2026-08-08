package lang

import (
	"fmt"
	"strings"
)

const stdModulePrefix = "volstd:"

// nativeHandle is an opaque host resource (DB connection, HTTP reply, …).
type nativeHandle struct {
	kind  string
	value any
}

func (h *nativeHandle) String() string {
	return "<" + h.kind + ">"
}

type stdModule struct {
	exports map[string]any
}

func stdModulePath(importPath string) string {
	return stdModulePrefix + importPath
}

func isStdModulePath(path string) bool {
	return strings.HasPrefix(path, stdModulePrefix)
}

func stdImportFromPath(path string) string {
	return strings.TrimPrefix(path, stdModulePrefix)
}

func isReservedStdImport(importPath string) bool {
	return importPath == "@std" || strings.HasPrefix(importPath, "@std/")
}

func knownStdModules() map[string]func() *stdModule {
	return map[string]func() *stdModule{
		"@std/math":    stdMathModule,
		"@std/strings": stdStringsModule,
		"@std/fs":      stdFSModule,
		"@std/path":    stdPathModule,
		"@std/env":     stdEnvModule,
		"@std/time":    stdTimeModule,
		"@std/url":     stdURLModule,
		"@std/json":    stdJSONModule,
		"@std/yaml":    stdYAMLModule,
		"@std/http":    stdHTTPModule,
		"@std/process": stdProcessModule,
		"@std/db":      stdDBModule,
	}
}

func loadStdModule(importPath string) (*stdModule, *Diagnostic) {
	if importPath == "@std" {
		return nil, &Diagnostic{
			Code:    "S031",
			Message: "Module `@std` is a reserved root; import a submodule such as `@std/math`.",
			Fix:     "Use `import \"@std/<module>\"` (for example `@std/http`).",
		}
	}
	factory, ok := knownStdModules()[importPath]
	if !ok {
		return nil, &Diagnostic{
			Code:    "S031",
			Message: "Unknown standard module `" + importPath + "`.",
			Fix:     "Supported `@std` modules: math, strings, fs, path, env, time, url, json, yaml, http, process, db.",
		}
	}
	return factory(), nil
}

func stdSymbols(mod *stdModule) map[string]symbol {
	out := map[string]symbol{}
	for name, value := range mod.exports {
		out[name] = symbolFromValue(value)
	}
	return out
}

func stdBuiltin(name string, minArgs, maxArgs int, call func([]any, Position) (any, *Diagnostic)) *builtinFunction {
	return &builtinFunction{name: name, call: func(values []any, pos Position) (any, *Diagnostic) {
		if len(values) < minArgs || (maxArgs >= 0 && len(values) > maxArgs) {
			want := fmt.Sprintf("%d", minArgs)
			if maxArgs != minArgs {
				if maxArgs < 0 {
					want = fmt.Sprintf("at least %d", minArgs)
				} else {
					want = fmt.Sprintf("%d or %d", minArgs, maxArgs)
				}
			}
			return nil, &Diagnostic{
				Code:    "R018",
				Message: fmt.Sprintf("Function `%s` expects %s arguments, got %d.", name, want, len(values)),
				Pos:     pos,
			}
		}
		return call(values, pos)
	}}
}

func okResult(value any) resultValue {
	return resultValue{ok: true, value: value}
}

func errResult(message string) resultValue {
	return resultValue{ok: false, value: message}
}

func someOption(value any) optionValue {
	return optionValue{present: true, value: value}
}

func noneOption() optionValue {
	return optionValue{}
}

func asString(value any) (string, bool) {
	s, ok := value.(string)
	return s, ok
}

func asInt64(value any) (int64, bool) {
	n, ok := value.(int64)
	return n, ok
}

func asFloat64(value any) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case int64:
		return float64(v), true
	default:
		return 0, false
	}
}

func requireStringArg(name string, values []any, index int, pos Position) (string, *Diagnostic) {
	s, ok := values[index].(string)
	if !ok {
		return "", &Diagnostic{
			Code:    "R047",
			Message: fmt.Sprintf("`%s` argument %d must be a string, got %s.", name, index+1, typeName(values[index])),
			Pos:     pos,
		}
	}
	return s, nil
}

func requireIntArg(name string, values []any, index int, pos Position) (int64, *Diagnostic) {
	n, ok := values[index].(int64)
	if !ok {
		return 0, &Diagnostic{
			Code:    "R047",
			Message: fmt.Sprintf("`%s` argument %d must be an integer, got %s.", name, index+1, typeName(values[index])),
			Pos:     pos,
		}
	}
	return n, nil
}

func numberArg(name string, values []any, index int, pos Position) (float64, *Diagnostic) {
	n, ok := asFloat64(values[index])
	if !ok {
		return 0, &Diagnostic{
			Code:    "R047",
			Message: fmt.Sprintf("`%s` argument %d must be a number, got %s.", name, index+1, typeName(values[index])),
			Pos:     pos,
		}
	}
	return n, nil
}

func structOf(name string, fields []string, values map[string]any) *structValue {
	return &structValue{typ: &structType{name: name, fields: fields}, fields: values}
}

func optsString(d *dictValue, key string) (string, bool) {
	if d == nil {
		return "", false
	}
	value, ok := d.get(key)
	if !ok {
		return "", false
	}
	s, ok := value.(string)
	return s, ok
}

package lang

import "os"

func stdEnvModule() *stdModule {
	return &stdModule{exports: map[string]any{
		"get": stdBuiltin("get", 1, 1, stdEnvGet),
		"set": stdBuiltin("set", 2, 2, stdEnvSet),
	}}
}

func stdEnvGet(values []any, pos Position) (any, *Diagnostic) {
	key, d := requireStringArg("get", values, 0, pos)
	if d != nil {
		return nil, d
	}
	value, ok := os.LookupEnv(key)
	if !ok {
		return noneOption(), nil
	}
	return someOption(value), nil
}

func stdEnvSet(values []any, pos Position) (any, *Diagnostic) {
	key, d := requireStringArg("set", values, 0, pos)
	if d != nil {
		return nil, d
	}
	value, d := requireStringArg("set", values, 1, pos)
	if d != nil {
		return nil, d
	}
	if err := os.Setenv(key, value); err != nil {
		return errResult(err.Error()), nil
	}
	return okResult(true), nil
}

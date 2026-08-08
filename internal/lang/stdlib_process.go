package lang

import (
	"bytes"
	"os"
	"os/exec"
)

func stdProcessModule() *stdModule {
	return &stdModule{exports: map[string]any{
		"run": stdBuiltin("run", 1, 2, stdProcessRun),
	}}
}

func stdProcessRun(values []any, pos Position) (any, *Diagnostic) {
	argvVal, ok := values[0].([]any)
	if !ok {
		return nil, &Diagnostic{Code: "R047", Message: "`run` argument 1 must be an array of strings, got " + typeName(values[0]) + ".", Pos: pos}
	}
	if len(argvVal) == 0 {
		return errResult("run argv must not be empty"), nil
	}
	argv := make([]string, len(argvVal))
	for i, item := range argvVal {
		s, ok := item.(string)
		if !ok {
			return nil, &Diagnostic{Code: "R047", Message: "`run` argv elements must be strings, got " + typeName(item) + ".", Pos: pos}
		}
		argv[i] = s
	}
	opts, d := requireDictOpts(values, 1, pos, "")
	if d != nil {
		return nil, d
	}
	cmd := exec.Command(argv[0], argv[1:]...)
	if cwd, ok := optsString(opts, "cwd"); ok {
		cmd.Dir = cwd
	}
	if envVal, ok := opts.get("env"); ok {
		envDict, ok := envVal.(*dictValue)
		if !ok {
			return nil, &Diagnostic{Code: "R047", Message: "`run` opts.env must be a dict, got " + typeName(envVal) + ".", Pos: pos}
		}
		env := os.Environ()
		for _, key := range envDict.keys() {
			k := key.(string)
			v, _ := envDict.get(k)
			text, ok := v.(string)
			if !ok {
				return nil, &Diagnostic{Code: "R047", Message: "`run` opts.env values must be strings.", Pos: pos}
			}
			env = append(env, k+"="+text)
		}
		cmd.Env = env
	}
	if stdin, ok := optsString(opts, "stdin"); ok {
		cmd.Stdin = bytes.NewBufferString(stdin)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	status := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			status = exitErr.ExitCode()
		} else {
			return errResult(err.Error()), nil
		}
	}
	return okResult(structOf("Proc", []string{"status", "stdout", "stderr"}, map[string]any{
		"status": int64(status),
		"stdout": stdout.String(),
		"stderr": stderr.String(),
	})), nil
}

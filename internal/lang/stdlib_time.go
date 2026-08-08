package lang

import (
	"time"
)

func stdTimeModule() *stdModule {
	return &stdModule{exports: map[string]any{
		"now":    stdBuiltin("now", 0, 0, stdTimeNow),
		"sleep":  stdBuiltin("sleep", 1, 1, stdTimeSleep),
		"format": stdBuiltin("format", 1, 2, stdTimeFormat),
	}}
}

func stdTimeNow(values []any, pos Position) (any, *Diagnostic) {
	return time.Now().UnixMilli(), nil
}

func stdTimeSleep(values []any, pos Position) (any, *Diagnostic) {
	ms, d := requireIntArg("sleep", values, 0, pos)
	if d != nil {
		return nil, d
	}
	if ms < 0 {
		return errResult("sleep duration must be non-negative"), nil
	}
	time.Sleep(time.Duration(ms) * time.Millisecond)
	return okResult(true), nil
}

func stdTimeFormat(values []any, pos Position) (any, *Diagnostic) {
	ms, d := requireIntArg("format", values, 0, pos)
	if d != nil {
		return nil, d
	}
	layout := time.RFC3339Nano
	if len(values) == 2 {
		text, d := requireStringArg("format", values, 1, pos)
		if d != nil {
			return nil, d
		}
		layout = text
	}
	t := time.UnixMilli(ms).UTC()
	return t.Format(layout), nil
}

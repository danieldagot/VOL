package lang

import (
	"fmt"
	"sort"
	"strings"
)

// dictValue is a mutable string-key map (SF-3). Construct with dict() or
// dict("k", v, …); index with ["k"].
type dictValue struct {
	entries map[string]any
}

func newDict() *dictValue {
	return &dictValue{entries: map[string]any{}}
}

// builtinDict builds an empty dict (0 args) or a dict from alternating string
// keys and values. Odd arity or non-string keys fail with R018 / R045.
func builtinDict(values []any, pos Position, file string) (any, *Diagnostic) {
	if len(values)%2 != 0 {
		return nil, &Diagnostic{
			Code:    "R018",
			Message: "Function `dict` expects an even number of arguments (key, value pairs), got " + fmt.Sprint(len(values)) + ".",
			File:    file,
			Pos:     pos,
			Fix:     "Use `dict()` for empty, or `dict(\"k\", v, …)` with alternating string keys and values.",
		}
	}
	d := newDict()
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, &Diagnostic{
				Code:    "R045",
				Message: "Dict key must be a string, got " + typeName(values[i]) + ".",
				File:    file,
				Pos:     pos,
			}
		}
		d.set(key, values[i+1])
	}
	return d, nil
}

func (d *dictValue) len() int64 {
	return int64(len(d.entries))
}

func (d *dictValue) keys() []any {
	keys := make([]string, 0, len(d.entries))
	for key := range d.entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]any, len(keys))
	for i, key := range keys {
		out[i] = key
	}
	return out
}

func (d *dictValue) get(key string) (any, bool) {
	value, ok := d.entries[key]
	return value, ok
}

func (d *dictValue) set(key string, value any) {
	d.entries[key] = value
}

func displayDict(d *dictValue) string {
	keys := make([]string, 0, len(d.entries))
	for key := range d.entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s: %s", key, display(d.entries[key])))
	}
	return "dict { " + strings.Join(parts, ", ") + " }"
}

func deepCopyDict(d *dictValue) *dictValue {
	out := newDict()
	for key, value := range d.entries {
		out.entries[key] = deepCopyValue(value)
	}
	return out
}

func dictFromStringMap(m map[string]string) *dictValue {
	d := newDict()
	for key, value := range m {
		d.entries[key] = value
	}
	return d
}

func asDict(value any) (*dictValue, bool) {
	d, ok := value.(*dictValue)
	return d, ok
}

func requireDictOpts(values []any, index int, pos Position, file string) (*dictValue, *Diagnostic) {
	if index >= len(values) {
		return newDict(), nil
	}
	d, ok := values[index].(*dictValue)
	if !ok {
		return nil, &Diagnostic{
			Code:    "R046",
			Message: "Options must be a dict, got " + typeName(values[index]) + ".",
			File:    file,
			Pos:     pos,
		}
	}
	return d, nil
}

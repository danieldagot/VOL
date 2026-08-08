package lang

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func stdYAMLModule() *stdModule {
	return &stdModule{exports: map[string]any{
		"parse": stdBuiltin("parse", 1, 1, stdYAMLParse),
		"dump":  stdBuiltin("dump", 1, 1, stdYAMLDump),
	}}
}

func stdYAMLParse(values []any, pos Position) (any, *Diagnostic) {
	text, d := requireStringArg("parse", values, 0, pos)
	if d != nil {
		return nil, d
	}
	var raw any
	if err := yaml.Unmarshal([]byte(text), &raw); err != nil {
		return errResult(err.Error()), nil
	}
	value, err := fromYAML(raw)
	if err != nil {
		return errResult(err.Error()), nil
	}
	return okResult(value), nil
}

func stdYAMLDump(values []any, pos Position) (any, *Diagnostic) {
	data, err := toJSON(values[0]) // same VOL→host mapping as JSON
	if err != nil {
		return errResult(err.Error()), nil
	}
	encoded, err := yaml.Marshal(data)
	if err != nil {
		return errResult(err.Error()), nil
	}
	return okResult(string(encoded)), nil
}

func fromYAML(raw any) (any, error) {
	switch v := raw.(type) {
	case nil:
		return noneOption(), nil
	case bool:
		return v, nil
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case uint64:
		return int64(v), nil
	case float64:
		if float64(int64(v)) == v {
			return int64(v), nil
		}
		return v, nil
	case string:
		return v, nil
	case []any:
		out := make([]any, len(v))
		for i, item := range v {
			converted, err := fromYAML(item)
			if err != nil {
				return nil, err
			}
			out[i] = converted
		}
		return out, nil
	case map[string]any:
		d := newDict()
		for key, item := range v {
			converted, err := fromYAML(item)
			if err != nil {
				return nil, err
			}
			d.set(key, converted)
		}
		return d, nil
	case map[any]any:
		d := newDict()
		for key, item := range v {
			text, ok := key.(string)
			if !ok {
				return nil, fmt.Errorf("YAML object keys must be strings")
			}
			converted, err := fromYAML(item)
			if err != nil {
				return nil, err
			}
			d.set(text, converted)
		}
		return d, nil
	default:
		return nil, fmt.Errorf("unsupported YAML value %T", raw)
	}
}

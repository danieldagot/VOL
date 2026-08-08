package lang

import (
	"encoding/json"
	"fmt"
)

func stdJSONModule() *stdModule {
	return &stdModule{exports: map[string]any{
		"parse": stdBuiltin("parse", 1, 1, stdJSONParse),
		"dump":  stdBuiltin("dump", 1, 1, stdJSONDump),
	}}
}

func stdJSONParse(values []any, pos Position) (any, *Diagnostic) {
	text, d := requireStringArg("parse", values, 0, pos)
	if d != nil {
		return nil, d
	}
	var raw any
	if err := json.Unmarshal([]byte(text), &raw); err != nil {
		return errResult(err.Error()), nil
	}
	value, err := fromJSON(raw)
	if err != nil {
		return errResult(err.Error()), nil
	}
	return okResult(value), nil
}

func stdJSONDump(values []any, pos Position) (any, *Diagnostic) {
	data, err := toJSON(values[0])
	if err != nil {
		return errResult(err.Error()), nil
	}
	encoded, err := json.Marshal(data)
	if err != nil {
		return errResult(err.Error()), nil
	}
	return okResult(string(encoded)), nil
}

func fromJSON(raw any) (any, error) {
	switch v := raw.(type) {
	case nil:
		return noneOption(), nil
	case bool:
		return v, nil
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
			converted, err := fromJSON(item)
			if err != nil {
				return nil, err
			}
			out[i] = converted
		}
		return out, nil
	case map[string]any:
		d := newDict()
		for key, item := range v {
			converted, err := fromJSON(item)
			if err != nil {
				return nil, err
			}
			d.set(key, converted)
		}
		return d, nil
	default:
		return nil, fmt.Errorf("unsupported JSON value %T", raw)
	}
}

func toJSON(value any) (any, error) {
	switch v := value.(type) {
	case nil:
		return nil, fmt.Errorf("cannot dump nothing")
	case optionValue:
		if !v.present {
			return nil, nil
		}
		return toJSON(v.value)
	case resultValue:
		return nil, fmt.Errorf("cannot dump Result")
	case bool, string:
		return v, nil
	case int64:
		return v, nil
	case float64:
		return v, nil
	case []any:
		out := make([]any, len(v))
		for i, item := range v {
			converted, err := toJSON(item)
			if err != nil {
				return nil, err
			}
			out[i] = converted
		}
		return out, nil
	case *dictValue:
		out := map[string]any{}
		for key, item := range v.entries {
			converted, err := toJSON(item)
			if err != nil {
				return nil, err
			}
			out[key] = converted
		}
		return out, nil
	case *structValue:
		out := map[string]any{}
		for _, key := range v.typ.fields {
			converted, err := toJSON(v.fields[key])
			if err != nil {
				return nil, err
			}
			out[key] = converted
		}
		return out, nil
	default:
		return nil, fmt.Errorf("cannot dump %s", typeName(value))
	}
}

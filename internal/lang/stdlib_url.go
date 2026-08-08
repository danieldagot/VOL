package lang

import (
	"net/url"
	"strconv"
)

func stdURLModule() *stdModule {
	return &stdModule{exports: map[string]any{
		"parse": stdBuiltin("parse", 1, 1, stdURLParse),
	}}
}

func stdURLParse(values []any, pos Position) (any, *Diagnostic) {
	raw, d := requireStringArg("parse", values, 0, pos)
	if d != nil {
		return nil, d
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return errResult(err.Error()), nil
	}
	if parsed.Scheme == "" {
		return errResult("url missing scheme"), nil
	}
	query := newDict()
	for key, vals := range parsed.Query() {
		if len(vals) == 0 {
			query.set(key, "")
			continue
		}
		query.set(key, vals[0])
	}
	var port int64
	if parsed.Port() != "" {
		n, err := strconv.ParseInt(parsed.Port(), 10, 64)
		if err != nil {
			return errResult("invalid port"), nil
		}
		port = n
	}
	host := parsed.Hostname()
	return okResult(structOf("Url", []string{"scheme", "host", "port", "path", "query"}, map[string]any{
		"scheme": parsed.Scheme,
		"host":   host,
		"port":   port,
		"path":   parsed.Path,
		"query":  query,
	})), nil
}

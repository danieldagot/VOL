package lang

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type httpReplyValue struct {
	status      int
	body        string
	contentType string
}

func stdHTTPModule() *stdModule {
	return &stdModule{exports: map[string]any{
		"fetch":  stdBuiltin("fetch", 1, 2, stdHTTPFetch),
		"listen": stdBuiltin("listen", 2, 3, stdHTTPListen),
		"reply":  stdBuiltin("reply", 2, 2, stdHTTPReply),
	}}
}

func stdHTTPReply(values []any, pos Position) (any, *Diagnostic) {
	status, d := requireIntArg("reply", values, 0, pos)
	if d != nil {
		return nil, d
	}
	body := values[1]
	contentType := "text/plain; charset=utf-8"
	text := ""
	switch v := body.(type) {
	case string:
		text = v
	case *dictValue, []any:
		data, err := toJSON(v)
		if err != nil {
			return errResult(err.Error()), nil
		}
		encoded, err := json.Marshal(data)
		if err != nil {
			return errResult(err.Error()), nil
		}
		text = string(encoded)
		contentType = "application/json; charset=utf-8"
	default:
		text = display(body)
	}
	return &httpReplyValue{status: int(status), body: text, contentType: contentType}, nil
}

func stdHTTPFetch(values []any, pos Position) (any, *Diagnostic) {
	rawURL, d := requireStringArg("fetch", values, 0, pos)
	if d != nil {
		return nil, d
	}
	opts, d := requireDictOpts(values, 1, pos, "")
	if d != nil {
		return nil, d
	}
	method := "GET"
	if m, ok := optsString(opts, "method"); ok {
		method = strings.ToUpper(m)
	}
	var bodyReader io.Reader
	if body, ok := optsString(opts, "body"); ok {
		bodyReader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, rawURL, bodyReader)
	if err != nil {
		return errResult(err.Error()), nil
	}
	if headersVal, ok := opts.get("headers"); ok {
		headers, ok := headersVal.(*dictValue)
		if !ok {
			return nil, &Diagnostic{Code: "R047", Message: "`fetch` opts.headers must be a dict.", Pos: pos}
		}
		for _, key := range headers.keys() {
			k := key.(string)
			v, _ := headers.get(k)
			text, ok := v.(string)
			if !ok {
				return nil, &Diagnostic{Code: "R047", Message: "`fetch` header values must be strings.", Pos: pos}
			}
			req.Header.Set(k, text)
		}
	}
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return errResult(err.Error()), nil
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return errResult(err.Error()), nil
	}
	headers := newDict()
	for key, vals := range resp.Header {
		if len(vals) > 0 {
			headers.set(key, vals[0])
		}
	}
	return okResult(structOf("Response", []string{"status", "body", "headers"}, map[string]any{
		"status":  int64(resp.StatusCode),
		"body":    string(data),
		"headers": headers,
	})), nil
}

func stdHTTPListen(values []any, pos Position) (any, *Diagnostic) {
	addr, d := requireStringArg("listen", values, 0, pos)
	if d != nil {
		return nil, d
	}
	handler := values[1]
	opts, d := requireDictOpts(values, 2, pos, "")
	if d != nil {
		return nil, d
	}
	cert, hasCert := optsString(opts, "cert")
	key, hasKey := optsString(opts, "key")
	if hasCert != hasKey {
		return errResult("listen TLS requires both cert and key"), nil
	}

	server := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req := requestFromHTTP(r)
			reply, err := callHTTPHandler(handler, req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if reply.contentType != "" {
				w.Header().Set("Content-Type", reply.contentType)
			}
			w.WriteHeader(reply.status)
			_, _ = io.WriteString(w, reply.body)
		}),
		ReadHeaderTimeout: 10 * time.Second,
	}

	var err error
	if hasCert {
		server.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		err = server.ListenAndServeTLS(cert, key)
	} else {
		err = server.ListenAndServe()
	}
	if err != nil && err != http.ErrServerClosed {
		return errResult(err.Error()), nil
	}
	return okResult(true), nil
}

func requestFromHTTP(r *http.Request) *structValue {
	query := newDict()
	for key, vals := range r.URL.Query() {
		if len(vals) == 0 {
			query.set(key, "")
		} else {
			query.set(key, vals[0])
		}
	}
	headers := newDict()
	for key, vals := range r.Header {
		if len(vals) > 0 {
			headers.set(key, vals[0])
		}
	}
	bodyBytes, _ := io.ReadAll(r.Body)
	_ = r.Body.Close()
	return structOf("Request", []string{"method", "path", "query", "headers", "body"}, map[string]any{
		"method":  r.Method,
		"path":    r.URL.Path,
		"query":   query,
		"headers": headers,
		"body":    string(bodyBytes),
	})
}

func callHTTPHandler(handler any, req *structValue) (*httpReplyValue, error) {
	switch fn := handler.(type) {
	case *builtinFunction:
		value, d := fn.call([]any{req}, Position{})
		if d != nil {
			return nil, errFromDiagnostic(d)
		}
		return coerceReply(value)
	case *function:
		// Handlers need an interpreter; listen installs a closure wrapper at call time.
		return nil, errString("internal: VOL function handler requires listen wrapper")
	default:
		return nil, errString("listen handler must be a function")
	}
}

type simpleError string

func (e simpleError) Error() string { return string(e) }
func errString(msg string) error    { return simpleError(msg) }
func errFromDiagnostic(d *Diagnostic) error {
	return simpleError(d.Message)
}

func coerceReply(value any) (*httpReplyValue, error) {
	if value == nil {
		return &httpReplyValue{status: 200, body: "", contentType: "text/plain; charset=utf-8"}, nil
	}
	if reply, ok := value.(*httpReplyValue); ok {
		return reply, nil
	}
	if result, ok := value.(resultValue); ok {
		if !result.ok {
			msg, _ := result.value.(string)
			return &httpReplyValue{status: 500, body: msg, contentType: "text/plain; charset=utf-8"}, nil
		}
		return coerceReply(result.value)
	}
	return &httpReplyValue{status: 200, body: display(value), contentType: "text/plain; charset=utf-8"}, nil
}

// listenWithInterpreter is used by the Call path when the handler is a VOL function.
func listenWithInterpreter(i *interpreter, addr string, handler any, opts *dictValue, pos Position) (any, *Diagnostic) {
	cert, hasCert := optsString(opts, "cert")
	key, hasKey := optsString(opts, "key")
	if hasCert != hasKey {
		return errResult("listen TLS requires both cert and key"), nil
	}
	server := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req := requestFromHTTP(r)
			reply, errMsg := invokeHandler(i, handler, req)
			if errMsg != "" {
				http.Error(w, errMsg, http.StatusInternalServerError)
				return
			}
			if reply.contentType != "" {
				w.Header().Set("Content-Type", reply.contentType)
			}
			w.WriteHeader(reply.status)
			_, _ = io.WriteString(w, reply.body)
		}),
		ReadHeaderTimeout: 10 * time.Second,
	}
	var err error
	if hasCert {
		server.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		err = server.ListenAndServeTLS(cert, key)
	} else {
		err = server.ListenAndServe()
	}
	if err != nil && err != http.ErrServerClosed {
		return errResult(err.Error()), nil
	}
	return okResult(true), nil
}

func invokeHandler(i *interpreter, handler any, req *structValue) (*httpReplyValue, string) {
	switch fn := handler.(type) {
	case *builtinFunction:
		value, d := fn.call([]any{req}, Position{})
		if d != nil {
			return nil, d.Message
		}
		reply, err := coerceReply(value)
		if err != nil {
			return nil, err.Error()
		}
		return reply, ""
	case *function:
		if len(fn.declaration.Parameters) != 1 {
			return nil, "listen handler must take 1 parameter"
		}
		value, d := i.call(fn, []any{req})
		if d != nil {
			return nil, d.Message
		}
		reply, err := coerceReply(value)
		if err != nil {
			return nil, err.Error()
		}
		return reply, ""
	default:
		return nil, "listen handler must be a function"
	}
}

package lang

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func stdDBModule() *stdModule {
	return &stdModule{exports: map[string]any{
		"open":  stdBuiltin("open", 1, 1, stdDBOpen),
		"exec":  stdBuiltin("exec", 2, -1, stdDBExec),
		"query": stdBuiltin("query", 2, -1, stdDBQuery),
		"close": stdBuiltin("close", 1, 1, stdDBClose),
	}}
}

func stdDBOpen(values []any, pos Position) (any, *Diagnostic) {
	path, d := requireStringArg("open", values, 0, pos)
	if d != nil {
		return nil, d
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return errResult(err.Error()), nil
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return errResult(err.Error()), nil
	}
	return okResult(&nativeHandle{kind: "db", value: db}), nil
}

func dbFromArg(value any, pos Position) (*sql.DB, *Diagnostic) {
	handle, ok := value.(*nativeHandle)
	if !ok || handle.kind != "db" {
		return nil, &Diagnostic{Code: "R047", Message: "DB handle required, got " + typeName(value) + ".", Pos: pos}
	}
	db, ok := handle.value.(*sql.DB)
	if !ok || db == nil {
		return nil, &Diagnostic{Code: "R047", Message: "DB handle is closed.", Pos: pos}
	}
	return db, nil
}

func stdDBArgs(values []any) []any {
	if len(values) <= 2 {
		return nil
	}
	return values[2:]
}

func stdDBExec(values []any, pos Position) (any, *Diagnostic) {
	db, d := dbFromArg(values[0], pos)
	if d != nil {
		return nil, d
	}
	query, d := requireStringArg("exec", values, 1, pos)
	if d != nil {
		return nil, d
	}
	args := toSQLArgs(stdDBArgs(values))
	result, err := db.Exec(query, args...)
	if err != nil {
		return errResult(err.Error()), nil
	}
	affected, _ := result.RowsAffected()
	return okResult(int64(affected)), nil
}

func stdDBQuery(values []any, pos Position) (any, *Diagnostic) {
	db, d := dbFromArg(values[0], pos)
	if d != nil {
		return nil, d
	}
	query, d := requireStringArg("query", values, 1, pos)
	if d != nil {
		return nil, d
	}
	args := toSQLArgs(stdDBArgs(values))
	rows, err := db.Query(query, args...)
	if err != nil {
		return errResult(err.Error()), nil
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return errResult(err.Error()), nil
	}
	out := []any{}
	for rows.Next() {
		raw := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range raw {
			ptrs[i] = &raw[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return errResult(err.Error()), nil
		}
		row := newDict()
		for i, col := range cols {
			row.set(col, fromSQLValue(raw[i]))
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		return errResult(err.Error()), nil
	}
	return okResult(out), nil
}

func stdDBClose(values []any, pos Position) (any, *Diagnostic) {
	handle, ok := values[0].(*nativeHandle)
	if !ok || handle.kind != "db" {
		return nil, &Diagnostic{Code: "R047", Message: "`close` requires a DB handle.", Pos: pos}
	}
	db, ok := handle.value.(*sql.DB)
	if !ok || db == nil {
		return okResult(true), nil
	}
	if err := db.Close(); err != nil {
		return errResult(err.Error()), nil
	}
	handle.value = (*sql.DB)(nil)
	return okResult(true), nil
}

func toSQLArgs(values []any) []any {
	out := make([]any, len(values))
	for i, value := range values {
		switch v := value.(type) {
		case optionValue:
			if !v.present {
				out[i] = nil
			} else {
				out[i] = v.value
			}
		default:
			out[i] = value
		}
	}
	return out
}

func fromSQLValue(value any) any {
	if value == nil {
		return noneOption()
	}
	switch v := value.(type) {
	case int64:
		return v
	case float64:
		return v
	case bool:
		return v
	case []byte:
		return string(v)
	case string:
		return v
	default:
		return fmt.Sprint(v)
	}
}

package lang

import (
	"math"
)

func stdMathModule() *stdModule {
	return &stdModule{exports: map[string]any{
		"abs":   stdBuiltin("abs", 1, 1, stdAbs),
		"min":   stdBuiltin("min", 2, 2, stdMin),
		"max":   stdBuiltin("max", 2, 2, stdMax),
		"clamp": stdBuiltin("clamp", 3, 3, stdClamp),
		"floor": stdBuiltin("floor", 1, 1, stdFloor),
		"ceil":  stdBuiltin("ceil", 1, 1, stdCeil),
		"sqrt":  stdBuiltin("sqrt", 1, 1, stdSqrt),
		"pow":   stdBuiltin("pow", 2, 2, stdPow),
	}}
}

func stdAbs(values []any, pos Position) (any, *Diagnostic) {
	switch v := values[0].(type) {
	case int64:
		if v == math.MinInt64 {
			return errResult("integer overflow in abs"), nil
		}
		if v < 0 {
			return okResult(-v), nil
		}
		return okResult(v), nil
	case float64:
		return okResult(math.Abs(v)), nil
	default:
		return nil, &Diagnostic{Code: "R047", Message: "`abs` requires a number, got " + typeName(values[0]) + ".", Pos: pos}
	}
}

func stdMin(values []any, pos Position) (any, *Diagnostic) {
	a, d := numberArg("min", values, 0, pos)
	if d != nil {
		return nil, d
	}
	b, d := numberArg("min", values, 1, pos)
	if d != nil {
		return nil, d
	}
	_, aInt := values[0].(int64)
	_, bInt := values[1].(int64)
	if aInt && bInt {
		if values[0].(int64) < values[1].(int64) {
			return okResult(values[0]), nil
		}
		return okResult(values[1]), nil
	}
	return okResult(math.Min(a, b)), nil
}

func stdMax(values []any, pos Position) (any, *Diagnostic) {
	a, d := numberArg("max", values, 0, pos)
	if d != nil {
		return nil, d
	}
	b, d := numberArg("max", values, 1, pos)
	if d != nil {
		return nil, d
	}
	_, aInt := values[0].(int64)
	_, bInt := values[1].(int64)
	if aInt && bInt {
		if values[0].(int64) > values[1].(int64) {
			return okResult(values[0]), nil
		}
		return okResult(values[1]), nil
	}
	return okResult(math.Max(a, b)), nil
}

func stdClamp(values []any, pos Position) (any, *Diagnostic) {
	x, d := numberArg("clamp", values, 0, pos)
	if d != nil {
		return nil, d
	}
	lo, d := numberArg("clamp", values, 1, pos)
	if d != nil {
		return nil, d
	}
	hi, d := numberArg("clamp", values, 2, pos)
	if d != nil {
		return nil, d
	}
	if lo > hi {
		return errResult("clamp lower bound greater than upper bound"), nil
	}
	if x < lo {
		x = lo
	}
	if x > hi {
		x = hi
	}
	_, xInt := values[0].(int64)
	_, loInt := values[1].(int64)
	_, hiInt := values[2].(int64)
	if xInt && loInt && hiInt {
		return okResult(int64(x)), nil
	}
	return okResult(x), nil
}

func stdFloor(values []any, pos Position) (any, *Diagnostic) {
	n, d := numberArg("floor", values, 0, pos)
	if d != nil {
		return nil, d
	}
	return okResult(math.Floor(n)), nil
}

func stdCeil(values []any, pos Position) (any, *Diagnostic) {
	n, d := numberArg("ceil", values, 0, pos)
	if d != nil {
		return nil, d
	}
	return okResult(math.Ceil(n)), nil
}

func stdSqrt(values []any, pos Position) (any, *Diagnostic) {
	n, d := numberArg("sqrt", values, 0, pos)
	if d != nil {
		return nil, d
	}
	if n < 0 {
		return errResult("sqrt of negative number"), nil
	}
	return okResult(math.Sqrt(n)), nil
}

func stdPow(values []any, pos Position) (any, *Diagnostic) {
	base, d := numberArg("pow", values, 0, pos)
	if d != nil {
		return nil, d
	}
	exp, d := numberArg("pow", values, 1, pos)
	if d != nil {
		return nil, d
	}
	result := math.Pow(base, exp)
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return errResult("pow overflow or undefined"), nil
	}
	return okResult(result), nil
}

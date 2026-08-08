package lang

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
)

const overflowFix = "Use smaller values, switch to floating-point, or wait for planned wrapping arithmetic / build modes."
const nothingFix = "Return a value from the function, or call it as a statement without assigning or using the result."

type environment struct {
	parent *environment
	values map[string]any
	consts map[string]bool
}

func newEnvironment(parent *environment) *environment {
	return &environment{parent: parent, values: map[string]any{}, consts: map[string]bool{}}
}
func (e *environment) define(name string, value any) { e.values[name] = value }
func (e *environment) defineConst(name string, value any) {
	e.values[name] = value
	e.consts[name] = true
}
func (e *environment) isConst(name string) bool {
	if _, ok := e.values[name]; ok {
		return e.consts[name]
	}
	if e.parent != nil {
		return e.parent.isConst(name)
	}
	return false
}
func (e *environment) get(name string) (any, bool) {
	if value, ok := e.values[name]; ok {
		return value, true
	}
	if e.parent != nil {
		return e.parent.get(name)
	}
	return nil, false
}
func (e *environment) assign(name string, value any) bool {
	if _, ok := e.values[name]; ok {
		e.values[name] = value
		return true
	}
	if e.parent != nil {
		return e.parent.assign(name, value)
	}
	return false
}

type interpreter struct {
	output io.Writer
	input  *bufio.Reader
	env    *environment
	file   string
}

type function struct {
	declaration *FunctionDeclaration
	closure     *environment
}

// optionValue is Option: some(x) when present, none when !present.
type optionValue struct {
	present bool
	value   any
}

// resultValue is Result: ok(x) when ok, err(x) when !ok.
type resultValue struct {
	ok    bool
	value any
}

type structType struct {
	name   string
	fields []string
}

type structValue struct {
	typ    *structType
	fields map[string]any
}

type returnValue struct{ value any }

type ExecuteOptions struct {
	Input io.Reader
	Args  []string
}

type builtinFunction struct {
	name string
	call func([]any, Position) (any, *Diagnostic)
}

func Execute(program *Program, output io.Writer) *Diagnostic {
	return ExecuteWithOptions(program, output, ExecuteOptions{Input: strings.NewReader("")})
}

func ExecuteWithOptions(program *Program, output io.Writer, options ExecuteOptions) *Diagnostic {
	for _, statement := range program.Statements {
		if _, ok := statement.(*ImportStatement); ok {
			return &Diagnostic{
				Code:    "S031",
				Message: "Programs with `import` must be run via the file loader.",
				File:    program.File,
				Fix:     "Use `vol run <file.vol>` so imports can be resolved.",
			}
		}
	}
	if d := Resolve(program); d != nil {
		return d
	}
	_, d := executeModule(program, output, options, nil)
	return d
}

func executeModule(program *Program, output io.Writer, options ExecuteOptions, imported map[string]any) (*environment, *Diagnostic) {
	if options.Input == nil {
		options.Input = strings.NewReader("")
	}
	i := &interpreter{output: output, input: bufio.NewReader(options.Input), env: newEnvironment(nil), file: program.File}
	i.installBuiltins(options.Args)
	for name, value := range imported {
		if _, exists := i.env.values[name]; exists {
			return nil, &Diagnostic{
				Code:    "S034",
				Message: "Imported name `" + name + "` collides with a built-in.",
				File:    program.File,
			}
		}
		i.env.define(name, value)
	}
	for _, statement := range program.Statements {
		switch node := statement.(type) {
		case *FunctionDeclaration:
			i.env.define(node.Name.Lexeme, &function{declaration: node, closure: i.env})
		case *StructDeclaration:
			fields := make([]string, len(node.Fields))
			for idx, field := range node.Fields {
				fields[idx] = field.Lexeme
			}
			i.env.define(node.Name.Lexeme, &structType{name: node.Name.Lexeme, fields: fields})
		}
	}
	for _, statement := range program.Statements {
		if d := i.execute(statement); d != nil {
			return nil, d
		}
	}
	return i.env, nil
}

func (i *interpreter) execute(statement Statement) *Diagnostic {
	switch node := statement.(type) {
	case *BlockStatement:
		return i.executeBlock(node.Body, newEnvironment(i.env))
	case *ExportStatement, *ImportStatement:
		// Module metadata — imports are installed by the loader; exports collected after run.
	case *StructDeclaration:
		// Struct types are installed before execution (module scope).
	case *FunctionDeclaration:
		// Module functions are installed before execution so forward calls are valid.
		// A nested function is installed when execution reaches its declaration.
		if _, exists := i.env.values[node.Name.Lexeme]; !exists {
			i.env.define(node.Name.Lexeme, &function{declaration: node, closure: i.env})
		}
	case *ReturnStatement:
		value, d := i.evaluate(node.Value)
		if d != nil {
			return d
		}
		panic(returnValue{value: value})
	case *Declaration:
		value, d := i.evaluate(node.Value)
		if d != nil {
			return d
		}
		if d := i.requireValue(value, node.Value.Position()); d != nil {
			return d
		}
		if _, exists := i.env.values[node.Name.Lexeme]; exists {
			return i.runtime(node.Name.Pos, "R001", "Variable `"+node.Name.Lexeme+"` is already declared in this scope.")
		}
		if node.Const {
			i.env.defineConst(node.Name.Lexeme, value)
		} else {
			i.env.define(node.Name.Lexeme, value)
		}
	case *MultiDeclaration:
		values := make([]any, len(node.Values))
		for idx, expr := range node.Values {
			value, d := i.evaluate(expr)
			if d != nil {
				return d
			}
			if d := i.requireValue(value, expr.Position()); d != nil {
				return d
			}
			values[idx] = value
		}
		for idx, name := range node.Names {
			if _, exists := i.env.values[name.Lexeme]; exists {
				return i.runtime(name.Pos, "R001", "Variable `"+name.Lexeme+"` is already declared in this scope.")
			}
			if node.Const {
				i.env.defineConst(name.Lexeme, values[idx])
			} else {
				i.env.define(name.Lexeme, values[idx])
			}
		}
	case *MultiAssignment:
		values := make([]any, len(node.Values))
		for idx, expr := range node.Values {
			value, d := i.evaluate(expr)
			if d != nil {
				return d
			}
			if d := i.requireValue(value, expr.Position()); d != nil {
				return d
			}
			values[idx] = value
		}
		for idx, name := range node.Names {
			if i.env.isConst(name.Lexeme) {
				return i.runtimeWithFix(name.Pos, "R030",
					"Cannot assign to const binding `"+name.Lexeme+"`.",
					"Declare a new binding instead, or remove `const` from the declaration.")
			}
			if !i.env.assign(name.Lexeme, values[idx]) {
				return i.runtime(name.Pos, "R002", "Unknown variable `"+name.Lexeme+"`.")
			}
		}
	case *Assignment:
		value, d := i.evaluate(node.Value)
		if d != nil {
			return d
		}
		if d := i.requireValue(value, node.Value.Position()); d != nil {
			return d
		}
		switch target := node.Target.(type) {
		case *Variable:
			if i.env.isConst(target.Name.Lexeme) {
				return &Diagnostic{
					Code:    "R030",
					Message: "Cannot assign to const binding `" + target.Name.Lexeme + "`.",
					File:    i.file,
					Pos:     target.Position(),
					Fix:     "Declare a new binding instead, or remove `const` from the declaration.",
				}
			}
			if !i.env.assign(target.Name.Lexeme, value) {
				return i.runtime(target.Position(), "R002", "Unknown variable `"+target.Name.Lexeme+"`.")
			}
		case *Index:
			collection, d := i.evaluate(target.Collection)
			if d != nil {
				return d
			}
			if d := i.requireValue(collection, target.Collection.Position()); d != nil {
				return d
			}
			array, ok := collection.([]any)
			if !ok {
				return i.runtime(target.Position(), "R003", "Only arrays can be indexed.")
			}
			index, d := i.arrayIndex(target.At, len(array))
			if d != nil {
				return d
			}
			array[index] = value
		case *Property:
			object, d := i.evaluate(target.Object)
			if d != nil {
				return d
			}
			if d := i.requireValue(object, target.Object.Position()); d != nil {
				return d
			}
			instance, ok := object.(*structValue)
			if !ok {
				return i.runtime(target.Position(), "R035", "Only struct fields can be assigned through `.`.")
			}
			if _, exists := instance.fields[target.Name.Lexeme]; !exists {
				return i.runtime(target.Name.Pos, "R036", "Struct `"+instance.typ.name+"` has no field `"+target.Name.Lexeme+"`.")
			}
			instance.fields[target.Name.Lexeme] = value
		}
	case *PrintStatement:
		parts := make([]string, 0, len(node.Values))
		for _, expr := range node.Values {
			value, d := i.evaluate(expr)
			if d != nil {
				return d
			}
			if d := i.requireValue(value, expr.Position()); d != nil {
				return d
			}
			parts = append(parts, display(value))
		}
		fmt.Fprintln(i.output, strings.Join(parts, " "))
	case *ExpressionStatement:
		// Discarding a call result is allowed, including `nothing` from a missing return.
		_, d := i.evaluate(node.Value)
		return d
	case *IfLetStatement:
		value, d := i.evaluate(node.Value)
		if d != nil {
			return d
		}
		if d := i.requireValue(value, node.Value.Position()); d != nil {
			return d
		}
		if node.IsOption() {
			option, ok := value.(optionValue)
			if !ok {
				return i.runtimeWithFix(node.Value.Position(), "R034",
					"Option if-let requires an Option (`some`/`none`), got "+typeName(value)+".",
					"Pass a `some(...)` / `none` value, or use a Boolean `if`.")
			}
			if option.present {
				scope := newEnvironment(i.env)
				scope.define(node.Name.Lexeme, option.value)
				return i.executeBlock(node.Then, scope)
			}
			return i.executeBlock(node.Else, newEnvironment(i.env))
		}
		result, ok := value.(resultValue)
		if !ok {
			return i.runtimeWithFix(node.Value.Position(), "R037",
				"Result if-let requires a Result (`ok`/`err`), got "+typeName(value)+".",
				"Pass an `ok(...)` / `err(...)` value, or use a Boolean `if`.")
		}
		if result.ok {
			scope := newEnvironment(i.env)
			scope.define(node.Name.Lexeme, result.value)
			return i.executeBlock(node.Then, scope)
		}
		scope := newEnvironment(i.env)
		scope.define(node.ErrName.Lexeme, result.value)
		return i.executeBlock(node.ErrBody, scope)
	case *IfStatement:
		condition, d := i.evaluate(node.Condition)
		if d != nil {
			return d
		}
		if d := i.requireValue(condition, node.Condition.Position()); d != nil {
			return d
		}
		truth, ok := condition.(bool)
		if !ok {
			return i.runtime(node.Condition.Position(), "R004", "Condition must be Boolean.")
		}
		if truth {
			return i.executeBlock(node.Then, newEnvironment(i.env))
		}
		for _, clause := range node.ElseIfs {
			condition, d := i.evaluate(clause.Condition)
			if d != nil {
				return d
			}
			if d := i.requireValue(condition, clause.Condition.Position()); d != nil {
				return d
			}
			truth, ok := condition.(bool)
			if !ok {
				return i.runtime(clause.Condition.Position(), "R004", "Condition must be Boolean.")
			}
			if truth {
				return i.executeBlock(clause.Then, newEnvironment(i.env))
			}
		}
		if node.Else != nil {
			return i.executeBlock(node.Else, newEnvironment(i.env))
		}
	case *RepeatStatement:
		count, d := i.evaluate(node.Count)
		if d != nil {
			return d
		}
		if d := i.requireValue(count, node.Count.Position()); d != nil {
			return d
		}
		times, ok := count.(int64)
		if !ok || times < 0 {
			return i.runtime(node.Count.Position(), "R005", "Repeat count must be a non-negative integer.")
		}
		for n := int64(0); n < times; n++ {
			if d := i.executeBlock(node.Body, newEnvironment(i.env)); d != nil {
				return d
			}
		}
	case *WhileStatement:
		for {
			condition, d := i.evaluate(node.Condition)
			if d != nil {
				return d
			}
			if d := i.requireValue(condition, node.Condition.Position()); d != nil {
				return d
			}
			truth, ok := condition.(bool)
			if !ok {
				return i.runtime(node.Condition.Position(), "R004", "Condition must be Boolean.")
			}
			if !truth {
				break
			}
			if d := i.executeBlock(node.Body, newEnvironment(i.env)); d != nil {
				return d
			}
		}
	case *EachStatement:
		collection, d := i.evaluate(node.Collection)
		if d != nil {
			return d
		}
		if d := i.requireValue(collection, node.Collection.Position()); d != nil {
			return d
		}
		array, ok := collection.([]any)
		if !ok {
			return i.runtime(node.Collection.Position(), "R006", "`.each` requires an array.")
		}
		for _, value := range array {
			scope := newEnvironment(i.env)
			scope.define(node.Name.Lexeme, value)
			if d := i.executeBlock(node.Body, scope); d != nil {
				return d
			}
		}
	}
	return nil
}

func (i *interpreter) executeBlock(block *Block, scope *environment) *Diagnostic {
	previous := i.env
	i.env = scope
	defer func() { i.env = previous }()
	for _, statement := range block.Statements {
		if d := i.execute(statement); d != nil {
			return d
		}
	}
	return nil
}

func (i *interpreter) evaluate(expression Expression) (any, *Diagnostic) {
	switch node := expression.(type) {
	case *Literal:
		return node.Value, nil
	case *SomeExpression:
		inner, d := i.evaluate(node.Value)
		if d != nil {
			return nil, d
		}
		if d := i.requireValue(inner, node.Value.Position()); d != nil {
			return nil, d
		}
		return optionValue{present: true, value: inner}, nil
	case *OkExpression:
		inner, d := i.evaluate(node.Value)
		if d != nil {
			return nil, d
		}
		if d := i.requireValue(inner, node.Value.Position()); d != nil {
			return nil, d
		}
		return resultValue{ok: true, value: inner}, nil
	case *ErrExpression:
		inner, d := i.evaluate(node.Value)
		if d != nil {
			return nil, d
		}
		if d := i.requireValue(inner, node.Value.Position()); d != nil {
			return nil, d
		}
		return resultValue{ok: false, value: inner}, nil
	case *StructLiteral:
		typeValue, found := i.env.get(node.Type.Lexeme)
		if !found {
			return nil, i.runtime(node.Type.Pos, "R038", "Unknown struct type `"+node.Type.Lexeme+"`.")
		}
		typ, ok := typeValue.(*structType)
		if !ok {
			return nil, i.runtime(node.Type.Pos, "R038", "`"+node.Type.Lexeme+"` is not a struct type.")
		}
		values := map[string]any{}
		if len(node.Positional) > 0 {
			if len(node.Positional) != len(typ.fields) {
				return nil, i.runtime(node.Open.Pos, "R043",
					fmt.Sprintf("Struct `%s` expects %d positional fields, got %d.", typ.name, len(typ.fields), len(node.Positional)))
			}
			for idx, expr := range node.Positional {
				value, d := i.evaluate(expr)
				if d != nil {
					return nil, d
				}
				if d := i.requireValue(value, expr.Position()); d != nil {
					return nil, d
				}
				values[typ.fields[idx]] = value
			}
			return &structValue{typ: typ, fields: values}, nil
		}
		provided := map[string]Expression{}
		for _, field := range node.Fields {
			provided[field.Name.Lexeme] = field.Value
		}
		for _, name := range typ.fields {
			expr, ok := provided[name]
			if !ok {
				return nil, i.runtime(node.Open.Pos, "R039", "Struct literal for `"+typ.name+"` missing field `"+name+"`.")
			}
			value, d := i.evaluate(expr)
			if d != nil {
				return nil, d
			}
			if d := i.requireValue(value, expr.Position()); d != nil {
				return nil, d
			}
			values[name] = value
			delete(provided, name)
		}
		for name, expr := range provided {
			return nil, i.runtime(expr.Position(), "R040", "Struct `"+typ.name+"` has no field `"+name+"`.")
		}
		return &structValue{typ: typ, fields: values}, nil
	case *FunctionExpression:
		declaration := &FunctionDeclaration{
			Keyword:    node.Keyword,
			Name:       Token{Kind: TokenIdentifier, Lexeme: "fn", Pos: node.Keyword.Pos},
			Parameters: node.Parameters,
			Body:       node.Body,
		}
		return &function{declaration: declaration, closure: i.env}, nil
	case *Variable:
		value, ok := i.env.get(node.Name.Lexeme)
		if !ok {
			return nil, i.runtime(node.Position(), "R002", "Unknown variable `"+node.Name.Lexeme+"`.")
		}
		return value, nil
	case *ArrayLiteral:
		values := make([]any, 0, len(node.Elements))
		for _, element := range node.Elements {
			value, d := i.evaluate(element)
			if d != nil {
				return nil, d
			}
			if d := i.requireValue(value, element.Position()); d != nil {
				return nil, d
			}
			values = append(values, value)
		}
		return values, nil
	case *Index:
		collection, d := i.evaluate(node.Collection)
		if d != nil {
			return nil, d
		}
		if d := i.requireValue(collection, node.Collection.Position()); d != nil {
			return nil, d
		}
		array, ok := collection.([]any)
		if !ok {
			return nil, i.runtime(node.Position(), "R003", "Only arrays can be indexed.")
		}
		index, d := i.arrayIndex(node.At, len(array))
		if d != nil {
			return nil, d
		}
		return array[index], nil
	case *Property:
		object, d := i.evaluate(node.Object)
		if d != nil {
			return nil, d
		}
		if d := i.requireValue(object, node.Object.Position()); d != nil {
			return nil, d
		}
		if instance, ok := object.(*structValue); ok {
			value, exists := instance.fields[node.Name.Lexeme]
			if !exists {
				return nil, i.runtime(node.Name.Pos, "R036", "Struct `"+instance.typ.name+"` has no field `"+node.Name.Lexeme+"`.")
			}
			return value, nil
		}
		if node.Name.Lexeme == "len" {
			if array, ok := object.([]any); ok {
				return int64(len(array)), nil
			}
			if text, ok := object.(string); ok {
				return int64(len([]rune(text))), nil
			}
		}
		if node.Name.Lexeme == "length" {
			return nil, &Diagnostic{
				Code:    "R007",
				Message: "Unknown property `length`.",
				File:    i.file,
				Pos:     node.Name.Pos,
				Fix:     "Use `.len` for array element count or string Unicode scalar count.",
			}
		}
		if node.Name.Lexeme == "each" {
			return nil, i.runtimeWithFix(
				node.Name.Pos,
				"R007",
				"Unknown property `each`.",
				"Use statement form `items.each item { ... }`, not `.each(fn...)`.",
			)
		}
		if node.Name.Lexeme == "byte_len" {
			if text, ok := object.(string); ok {
				return int64(len(text)), nil
			}
			return nil, i.runtime(node.Position(), "R033", "`.byte_len` requires a string.")
		}
		return nil, i.runtime(node.Position(), "R007", "Unknown property `"+node.Name.Lexeme+"`.")
	case *Unary:
		right, d := i.evaluate(node.Right)
		if d != nil {
			return nil, d
		}
		if d := i.requireValue(right, node.Right.Position()); d != nil {
			return nil, d
		}
		if node.Operator.Kind == TokenNot {
			value, ok := right.(bool)
			if !ok {
				return nil, i.runtime(node.Position(), "R008", "`not` requires a Boolean value.")
			}
			return !value, nil
		}
		if value, ok := right.(int64); ok {
			if value == math.MinInt64 {
				return nil, i.integerOverflow(node.Operator.Pos, "-")
			}
			return -value, nil
		}
		if value, ok := right.(float64); ok {
			return -value, nil
		}
		return nil, i.runtime(node.Position(), "R009", "Unary `-` requires a number.")
	case *Binary:
		return i.evaluateBinary(node)
	case *Coalesce:
		left, d := i.evaluate(node.Left)
		if d != nil {
			return nil, d
		}
		if d := i.requireValue(left, node.Left.Position()); d != nil {
			return nil, d
		}
		if _, isResult := left.(resultValue); isResult {
			return nil, i.runtimeWithFix(node.Op.Pos, "R042",
				"`??` does not apply to Result values.",
				"Use `if ok x := res { ... } else err e { ... }` to handle errors explicitly.")
		}
		option, ok := left.(optionValue)
		if !ok {
			return nil, i.runtimeWithFix(node.Op.Pos, "R041",
				"`??` requires an Option on the left, got "+typeName(left)+".",
				"Wrap with `some(...)` / use `none`, or pick another expression.")
		}
		if option.present {
			return option.value, nil
		}
		right, d := i.evaluate(node.Right)
		if d != nil {
			return nil, d
		}
		if d := i.requireValue(right, node.Right.Position()); d != nil {
			return nil, d
		}
		return right, nil
	case *TryPropagate:
		value, d := i.evaluate(node.Value)
		if d != nil {
			return nil, d
		}
		if d := i.requireValue(value, node.Value.Position()); d != nil {
			return nil, d
		}
		result, ok := value.(resultValue)
		if !ok {
			return nil, i.runtimeWithFix(node.Op.Pos, "R044",
				"`?` requires a Result on the left, got "+typeName(value)+".",
				"Use `ok(...)` / `err(...)`, or Option `??` / if-let for Option values.")
		}
		if result.ok {
			return result.value, nil
		}
		panic(returnValue{value: resultValue{ok: false, value: result.value}})
	case *Conditional:
		condition, d := i.evaluate(node.Condition)
		if d != nil {
			return nil, d
		}
		if d := i.requireValue(condition, node.Condition.Position()); d != nil {
			return nil, d
		}
		truth, ok := condition.(bool)
		if !ok {
			return nil, i.runtime(node.Condition.Position(), "R004", "Condition must be Boolean.")
		}
		if truth {
			value, d := i.evaluate(node.Then)
			if d != nil {
				return nil, d
			}
			if d := i.requireValue(value, node.Then.Position()); d != nil {
				return nil, d
			}
			return value, nil
		}
		value, d := i.evaluate(node.Else)
		if d != nil {
			return nil, d
		}
		if d := i.requireValue(value, node.Else.Position()); d != nil {
			return nil, d
		}
		return value, nil
	case *Call:
		if property, ok := node.Callee.(*Property); ok && property.Name.Lexeme == "where" {
			return i.evaluateWhere(property, node.Arguments)
		}
		if property, ok := node.Callee.(*Property); ok && property.Name.Lexeme == "map" {
			return i.evaluateMap(property, node.Arguments)
		}
		if property, ok := node.Callee.(*Property); ok && property.Name.Lexeme == "count" {
			return i.evaluateCount(property, node.Arguments)
		}
		if property, ok := node.Callee.(*Property); ok {
			switch property.Name.Lexeme {
			case "sum", "copy", "deep_copy":
				return i.evaluateArrayOperation(property, node.Arguments)
			}
		}
		callee, d := i.evaluate(node.Callee)
		if d != nil {
			return nil, d
		}
		if d := i.requireValue(callee, node.Callee.Position()); d != nil {
			return nil, d
		}
		if builtin, ok := callee.(*builtinFunction); ok {
			arguments := make([]any, len(node.Arguments))
			for index, argument := range node.Arguments {
				arguments[index], d = i.evaluate(argument)
				if d != nil {
					return nil, d
				}
				if d := i.requireValue(arguments[index], argument.Position()); d != nil {
					return nil, d
				}
			}
			return builtin.call(arguments, node.Position())
		}
		fn, ok := callee.(*function)
		if !ok {
			return nil, i.runtime(node.Position(), "R017", "Only functions can be called.")
		}
		if len(node.Arguments) != len(fn.declaration.Parameters) {
			return nil, i.runtime(node.Position(), "R018", fmt.Sprintf("Function `%s` expects %d arguments, got %d.", fn.declaration.Name.Lexeme, len(fn.declaration.Parameters), len(node.Arguments)))
		}
		arguments := make([]any, len(node.Arguments))
		for index, argument := range node.Arguments {
			arguments[index], d = i.evaluate(argument)
			if d != nil {
				return nil, d
			}
			if d := i.requireValue(arguments[index], argument.Position()); d != nil {
				return nil, d
			}
		}
		return i.call(fn, arguments)
	}
	return nil, i.runtime(expression.Position(), "R999", "Unsupported expression.")
}

func (i *interpreter) evaluateArrayOperation(property *Property, arguments []Expression) (any, *Diagnostic) {
	value, d := i.evaluate(property.Object)
	if d != nil {
		return nil, d
	}
	if d := i.requireValue(value, property.Object.Position()); d != nil {
		return nil, d
	}
	array, ok := value.([]any)
	if !ok {
		switch property.Name.Lexeme {
		case "sum":
			return nil, i.runtime(property.Position(), "R019", "`.sum()` requires an array.")
		case "copy":
			return nil, i.runtime(property.Position(), "R031", "`.copy()` requires an array.")
		default:
			return nil, i.runtime(property.Position(), "R032", "`.deep_copy()` requires an array.")
		}
	}
	switch property.Name.Lexeme {
	case "sum":
		var sum any = int64(0)
		for _, item := range array {
			operator := Token{Kind: TokenPlus, Lexeme: "+", Pos: property.Name.Pos}
			sum, d = i.evaluateBinaryValues(operator, sum, item)
			if d != nil {
				return nil, d
			}
		}
		return sum, nil
	case "copy":
		return shallowCopyArray(array), nil
	default:
		return deepCopyArray(array), nil
	}
}

func (i *interpreter) evaluateMap(property *Property, arguments []Expression) (any, *Diagnostic) {
	value, d := i.evaluate(property.Object)
	if d != nil {
		return nil, d
	}
	if d := i.requireValue(value, property.Object.Position()); d != nil {
		return nil, d
	}
	array, ok := value.([]any)
	if !ok {
		return nil, i.runtime(property.Position(), "R021", "`.map` requires an array.")
	}
	if len(arguments) != 1 {
		return nil, i.runtime(property.Name.Pos, "R020", "`.map` expects exactly 1 argument.")
	}
	mapped := make([]any, 0, len(array))
	for _, item := range array {
		scope := newEnvironment(i.env)
		scope.define("_", item)
		previous := i.env
		i.env = scope
		element, diagnostic := i.evaluate(arguments[0])
		i.env = previous
		if diagnostic != nil {
			return nil, diagnostic
		}
		if diagnostic := i.requireValue(element, arguments[0].Position()); diagnostic != nil {
			return nil, diagnostic
		}
		mapped = append(mapped, element)
	}
	return mapped, nil
}

func (i *interpreter) evaluateCount(property *Property, arguments []Expression) (any, *Diagnostic) {
	value, d := i.evaluate(property.Object)
	if d != nil {
		return nil, d
	}
	if d := i.requireValue(value, property.Object.Position()); d != nil {
		return nil, d
	}
	if len(arguments) == 0 {
		if array, ok := value.([]any); ok {
			return int64(len(array)), nil
		}
		if text, ok := value.(string); ok {
			return int64(len([]rune(text))), nil
		}
		return nil, i.runtime(property.Position(), "R021", "`.count` requires an array or string.")
	}
	if len(arguments) != 1 {
		return nil, i.runtime(property.Name.Pos, "R020", "`.count` expects 0 or 1 arguments.")
	}
	array, ok := value.([]any)
	if !ok {
		return nil, i.runtime(property.Position(), "R021", "`.count` requires an array.")
	}
	var n int64
	for _, item := range array {
		scope := newEnvironment(i.env)
		scope.define("_", item)
		previous := i.env
		i.env = scope
		condition, diagnostic := i.evaluate(arguments[0])
		i.env = previous
		if diagnostic != nil {
			return nil, diagnostic
		}
		if diagnostic := i.requireValue(condition, arguments[0].Position()); diagnostic != nil {
			return nil, diagnostic
		}
		matches, ok := condition.(bool)
		if !ok {
			return nil, i.runtimeWithFix(
				arguments[0].Position(),
				"R022",
				"`.count` condition must be Boolean.",
				"Use a Boolean `_` expression, e.g. `.count(_ > 5)`, not a `fn` value.",
			)
		}
		if matches {
			n++
		}
	}
	return n, nil
}

func (i *interpreter) evaluateWhere(property *Property, arguments []Expression) (any, *Diagnostic) {
	if len(arguments) != 1 {
		return nil, i.runtime(property.Position(), "R020", "`.where` expects one condition.")
	}
	value, d := i.evaluate(property.Object)
	if d != nil {
		return nil, d
	}
	if d := i.requireValue(value, property.Object.Position()); d != nil {
		return nil, d
	}
	array, ok := value.([]any)
	if !ok {
		return nil, i.runtime(property.Position(), "R021", "`.where` requires an array.")
	}
	filtered := make([]any, 0, len(array))
	for _, item := range array {
		scope := newEnvironment(i.env)
		scope.define("_", item)
		previous := i.env
		i.env = scope
		condition, diagnostic := i.evaluate(arguments[0])
		i.env = previous
		if diagnostic != nil {
			return nil, diagnostic
		}
		if diagnostic := i.requireValue(condition, arguments[0].Position()); diagnostic != nil {
			return nil, diagnostic
		}
		matches, ok := condition.(bool)
		if !ok {
			return nil, i.runtimeWithFix(
				arguments[0].Position(),
				"R022",
				"`.where` condition must be Boolean.",
				"Use a Boolean `_` expression, e.g. `.where(_ > 5)`, not a `fn` value.",
			)
		}
		if matches {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

func (i *interpreter) call(fn *function, arguments []any) (value any, diagnostic *Diagnostic) {
	scope := newEnvironment(fn.closure)
	for index, parameter := range fn.declaration.Parameters {
		scope.define(parameter.Lexeme, arguments[index])
	}
	returned := false
	func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				if result, ok := recovered.(returnValue); ok {
					value, returned = result.value, true
					return
				}
				panic(recovered)
			}
		}()
		diagnostic = i.executeBlock(fn.declaration.Body, scope)
	}()
	if returned {
		return value, nil
	}
	return nil, diagnostic
}

func (i *interpreter) evaluateBinary(node *Binary) (any, *Diagnostic) {
	left, d := i.evaluate(node.Left)
	if d != nil {
		return nil, d
	}
	if d := i.requireValue(left, node.Left.Position()); d != nil {
		return nil, d
	}
	if node.Operator.Kind == TokenAnd {
		value, ok := left.(bool)
		if !ok {
			return nil, i.runtime(node.Position(), "R010", "`and` requires Boolean values.")
		}
		if !value {
			return false, nil
		}
	}
	if node.Operator.Kind == TokenOr {
		value, ok := left.(bool)
		if !ok {
			return nil, i.runtime(node.Position(), "R011", "`or` requires Boolean values.")
		}
		if value {
			return true, nil
		}
	}
	right, d := i.evaluate(node.Right)
	if d != nil {
		return nil, d
	}
	if d := i.requireValue(right, node.Right.Position()); d != nil {
		return nil, d
	}
	return i.evaluateBinaryValues(node.Operator, left, right)
}

func (i *interpreter) evaluateBinaryValues(operator Token, left, right any) (any, *Diagnostic) {
	switch operator.Kind {
	case TokenEqualEqual:
		if equal, ok := numbersEqual(left, right); ok {
			return equal, nil
		}
		return reflect.DeepEqual(left, right), nil
	case TokenBangEqual:
		if equal, ok := numbersEqual(left, right); ok {
			return !equal, nil
		}
		return !reflect.DeepEqual(left, right), nil
	case TokenAnd, TokenOr:
		rightBool, ok := right.(bool)
		if !ok {
			return nil, i.runtime(operator.Pos, "R012", "Boolean operator requires Boolean values.")
		}
		return rightBool, nil
	case TokenPlus:
		if a, ok := left.(string); ok {
			if b, ok := right.(string); ok {
				return a + b, nil
			}
			// String-context coercion: string + displayable → concat.
			return a + display(right), nil
		}
	}
	if a, ok := left.(int64); ok {
		if b, ok := right.(int64); ok {
			return i.integerBinary(operator, a, b)
		}
	}
	a, aok := number(left)
	b, bok := number(right)
	if aok && bok {
		return i.floatBinary(operator, a, b)
	}
	return nil, i.runtime(operator.Pos, "R013", fmt.Sprintf("Operator `%s` cannot be applied to %s and %s.", operator.Lexeme, typeName(left), typeName(right)))
}

func (i *interpreter) integerBinary(operator Token, a, b int64) (any, *Diagnostic) {
	switch operator.Kind {
	case TokenPlus:
		sum := a + b
		// Two negatives never sum to zero mathematically; MinInt64+MinInt64 wraps to 0.
		if (a > 0 && b > 0 && sum < 0) || (a < 0 && b < 0 && sum >= 0) {
			return nil, i.integerOverflow(operator.Pos, "+")
		}
		return sum, nil
	case TokenMinus:
		diff := a - b
		if (b > 0 && a < math.MinInt64+b) || (b < 0 && a > math.MaxInt64+b) {
			return nil, i.integerOverflow(operator.Pos, "-")
		}
		return diff, nil
	case TokenStar:
		if a == 0 || b == 0 {
			return int64(0), nil
		}
		if (a == math.MinInt64 && b == -1) || (b == math.MinInt64 && a == -1) {
			return nil, i.integerOverflow(operator.Pos, "*")
		}
		product := a * b
		if product/a != b {
			return nil, i.integerOverflow(operator.Pos, "*")
		}
		return product, nil
	case TokenSlash:
		if b == 0 {
			return nil, i.runtime(operator.Pos, "R014", "Division by zero.")
		}
		if a == math.MinInt64 && b == -1 {
			return nil, i.integerOverflow(operator.Pos, "/")
		}
		return a / b, nil
	case TokenLess:
		return a < b, nil
	case TokenLessEqual:
		return a <= b, nil
	case TokenGreater:
		return a > b, nil
	case TokenGreaterEqual:
		return a >= b, nil
	}
	return nil, i.runtime(operator.Pos, "R013", "Unsupported integer operator.")
}
func (i *interpreter) floatBinary(operator Token, a, b float64) (any, *Diagnostic) {
	switch operator.Kind {
	case TokenPlus:
		return a + b, nil
	case TokenMinus:
		return a - b, nil
	case TokenStar:
		return a * b, nil
	case TokenSlash:
		if b == 0 {
			return nil, i.runtime(operator.Pos, "R014", "Division by zero.")
		}
		return a / b, nil
	case TokenLess:
		return a < b, nil
	case TokenLessEqual:
		return a <= b, nil
	case TokenGreater:
		return a > b, nil
	case TokenGreaterEqual:
		return a >= b, nil
	}
	return nil, i.runtime(operator.Pos, "R013", "Unsupported numeric operator.")
}
func number(value any) (float64, bool) {
	switch n := value.(type) {
	case int64:
		return float64(n), true
	case float64:
		return n, true
	}
	return 0, false
}

func numbersEqual(left, right any) (bool, bool) {
	switch a := left.(type) {
	case int64:
		switch b := right.(type) {
		case int64:
			return a == b, true
		case float64:
			return integerEqualsFloat(a, b), true
		}
	case float64:
		switch b := right.(type) {
		case int64:
			return integerEqualsFloat(b, a), true
		case float64:
			return a == b, true
		}
	}
	return false, false
}

func integerEqualsFloat(integer int64, floating float64) bool {
	const (
		minimumIntegerFloat = -9223372036854775808.0
		maximumIntegerFloat = 9223372036854775808.0
	)
	return floating >= minimumIntegerFloat && floating < maximumIntegerFloat &&
		math.Trunc(floating) == floating && int64(floating) == integer
}

func (i *interpreter) arrayIndex(expr Expression, length int) (int, *Diagnostic) {
	value, d := i.evaluate(expr)
	if d != nil {
		return 0, d
	}
	if d := i.requireValue(value, expr.Position()); d != nil {
		return 0, d
	}
	index, ok := value.(int64)
	if !ok {
		return 0, i.runtime(expr.Position(), "R015", "Array index must be an integer.")
	}
	if index < 0 || index >= int64(length) {
		return 0, i.runtime(expr.Position(), "R016", fmt.Sprintf("Array index %d is outside length %d.", index, length))
	}
	return int(index), nil
}
func (i *interpreter) runtime(pos Position, code, message string) *Diagnostic {
	return &Diagnostic{Code: code, Message: message, File: i.file, Pos: pos}
}

func (i *interpreter) runtimeWithFix(pos Position, code, message, fix string) *Diagnostic {
	return &Diagnostic{Code: code, Message: message, File: i.file, Pos: pos, Fix: fix}
}

func (i *interpreter) integerOverflow(pos Position, op string) *Diagnostic {
	return &Diagnostic{
		Code:    "R028",
		Message: fmt.Sprintf("Integer overflow in `%s`.", op),
		File:    i.file,
		Pos:     pos,
		Fix:     overflowFix,
	}
}

func (i *interpreter) requireValue(value any, pos Position) *Diagnostic {
	if value != nil {
		return nil
	}
	return &Diagnostic{
		Code:    "R029",
		Message: "Expected a value, got `nothing`.",
		File:    i.file,
		Pos:     pos,
		Fix:     nothingFix,
	}
}
func display(value any) string {
	if value == nil {
		return "nothing"
	}
	if option, ok := value.(optionValue); ok {
		if !option.present {
			return "none"
		}
		return "some(" + display(option.value) + ")"
	}
	if result, ok := value.(resultValue); ok {
		if result.ok {
			return "ok(" + display(result.value) + ")"
		}
		return "err(" + display(result.value) + ")"
	}
	if instance, ok := value.(*structValue); ok {
		parts := make([]string, 0, len(instance.typ.fields))
		for _, name := range instance.typ.fields {
			parts = append(parts, name+": "+display(instance.fields[name]))
		}
		return instance.typ.name + " { " + strings.Join(parts, ", ") + " }"
	}
	if typ, ok := value.(*structType); ok {
		return "struct " + typ.name
	}
	if array, ok := value.([]any); ok {
		result := "["
		for index, item := range array {
			if index > 0 {
				result += ", "
			}
			result += display(item)
		}
		return result + "]"
	}
	return fmt.Sprint(value)
}

func typeName(value any) string {
	switch value.(type) {
	case nil:
		return "nothing"
	case optionValue:
		return "option"
	case resultValue:
		return "result"
	case *structValue:
		return "struct"
	case *structType:
		return "struct type"
	case int64:
		return "integer"
	case float64:
		return "float"
	case bool:
		return "Boolean"
	case string:
		return "string"
	case []any:
		return "array"
	case *function, *builtinFunction:
		return "function"
	}
	return "value"
}

func shallowCopyArray(array []any) []any {
	result := make([]any, len(array))
	copy(result, array)
	return result
}

func deepCopyArray(array []any) []any {
	result := make([]any, len(array))
	for idx, item := range array {
		result[idx] = deepCopyValue(item)
	}
	return result
}

func deepCopyValue(value any) any {
	if array, ok := value.([]any); ok {
		return deepCopyArray(array)
	}
	if option, ok := value.(optionValue); ok {
		if !option.present {
			return optionValue{}
		}
		return optionValue{present: true, value: deepCopyValue(option.value)}
	}
	if result, ok := value.(resultValue); ok {
		return resultValue{ok: result.ok, value: deepCopyValue(result.value)}
	}
	if instance, ok := value.(*structValue); ok {
		fields := map[string]any{}
		for name, field := range instance.fields {
			fields[name] = deepCopyValue(field)
		}
		return &structValue{typ: instance.typ, fields: fields}
	}
	return value
}

func (i *interpreter) installBuiltins(arguments []string) {
	args := make([]any, len(arguments))
	for index, argument := range arguments {
		args[index] = argument
	}
	i.env.define("args", args)
	i.env.define("string", &builtinFunction{name: "string", call: func(values []any, pos Position) (any, *Diagnostic) { return display(values[0]), nil }})
	i.env.define("input", &builtinFunction{name: "input", call: func(values []any, pos Position) (any, *Diagnostic) {
		if len(values) == 1 {
			prompt, ok := values[0].(string)
			if !ok {
				return nil, i.runtime(pos, "R023", "`input` prompt must be a string, got "+typeName(values[0])+".")
			}
			fmt.Fprint(i.output, prompt)
		}
		line, err := i.input.ReadString('\n')
		if err != nil && len(line) == 0 {
			if err == io.EOF {
				return "", nil
			}
			return nil, i.runtime(pos, "R024", "Could not read input: "+err.Error())
		}
		return strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r"), nil
	}})
	i.env.define("assert", &builtinFunction{name: "assert", call: func(values []any, pos Position) (any, *Diagnostic) {
		condition, ok := values[0].(bool)
		if !ok {
			return nil, i.runtime(pos, "R025", "`assert` condition must be Boolean, got "+typeName(values[0])+".")
		}
		if condition {
			return nil, nil
		}
		message := "Assertion failed."
		if len(values) == 2 {
			text, ok := values[1].(string)
			if !ok {
				return nil, i.runtime(pos, "R026", "`assert` message must be a string, got "+typeName(values[1])+".")
			}
			message = text
		}
		return nil, i.runtime(pos, "R027", message)
	}})
}

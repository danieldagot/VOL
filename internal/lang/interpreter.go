package lang

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
)

type environment struct {
	parent *environment
	values map[string]any
}

func newEnvironment(parent *environment) *environment {
	return &environment{parent: parent, values: map[string]any{}}
}
func (e *environment) define(name string, value any) { e.values[name] = value }
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
	if d := Resolve(program); d != nil {
		return d
	}
	if options.Input == nil {
		options.Input = strings.NewReader("")
	}
	i := &interpreter{output: output, input: bufio.NewReader(options.Input), env: newEnvironment(nil), file: program.File}
	i.installBuiltins(options.Args)
	for _, statement := range program.Statements {
		if fn, ok := statement.(*FunctionDeclaration); ok {
			i.env.define(fn.Name.Lexeme, &function{declaration: fn, closure: i.env})
		}
	}
	for _, statement := range program.Statements {
		if d := i.execute(statement); d != nil {
			return d
		}
	}
	return nil
}

func (i *interpreter) execute(statement Statement) *Diagnostic {
	switch node := statement.(type) {
	case *BlockStatement:
		return i.executeBlock(node.Body, newEnvironment(i.env))
	case *ExportStatement:
		// Export declarations are module metadata and have no runtime behavior.
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
		if _, exists := i.env.values[node.Name.Lexeme]; exists {
			return i.runtime(node.Name.Pos, "R001", "Variable `"+node.Name.Lexeme+"` is already declared in this scope.")
		}
		i.env.define(node.Name.Lexeme, value)
	case *Assignment:
		value, d := i.evaluate(node.Value)
		if d != nil {
			return d
		}
		switch target := node.Target.(type) {
		case *Variable:
			if !i.env.assign(target.Name.Lexeme, value) {
				return i.runtime(target.Position(), "R002", "Unknown variable `"+target.Name.Lexeme+"`.")
			}
		case *Index:
			collection, d := i.evaluate(target.Collection)
			if d != nil {
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
		}
	case *PrintStatement:
		value, d := i.evaluate(node.Value)
		if d != nil {
			return d
		}
		fmt.Fprintln(i.output, display(value))
	case *ExpressionStatement:
		_, d := i.evaluate(node.Value)
		return d
	case *IfStatement:
		condition, d := i.evaluate(node.Condition)
		if d != nil {
			return d
		}
		truth, ok := condition.(bool)
		if !ok {
			return i.runtime(node.Condition.Position(), "R004", "Condition must be Boolean.")
		}
		if truth {
			return i.executeBlock(node.Then, newEnvironment(i.env))
		}
		if node.Else != nil {
			return i.executeBlock(node.Else, newEnvironment(i.env))
		}
	case *RepeatStatement:
		count, d := i.evaluate(node.Count)
		if d != nil {
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
			values = append(values, value)
		}
		return values, nil
	case *Index:
		collection, d := i.evaluate(node.Collection)
		if d != nil {
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
		if node.Name.Lexeme == "length" {
			if array, ok := object.([]any); ok {
				return int64(len(array)), nil
			}
			if text, ok := object.(string); ok {
				return int64(len([]rune(text))), nil
			}
		}
		if node.Name.Lexeme == "sum" {
			array, ok := object.([]any)
			if !ok {
				return nil, i.runtime(node.Position(), "R019", "`.sum` requires an array.")
			}
			var sum any = int64(0)
			for _, item := range array {
				operator := Token{Kind: TokenPlus, Lexeme: "+", Pos: node.Name.Pos}
				sum, d = i.evaluateBinaryValues(operator, sum, item)
				if d != nil {
					return nil, d
				}
			}
			return sum, nil
		}
		return nil, i.runtime(node.Position(), "R007", "Unknown property `"+node.Name.Lexeme+"`.")
	case *Unary:
		right, d := i.evaluate(node.Right)
		if d != nil {
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
			return -value, nil
		}
		if value, ok := right.(float64); ok {
			return -value, nil
		}
		return nil, i.runtime(node.Position(), "R009", "Unary `-` requires a number.")
	case *Binary:
		return i.evaluateBinary(node)
	case *Call:
		if property, ok := node.Callee.(*Property); ok && property.Name.Lexeme == "where" {
			return i.evaluateWhere(property, node.Arguments)
		}
		callee, d := i.evaluate(node.Callee)
		if d != nil {
			return nil, d
		}
		if builtin, ok := callee.(*builtinFunction); ok {
			arguments := make([]any, len(node.Arguments))
			for index, argument := range node.Arguments {
				arguments[index], d = i.evaluate(argument)
				if d != nil {
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
		}
		return i.call(fn, arguments)
	}
	return nil, i.runtime(expression.Position(), "R999", "Unsupported expression.")
}

func (i *interpreter) evaluateWhere(property *Property, arguments []Expression) (any, *Diagnostic) {
	if len(arguments) != 1 {
		return nil, i.runtime(property.Position(), "R020", "`.where` expects one condition.")
	}
	value, d := i.evaluate(property.Object)
	if d != nil {
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
		matches, ok := condition.(bool)
		if !ok {
			return nil, i.runtime(arguments[0].Position(), "R022", "`.where` condition must be Boolean.")
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
func display(value any) string {
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

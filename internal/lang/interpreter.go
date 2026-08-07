package lang

import (
	"fmt"
	"io"
	"reflect"
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
	env    *environment
	file   string
}

func Execute(program *Program, output io.Writer) *Diagnostic {
	i := &interpreter{output: output, env: newEnvironment(nil), file: program.File}
	for _, statement := range program.Statements {
		if d := i.execute(statement); d != nil {
			return d
		}
	}
	return nil
}

func (i *interpreter) execute(statement Statement) *Diagnostic {
	switch node := statement.(type) {
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
	}
	return nil, i.runtime(expression.Position(), "R999", "Unsupported expression.")
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
	switch node.Operator.Kind {
	case TokenEqualEqual:
		return reflect.DeepEqual(left, right), nil
	case TokenBangEqual:
		return !reflect.DeepEqual(left, right), nil
	case TokenAnd, TokenOr:
		rightBool, ok := right.(bool)
		if !ok {
			return nil, i.runtime(node.Position(), "R012", "Boolean operator requires Boolean values.")
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
			return i.integerBinary(node, a, b)
		}
	}
	a, aok := number(left)
	b, bok := number(right)
	if aok && bok {
		return i.floatBinary(node, a, b)
	}
	return nil, i.runtime(node.Position(), "R013", "Operator `"+node.Operator.Lexeme+"` cannot be applied to these values.")
}

func (i *interpreter) integerBinary(node *Binary, a, b int64) (any, *Diagnostic) {
	switch node.Operator.Kind {
	case TokenPlus:
		return a + b, nil
	case TokenMinus:
		return a - b, nil
	case TokenStar:
		return a * b, nil
	case TokenSlash:
		if b == 0 {
			return nil, i.runtime(node.Position(), "R014", "Division by zero.")
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
	return nil, i.runtime(node.Position(), "R013", "Unsupported integer operator.")
}
func (i *interpreter) floatBinary(node *Binary, a, b float64) (any, *Diagnostic) {
	switch node.Operator.Kind {
	case TokenPlus:
		return a + b, nil
	case TokenMinus:
		return a - b, nil
	case TokenStar:
		return a * b, nil
	case TokenSlash:
		if b == 0 {
			return nil, i.runtime(node.Position(), "R014", "Division by zero.")
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
	return nil, i.runtime(node.Position(), "R013", "Unsupported numeric operator.")
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

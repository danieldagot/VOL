package lang

import "fmt"

type symbol struct {
	kind   string
	arity  int
	const_ bool
}

type resolver struct {
	file   string
	scopes []map[string]symbol
}

var builtins = map[string]symbol{
	"args":   {kind: "value", arity: -1},
	"input":  {kind: "builtin", arity: -1},
	"assert": {kind: "builtin", arity: -1},
	"string": {kind: "builtin", arity: 1},
}

// Resolve validates lexical names and statically-known call arities before execution.
func Resolve(program *Program) *Diagnostic {
	r := &resolver{file: program.File, scopes: []map[string]symbol{{}}}
	for name, value := range builtins {
		r.scopes[0][name] = value
	}
	// Module functions are visible throughout the module, including before their declaration.
	for _, statement := range program.Statements {
		if fn, ok := statement.(*FunctionDeclaration); ok {
			if d := r.declare(fn.Name, symbol{kind: "function", arity: len(fn.Parameters)}); d != nil {
				return d
			}
		}
	}
	// Resolve module execution first so every global is known when function bodies
	// are checked, while ordinary local declarations remain order-sensitive.
	for _, statement := range program.Statements {
		if _, ok := statement.(*FunctionDeclaration); ok {
			continue
		}
		if d := r.statement(statement); d != nil {
			return d
		}
	}
	for _, statement := range program.Statements {
		if _, ok := statement.(*FunctionDeclaration); !ok {
			continue
		}
		if d := r.statement(statement); d != nil {
			return d
		}
	}
	return nil
}

func (r *resolver) statement(statement Statement) *Diagnostic {
	switch node := statement.(type) {
	case *BlockStatement:
		return r.block(node.Body)
	case *ExportStatement:
		return nil
	case *FunctionDeclaration:
		if len(r.scopes) > 1 {
			if d := r.declare(node.Name, symbol{kind: "function", arity: len(node.Parameters)}); d != nil {
				return d
			}
		}
		r.begin()
		for _, parameter := range node.Parameters {
			if d := r.declare(parameter, symbol{kind: "value", arity: -1}); d != nil {
				return d
			}
		}
		d := r.blockContents(node.Body)
		r.end()
		return d
	case *Declaration:
		if d := r.expression(node.Value); d != nil {
			return d
		}
		return r.declare(node.Name, symbol{kind: "value", arity: -1, const_: node.Const})
	case *Assignment:
		if target, ok := node.Target.(*Variable); ok {
			if sym, found := r.lookup(target.Name.Lexeme); found && sym.const_ {
				return r.error(target.Name, "S030", "Cannot assign to const binding `"+target.Name.Lexeme+"`.", "Declare a new binding instead, or remove `const` from the declaration.")
			}
		}
		if d := r.expression(node.Target); d != nil {
			return d
		}
		return r.expression(node.Value)
	case *ReturnStatement:
		return r.expression(node.Value)
	case *PrintStatement:
		return r.expression(node.Value)
	case *ExpressionStatement:
		return r.expression(node.Value)
	case *IfStatement:
		if d := r.expression(node.Condition); d != nil {
			return d
		}
		if d := r.block(node.Then); d != nil {
			return d
		}
		for _, clause := range node.ElseIfs {
			if d := r.expression(clause.Condition); d != nil {
				return d
			}
			if d := r.block(clause.Then); d != nil {
				return d
			}
		}
		if node.Else != nil {
			return r.block(node.Else)
		}
	case *RepeatStatement:
		if d := r.expression(node.Count); d != nil {
			return d
		}
		return r.block(node.Body)
	case *WhileStatement:
		if d := r.expression(node.Condition); d != nil {
			return d
		}
		return r.block(node.Body)
	case *EachStatement:
		if d := r.expression(node.Collection); d != nil {
			return d
		}
		r.begin()
		if d := r.declare(node.Name, symbol{kind: "value", arity: -1}); d != nil {
			r.end()
			return d
		}
		d := r.blockContents(node.Body)
		r.end()
		return d
	}
	return nil
}

func (r *resolver) expression(expression Expression) *Diagnostic {
	switch node := expression.(type) {
	case *Literal:
		return nil
	case *Variable:
		if _, ok := r.lookup(node.Name.Lexeme); !ok {
			return r.error(node.Name, "S002", "Undefined name `"+node.Name.Lexeme+"`.", "Declare the name before using it.")
		}
	case *Unary:
		return r.expression(node.Right)
	case *Binary:
		if d := r.expression(node.Left); d != nil {
			return d
		}
		return r.expression(node.Right)
	case *Conditional:
		if d := r.expression(node.Condition); d != nil {
			return d
		}
		if d := r.expression(node.Then); d != nil {
			return d
		}
		return r.expression(node.Else)
	case *ArrayLiteral:
		for _, item := range node.Elements {
			if d := r.expression(item); d != nil {
				return d
			}
		}
	case *Index:
		if d := r.expression(node.Collection); d != nil {
			return d
		}
		return r.expression(node.At)
	case *Property:
		return r.expression(node.Object)
	case *Call:
		if property, ok := node.Callee.(*Property); ok && property.Name.Lexeme == "where" {
			if d := r.expression(property.Object); d != nil {
				return d
			}
			if len(node.Arguments) != 1 {
				return r.arity(property.Name, ".where", 1, len(node.Arguments))
			}
			r.begin()
			r.scopes[len(r.scopes)-1]["_"] = symbol{kind: "value", arity: -1}
			d := r.expression(node.Arguments[0])
			r.end()
			return d
		}
		if property, ok := node.Callee.(*Property); ok {
			switch property.Name.Lexeme {
			case "sum", "copy", "deep_copy":
				if d := r.expression(property.Object); d != nil {
					return d
				}
				if len(node.Arguments) != 0 {
					return r.arity(property.Name, "."+property.Name.Lexeme, 0, len(node.Arguments))
				}
				return nil
			}
		}
		if d := r.expression(node.Callee); d != nil {
			return d
		}
		if variable, ok := node.Callee.(*Variable); ok {
			sym, _ := r.lookup(variable.Name.Lexeme)
			got := len(node.Arguments)
			if sym.kind == "builtin" {
				switch variable.Name.Lexeme {
				case "input":
					if got > 1 {
						return r.rangeArity(variable.Name, "input", "zero or one", got)
					}
				case "assert":
					if got < 1 || got > 2 {
						return r.rangeArity(variable.Name, "assert", "one or two", got)
					}
				}
			}
			if sym.arity >= 0 && got != sym.arity {
				return r.arity(variable.Name, variable.Name.Lexeme, sym.arity, got)
			}
		}
		for _, argument := range node.Arguments {
			if d := r.expression(argument); d != nil {
				return d
			}
		}
	}
	return nil
}

func (r *resolver) block(block *Block) *Diagnostic {
	r.begin()
	d := r.blockContents(block)
	r.end()
	return d
}
func (r *resolver) blockContents(block *Block) *Diagnostic {
	for _, statement := range block.Statements {
		if d := r.statement(statement); d != nil {
			return d
		}
	}
	return nil
}
func (r *resolver) begin() { r.scopes = append(r.scopes, map[string]symbol{}) }
func (r *resolver) end()   { r.scopes = r.scopes[:len(r.scopes)-1] }
func (r *resolver) declare(name Token, value symbol) *Diagnostic {
	scope := r.scopes[len(r.scopes)-1]
	if _, exists := scope[name.Lexeme]; exists {
		return r.error(name, "S001", "Name `"+name.Lexeme+"` is already declared in this scope.", "Choose a different name or remove the duplicate declaration.")
	}
	scope[name.Lexeme] = value
	return nil
}
func (r *resolver) lookup(name string) (symbol, bool) {
	for index := len(r.scopes) - 1; index >= 0; index-- {
		if value, ok := r.scopes[index][name]; ok {
			return value, true
		}
	}
	return symbol{}, false
}
func (r *resolver) arity(token Token, name string, want, got int) *Diagnostic {
	return r.error(token, "S003", fmt.Sprintf("`%s` expects %d arguments, got %d.", name, want, got), "Pass the required number of arguments.")
}
func (r *resolver) rangeArity(token Token, name, want string, got int) *Diagnostic {
	return r.error(token, "S003", fmt.Sprintf("`%s` expects %s arguments, got %d.", name, want, got), "Pass the required number of arguments.")
}
func (r *resolver) error(token Token, code, message, fix string) *Diagnostic {
	return &Diagnostic{Code: code, Message: message, File: r.file, Pos: token.Pos, Fix: fix}
}

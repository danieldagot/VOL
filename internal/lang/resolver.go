package lang

import "fmt"

type symbol struct {
	kind   string
	arity  int
	const_ bool
}

type resolver struct {
	file          string
	scopes        []map[string]symbol
	functionDepth int
}

var builtins = map[string]symbol{
	"args":   {kind: "value", arity: -1},
	"dict":   {kind: "builtin", arity: -1},
	"input":  {kind: "builtin", arity: -1},
	"assert": {kind: "builtin", arity: -1},
	"string": {kind: "builtin", arity: 1},
}

// Resolve validates lexical names and statically-known call arities before execution.
func Resolve(program *Program) *Diagnostic {
	return resolveModule(program, nil)
}

func resolveModule(program *Program, imported map[string]symbol) *Diagnostic {
	r := &resolver{file: program.File, scopes: []map[string]symbol{{}}}
	for name, value := range builtins {
		r.scopes[0][name] = value
	}
	for name, sym := range imported {
		if _, exists := r.scopes[0][name]; exists {
			return &Diagnostic{
				Code:    "S034",
				Message: "Imported name `" + name + "` collides with a built-in.",
				File:    program.File,
				Fix:     "Rename the export in the imported module.",
			}
		}
		r.scopes[0][name] = sym
	}
	// Module functions and struct types are visible throughout the module.
	for _, statement := range program.Statements {
		switch node := statement.(type) {
		case *FunctionDeclaration:
			if d := r.declare(node.Name, symbol{kind: "function", arity: len(node.Parameters)}); d != nil {
				return d
			}
		case *StructDeclaration:
			if d := r.declare(node.Name, symbol{kind: "struct", arity: -1}); d != nil {
				return d
			}
		}
	}
	// Resolve module execution first so every global is known when function bodies
	// are checked, while ordinary local declarations remain order-sensitive.
	for _, statement := range program.Statements {
		switch statement.(type) {
		case *FunctionDeclaration, *StructDeclaration:
			continue
		}
		if d := r.statement(statement); d != nil {
			return d
		}
	}
	for _, statement := range program.Statements {
		switch statement.(type) {
		case *FunctionDeclaration, *StructDeclaration:
			if d := r.statement(statement); d != nil {
				return d
			}
		}
	}
	return nil
}

func (r *resolver) statement(statement Statement) *Diagnostic {
	switch node := statement.(type) {
	case *BlockStatement:
		return r.block(node.Body)
	case *ExportStatement, *ImportStatement:
		return nil
	case *StructDeclaration:
		if len(r.scopes) > 1 {
			return r.error(node.Name, "S037", "`struct` declarations are only allowed at module scope.", "Move the struct to module scope.")
		}
		return nil
	case *FunctionDeclaration:
		if len(r.scopes) > 1 {
			if d := r.declare(node.Name, symbol{kind: "function", arity: len(node.Parameters)}); d != nil {
				return d
			}
		}
		r.begin()
		r.functionDepth++
		for _, parameter := range node.Parameters {
			if d := r.declare(parameter, symbol{kind: "value", arity: -1}); d != nil {
				r.functionDepth--
				r.end()
				return d
			}
		}
		d := r.blockContents(node.Body)
		r.functionDepth--
		r.end()
		return d
	case *Declaration:
		if d := r.expression(node.Value); d != nil {
			return d
		}
		return r.declare(node.Name, symbol{kind: "value", arity: -1, const_: node.Const})
	case *MultiDeclaration:
		for _, value := range node.Values {
			if d := r.expression(value); d != nil {
				return d
			}
		}
		for _, name := range node.Names {
			if d := r.declare(name, symbol{kind: "value", arity: -1, const_: node.Const}); d != nil {
				return d
			}
		}
		return nil
	case *MultiAssignment:
		for _, name := range node.Names {
			if sym, found := r.lookup(name.Lexeme); found && sym.const_ {
				return r.error(name, "S030", "Cannot assign to const binding `"+name.Lexeme+"`.", "Declare a new binding instead, or remove `const` from the declaration.")
			}
			if _, found := r.lookup(name.Lexeme); !found {
				return r.error(name, "S002", "Undefined name `"+name.Lexeme+"`.", undefinedNameFix(name.Lexeme))
			}
		}
		for _, value := range node.Values {
			if d := r.expression(value); d != nil {
				return d
			}
		}
		return nil
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
		for _, value := range node.Values {
			if d := r.expression(value); d != nil {
				return d
			}
		}
		return nil
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
	case *IfLetStatement:
		if d := r.expression(node.Value); d != nil {
			return d
		}
		r.begin()
		if d := r.declare(node.Name, symbol{kind: "value", arity: -1}); d != nil {
			r.end()
			return d
		}
		if d := r.blockContents(node.Then); d != nil {
			r.end()
			return d
		}
		r.end()
		if node.IsOption() {
			return r.block(node.Else)
		}
		r.begin()
		if d := r.declare(node.ErrName, symbol{kind: "value", arity: -1}); d != nil {
			r.end()
			return d
		}
		d := r.blockContents(node.ErrBody)
		r.end()
		return d
	}
	return nil
}

func (r *resolver) expression(expression Expression) *Diagnostic {
	switch node := expression.(type) {
	case *Literal:
		return nil
	case *SomeExpression:
		return r.expression(node.Value)
	case *OkExpression:
		return r.expression(node.Value)
	case *ErrExpression:
		return r.expression(node.Value)
	case *StructLiteral:
		if node.Module != nil {
			if d := r.expression(node.Module); d != nil {
				return d
			}
			// Qualified `mod.Type { }` — export kind checked at runtime.
		} else if sym, ok := r.lookup(node.Type.Lexeme); !ok || sym.kind != "struct" {
			return r.error(node.Type, "S038", "`"+node.Type.Lexeme+"` is not a struct type.", "Declare `struct "+node.Type.Lexeme+" { ... }` before constructing it.")
		}
		for _, field := range node.Fields {
			if d := r.expression(field.Value); d != nil {
				return d
			}
		}
		for _, value := range node.Positional {
			if d := r.expression(value); d != nil {
				return d
			}
		}
		return nil
	case *DictLiteral:
		for _, entry := range node.Entries {
			if d := r.expression(entry.Value); d != nil {
				return d
			}
		}
		return nil
	case *FunctionExpression:
		r.begin()
		r.functionDepth++
		for _, parameter := range node.Parameters {
			if d := r.declare(parameter, symbol{kind: "value", arity: -1}); d != nil {
				r.functionDepth--
				r.end()
				return d
			}
		}
		d := r.blockContents(node.Body)
		r.functionDepth--
		r.end()
		return d
	case *Variable:
		if _, ok := r.lookup(node.Name.Lexeme); !ok {
			return r.error(node.Name, "S002", "Undefined name `"+node.Name.Lexeme+"`.", undefinedNameFix(node.Name.Lexeme))
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
	case *Coalesce:
		if d := r.expression(node.Left); d != nil {
			return d
		}
		return r.expression(node.Right)
	case *TryPropagate:
		return r.expression(node.Value)
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
		if property, ok := node.Callee.(*Property); ok {
			switch property.Name.Lexeme {
			case "where", "map":
				if d := r.expression(property.Object); d != nil {
					return d
				}
				if len(node.Arguments) != 1 {
					return r.arity(property.Name, "."+property.Name.Lexeme, 1, len(node.Arguments))
				}
				r.begin()
				r.scopes[len(r.scopes)-1]["_"] = symbol{kind: "value", arity: -1}
				d := r.expression(node.Arguments[0])
				r.end()
				return d
			case "count":
				if d := r.expression(property.Object); d != nil {
					return d
				}
				if len(node.Arguments) != 1 {
					return r.error(property.Name, "S003",
						fmt.Sprintf("`.count` expects 1 arguments, got %d.", len(node.Arguments)),
						"Use `.len` for length, or `.count(condition)` to count matches (e.g. `.count(_ > 5)`).")
				}
				r.begin()
				r.scopes[len(r.scopes)-1]["_"] = symbol{kind: "value", arity: -1}
				d := r.expression(node.Arguments[0])
				r.end()
				return d
			}
		}
		if property, ok := node.Callee.(*Property); ok {
			switch property.Name.Lexeme {
			case "sum", "copy", "deep_copy", "keys":
				if d := r.expression(property.Object); d != nil {
					return d
				}
				if len(node.Arguments) != 0 {
					fix := "Pass the required number of arguments."
					if property.Name.Lexeme == "sum" {
						fix = "`.sum()` takes no filter; use `.where(condition).sum()`."
					}
					return r.error(property.Name, "S003",
						fmt.Sprintf("`.%s` expects %d arguments, got %d.", property.Name.Lexeme, 0, len(node.Arguments)),
						fix)
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
	d := &Diagnostic{Code: code, Message: message, File: r.file, Pos: token.Pos, Fix: fix}
	if fix != "" {
		d.Repairs = []Repair{{Description: fix}}
	}
	if code == "S003" {
		d.Operation = "arity"
	}
	return d
}

// undefinedNameFix steers common LLM / Python aliases toward Supported @std names.
func undefinedNameFix(name string) string {
	switch name {
	case "contains":
		return "Import `@std/strings` and use `strings.has` (there is no `contains`)."
	case "stringify":
		return "Import `@std/json` or `@std/yaml` and use `json.dump` / `yaml.dump`."
	case "json", "yaml", "http", "fs", "path", "env", "db", "process", "math", "strings", "time", "url":
		return "Import `@std/" + name + "` then call `" + name + ".…` (imports bind module names)."
	case "trim", "split", "join", "has", "prefix", "suffix", "replace":
		return "Import `@std/strings` and call `strings." + name + "(...)` (imports bind module names)."
	case "parse", "dump":
		return "Import `@std/json` (or `@std/yaml` / `@std/url`) and call `json." + name + "(...)` (or `yaml.` / `url.`)."
	case "abs", "min", "max", "clamp", "floor", "ceil", "sqrt", "pow":
		return "Import `@std/math` and call `math." + name + "(...)`."
	case "fetch", "listen", "reply":
		return "Import `@std/http` and call `http." + name + "(...)`."
	case "read", "write", "list", "exists":
		return "Import `@std/fs` and call `fs." + name + "(...)`."
	case "get", "set":
		return "Import `@std/env` and call `env." + name + "(...)`."
	case "open", "exec", "query", "close":
		return "Import `@std/db` and call `db." + name + "(...)`."
	case "run":
		return "Import `@std/process` and call `process.run(...)`."
	case "now", "sleep", "format":
		return "Import `@std/time` and call `time." + name + "(...)`."
	case "base", "dir", "ext":
		return "Import `@std/path` and call `path." + name + "(...)`."
	default:
		return "Declare the name before using it, or `import` the module that exports it."
	}
}

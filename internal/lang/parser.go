package lang

import "fmt"

type parser struct {
	file          string
	tokens        []Token
	current       int
	functionDepth int
}

func Parse(file, source string) (*Program, *Diagnostic) {
	tokens, diagnostic := lex(file, source)
	if diagnostic != nil {
		return nil, diagnostic
	}
	p := parser{file: file, tokens: tokens}
	return p.program()
}

func (p *parser) program() (*Program, *Diagnostic) {
	program := &Program{File: p.file}
	p.skipNewlines()
	for !p.check(TokenEOF) {
		statement, diagnostic := p.statement()
		if diagnostic != nil {
			return nil, diagnostic
		}
		program.Statements = append(program.Statements, statement)
		if declaration, ok := statement.(*ExportStatement); ok {
			program.Exports = append(program.Exports, declaration.Names...)
		}
		if diagnostic := p.finishStatement(TokenEOF); diagnostic != nil {
			return nil, diagnostic
		}
	}
	declared := map[string]bool{}
	for _, statement := range program.Statements {
		switch declaration := statement.(type) {
		case *FunctionDeclaration:
			declared[declaration.Name.Lexeme] = true
		case *Declaration:
			declared[declaration.Name.Lexeme] = true
		case *StructDeclaration:
			declared[declaration.Name.Lexeme] = true
		}
	}
	exported := map[string]bool{}
	for _, name := range program.Exports {
		if exported[name.Lexeme] {
			return nil, p.error(name, "E117", "Name `"+name.Lexeme+"` is exported more than once.")
		}
		if !declared[name.Lexeme] {
			return nil, p.error(name, "E118", "Cannot export unknown name `"+name.Lexeme+"`.")
		}
		exported[name.Lexeme] = true
	}
	for _, statement := range program.Statements {
		switch declaration := statement.(type) {
		case *FunctionDeclaration:
			declaration.Public = exported[declaration.Name.Lexeme]
		case *Declaration:
			declaration.Public = exported[declaration.Name.Lexeme]
		case *StructDeclaration:
			declaration.Public = exported[declaration.Name.Lexeme]
		}
	}
	return program, nil
}

func (p *parser) statement() (Statement, *Diagnostic) {
	if p.match(TokenLeftBrace) {
		body, diagnostic := p.blockAfterOpen(p.previous())
		if diagnostic != nil {
			return nil, diagnostic
		}
		return &BlockStatement{Body: body}, nil
	}
	if p.match(TokenExport) {
		return p.exportStatement(p.previous())
	}
	if p.match(TokenImport) {
		return p.importStatement(p.previous())
	}
	if p.match(TokenStruct) {
		return p.structDeclaration(p.previous())
	}
	// Named `fn name(...)` is a statement; anonymous `fn(...)` is an expression.
	if p.check(TokenFn) && p.peekNext().Kind == TokenIdentifier {
		p.advance()
		return p.functionDeclaration(p.previous())
	}
	if p.match(TokenReturn) {
		keyword := p.previous()
		if p.functionDepth == 0 {
			return nil, p.error(keyword, "E115", "`return` can only be used inside a function.")
		}
		value, d := p.expression()
		if d != nil {
			return nil, d
		}
		return &ReturnStatement{Keyword: keyword, Value: value}, nil
	}
	if p.match(TokenPrint) {
		keyword := p.previous()
		value, d := p.expression()
		if d != nil {
			return nil, d
		}
		values := []Expression{value}
		for p.match(TokenComma) {
			next, d := p.expression()
			if d != nil {
				return nil, d
			}
			values = append(values, next)
		}
		return &PrintStatement{Keyword: keyword, Values: values}, nil
	}
	if p.match(TokenIf) {
		return p.ifStatement(p.previous())
	}
	if p.match(TokenMatch) {
		return nil, p.errorWithFix(p.previous(), "E153",
			"`match` was removed; unwrap with if-let or `??`.",
			"Use `if some x := opt { ... } else { ... }`, `if ok x := res { ... } else err e { ... }`, or `opt ?? default`.")
	}
	if p.match(TokenRepeat) {
		return p.repeatStatement(p.previous())
	}
	if p.match(TokenWhile) {
		return p.whileStatement(p.previous())
	}

	if p.match(TokenConst) {
		keyword := p.previous()
		if !p.check(TokenIdentifier) {
			return nil, p.error(p.peek(), "E120", "Expected a name after `const`.")
		}
		names, d := p.identifierList()
		if d != nil {
			return nil, d
		}
		if !p.match(TokenColonEqual) {
			return nil, p.error(p.peek(), "E121", "Expected `:=` after const name.")
		}
		values, d := p.expressionList()
		if d != nil {
			return nil, d
		}
		if len(names) == 1 && len(values) == 1 {
			return &Declaration{Keyword: keyword, Name: names[0], Value: values[0], Const: true}, nil
		}
		if len(names) != len(values) {
			return nil, p.error(keyword, "E158", fmt.Sprintf("Multi-declare expects %d values, got %d.", len(names), len(values)))
		}
		return &MultiDeclaration{Keyword: keyword, Names: names, Values: values, Const: true}, nil
	}

	if p.check(TokenIdentifier) && (p.peekNext().Kind == TokenColonEqual || p.peekNext().Kind == TokenComma) {
		names, d := p.identifierList()
		if d != nil {
			return nil, d
		}
		if p.match(TokenColonEqual) {
			values, d := p.expressionList()
			if d != nil {
				return nil, d
			}
			if len(names) == 1 && len(values) == 1 {
				return &Declaration{Name: names[0], Value: values[0]}, nil
			}
			if len(names) != len(values) {
				return nil, p.error(names[0], "E158", fmt.Sprintf("Multi-declare expects %d values, got %d.", len(names), len(values)))
			}
			return &MultiDeclaration{Names: names, Values: values}, nil
		}
		if p.match(TokenEqual) {
			equal := p.previous()
			values, d := p.expressionList()
			if d != nil {
				return nil, d
			}
			if len(names) != len(values) {
				return nil, p.error(equal, "E159", fmt.Sprintf("Multi-assign expects %d values, got %d.", len(names), len(values)))
			}
			return &MultiAssignment{Names: names, Equal: equal, Values: values}, nil
		}
		return nil, p.error(p.peek(), "E121", "Expected `:=` or `=` after names.")
	}

	expr, diagnostic := p.expression()
	if diagnostic != nil {
		return nil, diagnostic
	}
	if p.match(TokenEqual) {
		equal := p.previous()
		value, d := p.expression()
		if d != nil {
			return nil, d
		}
		switch expr.(type) {
		case *Variable, *Index, *Property:
		default:
			return nil, p.error(equal, "E102", "Invalid assignment target.")
		}
		return &Assignment{Target: expr, Equal: equal, Value: value}, nil
	}
	if property, ok := expr.(*Property); ok && property.Name.Lexeme == "each" {
		if !p.check(TokenIdentifier) {
			return nil, p.errorWithFix(
				p.peek(),
				"E103",
				"Expected an item name after `.each`.",
				"Use statement form `items.each item { ... }`, not `.each(fn...)`.",
			)
		}
		name := p.advance()
		body, d := p.block()
		if d != nil {
			return nil, d
		}
		return &EachStatement{Collection: property.Object, Name: name, Body: body}, nil
	}
	return &ExpressionStatement{Value: expr}, nil
}

func (p *parser) exportStatement(keyword Token) (Statement, *Diagnostic) {
	declaration := &ExportStatement{Keyword: keyword}
	for {
		if !p.check(TokenIdentifier) {
			return nil, p.error(p.peek(), "E116", "Expected a name after `export`.")
		}
		declaration.Names = append(declaration.Names, p.advance())
		if !p.match(TokenComma) {
			break
		}
	}
	return declaration, nil
}

func (p *parser) importStatement(keyword Token) (Statement, *Diagnostic) {
	if !p.check(TokenString) {
		return nil, p.error(p.peek(), "E131", "Expected a string path after `import`.")
	}
	return &ImportStatement{Keyword: keyword, Path: p.advance()}, nil
}

func (p *parser) structDeclaration(keyword Token) (Statement, *Diagnostic) {
	if !p.check(TokenIdentifier) {
		return nil, p.error(p.peek(), "E132", "Expected a type name after `struct`.")
	}
	name := p.advance()
	if !p.match(TokenLeftBrace) {
		return nil, p.error(p.peek(), "E133", "Expected `{` after struct name.")
	}
	declaration := &StructDeclaration{Keyword: keyword, Name: name}
	p.skipNewlines()
	seen := map[string]bool{}
	for !p.check(TokenRightBrace) && !p.check(TokenEOF) {
		if !p.check(TokenIdentifier) {
			return nil, p.error(p.peek(), "E134", "Expected a field name in `struct`.")
		}
		field := p.advance()
		if seen[field.Lexeme] {
			return nil, p.error(field, "E135", "Duplicate struct field `"+field.Lexeme+"`.")
		}
		seen[field.Lexeme] = true
		declaration.Fields = append(declaration.Fields, field)
		if d := p.finishStatement(TokenRightBrace); d != nil {
			return nil, d
		}
	}
	if !p.match(TokenRightBrace) {
		return nil, p.error(p.peek(), "E136", "Expected `}` to close `struct`.")
	}
	if len(declaration.Fields) == 0 {
		return nil, p.error(name, "E137", "Struct `"+name.Lexeme+"` must declare at least one field.")
	}
	return declaration, nil
}

func (p *parser) functionDeclaration(keyword Token) (Statement, *Diagnostic) {
	if !p.check(TokenIdentifier) {
		return nil, p.error(p.peek(), "E110", "Expected a function name after `fn`.")
	}
	name := p.advance()
	if !p.match(TokenLeftParen) {
		return nil, p.error(p.peek(), "E111", "Expected `(` after the function name.")
	}
	parameters, d := p.parameterList()
	if d != nil {
		return nil, d
	}
	p.functionDepth++
	body, d := p.functionBody()
	p.functionDepth--
	if d != nil {
		return nil, d
	}
	return &FunctionDeclaration{Keyword: keyword, Name: name, Parameters: parameters, Body: body}, nil
}

func (p *parser) functionExpression(keyword Token) (Expression, *Diagnostic) {
	if !p.match(TokenLeftParen) {
		return nil, p.error(p.peek(), "E111", "Expected `(` after `fn`.")
	}
	parameters, d := p.parameterList()
	if d != nil {
		return nil, d
	}
	p.functionDepth++
	body, d := p.functionBody()
	p.functionDepth--
	if d != nil {
		return nil, d
	}
	return &FunctionExpression{Keyword: keyword, Parameters: parameters, Body: body}, nil
}

// functionBody is `{ ... }` or a single expression (implicit return).
func (p *parser) functionBody() (*Block, *Diagnostic) {
	if p.check(TokenLeftBrace) {
		return p.block()
	}
	value, d := p.expression()
	if d != nil {
		return nil, d
	}
	ret := &ReturnStatement{Keyword: Token{Kind: TokenReturn, Lexeme: "return", Pos: value.Position()}, Value: value}
	return &Block{Open: Token{Kind: TokenLeftBrace, Lexeme: "{", Pos: value.Position()}, Statements: []Statement{ret}}, nil
}

func (p *parser) identifierList() ([]Token, *Diagnostic) {
	if !p.check(TokenIdentifier) {
		return nil, p.error(p.peek(), "E120", "Expected a name.")
	}
	names := []Token{p.advance()}
	for p.match(TokenComma) {
		if !p.check(TokenIdentifier) {
			return nil, p.error(p.peek(), "E120", "Expected a name after `,`.")
		}
		names = append(names, p.advance())
	}
	return names, nil
}

func (p *parser) expressionList() ([]Expression, *Diagnostic) {
	first, d := p.expression()
	if d != nil {
		return nil, d
	}
	values := []Expression{first}
	for p.match(TokenComma) {
		value, d := p.expression()
		if d != nil {
			return nil, d
		}
		values = append(values, value)
	}
	return values, nil
}

func (p *parser) parameterList() ([]Token, *Diagnostic) {
	var parameters []Token
	if !p.check(TokenRightParen) {
		for {
			if !p.check(TokenIdentifier) {
				return nil, p.error(p.peek(), "E112", "Expected a parameter name.")
			}
			parameters = append(parameters, p.advance())
			if !p.match(TokenComma) {
				break
			}
		}
	}
	if !p.match(TokenRightParen) {
		return nil, p.error(p.peek(), "E113", "Expected `)` after function parameters.")
	}
	return parameters, nil
}

func (p *parser) looksLikeIfLet() bool {
	if !p.check(TokenSome) && !p.check(TokenOk) {
		return false
	}
	if p.current+2 >= len(p.tokens) {
		return false
	}
	return p.tokens[p.current+1].Kind == TokenIdentifier && p.tokens[p.current+2].Kind == TokenColonEqual
}

func (p *parser) ifLetStatement(keyword Token) (Statement, *Diagnostic) {
	tag := p.advance() // some or ok
	name := p.advance()
	p.advance() // :=
	value, d := p.expression()
	if d != nil {
		return nil, d
	}
	then, d := p.block()
	if d != nil {
		return nil, d
	}
	checkpoint := p.current
	p.skipNewlines()
	if !p.match(TokenElse) {
		p.current = checkpoint
		if tag.Kind == TokenSome {
			return nil, p.errorWithFix(keyword, "E154",
				"Option if-let requires an `else` branch.",
				"Write `if some x := opt { ... } else { ... }`.")
		}
		return nil, p.errorWithFix(keyword, "E155",
			"Result if-let requires `else err <name> { ... }`.",
			"Write `if ok x := res { ... } else err e { ... }`.")
	}
	if tag.Kind == TokenSome {
		otherwise, d := p.block()
		if d != nil {
			return nil, d
		}
		return &IfLetStatement{Keyword: keyword, Tag: tag, Name: name, Value: value, Then: then, Else: otherwise}, nil
	}
	// Result: else err name { ... }
	if !p.match(TokenErr) {
		return nil, p.errorWithFix(p.peek(), "E156",
			"Expected `err` after `else` in Result if-let.",
			"Write `else err e { ... }` to bind the error payload.")
	}
	errKeyword := p.previous()
	if !p.check(TokenIdentifier) {
		return nil, p.error(p.peek(), "E157", "Expected a binding name after `err`.")
	}
	errName := p.advance()
	errBody, d := p.block()
	if d != nil {
		return nil, d
	}
	return &IfLetStatement{
		Keyword: keyword, Tag: tag, Name: name, Value: value, Then: then,
		ErrKeyword: errKeyword, ErrName: errName, ErrBody: errBody,
	}, nil
}

func (p *parser) ifStatement(keyword Token) (Statement, *Diagnostic) {
	if p.looksLikeIfLet() {
		return p.ifLetStatement(keyword)
	}
	condition, d := p.expression()
	if d != nil {
		return nil, d
	}
	then, d := p.block()
	if d != nil {
		return nil, d
	}
	var elseIfs []ElseIfClause
	var otherwise *Block
	for {
		checkpoint := p.current
		p.skipNewlines()
		if p.match(TokenElif) {
			elifKeyword := p.previous()
			elifCondition, d := p.expression()
			if d != nil {
				return nil, d
			}
			elifThen, d := p.block()
			if d != nil {
				return nil, d
			}
			elseIfs = append(elseIfs, ElseIfClause{Keyword: elifKeyword, Condition: elifCondition, Then: elifThen})
			continue
		}
		if p.match(TokenElse) {
			otherwise, d = p.block()
			if d != nil {
				return nil, d
			}
			break
		}
		p.current = checkpoint
		break
	}
	return &IfStatement{Keyword: keyword, Condition: condition, Then: then, ElseIfs: elseIfs, Else: otherwise}, nil
}

func (p *parser) repeatStatement(keyword Token) (Statement, *Diagnostic) {
	count, d := p.expression()
	if d != nil {
		return nil, d
	}
	body, d := p.block()
	if d != nil {
		return nil, d
	}
	return &RepeatStatement{Keyword: keyword, Count: count, Body: body}, nil
}
func (p *parser) whileStatement(keyword Token) (Statement, *Diagnostic) {
	condition, d := p.expression()
	if d != nil {
		return nil, d
	}
	body, d := p.block()
	if d != nil {
		return nil, d
	}
	return &WhileStatement{Keyword: keyword, Condition: condition, Body: body}, nil
}

func (p *parser) block() (*Block, *Diagnostic) {
	if !p.match(TokenLeftBrace) {
		return nil, p.error(p.peek(), "E104", "Expected `{` to begin a block.")
	}
	return p.blockAfterOpen(p.previous())
}

func (p *parser) blockAfterOpen(open Token) (*Block, *Diagnostic) {
	block := &Block{Open: open}
	p.skipNewlines()
	for !p.check(TokenRightBrace) && !p.check(TokenEOF) {
		statement, d := p.statement()
		if d != nil {
			return nil, d
		}
		block.Statements = append(block.Statements, statement)
		if d := p.finishStatement(TokenRightBrace); d != nil {
			return nil, d
		}
	}
	if !p.match(TokenRightBrace) {
		return nil, p.error(p.peek(), "E105", "Expected `}` to close the block.")
	}
	return block, nil
}

func (p *parser) expression() (Expression, *Diagnostic) { return p.ternary() }

func (p *parser) ternary() (Expression, *Diagnostic) {
	condition, d := p.coalesce()
	if d != nil {
		return nil, d
	}
	if !p.check(TokenQuestion) {
		return condition, nil
	}
	// Postfix Result `?` when `?` is not starting `? :`.
	if p.isPostfixTry() {
		p.advance()
		return &TryPropagate{Value: condition, Op: p.previous()}, nil
	}
	p.advance()
	question := p.previous()
	thenExpr, d := p.ternary()
	if d != nil {
		return nil, d
	}
	if !p.match(TokenColon) {
		return nil, p.error(p.peek(), "E001", "Expected `:` after the true branch of `?`.")
	}
	colon := p.previous()
	elseExpr, d := p.ternary()
	if d != nil {
		return nil, d
	}
	return &Conditional{Condition: condition, Question: question, Then: thenExpr, Colon: colon, Else: elseExpr}, nil
}

func (p *parser) isPostfixTry() bool {
	if !p.check(TokenQuestion) {
		return false
	}
	if p.current+1 >= len(p.tokens) {
		return true
	}
	next := p.tokens[p.current+1].Kind
	switch next {
	case TokenNewline, TokenEOF, TokenRightParen, TokenRightBracket, TokenRightBrace,
		TokenComma, TokenPlus, TokenMinus, TokenStar, TokenSlash,
		TokenEqualEqual, TokenBangEqual, TokenLess, TokenLessEqual, TokenGreater, TokenGreaterEqual,
		TokenAnd, TokenOr, TokenQuestionQuestion, TokenEqual, TokenColonEqual,
		TokenDot:
		return true
	default:
		return false
	}
}

// coalesce parses Option `??` (tighter than `? :`, looser than `or`), right-associative.
func (p *parser) coalesce() (Expression, *Diagnostic) {
	left, d := p.or()
	if d != nil {
		return nil, d
	}
	if !p.match(TokenQuestionQuestion) {
		return left, nil
	}
	op := p.previous()
	p.skipNewlines()
	right, d := p.coalesce()
	if d != nil {
		return nil, d
	}
	return &Coalesce{Left: left, Op: op, Right: right}, nil
}

func (p *parser) or() (Expression, *Diagnostic)  { return p.binary(p.and, TokenOr) }
func (p *parser) and() (Expression, *Diagnostic) { return p.binary(p.not, TokenAnd) }
func (p *parser) not() (Expression, *Diagnostic) {
	if p.match(TokenNot) {
		operator := p.previous()
		right, d := p.not()
		if d != nil {
			return nil, d
		}
		return &Unary{Operator: operator, Right: right}, nil
	}
	return p.equality()
}
func (p *parser) equality() (Expression, *Diagnostic) {
	return p.binary(p.comparison, TokenEqualEqual, TokenBangEqual)
}
func (p *parser) comparison() (Expression, *Diagnostic) {
	return p.binary(p.term, TokenLess, TokenLessEqual, TokenGreater, TokenGreaterEqual)
}
func (p *parser) term() (Expression, *Diagnostic)   { return p.binary(p.factor, TokenPlus, TokenMinus) }
func (p *parser) factor() (Expression, *Diagnostic) { return p.binary(p.unary, TokenStar, TokenSlash) }

func (p *parser) binary(next func() (Expression, *Diagnostic), kinds ...TokenKind) (Expression, *Diagnostic) {
	expr, d := next()
	if d != nil {
		return nil, d
	}
	for p.match(kinds...) {
		operator := p.previous()
		p.skipNewlines() // continue after operator: `1 +\n2`
		right, d := next()
		if d != nil {
			return nil, d
		}
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *parser) unary() (Expression, *Diagnostic) {
	if p.match(TokenMinus) {
		operator := p.previous()
		right, d := p.unary()
		if d != nil {
			return nil, d
		}
		return &Unary{Operator: operator, Right: right}, nil
	}
	return p.postfix()
}

func (p *parser) postfix() (Expression, *Diagnostic) {
	expr, d := p.primary()
	if d != nil {
		return nil, d
	}
	for {
		// Same-line call / index (do not cross newlines — that would glue the
		// next statement's `[` / `(` onto this expression).
		if p.match(TokenLeftParen) {
			call := &Call{Callee: expr, Open: p.previous()}
			p.skipNewlines()
			if !p.check(TokenRightParen) {
				for {
					argument, d := p.expression()
					if d != nil {
						return nil, d
					}
					call.Arguments = append(call.Arguments, argument)
					p.skipNewlines()
					if !p.match(TokenComma) {
						break
					}
					p.skipNewlines()
				}
			}
			if !p.match(TokenRightParen) {
				return nil, p.error(p.peek(), "E114", "Expected `)` after function arguments.")
			}
			expr = call
			continue
		}
		if p.match(TokenLeftBracket) {
			open := p.previous()
			p.skipNewlines()
			at, d := p.expression()
			if d != nil {
				return nil, d
			}
			p.skipNewlines()
			if !p.match(TokenRightBracket) {
				return nil, p.error(p.peek(), "E106", "Expected `]` after array index.")
			}
			expr = &Index{Collection: expr, Open: open, At: at}
			continue
		}
		// Multiline method chains: allow newline before `.` only.
		saved := p.current
		p.skipNewlines()
		if p.match(TokenDot) {
			if !p.check(TokenIdentifier) {
				return nil, p.error(p.peek(), "E107", "Expected a property name after `.`.")
			}
			expr = &Property{Object: expr, Name: p.advance()}
			continue
		}
		p.current = saved
		if variable, ok := expr.(*Variable); ok && variable.Name.Lexeme == "dict" && p.looksLikeDictLiteral() {
			literal, d := p.dictLiteral(variable.Name)
			if d != nil {
				return nil, d
			}
			expr = literal
			continue
		}
		if prop, ok := expr.(*Property); ok && p.looksLikeStructLiteral() {
			literal, d := p.structLiteral(prop.Name)
			if d != nil {
				return nil, d
			}
			literal.(*StructLiteral).Module = prop.Object
			expr = literal
			continue
		}
		if variable, ok := expr.(*Variable); ok && p.looksLikeStructLiteral() {
			literal, d := p.structLiteral(variable.Name)
			if d != nil {
				return nil, d
			}
			expr = literal
			continue
		}
		break
	}
	return expr, nil
}

// looksLikeStructLiteral distinguishes struct literals from block-bodied statements.
func (p *parser) looksLikeStructLiteral() bool {
	if !p.check(TokenLeftBrace) {
		return false
	}
	index := p.current + 1
	for index < len(p.tokens) && p.tokens[index].Kind == TokenNewline {
		index++
	}
	if index >= len(p.tokens) {
		return false
	}
	kind := p.tokens[index].Kind
	// Bare `{ }` is a block, not an empty struct literal.
	if kind == TokenRightBrace {
		return false
	}
	if kind == TokenIdentifier {
		index++
		if index >= len(p.tokens) {
			return false
		}
		next := p.tokens[index].Kind
		return next == TokenColon || next == TokenComma || next == TokenRightBrace
	}
	switch kind {
	case TokenInteger, TokenFloat, TokenString, TokenTrue, TokenFalse, TokenNone,
		TokenSome, TokenOk, TokenErr, TokenFn, TokenLeftParen, TokenLeftBracket, TokenMinus:
		return true
	default:
		return false
	}
}

func (p *parser) structLiteral(typeName Token) (Expression, *Diagnostic) {
	if !p.match(TokenLeftBrace) {
		return nil, p.error(p.peek(), "E144", "Expected `{` after struct type name.")
	}
	literal := &StructLiteral{Type: typeName, Open: p.previous()}
	p.skipNewlines()
	if p.check(TokenRightBrace) {
		p.advance()
		return literal, nil
	}
	// Named if first field is `name:`; otherwise positional.
	named := p.check(TokenIdentifier) && p.current+1 < len(p.tokens) && p.tokens[p.current+1].Kind == TokenColon
	if named {
		seen := map[string]bool{}
		for {
			if !p.check(TokenIdentifier) {
				return nil, p.error(p.peek(), "E145", "Expected a field name in struct literal.")
			}
			name := p.advance()
			if seen[name.Lexeme] {
				return nil, p.error(name, "E146", "Duplicate field `"+name.Lexeme+"` in struct literal.")
			}
			seen[name.Lexeme] = true
			if !p.match(TokenColon) {
				return nil, p.error(p.peek(), "E147", "Expected `:` after field name in struct literal.")
			}
			value, d := p.expression()
			if d != nil {
				return nil, d
			}
			literal.Fields = append(literal.Fields, StructFieldInit{Name: name, Value: value})
			p.skipNewlines()
			if !p.match(TokenComma) {
				break
			}
			p.skipNewlines()
		}
	} else {
		for {
			value, d := p.expression()
			if d != nil {
				return nil, d
			}
			literal.Positional = append(literal.Positional, value)
			p.skipNewlines()
			if !p.match(TokenComma) {
				break
			}
			p.skipNewlines()
		}
	}
	if !p.match(TokenRightBrace) {
		return nil, p.error(p.peek(), "E148", "Expected `}` to close struct literal.")
	}
	return literal, nil
}

func (p *parser) primary() (Expression, *Diagnostic) {
	if p.match(TokenInteger) {
		token := p.previous()
		return &Literal{Token: token, Value: parseInteger(token)}, nil
	}
	if p.match(TokenFloat) {
		token := p.previous()
		return &Literal{Token: token, Value: parseFloat(token)}, nil
	}
	if p.match(TokenString) {
		token := p.previous()
		return &Literal{Token: token, Value: token.Lexeme}, nil
	}
	if p.match(TokenTrue) {
		token := p.previous()
		return &Literal{Token: token, Value: true}, nil
	}
	if p.match(TokenFalse) {
		token := p.previous()
		return &Literal{Token: token, Value: false}, nil
	}
	if p.match(TokenNone) {
		token := p.previous()
		return &Literal{Token: token, Value: optionValue{}}, nil
	}
	if p.match(TokenSome) {
		keyword := p.previous()
		if !p.match(TokenLeftParen) {
			return nil, p.error(p.peek(), "E129", "Expected `(` after `some`.")
		}
		value, d := p.expression()
		if d != nil {
			return nil, d
		}
		if !p.match(TokenRightParen) {
			return nil, p.error(p.peek(), "E130", "Expected `)` after `some` value.")
		}
		return &SomeExpression{Keyword: keyword, Value: value}, nil
	}
	if p.match(TokenOk) {
		keyword := p.previous()
		if !p.match(TokenLeftParen) {
			return nil, p.error(p.peek(), "E149", "Expected `(` after `ok`.")
		}
		value, d := p.expression()
		if d != nil {
			return nil, d
		}
		if !p.match(TokenRightParen) {
			return nil, p.error(p.peek(), "E150", "Expected `)` after `ok` value.")
		}
		return &OkExpression{Keyword: keyword, Value: value}, nil
	}
	if p.match(TokenErr) {
		keyword := p.previous()
		if !p.match(TokenLeftParen) {
			return nil, p.error(p.peek(), "E151", "Expected `(` after `err`.")
		}
		value, d := p.expression()
		if d != nil {
			return nil, d
		}
		if !p.match(TokenRightParen) {
			return nil, p.error(p.peek(), "E152", "Expected `)` after `err` value.")
		}
		return &ErrExpression{Keyword: keyword, Value: value}, nil
	}
	if p.match(TokenFn) {
		return p.functionExpression(p.previous())
	}
	if p.match(TokenIdentifier) {
		return &Variable{Name: p.previous()}, nil
	}
	if p.match(TokenLeftParen) {
		expr, d := p.expression()
		if d != nil {
			return nil, d
		}
		if !p.match(TokenRightParen) {
			return nil, p.error(p.peek(), "E108", "Expected `)` after expression.")
		}
		return expr, nil
	}
	if p.match(TokenLeftBracket) {
		array := &ArrayLiteral{Open: p.previous()}
		p.skipNewlines()
		if !p.check(TokenRightBracket) {
			for {
				value, d := p.expression()
				if d != nil {
					return nil, d
				}
				array.Elements = append(array.Elements, value)
				p.skipNewlines()
				if !p.match(TokenComma) {
					break
				}
				p.skipNewlines()
			}
		}
		if !p.match(TokenRightBracket) {
			return nil, p.error(p.peek(), "E109", "Expected `]` after array elements.")
		}
		return array, nil
	}
	return nil, p.error(p.peek(), "E101", "Expected an expression, found "+p.peek().String()+".")
}

func (p *parser) isDictKeyToken(kind TokenKind) bool {
	switch kind {
	case TokenIdentifier, TokenString,
		TokenOk, TokenErr, TokenSome, TokenNone, TokenTrue, TokenFalse,
		TokenFn, TokenIf, TokenElse, TokenElif, TokenWhile, TokenRepeat,
		TokenReturn, TokenStruct, TokenExport, TokenImport, TokenConst,
		TokenAnd, TokenOr, TokenNot, TokenPrint:
		return true
	default:
		return false
	}
}

// looksLikeDictLiteral is `dict {` with empty body or `key: value` entries.
func (p *parser) looksLikeDictLiteral() bool {
	if !p.check(TokenLeftBrace) {
		return false
	}
	index := p.current + 1
	for index < len(p.tokens) && p.tokens[index].Kind == TokenNewline {
		index++
	}
	if index >= len(p.tokens) {
		return false
	}
	kind := p.tokens[index].Kind
	if kind == TokenRightBrace {
		return true
	}
	if !p.isDictKeyToken(kind) {
		return false
	}
	index++
	for index < len(p.tokens) && p.tokens[index].Kind == TokenNewline {
		index++
	}
	return index < len(p.tokens) && p.tokens[index].Kind == TokenColon
}

func (p *parser) dictLiteral(keyword Token) (Expression, *Diagnostic) {
	if !p.match(TokenLeftBrace) {
		return nil, p.error(p.peek(), "E160", "Expected `{` after `dict`.")
	}
	literal := &DictLiteral{Keyword: keyword, Open: p.previous()}
	p.skipNewlines()
	if p.check(TokenRightBrace) {
		p.advance()
		return literal, nil
	}
	for {
		if !p.isDictKeyToken(p.peek().Kind) {
			return nil, p.error(p.peek(), "E161", "Expected a dict key (identifier or string).")
		}
		key := p.advance()
		p.skipNewlines()
		if !p.match(TokenColon) {
			return nil, p.error(p.peek(), "E162", "Expected `:` after dict key.")
		}
		p.skipNewlines()
		value, d := p.expression()
		if d != nil {
			return nil, d
		}
		literal.Entries = append(literal.Entries, DictEntry{Key: key, Value: value})
		p.skipNewlines()
		if p.match(TokenComma) {
			p.skipNewlines()
			if p.check(TokenRightBrace) {
				break
			}
			continue
		}
		break
	}
	if !p.match(TokenRightBrace) {
		return nil, p.error(p.peek(), "E163", "Expected `}` to close dict literal.")
	}
	return literal, nil
}

func (p *parser) match(kinds ...TokenKind) bool {
	for _, kind := range kinds {
		if p.check(kind) {
			p.advance()
			return true
		}
	}
	return false
}
func (p *parser) check(kind TokenKind) bool { return p.peek().Kind == kind }
func (p *parser) advance() Token {
	if !p.check(TokenEOF) {
		p.current++
	}
	return p.previous()
}
func (p *parser) peek() Token { return p.tokens[p.current] }
func (p *parser) peekNext() Token {
	if p.current+1 >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[p.current+1]
}
func (p *parser) previous() Token { return p.tokens[p.current-1] }
func (p *parser) skipNewlines() {
	for p.match(TokenNewline) {
	}
}
func (p *parser) finishStatement(terminator TokenKind) *Diagnostic {
	if p.check(terminator) || p.check(TokenEOF) {
		return nil
	}
	if p.match(TokenNewline) {
		p.skipNewlines()
		return nil
	}
	// A closing brace makes block-bodied statements self-delimiting.
	if p.previous().Kind == TokenRightBrace {
		return nil
	}
	return p.error(p.peek(), "E119", "Expected a newline after the statement.")
}
func (p *parser) error(token Token, code, message string) *Diagnostic {
	return &Diagnostic{Code: code, Message: message, File: p.file, Pos: token.Pos}
}

func (p *parser) errorWithFix(token Token, code, message, fix string) *Diagnostic {
	return &Diagnostic{Code: code, Message: message, File: p.file, Pos: token.Pos, Fix: fix}
}

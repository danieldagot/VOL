package lang

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
		if function, ok := statement.(*FunctionDeclaration); ok {
			function.Public = exported[function.Name.Lexeme]
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
	if p.match(TokenFn) {
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
		return &PrintStatement{Keyword: keyword, Value: value}, nil
	}
	if p.match(TokenIf) {
		return p.ifStatement(p.previous())
	}
	if p.match(TokenRepeat) {
		return p.repeatStatement(p.previous())
	}
	if p.match(TokenWhile) {
		return p.whileStatement(p.previous())
	}

	if p.check(TokenIdentifier) && p.peekNext().Kind == TokenColonEqual {
		name := p.advance()
		p.advance()
		value, d := p.expression()
		if d != nil {
			return nil, d
		}
		return &Declaration{Name: name, Value: value}, nil
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
		case *Variable, *Index:
		default:
			return nil, p.error(equal, "E102", "Invalid assignment target.")
		}
		return &Assignment{Target: expr, Equal: equal, Value: value}, nil
	}
	if property, ok := expr.(*Property); ok && property.Name.Lexeme == "each" {
		if !p.check(TokenIdentifier) {
			return nil, p.error(p.peek(), "E103", "Expected an item name after `.each`.")
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

func (p *parser) functionDeclaration(keyword Token) (Statement, *Diagnostic) {
	if !p.check(TokenIdentifier) {
		return nil, p.error(p.peek(), "E110", "Expected a function name after `fn`.")
	}
	name := p.advance()
	if !p.match(TokenLeftParen) {
		return nil, p.error(p.peek(), "E111", "Expected `(` after the function name.")
	}
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
	p.functionDepth++
	body, d := p.block()
	p.functionDepth--
	if d != nil {
		return nil, d
	}
	return &FunctionDeclaration{Keyword: keyword, Name: name, Parameters: parameters, Body: body}, nil
}

func (p *parser) ifStatement(keyword Token) (Statement, *Diagnostic) {
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
	condition, d := p.or()
	if d != nil {
		return nil, d
	}
	if !p.match(TokenQuestion) {
		return condition, nil
	}
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

func (p *parser) or() (Expression, *Diagnostic)  { return p.binary(p.and, TokenOr) }
func (p *parser) and() (Expression, *Diagnostic) { return p.binary(p.equality, TokenAnd) }
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
		right, d := next()
		if d != nil {
			return nil, d
		}
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *parser) unary() (Expression, *Diagnostic) {
	if p.match(TokenNot, TokenMinus) {
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
		if p.match(TokenLeftParen) {
			call := &Call{Callee: expr, Open: p.previous()}
			if !p.check(TokenRightParen) {
				for {
					argument, d := p.expression()
					if d != nil {
						return nil, d
					}
					call.Arguments = append(call.Arguments, argument)
					if !p.match(TokenComma) {
						break
					}
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
			at, d := p.expression()
			if d != nil {
				return nil, d
			}
			if !p.match(TokenRightBracket) {
				return nil, p.error(p.peek(), "E106", "Expected `]` after array index.")
			}
			expr = &Index{Collection: expr, Open: open, At: at}
			continue
		}
		if p.match(TokenDot) {
			if !p.check(TokenIdentifier) {
				return nil, p.error(p.peek(), "E107", "Expected a property name after `.`.")
			}
			expr = &Property{Object: expr, Name: p.advance()}
			continue
		}
		break
	}
	return expr, nil
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
		if !p.check(TokenRightBracket) {
			for {
				value, d := p.expression()
				if d != nil {
					return nil, d
				}
				array.Elements = append(array.Elements, value)
				if !p.match(TokenComma) {
					break
				}
			}
		}
		if !p.match(TokenRightBracket) {
			return nil, p.error(p.peek(), "E109", "Expected `]` after array elements.")
		}
		return array, nil
	}
	return nil, p.error(p.peek(), "E101", "Expected an expression, found "+p.peek().String()+".")
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

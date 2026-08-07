package lang

type parser struct {
	file    string
	tokens  []Token
	current int
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
	for !p.check(TokenEOF) {
		statement, diagnostic := p.statement()
		if diagnostic != nil {
			return nil, diagnostic
		}
		program.Statements = append(program.Statements, statement)
	}
	return program, nil
}

func (p *parser) statement() (Statement, *Diagnostic) {
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

func (p *parser) ifStatement(keyword Token) (Statement, *Diagnostic) {
	condition, d := p.expression()
	if d != nil {
		return nil, d
	}
	then, d := p.block()
	if d != nil {
		return nil, d
	}
	var otherwise *Block
	if p.match(TokenElse) {
		otherwise, d = p.block()
		if d != nil {
			return nil, d
		}
	}
	return &IfStatement{Keyword: keyword, Condition: condition, Then: then, Else: otherwise}, nil
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
	block := &Block{Open: p.previous()}
	for !p.check(TokenRightBrace) && !p.check(TokenEOF) {
		statement, d := p.statement()
		if d != nil {
			return nil, d
		}
		block.Statements = append(block.Statements, statement)
	}
	if !p.match(TokenRightBrace) {
		return nil, p.error(p.peek(), "E105", "Expected `}` to close the block.")
	}
	return block, nil
}

func (p *parser) expression() (Expression, *Diagnostic) { return p.or() }
func (p *parser) or() (Expression, *Diagnostic)         { return p.binary(p.and, TokenOr) }
func (p *parser) and() (Expression, *Diagnostic)        { return p.binary(p.equality, TokenAnd) }
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
func (p *parser) error(token Token, code, message string) *Diagnostic {
	return &Diagnostic{Code: code, Message: message, File: p.file, Pos: token.Pos}
}

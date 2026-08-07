package lang

import "testing"

func TestEveryASTNodeReportsItsSourcePosition(t *testing.T) {
	position := Position{Offset: 4, Line: 2, Column: 3}
	token := Token{Kind: TokenIdentifier, Lexeme: "node", Pos: position}
	literal := &Literal{Token: token, Value: int64(1)}
	block := &Block{Open: token}

	statements := []Statement{
		&BlockStatement{Body: block},
		&Declaration{Name: token, Value: literal},
		&ExportStatement{Keyword: token},
		&FunctionDeclaration{Keyword: token, Name: token, Body: block},
		&ReturnStatement{Keyword: token, Value: literal},
		&Assignment{Target: literal, Equal: token, Value: literal},
		&PrintStatement{Keyword: token, Value: literal},
		&ExpressionStatement{Value: literal},
		&IfStatement{Keyword: token, Condition: literal, Then: block},
		&RepeatStatement{Keyword: token, Count: literal, Body: block},
		&WhileStatement{Keyword: token, Condition: literal, Body: block},
		&EachStatement{Collection: literal, Name: token, Body: block},
	}
	for _, statement := range statements {
		statement.statement()
		if got := statement.Position(); got != position {
			t.Errorf("%T position = %#v", statement, got)
		}
	}
	if got := block.Position(); got != position {
		t.Fatalf("block position = %#v", got)
	}

	expressions := []Expression{
		literal,
		&Variable{Name: token},
		&Unary{Operator: token, Right: literal},
		&Binary{Left: literal, Operator: token, Right: literal},
		&ArrayLiteral{Open: token},
		&Index{Collection: literal, Open: token, At: literal},
		&Property{Object: literal, Name: token},
		&Call{Callee: literal, Open: token},
	}
	for _, expression := range expressions {
		expression.expression()
		if got := expression.Position(); got != position {
			t.Errorf("%T position = %#v", expression, got)
		}
	}
}

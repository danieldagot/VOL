package lang

type Program struct {
	File       string
	Statements []Statement
	Exports    []Token
}

type Node interface{ Position() Position }
type Statement interface {
	Node
	statement()
}
type Expression interface {
	Node
	expression()
}

type Block struct {
	Open       Token
	Statements []Statement
}

func (n *Block) Position() Position { return n.Open.Pos }

type BlockStatement struct {
	Body *Block
}

func (*BlockStatement) statement()           {}
func (n *BlockStatement) Position() Position { return n.Body.Position() }

type Declaration struct {
	Keyword Token
	Name    Token
	Value   Expression
	Const   bool
	Public  bool
}

func (*Declaration) statement()           {}
func (n *Declaration) Position() Position { return n.Name.Pos }

// MultiDeclaration is `a, b := x, y` (or const).
type MultiDeclaration struct {
	Keyword Token // optional const
	Names   []Token
	Values  []Expression
	Const   bool
}

func (*MultiDeclaration) statement()           {}
func (n *MultiDeclaration) Position() Position { return n.Names[0].Pos }

// MultiAssignment is `a, b = x, y` (RHS fully evaluated before assigns).
type MultiAssignment struct {
	Names  []Token
	Equal  Token
	Values []Expression
}

func (*MultiAssignment) statement()           {}
func (n *MultiAssignment) Position() Position { return n.Equal.Pos }

type ExportStatement struct {
	Keyword Token
	Names   []Token
}

func (*ExportStatement) statement()           {}
func (n *ExportStatement) Position() Position { return n.Keyword.Pos }

type ImportStatement struct {
	Keyword Token
	Path    Token // string literal
}

func (*ImportStatement) statement()           {}
func (n *ImportStatement) Position() Position { return n.Keyword.Pos }

type FunctionDeclaration struct {
	Keyword    Token
	Name       Token
	Parameters []Token
	Body       *Block
	Public     bool
}

func (*FunctionDeclaration) statement()           {}
func (n *FunctionDeclaration) Position() Position { return n.Keyword.Pos }

// FunctionExpression is an anonymous function value: fn(params) { ... } or fn(params) expr.
type FunctionExpression struct {
	Keyword    Token
	Parameters []Token
	Body       *Block
}

func (*FunctionExpression) expression()          {}
func (n *FunctionExpression) Position() Position { return n.Keyword.Pos }

// SomeExpression wraps a present Option value: some(expr).
type SomeExpression struct {
	Keyword Token
	Value   Expression
}

func (*SomeExpression) expression()          {}
func (n *SomeExpression) Position() Position { return n.Keyword.Pos }

// OkExpression wraps a successful Result value: ok(expr).
type OkExpression struct {
	Keyword Token
	Value   Expression
}

func (*OkExpression) expression()          {}
func (n *OkExpression) Position() Position { return n.Keyword.Pos }

// ErrExpression wraps a failed Result value: err(expr).
type ErrExpression struct {
	Keyword Token
	Value   Expression
}

func (*ErrExpression) expression()          {}
func (n *ErrExpression) Position() Position { return n.Keyword.Pos }

// IfLetStatement unwraps Option or Result.
// Option: if some name := value { then } else { else }
// Result: if ok name := value { then } else err errName { errBody }
type IfLetStatement struct {
	Keyword    Token // if
	Tag        Token // some or ok
	Name       Token
	Value      Expression
	Then       *Block
	Else       *Block // Option none branch
	ErrKeyword Token  // err (Result only)
	ErrName    Token  // Result only
	ErrBody    *Block
}

func (*IfLetStatement) statement()           {}
func (n *IfLetStatement) Position() Position { return n.Keyword.Pos }

func (n *IfLetStatement) IsOption() bool { return n.Tag.Kind == TokenSome }
func (n *IfLetStatement) IsResult() bool { return n.Tag.Kind == TokenOk }

type StructDeclaration struct {
	Keyword Token
	Name    Token
	Fields  []Token
	Public  bool
}

func (*StructDeclaration) statement()           {}
func (n *StructDeclaration) Position() Position { return n.Keyword.Pos }

type StructFieldInit struct {
	Name  Token
	Value Expression
}

// StructLiteral constructs a struct value: Type { field: expr, ... } or Type { expr, ... }.
// Module is set for qualified forms like `cart.Item { ... }` (SF-3.1).
type StructLiteral struct {
	Module     Expression // optional module namespace expression
	Type       Token
	Open       Token
	Fields     []StructFieldInit // named form
	Positional []Expression      // positional form (exclusive with Fields)
}

func (*StructLiteral) expression()          {}
func (n *StructLiteral) Position() Position { return n.Type.Pos }

// DictEntry is one `key: value` pair in a dict literal.
type DictEntry struct {
	Key   Token // identifier or string
	Value Expression
}

// DictLiteral is `dict { key: value, ... }` (SF-3.1).
type DictLiteral struct {
	Keyword Token
	Open    Token
	Entries []DictEntry
}

func (*DictLiteral) expression()          {}
func (n *DictLiteral) Position() Position { return n.Keyword.Pos }

type ReturnStatement struct {
	Keyword Token
	Value   Expression
}

func (*ReturnStatement) statement()           {}
func (n *ReturnStatement) Position() Position { return n.Keyword.Pos }

type Assignment struct {
	Target Expression
	Equal  Token
	Value  Expression
}

func (*Assignment) statement()           {}
func (n *Assignment) Position() Position { return n.Equal.Pos }

type PrintStatement struct {
	Keyword Token
	Values  []Expression
}

func (*PrintStatement) statement()           {}
func (n *PrintStatement) Position() Position { return n.Keyword.Pos }

type ExpressionStatement struct{ Value Expression }

func (*ExpressionStatement) statement()           {}
func (n *ExpressionStatement) Position() Position { return n.Value.Position() }

type ElseIfClause struct {
	Keyword   Token
	Condition Expression
	Then      *Block
}

type IfStatement struct {
	Keyword   Token
	Condition Expression
	Then      *Block
	ElseIfs   []ElseIfClause
	Else      *Block
}

func (*IfStatement) statement()           {}
func (n *IfStatement) Position() Position { return n.Keyword.Pos }

type Conditional struct {
	Condition Expression
	Question  Token
	Then      Expression
	Colon     Token
	Else      Expression
}

func (*Conditional) expression()          {}
func (n *Conditional) Position() Position { return n.Question.Pos }

// Coalesce is Option nullish coalesce: left ?? right.
type Coalesce struct {
	Left  Expression
	Op    Token
	Right Expression
}

func (*Coalesce) expression()          {}
func (n *Coalesce) Position() Position { return n.Op.Pos }

// TryPropagate unwraps Result: expr? → ok value; on err, return err from a
// function or abort the script at module top level (R049).
type TryPropagate struct {
	Value Expression
	Op    Token
}

func (*TryPropagate) expression()          {}
func (n *TryPropagate) Position() Position { return n.Op.Pos }

type RepeatStatement struct {
	Keyword Token
	Count   Expression
	Body    *Block
}

func (*RepeatStatement) statement()           {}
func (n *RepeatStatement) Position() Position { return n.Keyword.Pos }

type WhileStatement struct {
	Keyword   Token
	Condition Expression
	Body      *Block
}

func (*WhileStatement) statement()           {}
func (n *WhileStatement) Position() Position { return n.Keyword.Pos }

type EachStatement struct {
	Collection Expression
	Name       Token
	Body       *Block
}

func (*EachStatement) statement()           {}
func (n *EachStatement) Position() Position { return n.Collection.Position() }

type Literal struct {
	Token Token
	Value any
}

func (*Literal) expression()          {}
func (n *Literal) Position() Position { return n.Token.Pos }

type Variable struct{ Name Token }

func (*Variable) expression()          {}
func (n *Variable) Position() Position { return n.Name.Pos }

type Unary struct {
	Operator Token
	Right    Expression
}

func (*Unary) expression()          {}
func (n *Unary) Position() Position { return n.Operator.Pos }

type Binary struct {
	Left     Expression
	Operator Token
	Right    Expression
}

func (*Binary) expression()          {}
func (n *Binary) Position() Position { return n.Operator.Pos }

type ArrayLiteral struct {
	Open     Token
	Elements []Expression
}

func (*ArrayLiteral) expression()          {}
func (n *ArrayLiteral) Position() Position { return n.Open.Pos }

type Index struct {
	Collection Expression
	Open       Token
	At         Expression
}

func (*Index) expression()          {}
func (n *Index) Position() Position { return n.Open.Pos }

type Property struct {
	Object Expression
	Name   Token
}

func (*Property) expression()          {}
func (n *Property) Position() Position { return n.Name.Pos }

type Call struct {
	Callee    Expression
	Open      Token
	Arguments []Expression
}

func (*Call) expression()          {}
func (n *Call) Position() Position { return n.Open.Pos }

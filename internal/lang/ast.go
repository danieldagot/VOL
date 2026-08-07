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
	Name  Token
	Value Expression
}

func (*Declaration) statement()           {}
func (n *Declaration) Position() Position { return n.Name.Pos }

type ExportStatement struct {
	Keyword Token
	Names   []Token
}

func (*ExportStatement) statement()           {}
func (n *ExportStatement) Position() Position { return n.Keyword.Pos }

type FunctionDeclaration struct {
	Keyword    Token
	Name       Token
	Parameters []Token
	Body       *Block
	Public     bool
}

func (*FunctionDeclaration) statement()           {}
func (n *FunctionDeclaration) Position() Position { return n.Keyword.Pos }

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
	Value   Expression
}

func (*PrintStatement) statement()           {}
func (n *PrintStatement) Position() Position { return n.Keyword.Pos }

type ExpressionStatement struct{ Value Expression }

func (*ExpressionStatement) statement()           {}
func (n *ExpressionStatement) Position() Position { return n.Value.Position() }

type IfStatement struct {
	Keyword   Token
	Condition Expression
	Then      *Block
	Else      *Block
}

func (*IfStatement) statement()           {}
func (n *IfStatement) Position() Position { return n.Keyword.Pos }

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

package lang

import "fmt"

type TokenKind int

const (
	TokenEOF TokenKind = iota
	TokenNewline
	TokenIdentifier
	TokenInteger
	TokenFloat
	TokenString
	TokenTrue
	TokenFalse
	TokenIf
	TokenElif
	TokenElse
	TokenRepeat
	TokenWhile
	TokenPrint
	TokenFn
	TokenReturn
	TokenExport
	TokenImport
	TokenConst
	TokenAnd
	TokenOr
	TokenNot
	TokenSome
	TokenNone
	TokenMatch
	TokenStruct
	TokenOk
	TokenErr
	TokenColonEqual
	TokenColon
	TokenQuestion
	TokenQuestionQuestion
	TokenEqual
	TokenEqualEqual
	TokenBangEqual
	TokenLess
	TokenLessEqual
	TokenGreater
	TokenGreaterEqual
	TokenPlus
	TokenMinus
	TokenStar
	TokenSlash
	TokenLeftBrace
	TokenRightBrace
	TokenLeftBracket
	TokenRightBracket
	TokenLeftParen
	TokenRightParen
	TokenComma
	TokenDot
)

type Position struct {
	Offset int
	Line   int
	Column int
}

type Token struct {
	Kind   TokenKind
	Lexeme string
	Pos    Position
}

func (t Token) String() string {
	if t.Kind == TokenEOF {
		return "end of file"
	}
	return fmt.Sprintf("%q", t.Lexeme)
}

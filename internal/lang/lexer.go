package lang

import (
	"strconv"
	"unicode"
	"unicode/utf8"
)

type lexer struct {
	file   string
	source string
	offset int
	line   int
	column int
}

var keywords = map[string]TokenKind{
	"true": TokenTrue, "false": TokenFalse, "if": TokenIf, "elif": TokenElif, "else": TokenElse,
	"repeat": TokenRepeat, "while": TokenWhile, "print": TokenPrint,
	"fn": TokenFn, "return": TokenReturn, "export": TokenExport, "const": TokenConst,
	"and": TokenAnd, "or": TokenOr, "not": TokenNot,
}

func lex(file, source string) ([]Token, *Diagnostic) {
	l := lexer{file: file, source: source, line: 1, column: 1}
	var tokens []Token
	for {
		l.skipWhitespaceAndComments()
		pos := l.position()
		if l.atEnd() {
			return append(tokens, Token{Kind: TokenEOF, Pos: pos}), nil
		}

		r, _ := l.peek()
		if unicode.IsLetter(r) || r == '_' {
			tokens = append(tokens, l.identifier())
			continue
		}
		if unicode.IsDigit(r) {
			token := l.number()
			if !validNumber(token) {
				return nil, l.errorAt(token.Pos, "E006", "Numeric literal is outside the supported range.")
			}
			tokens = append(tokens, token)
			continue
		}

		switch r {
		case '\n':
			l.advance()
			tokens = append(tokens, Token{Kind: TokenNewline, Lexeme: "\n", Pos: pos})
		case '"':
			token, diagnostic := l.stringToken()
			if diagnostic != nil {
				return nil, diagnostic
			}
			tokens = append(tokens, token)
		case ':':
			l.advance()
			if l.match('=') {
				tokens = append(tokens, Token{Kind: TokenColonEqual, Lexeme: ":=", Pos: pos})
			} else {
				tokens = append(tokens, Token{Kind: TokenColon, Lexeme: ":", Pos: pos})
			}
		case '?':
			l.advance()
			tokens = append(tokens, Token{Kind: TokenQuestion, Lexeme: "?", Pos: pos})
		case '=':
			l.advance()
			kind, text := TokenEqual, "="
			if l.match('=') {
				kind, text = TokenEqualEqual, "=="
			}
			tokens = append(tokens, Token{Kind: kind, Lexeme: text, Pos: pos})
		case '!':
			l.advance()
			if !l.match('=') {
				return nil, l.errorAt(pos, "E002", "Expected `=` after `!`; use `not` for Boolean negation.")
			}
			tokens = append(tokens, Token{Kind: TokenBangEqual, Lexeme: "!=", Pos: pos})
		case '<', '>':
			l.advance()
			kind, text := TokenLess, string(r)
			if r == '>' {
				kind = TokenGreater
			}
			if l.match('=') {
				text += "="
				if r == '<' {
					kind = TokenLessEqual
				} else {
					kind = TokenGreaterEqual
				}
			}
			tokens = append(tokens, Token{Kind: kind, Lexeme: text, Pos: pos})
		default:
			kinds := map[rune]TokenKind{'+': TokenPlus, '-': TokenMinus, '*': TokenStar, '/': TokenSlash, '{': TokenLeftBrace, '}': TokenRightBrace, '[': TokenLeftBracket, ']': TokenRightBracket, '(': TokenLeftParen, ')': TokenRightParen, ',': TokenComma, '.': TokenDot}
			kind, ok := kinds[r]
			if !ok {
				return nil, l.errorAt(pos, "E003", "Unexpected character `"+string(r)+"`.")
			}
			l.advance()
			tokens = append(tokens, Token{Kind: kind, Lexeme: string(r), Pos: pos})
		}
	}
}

func (l *lexer) identifier() Token {
	pos, start := l.position(), l.offset
	for !l.atEnd() {
		r, _ := l.peek()
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			break
		}
		l.advance()
	}
	text := l.source[start:l.offset]
	kind := TokenIdentifier
	if keyword, ok := keywords[text]; ok {
		kind = keyword
	}
	return Token{Kind: kind, Lexeme: text, Pos: pos}
}

func (l *lexer) number() Token {
	pos, start, kind := l.position(), l.offset, TokenInteger
	for !l.atEnd() {
		r, _ := l.peek()
		if !unicode.IsDigit(r) {
			break
		}
		l.advance()
	}
	if r, _ := l.peek(); r == '.' {
		if next, ok := l.peekNext(); ok && unicode.IsDigit(next) {
			kind = TokenFloat
			l.advance()
			for !l.atEnd() {
				r, _ = l.peek()
				if !unicode.IsDigit(r) {
					break
				}
				l.advance()
			}
		}
	}
	return Token{Kind: kind, Lexeme: l.source[start:l.offset], Pos: pos}
}

func (l *lexer) stringToken() (Token, *Diagnostic) {
	pos := l.position()
	l.advance()
	var value []rune
	for !l.atEnd() {
		r, _ := l.peek()
		if r == '"' {
			l.advance()
			return Token{Kind: TokenString, Lexeme: string(value), Pos: pos}, nil
		}
		if r == '\n' {
			return Token{}, l.errorAt(pos, "E004", "String literal is not closed.")
		}
		l.advance()
		if r == '\\' {
			if l.atEnd() {
				break
			}
			escaped, _ := l.peek()
			l.advance()
			switch escaped {
			case 'n':
				value = append(value, '\n')
			case 'r':
				value = append(value, '\r')
			case 't':
				value = append(value, '\t')
			case '\\', '"':
				value = append(value, escaped)
			default:
				return Token{}, l.errorAt(pos, "E005", "Unknown string escape `\\"+string(escaped)+"`.")
			}
		} else {
			value = append(value, r)
		}
	}
	return Token{}, l.errorAt(pos, "E004", "String literal is not closed.")
}

func (l *lexer) skipWhitespaceAndComments() {
	for !l.atEnd() {
		r, _ := l.peek()
		if unicode.IsSpace(r) && r != '\n' {
			l.advance()
			continue
		}
		if r == '/' {
			next, ok := l.peekNext()
			if ok && next == '/' {
				for !l.atEnd() {
					r, _ = l.peek()
					if r == '\n' {
						break
					}
					l.advance()
				}
				continue
			}
		}
		break
	}
}

func (l *lexer) peek() (rune, int) {
	if l.atEnd() {
		return 0, 0
	}
	return utf8.DecodeRuneInString(l.source[l.offset:])
}
func (l *lexer) peekNext() (rune, bool) {
	_, size := l.peek()
	if size == 0 || l.offset+size >= len(l.source) {
		return 0, false
	}
	r, _ := utf8.DecodeRuneInString(l.source[l.offset+size:])
	return r, true
}
func (l *lexer) advance() rune {
	r, size := l.peek()
	l.offset += size
	if r == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	return r
}
func (l *lexer) match(want rune) bool {
	r, _ := l.peek()
	if r != want {
		return false
	}
	l.advance()
	return true
}
func (l *lexer) atEnd() bool { return l.offset >= len(l.source) }
func (l *lexer) position() Position {
	return Position{Offset: l.offset, Line: l.line, Column: l.column}
}
func (l *lexer) errorAt(pos Position, code, message string) *Diagnostic {
	return &Diagnostic{Code: code, Message: message, File: l.file, Pos: pos}
}

func parseInteger(token Token) int64 {
	value, _ := strconv.ParseInt(token.Lexeme, 10, 64)
	return value
}
func parseFloat(token Token) float64 { value, _ := strconv.ParseFloat(token.Lexeme, 64); return value }

func validNumber(token Token) bool {
	if token.Kind == TokenInteger {
		_, err := strconv.ParseInt(token.Lexeme, 10, 64)
		return err == nil
	}
	_, err := strconv.ParseFloat(token.Lexeme, 64)
	return err == nil
}

// Package lexer implements the lexical analyzer (tokenizer) for smog.
package lexer

import (
	"fmt"
	"unicode"
)

// TokenType represents the type of a token
type TokenType int

const (
	// Special tokens
	TokenEOF TokenType = iota
	TokenIllegal

	// Literals
	TokenInteger
	TokenFloat
	TokenString
	TokenSymbol

	// Keywords/Identifiers
	TokenIdentifier
	TokenTrue
	TokenFalse
	TokenNil

	// Delimiters
	TokenPeriod      // .
	TokenPipe        // |
	TokenColon       // :
	TokenAssign      // :=
	TokenCaret       // ^
	TokenLParen      // (
	TokenRParen      // )
	TokenLBracket    // [
	TokenRBracket    // ]
	TokenHash        // #
	TokenHashLParen  // #(

	// Operators (binary messages)
	TokenPlus     // +
	TokenMinus    // -
	TokenStar     // *
	TokenSlash    // /
	TokenPercent  // %
	TokenLess     // <
	TokenGreater  // >
	TokenLessEq   // <=
	TokenGreaterEq // >=
	TokenEqual    // =
	TokenNotEqual // ~=
)

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// String returns a string representation of the token type
func (tt TokenType) String() string {
	switch tt {
	case TokenEOF:
		return "EOF"
	case TokenIllegal:
		return "ILLEGAL"
	case TokenInteger:
		return "INTEGER"
	case TokenFloat:
		return "FLOAT"
	case TokenString:
		return "STRING"
	case TokenSymbol:
		return "SYMBOL"
	case TokenIdentifier:
		return "IDENTIFIER"
	case TokenTrue:
		return "TRUE"
	case TokenFalse:
		return "FALSE"
	case TokenNil:
		return "NIL"
	case TokenPeriod:
		return "PERIOD"
	case TokenPipe:
		return "PIPE"
	case TokenColon:
		return "COLON"
	case TokenAssign:
		return "ASSIGN"
	case TokenCaret:
		return "CARET"
	case TokenLParen:
		return "LPAREN"
	case TokenRParen:
		return "RPAREN"
	case TokenLBracket:
		return "LBRACKET"
	case TokenRBracket:
		return "RBRACKET"
	case TokenHash:
		return "HASH"
	case TokenHashLParen:
		return "HASH_LPAREN"
	case TokenPlus:
		return "PLUS"
	case TokenMinus:
		return "MINUS"
	case TokenStar:
		return "STAR"
	case TokenSlash:
		return "SLASH"
	case TokenPercent:
		return "PERCENT"
	case TokenLess:
		return "LESS"
	case TokenGreater:
		return "GREATER"
	case TokenLessEq:
		return "LESS_EQ"
	case TokenGreaterEq:
		return "GREATER_EQ"
	case TokenEqual:
		return "EQUAL"
	case TokenNotEqual:
		return "NOT_EQUAL"
	default:
		return "UNKNOWN"
	}
}

// Lexer represents the lexical analyzer
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int
	column       int
}

// New creates a new lexer for the given input
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// readChar reads the next character
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++
}

// peekChar returns the next character without advancing
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case 0:
		tok.Type = TokenEOF
		tok.Literal = ""
	case '"':
		l.skipComment()
		return l.NextToken()
	case '\'':
		tok.Type = TokenString
		tok.Literal = l.readString()
	case '#':
		if l.peekChar() == '(' {
			l.readChar()
			tok.Type = TokenHashLParen
			tok.Literal = "#("
		} else {
			tok.Type = TokenHash
			tok.Literal = "#"
		}
		l.readChar()
	case '.':
		tok.Type = TokenPeriod
		tok.Literal = "."
		l.readChar()
	case '|':
		tok.Type = TokenPipe
		tok.Literal = "|"
		l.readChar()
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = TokenAssign
			tok.Literal = string(ch) + string(l.ch)
			l.readChar()
		} else {
			tok.Type = TokenColon
			tok.Literal = ":"
			l.readChar()
		}
	case '^':
		tok.Type = TokenCaret
		tok.Literal = "^"
		l.readChar()
	case '(':
		tok.Type = TokenLParen
		tok.Literal = "("
		l.readChar()
	case ')':
		tok.Type = TokenRParen
		tok.Literal = ")"
		l.readChar()
	case '[':
		tok.Type = TokenLBracket
		tok.Literal = "["
		l.readChar()
	case ']':
		tok.Type = TokenRBracket
		tok.Literal = "]"
		l.readChar()
	case '+':
		tok.Type = TokenPlus
		tok.Literal = "+"
		l.readChar()
	case '-':
		// Could be negative number or minus operator
		if unicode.IsDigit(rune(l.peekChar())) {
			l.readChar() // consume the minus
			tok.Type, tok.Literal = l.readNumber()
			tok.Literal = "-" + tok.Literal
			return tok
		} else {
			tok.Type = TokenMinus
			tok.Literal = "-"
			l.readChar()
		}
	case '*':
		tok.Type = TokenStar
		tok.Literal = "*"
		l.readChar()
	case '/':
		tok.Type = TokenSlash
		tok.Literal = "/"
		l.readChar()
	case '%':
		tok.Type = TokenPercent
		tok.Literal = "%"
		l.readChar()
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = TokenLessEq
			tok.Literal = string(ch) + string(l.ch)
			l.readChar()
		} else {
			tok.Type = TokenLess
			tok.Literal = "<"
			l.readChar()
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = TokenGreaterEq
			tok.Literal = string(ch) + string(l.ch)
			l.readChar()
		} else {
			tok.Type = TokenGreater
			tok.Literal = ">"
			l.readChar()
		}
	case '=':
		tok.Type = TokenEqual
		tok.Literal = "="
		l.readChar()
	case '~':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = TokenNotEqual
			tok.Literal = string(ch) + string(l.ch)
			l.readChar()
		} else {
			tok.Type = TokenIllegal
			tok.Literal = string(l.ch)
			l.readChar()
		}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			return tok
		} else if unicode.IsDigit(rune(l.ch)) {
			tok.Type, tok.Literal = l.readNumber()
			return tok
		} else {
			tok.Type = TokenIllegal
			tok.Literal = string(l.ch)
			l.readChar()
		}
	}

	return tok
}

// skipWhitespace skips whitespace characters
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}
}

// skipComment skips a comment (enclosed in double quotes)
func (l *Lexer) skipComment() {
	l.readChar() // skip opening quote
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}
	l.readChar() // skip closing quote
}

// readString reads a string literal
func (l *Lexer) readString() string {
	l.readChar() // skip opening quote
	position := l.position
	for l.ch != '\'' && l.ch != 0 {
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}
	str := l.input[position:l.position]
	l.readChar() // skip closing quote
	return str
}

// readIdentifier reads an identifier or keyword
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || unicode.IsDigit(rune(l.ch)) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads a number (integer or float)
func (l *Lexer) readNumber() (TokenType, string) {
	position := l.position
	hasDecimal := false

	for unicode.IsDigit(rune(l.ch)) || l.ch == '.' {
		if l.ch == '.' {
			// Check if this is a decimal point or statement terminator
			if hasDecimal || !unicode.IsDigit(rune(l.peekChar())) {
				break
			}
			hasDecimal = true
		}
		l.readChar()
	}

	literal := l.input[position:l.position]
	if hasDecimal {
		return TokenFloat, literal
	}
	return TokenInteger, literal
}

// isLetter checks if a character is a letter
func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

// lookupIdent checks if an identifier is a keyword
func lookupIdent(ident string) TokenType {
	switch ident {
	case "true":
		return TokenTrue
	case "false":
		return TokenFalse
	case "nil":
		return TokenNil
	default:
		return TokenIdentifier
	}
}

// Tokenize returns all tokens from the input
func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			break
		}
		if tok.Type == TokenIllegal {
			return tokens, fmt.Errorf("illegal token '%s' at line %d, column %d", tok.Literal, tok.Line, tok.Column)
		}
	}
	return tokens, nil
}

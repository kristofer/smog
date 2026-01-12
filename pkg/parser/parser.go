// Package parser implements the smog language parser.
// It converts source code text into an Abstract Syntax Tree (AST).
package parser

import (
	"fmt"
	"strconv"

	"github.com/kristofer/smog/pkg/ast"
	"github.com/kristofer/smog/pkg/lexer"
)

// Parser represents the smog parser
type Parser struct {
	l      *lexer.Lexer
	curTok lexer.Token
	peekTok lexer.Token
	errors []string
}

// New creates a new parser for the given source code
func New(input string) *Parser {
	p := &Parser{
		l:      lexer.New(input),
		errors: []string{},
	}

	// Read two tokens, so curTok and peekTok are both set
	p.nextToken()
	p.nextToken()

	return p
}

// nextToken advances to the next token
func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.l.NextToken()
}

// Parse parses the source code and returns an AST
func (p *Parser) Parse() (*ast.Program, error) {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curTok.Type != lexer.TokenEOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	if len(p.errors) > 0 {
		return program, fmt.Errorf("parser errors: %v", p.errors)
	}

	return program, nil
}

// parseStatement parses a statement
func (p *Parser) parseStatement() ast.Statement {
	// For now, we only parse expression statements
	expr := p.parseExpression()
	if expr == nil {
		return nil
	}

	stmt := &ast.ExpressionStatement{Expression: expr}

	// Skip optional period at end of statement
	if p.peekTok.Type == lexer.TokenPeriod {
		p.nextToken()
	}

	return stmt
}

// parseExpression parses an expression
func (p *Parser) parseExpression() ast.Expression {
	return p.parsePrimaryExpression()
}

// parsePrimaryExpression parses a primary expression (literals, identifiers)
func (p *Parser) parsePrimaryExpression() ast.Expression {
	switch p.curTok.Type {
	case lexer.TokenInteger:
		return p.parseIntegerLiteral()
	case lexer.TokenFloat:
		return p.parseFloatLiteral()
	case lexer.TokenString:
		return p.parseStringLiteral()
	case lexer.TokenTrue:
		return &ast.BooleanLiteral{Value: true}
	case lexer.TokenFalse:
		return &ast.BooleanLiteral{Value: false}
	case lexer.TokenNil:
		return &ast.NilLiteral{}
	case lexer.TokenIdentifier:
		return &ast.Identifier{Name: p.curTok.Literal}
	default:
		p.addError(fmt.Sprintf("unexpected token: %s", p.curTok.Type))
		return nil
	}
}

// parseIntegerLiteral parses an integer literal
func (p *Parser) parseIntegerLiteral() ast.Expression {
	value, err := strconv.ParseInt(p.curTok.Literal, 10, 64)
	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as integer", p.curTok.Literal))
		return nil
	}
	return &ast.IntegerLiteral{Value: value}
}

// parseFloatLiteral parses a float literal
func (p *Parser) parseFloatLiteral() ast.Expression {
	value, err := strconv.ParseFloat(p.curTok.Literal, 64)
	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as float", p.curTok.Literal))
		return nil
	}
	return &ast.FloatLiteral{Value: value}
}

// parseStringLiteral parses a string literal
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Value: p.curTok.Literal}
}

// addError adds an error message
func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
}

// Errors returns the list of parsing errors
func (p *Parser) Errors() []string {
	return p.errors
}

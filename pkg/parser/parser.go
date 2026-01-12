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
	// Check for variable declarations
	if p.curTok.Type == lexer.TokenPipe {
		return p.parseVariableDeclaration()
	}

	// Otherwise, parse expression statement
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

// parseVariableDeclaration parses a variable declaration
func (p *Parser) parseVariableDeclaration() ast.Statement {
	// Skip opening pipe
	p.nextToken()

	var names []string
	for p.curTok.Type == lexer.TokenIdentifier {
		names = append(names, p.curTok.Literal)
		p.nextToken()
	}

	// Expect closing pipe
	if p.curTok.Type != lexer.TokenPipe {
		p.addError("expected closing | in variable declaration")
		return nil
	}

	return &ast.VariableDeclaration{Names: names}
}

// parseExpression parses an expression
func (p *Parser) parseExpression() ast.Expression {
	// Try to parse an assignment first
	if p.curTok.Type == lexer.TokenIdentifier && p.peekTok.Type == lexer.TokenAssign {
		return p.parseAssignment()
	}

	// Otherwise, parse a message send or primary expression
	return p.parseMessageSend()
}

// parseAssignment parses an assignment expression
func (p *Parser) parseAssignment() ast.Expression {
	name := p.curTok.Literal
	p.nextToken() // consume identifier

	if p.curTok.Type != lexer.TokenAssign {
		p.addError("expected := in assignment")
		return nil
	}
	p.nextToken() // consume :=

	value := p.parseMessageSend()
	if value == nil {
		return nil
	}

	return &ast.Assignment{
		Name:  name,
		Value: value,
	}
}

// parseMessageSend parses a message send expression
func (p *Parser) parseMessageSend() ast.Expression {
	receiver := p.parsePrimaryExpression()
	if receiver == nil {
		return nil
	}

	// Check for message send (unary or keyword)
	if p.peekTok.Type == lexer.TokenIdentifier {
		p.nextToken() // advance to the identifier
		// Now curTok is the identifier, check if it's followed by a colon
		if p.peekTok.Type == lexer.TokenColon {
			// It's a keyword message - put the identifier back as the first keyword
			keyword := p.curTok.Literal
			p.nextToken() // consume colon
			selector := keyword + ":"
			
			// Parse first argument
			p.nextToken()
			arg := p.parsePrimaryExpression()
			if arg == nil {
				p.addError("expected argument after keyword")
				return nil
			}
			args := []ast.Expression{arg}
			
			// Check for more keyword parts
			for p.peekTok.Type == lexer.TokenIdentifier {
				// Save position to check if it's another keyword part
				savedCur := p.curTok
				savedPeek := p.peekTok
				p.nextToken() // advance to identifier
				if p.peekTok.Type == lexer.TokenColon {
					// Yes, another keyword part
					keyword := p.curTok.Literal
					p.nextToken() // consume colon
					selector += keyword + ":"
					
					// Parse argument
					p.nextToken()
					arg := p.parsePrimaryExpression()
					if arg == nil {
						p.addError("expected argument after keyword")
						return nil
					}
					args = append(args, arg)
				} else {
					// Not a keyword part, restore
					p.curTok = savedCur
					p.peekTok = savedPeek
					break
				}
			}
			
			return &ast.MessageSend{
				Receiver: receiver,
				Selector: selector,
				Args:     args,
			}
		} else {
			// It's a unary message
			selector := p.curTok.Literal
			return &ast.MessageSend{
				Receiver: receiver,
				Selector: selector,
				Args:     []ast.Expression{},
			}
		}
	}

	// Check for binary message
	if p.isBinaryOperator(p.peekTok.Type) {
		p.nextToken()
		operator := p.curTok.Literal
		p.nextToken()
		arg := p.parsePrimaryExpression()
		if arg == nil {
			return nil
		}
		return &ast.MessageSend{
			Receiver: receiver,
			Selector: operator,
			Args:     []ast.Expression{arg},
		}
	}

	return receiver
}

// isBinaryOperator checks if a token type is a binary operator
func (p *Parser) isBinaryOperator(tt lexer.TokenType) bool {
	return tt == lexer.TokenPlus ||
		tt == lexer.TokenMinus ||
		tt == lexer.TokenStar ||
		tt == lexer.TokenSlash ||
		tt == lexer.TokenPercent ||
		tt == lexer.TokenLess ||
		tt == lexer.TokenGreater ||
		tt == lexer.TokenLessEq ||
		tt == lexer.TokenGreaterEq ||
		tt == lexer.TokenEqual ||
		tt == lexer.TokenNotEqual
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

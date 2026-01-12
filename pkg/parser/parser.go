// Package parser implements the smog language parser.
// It converts source code text into an Abstract Syntax Tree (AST).
package parser

import (
	"github.com/kristofer/smog/pkg/ast"
)

// Parser represents the smog parser
type Parser struct {
	source string
	pos    int
}

// New creates a new parser for the given source code
func New(source string) *Parser {
	return &Parser{
		source: source,
		pos:    0,
	}
}

// Parse parses the source code and returns an AST
func (p *Parser) Parse() (*ast.Program, error) {
	// TODO: Implement parser
	return &ast.Program{}, nil
}

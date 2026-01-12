// Package compiler compiles AST nodes into bytecode.
package compiler

import (
	"github.com/kristofer/smog/pkg/ast"
	"github.com/kristofer/smog/pkg/bytecode"
)

// Compiler represents the bytecode compiler
type Compiler struct {
	constants []interface{}
	symbols   map[string]int
}

// New creates a new compiler
func New() *Compiler {
	return &Compiler{
		constants: make([]interface{}, 0),
		symbols:   make(map[string]int),
	}
}

// Compile compiles an AST program into bytecode
func (c *Compiler) Compile(program *ast.Program) (*bytecode.Bytecode, error) {
	// TODO: Implement compiler
	return &bytecode.Bytecode{}, nil
}

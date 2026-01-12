// Package compiler compiles AST nodes into bytecode.
package compiler

import (
	"fmt"

	"github.com/kristofer/smog/pkg/ast"
	"github.com/kristofer/smog/pkg/bytecode"
)

// Compiler represents the bytecode compiler
type Compiler struct {
	instructions []bytecode.Instruction
	constants    []interface{}
	symbols      map[string]int
	localCount   int
}

// New creates a new compiler
func New() *Compiler {
	return &Compiler{
		instructions: make([]bytecode.Instruction, 0),
		constants:    make([]interface{}, 0),
		symbols:      make(map[string]int),
		localCount:   0,
	}
}

// Compile compiles an AST program into bytecode
func (c *Compiler) Compile(program *ast.Program) (*bytecode.Bytecode, error) {
	for _, stmt := range program.Statements {
		if err := c.compileStatement(stmt); err != nil {
			return nil, err
		}
	}

	// Add final return
	c.emit(bytecode.OpReturn, 0)

	return &bytecode.Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}, nil
}

// compileStatement compiles a statement
func (c *Compiler) compileStatement(stmt ast.Statement) error {
	switch s := stmt.(type) {
	case *ast.ExpressionStatement:
		return c.compileExpression(s.Expression)
	case *ast.VariableDeclaration:
		// Variable declarations don't generate code, they just reserve space
		for _, name := range s.Names {
			c.symbols[name] = c.localCount
			c.localCount++
		}
		return nil
	default:
		return fmt.Errorf("unknown statement type: %T", stmt)
	}
}

// compileExpression compiles an expression
func (c *Compiler) compileExpression(expr ast.Expression) error {
	switch e := expr.(type) {
	case *ast.IntegerLiteral:
		idx := c.addConstant(e.Value)
		c.emit(bytecode.OpPush, idx)
		return nil

	case *ast.FloatLiteral:
		idx := c.addConstant(e.Value)
		c.emit(bytecode.OpPush, idx)
		return nil

	case *ast.StringLiteral:
		idx := c.addConstant(e.Value)
		c.emit(bytecode.OpPush, idx)
		return nil

	case *ast.BooleanLiteral:
		if e.Value {
			c.emit(bytecode.OpPushTrue, 0)
		} else {
			c.emit(bytecode.OpPushFalse, 0)
		}
		return nil

	case *ast.NilLiteral:
		c.emit(bytecode.OpPushNil, 0)
		return nil

	case *ast.Identifier:
		// Look up the variable in the symbol table
		if idx, ok := c.symbols[e.Name]; ok {
			c.emit(bytecode.OpLoadLocal, idx)
		} else {
			// Try to load as global
			idx := c.addConstant(e.Name)
			c.emit(bytecode.OpLoadGlobal, idx)
		}
		return nil

	case *ast.Assignment:
		// Compile the value first
		if err := c.compileExpression(e.Value); err != nil {
			return err
		}

		// Store to variable
		if idx, ok := c.symbols[e.Name]; ok {
			c.emit(bytecode.OpStoreLocal, idx)
		} else {
			// Store as global
			nameIdx := c.addConstant(e.Name)
			c.emit(bytecode.OpStoreGlobal, nameIdx)
		}
		return nil

	case *ast.MessageSend:
		// Compile receiver
		if err := c.compileExpression(e.Receiver); err != nil {
			return err
		}

		// Compile arguments
		for _, arg := range e.Args {
			if err := c.compileExpression(arg); err != nil {
				return err
			}
		}

		// Emit send instruction
		selectorIdx := c.addConstant(e.Selector)
		argCount := len(e.Args)
		// Pack selector and arg count into operand using shared constants
		operand := (selectorIdx << bytecode.SelectorIndexShift) | argCount
		c.emit(bytecode.OpSend, operand)
		return nil

	default:
		return fmt.Errorf("unknown expression type: %T", expr)
	}
}

// emit adds an instruction to the bytecode
func (c *Compiler) emit(op bytecode.Opcode, operand int) {
	c.instructions = append(c.instructions, bytecode.Instruction{
		Op:      op,
		Operand: operand,
	})
}

// addConstant adds a constant to the constant pool and returns its index
func (c *Compiler) addConstant(obj interface{}) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

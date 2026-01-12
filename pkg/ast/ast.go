// Package ast defines the Abstract Syntax Tree nodes for smog.
package ast

// Node is the interface that all AST nodes implement
type Node interface {
	TokenLiteral() string
}

// Expression represents an expression node
type Expression interface {
	Node
	expressionNode()
}

// Statement represents a statement node
type Statement interface {
	Node
	statementNode()
}

// Program represents the root node of the AST
type Program struct {
	Statements []Statement
}

// TokenLiteral returns the token literal
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// Class represents a class definition
type Class struct {
	Name       string
	SuperClass string
	Methods    []*Method
	Fields     []string
}

// TokenLiteral returns the token literal
func (c *Class) TokenLiteral() string { return "class" }
func (c *Class) statementNode()       {}

// Method represents a method definition
type Method struct {
	Name       string
	Parameters []string
	Body       []Statement
}

// TokenLiteral returns the token literal
func (m *Method) TokenLiteral() string { return "method" }

// MessageSend represents a message send expression
type MessageSend struct {
	Receiver Expression
	Selector string
	Args     []Expression
}

// TokenLiteral returns the token literal
func (m *MessageSend) TokenLiteral() string { return m.Selector }
func (m *MessageSend) expressionNode()      {}

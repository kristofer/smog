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

// ExpressionStatement represents an expression used as a statement
type ExpressionStatement struct {
	Expression Expression
}

// TokenLiteral returns the token literal
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Expression.TokenLiteral()
}
func (es *ExpressionStatement) statementNode() {}

// VariableDeclaration represents a variable declaration
type VariableDeclaration struct {
	Names []string
}

// TokenLiteral returns the token literal
func (vd *VariableDeclaration) TokenLiteral() string { return "" }
func (vd *VariableDeclaration) statementNode()       {}

// Assignment represents a variable assignment
type Assignment struct {
	Name  string
	Value Expression
}

// TokenLiteral returns the token literal
func (a *Assignment) TokenLiteral() string { return a.Name }
func (a *Assignment) expressionNode()      {}

// IntegerLiteral represents an integer literal
type IntegerLiteral struct {
	Value int64
}

// TokenLiteral returns the token literal
func (il *IntegerLiteral) TokenLiteral() string { return "" }
func (il *IntegerLiteral) expressionNode()      {}

// FloatLiteral represents a float literal
type FloatLiteral struct {
	Value float64
}

// TokenLiteral returns the token literal
func (fl *FloatLiteral) TokenLiteral() string { return "" }
func (fl *FloatLiteral) expressionNode()      {}

// StringLiteral represents a string literal
type StringLiteral struct {
	Value string
}

// TokenLiteral returns the token literal
func (sl *StringLiteral) TokenLiteral() string { return sl.Value }
func (sl *StringLiteral) expressionNode()      {}

// BooleanLiteral represents a boolean literal
type BooleanLiteral struct {
	Value bool
}

// TokenLiteral returns the token literal
func (bl *BooleanLiteral) TokenLiteral() string {
	if bl.Value {
		return "true"
	}
	return "false"
}
func (bl *BooleanLiteral) expressionNode() {}

// NilLiteral represents the nil literal
type NilLiteral struct{}

// TokenLiteral returns the token literal
func (nl *NilLiteral) TokenLiteral() string { return "nil" }
func (nl *NilLiteral) expressionNode()      {}

// Identifier represents an identifier
type Identifier struct {
	Name string
}

// TokenLiteral returns the token literal
func (i *Identifier) TokenLiteral() string { return i.Name }
func (i *Identifier) expressionNode()      {}

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

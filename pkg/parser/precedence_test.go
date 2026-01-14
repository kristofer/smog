package parser

import (
	"testing"

	"github.com/kristofer/smog/pkg/ast"
)

// TestParseUnaryBinaryPrecedence tests that unary messages have higher precedence than binary
func TestParseUnaryBinaryPrecedence(t *testing.T) {
input := "arr size + 1"

p := New(input)
program, err := p.Parse()

if err != nil {
t.Fatalf("Parse returned error: %v", err)
}

stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
if !ok {
t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
}

// Should be: (arr size) + 1
// Top level is binary "+"
msg, ok := stmt.Expression.(*ast.MessageSend)
if !ok {
t.Fatalf("Expected MessageSend, got %T", stmt.Expression)
}

if msg.Selector != "+" {
t.Errorf("Expected top-level selector '+', got %s", msg.Selector)
}

// Receiver should be (arr size)
receiverMsg, ok := msg.Receiver.(*ast.MessageSend)
if !ok {
t.Fatalf("Expected MessageSend receiver, got %T", msg.Receiver)
}

if receiverMsg.Selector != "size" {
t.Errorf("Expected receiver selector 'size', got %s", receiverMsg.Selector)
}
}

// TestParseBinaryChaining tests that binary messages chain left-to-right
func TestParseBinaryChaining(t *testing.T) {
input := "3 + 4 * 2"

p := New(input)
program, err := p.Parse()

if err != nil {
t.Fatalf("Parse returned error: %v", err)
}

stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
if !ok {
t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
}

// Should be: (3 + 4) * 2
// Top level is binary "*"
msg, ok := stmt.Expression.(*ast.MessageSend)
if !ok {
t.Fatalf("Expected MessageSend, got %T", stmt.Expression)
}

if msg.Selector != "*" {
t.Errorf("Expected top-level selector '*', got %s", msg.Selector)
}

// Receiver should be (3 + 4)
receiverMsg, ok := msg.Receiver.(*ast.MessageSend)
if !ok {
t.Fatalf("Expected MessageSend receiver, got %T", msg.Receiver)
}

if receiverMsg.Selector != "+" {
t.Errorf("Expected receiver selector '+', got %s", receiverMsg.Selector)
}
}

// TestParseUnaryChaining tests that unary messages chain left-to-right
func TestParseUnaryChaining(t *testing.T) {
input := "x sqrt floor"

p := New(input)
program, err := p.Parse()

if err != nil {
t.Fatalf("Parse returned error: %v", err)
}

stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
if !ok {
t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
}

// Should be: (x sqrt) floor
// Top level is unary "floor"
msg, ok := stmt.Expression.(*ast.MessageSend)
if !ok {
t.Fatalf("Expected MessageSend, got %T", stmt.Expression)
}

if msg.Selector != "floor" {
t.Errorf("Expected top-level selector 'floor', got %s", msg.Selector)
}

// Receiver should be (x sqrt)
receiverMsg, ok := msg.Receiver.(*ast.MessageSend)
if !ok {
t.Fatalf("Expected MessageSend receiver, got %T", msg.Receiver)
}

if receiverMsg.Selector != "sqrt" {
t.Errorf("Expected receiver selector 'sqrt', got %s", receiverMsg.Selector)
}
}

// TestParseKeywordWithBinaryArg tests that keyword message arguments can be binary expressions
func TestParseKeywordWithBinaryArg(t *testing.T) {
input := "arr at: index + 1"

p := New(input)
program, err := p.Parse()

if err != nil {
t.Fatalf("Parse returned error: %v", err)
}

stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
if !ok {
t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
}

// Top level is keyword "at:"
msg, ok := stmt.Expression.(*ast.MessageSend)
if !ok {
t.Fatalf("Expected MessageSend, got %T", stmt.Expression)
}

if msg.Selector != "at:" {
t.Errorf("Expected selector 'at:', got %s", msg.Selector)
}

// Argument should be (index + 1)
if len(msg.Args) != 1 {
t.Fatalf("Expected 1 argument, got %d", len(msg.Args))
}

argMsg, ok := msg.Args[0].(*ast.MessageSend)
if !ok {
t.Fatalf("Expected MessageSend argument, got %T", msg.Args[0])
}

if argMsg.Selector != "+" {
t.Errorf("Expected argument selector '+', got %s", argMsg.Selector)
}
}

// TestParseComplexPrecedence tests a complex expression with all three precedence levels
func TestParseComplexPrecedence(t *testing.T) {
input := "point x: a + b y: c size"

p := New(input)
program, err := p.Parse()

if err != nil {
t.Fatalf("Parse returned error: %v", err)
}

stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
if !ok {
t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
}

// Top level is keyword "x:y:"
msg, ok := stmt.Expression.(*ast.MessageSend)
if !ok {
t.Fatalf("Expected MessageSend, got %T", stmt.Expression)
}

if msg.Selector != "x:y:" {
t.Errorf("Expected selector 'x:y:', got %s", msg.Selector)
}

// Should have 2 arguments
if len(msg.Args) != 2 {
t.Fatalf("Expected 2 arguments, got %d", len(msg.Args))
}

// First argument should be (a + b)
arg1Msg, ok := msg.Args[0].(*ast.MessageSend)
if !ok {
t.Fatalf("Expected MessageSend first argument, got %T", msg.Args[0])
}
if arg1Msg.Selector != "+" {
t.Errorf("Expected first argument selector '+', got %s", arg1Msg.Selector)
}

// Second argument should be (c size)
arg2Msg, ok := msg.Args[1].(*ast.MessageSend)
if !ok {
t.Fatalf("Expected MessageSend second argument, got %T", msg.Args[1])
}
if arg2Msg.Selector != "size" {
t.Errorf("Expected second argument selector 'size', got %s", arg2Msg.Selector)
}
}

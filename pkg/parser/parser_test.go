package parser

import (
	"testing"

	"github.com/kristofer/smog/pkg/ast"
)

func TestParseIntegerLiteral(t *testing.T) {
	input := "42"

	p := New(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	intLit, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expected IntegerLiteral, got %T", stmt.Expression)
	}

	if intLit.Value != 42 {
		t.Errorf("Expected value 42, got %d", intLit.Value)
	}
}

func TestParseFloatLiteral(t *testing.T) {
	input := "3.14"

	p := New(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	floatLit, ok := stmt.Expression.(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("Expected FloatLiteral, got %T", stmt.Expression)
	}

	if floatLit.Value != 3.14 {
		t.Errorf("Expected value 3.14, got %f", floatLit.Value)
	}
}

func TestParseStringLiteral(t *testing.T) {
	input := "'Hello, World!'"

	p := New(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	strLit, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("Expected StringLiteral, got %T", stmt.Expression)
	}

	if strLit.Value != "Hello, World!" {
		t.Errorf("Expected value 'Hello, World!', got %s", strLit.Value)
	}
}

func TestParseBooleanLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		p := New(tt.input)
		program, err := p.Parse()

		if err != nil {
			t.Fatalf("Parse returned error: %v", err)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
		}

		boolLit, ok := stmt.Expression.(*ast.BooleanLiteral)
		if !ok {
			t.Fatalf("Expected BooleanLiteral, got %T", stmt.Expression)
		}

		if boolLit.Value != tt.expected {
			t.Errorf("Expected value %v, got %v", tt.expected, boolLit.Value)
		}
	}
}

func TestParseNilLiteral(t *testing.T) {
	input := "nil"

	p := New(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	_, ok = stmt.Expression.(*ast.NilLiteral)
	if !ok {
		t.Fatalf("Expected NilLiteral, got %T", stmt.Expression)
	}
}

func TestParseIdentifier(t *testing.T) {
	input := "println"

	p := New(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier, got %T", stmt.Expression)
	}

	if ident.Name != "println" {
		t.Errorf("Expected identifier 'println', got %s", ident.Name)
	}
}

func TestParseMultipleStatements(t *testing.T) {
	input := `42.
'hello'.
true.`

	p := New(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Expected 3 statements, got %d", len(program.Statements))
	}

	// First statement: integer
	stmt1, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}
	if _, ok := stmt1.Expression.(*ast.IntegerLiteral); !ok {
		t.Errorf("Expected IntegerLiteral in first statement, got %T", stmt1.Expression)
	}

	// Second statement: string
	stmt2, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[1])
	}
	if _, ok := stmt2.Expression.(*ast.StringLiteral); !ok {
		t.Errorf("Expected StringLiteral in second statement, got %T", stmt2.Expression)
	}

	// Third statement: boolean
	stmt3, ok := program.Statements[2].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[2])
	}
	if _, ok := stmt3.Expression.(*ast.BooleanLiteral); !ok {
		t.Errorf("Expected BooleanLiteral in third statement, got %T", stmt3.Expression)
	}
}

func TestParseNegativeNumber(t *testing.T) {
	input := "-17"

	p := New(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	intLit, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expected IntegerLiteral, got %T", stmt.Expression)
	}

	if intLit.Value != -17 {
		t.Errorf("Expected value -17, got %d", intLit.Value)
	}
}

func TestParseWithComments(t *testing.T) {
	input := `" This is a comment "
42`

	p := New(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	intLit, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expected IntegerLiteral, got %T", stmt.Expression)
	}

	if intLit.Value != 42 {
		t.Errorf("Expected value 42, got %d", intLit.Value)
	}
}

func TestParseVariableDeclaration(t *testing.T) {
	input := `| x y z |`

	p := New(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	varDecl, ok := program.Statements[0].(*ast.VariableDeclaration)
	if !ok {
		t.Fatalf("Expected VariableDeclaration, got %T", program.Statements[0])
	}

	if len(varDecl.Names) != 3 {
		t.Fatalf("Expected 3 variable names, got %d", len(varDecl.Names))
	}

	expectedNames := []string{"x", "y", "z"}
	for i, name := range expectedNames {
		if varDecl.Names[i] != name {
			t.Errorf("Expected variable name %s, got %s", name, varDecl.Names[i])
		}
	}
}

func TestParseAssignment(t *testing.T) {
	input := `x := 42`

	p := New(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	assign, ok := stmt.Expression.(*ast.Assignment)
	if !ok {
		t.Fatalf("Expected Assignment, got %T", stmt.Expression)
	}

	if assign.Name != "x" {
		t.Errorf("Expected variable name 'x', got %s", assign.Name)
	}

	intLit, ok := assign.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expected IntegerLiteral value, got %T", assign.Value)
	}

	if intLit.Value != 42 {
		t.Errorf("Expected value 42, got %d", intLit.Value)
	}
}

func TestParseUnaryMessageSend(t *testing.T) {
	input := `'Hello' println`

	p := New(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	msg, ok := stmt.Expression.(*ast.MessageSend)
	if !ok {
		t.Fatalf("Expected MessageSend, got %T", stmt.Expression)
	}

	if msg.Selector != "println" {
		t.Errorf("Expected selector 'println', got %s", msg.Selector)
	}

	strLit, ok := msg.Receiver.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("Expected StringLiteral receiver, got %T", msg.Receiver)
	}

	if strLit.Value != "Hello" {
		t.Errorf("Expected receiver 'Hello', got %s", strLit.Value)
	}

	if len(msg.Args) != 0 {
		t.Errorf("Expected 0 arguments, got %d", len(msg.Args))
	}
}

func TestParseBinaryMessageSend(t *testing.T) {
	input := `3 + 4`

	p := New(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	msg, ok := stmt.Expression.(*ast.MessageSend)
	if !ok {
		t.Fatalf("Expected MessageSend, got %T", stmt.Expression)
	}

	if msg.Selector != "+" {
		t.Errorf("Expected selector '+', got %s", msg.Selector)
	}

	receiver, ok := msg.Receiver.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expected IntegerLiteral receiver, got %T", msg.Receiver)
	}

	if receiver.Value != 3 {
		t.Errorf("Expected receiver value 3, got %d", receiver.Value)
	}

	if len(msg.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(msg.Args))
	}

	arg, ok := msg.Args[0].(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expected IntegerLiteral argument, got %T", msg.Args[0])
	}

	if arg.Value != 4 {
		t.Errorf("Expected argument value 4, got %d", arg.Value)
	}
}

func TestParseKeywordMessageSend(t *testing.T) {
	input := `point x: 10 y: 20`

	p := New(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	msg, ok := stmt.Expression.(*ast.MessageSend)
	if !ok {
		t.Fatalf("Expected MessageSend, got %T", stmt.Expression)
	}

	if msg.Selector != "x:y:" {
		t.Errorf("Expected selector 'x:y:', got %s", msg.Selector)
	}

	receiver, ok := msg.Receiver.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier receiver, got %T", msg.Receiver)
	}

	if receiver.Name != "point" {
		t.Errorf("Expected receiver 'point', got %s", receiver.Name)
	}

	if len(msg.Args) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(msg.Args))
	}

	arg1, ok := msg.Args[0].(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expected IntegerLiteral first argument, got %T", msg.Args[0])
	}

	if arg1.Value != 10 {
		t.Errorf("Expected first argument 10, got %d", arg1.Value)
	}

	arg2, ok := msg.Args[1].(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expected IntegerLiteral second argument, got %T", msg.Args[1])
	}

	if arg2.Value != 20 {
		t.Errorf("Expected second argument 20, got %d", arg2.Value)
	}
}

func TestParseSimpleBlockLiteral(t *testing.T) {
	input := "[ 'Hello' println ]"

	p := New(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	block, ok := stmt.Expression.(*ast.BlockLiteral)
	if !ok {
		t.Fatalf("Expected BlockLiteral, got %T", stmt.Expression)
	}

	if len(block.Parameters) != 0 {
		t.Errorf("Expected 0 parameters, got %d", len(block.Parameters))
	}

	if len(block.Body) != 1 {
		t.Fatalf("Expected 1 statement in block body, got %d", len(block.Body))
	}
}

func TestParseBlockLiteralWithOneParameter(t *testing.T) {
	input := "[ :x | x * 2 ]"

	p := New(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	block, ok := stmt.Expression.(*ast.BlockLiteral)
	if !ok {
		t.Fatalf("Expected BlockLiteral, got %T", stmt.Expression)
	}

	if len(block.Parameters) != 1 {
		t.Fatalf("Expected 1 parameter, got %d", len(block.Parameters))
	}

	if block.Parameters[0] != "x" {
		t.Errorf("Expected parameter 'x', got '%s'", block.Parameters[0])
	}

	if len(block.Body) != 1 {
		t.Fatalf("Expected 1 statement in block body, got %d", len(block.Body))
	}
}

func TestParseBlockLiteralWithMultipleParameters(t *testing.T) {
	input := "[ :x :y | x + y ]"

	p := New(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	block, ok := stmt.Expression.(*ast.BlockLiteral)
	if !ok {
		t.Fatalf("Expected BlockLiteral, got %T", stmt.Expression)
	}

	if len(block.Parameters) != 2 {
		t.Fatalf("Expected 2 parameters, got %d", len(block.Parameters))
	}

	if block.Parameters[0] != "x" {
		t.Errorf("Expected first parameter 'x', got '%s'", block.Parameters[0])
	}

	if block.Parameters[1] != "y" {
		t.Errorf("Expected second parameter 'y', got '%s'", block.Parameters[1])
	}
}

func TestParseReturnStatement(t *testing.T) {
	input := "^42"

	p := New(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	ret, ok := program.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("Expected ReturnStatement, got %T", program.Statements[0])
	}

	intLit, ok := ret.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expected IntegerLiteral in return, got %T", ret.Value)
	}

	if intLit.Value != 42 {
		t.Errorf("Expected return value 42, got %d", intLit.Value)
	}
}

func TestParseArrayLiteral(t *testing.T) {
	input := "#(1 2 3 4 5)"

	p := New(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	arr, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("Expected ArrayLiteral, got %T", stmt.Expression)
	}

	if len(arr.Elements) != 5 {
		t.Fatalf("Expected 5 elements, got %d", len(arr.Elements))
	}

	expected := []int64{1, 2, 3, 4, 5}
	for i, elem := range arr.Elements {
		intLit, ok := elem.(*ast.IntegerLiteral)
		if !ok {
			t.Fatalf("Expected IntegerLiteral at index %d, got %T", i, elem)
		}
		if intLit.Value != expected[i] {
			t.Errorf("Expected element %d to be %d, got %d", i, expected[i], intLit.Value)
		}
	}
}

// TestParseSelfKeyword tests parsing the 'self' keyword
func TestParseSelfKeyword(t *testing.T) {
input := "self"

p := New(input)
program, err := p.Parse()

if err != nil {
t.Fatalf("Parse returned error: %v", err)
}

if len(program.Statements) != 1 {
t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
}

stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
if !ok {
t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
}

ident, ok := stmt.Expression.(*ast.Identifier)
if !ok {
t.Fatalf("Expected Identifier, got %T", stmt.Expression)
}

if ident.Name != "self" {
t.Errorf("Expected identifier 'self', got %s", ident.Name)
}
}

// TestParseSuperUnaryMessage tests parsing super with a unary message
func TestParseSuperUnaryMessage(t *testing.T) {
input := "super initialize"

p := New(input)
program, err := p.Parse()

if err != nil {
t.Fatalf("Parse returned error: %v", err)
}

if len(program.Statements) != 1 {
t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
}

stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
if !ok {
t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
}

msgSend, ok := stmt.Expression.(*ast.MessageSend)
if !ok {
t.Fatalf("Expected MessageSend, got %T", stmt.Expression)
}

if !msgSend.IsSuper {
t.Error("Expected IsSuper to be true")
}

if msgSend.Selector != "initialize" {
t.Errorf("Expected selector 'initialize', got %s", msgSend.Selector)
}

if len(msgSend.Args) != 0 {
t.Errorf("Expected 0 arguments, got %d", len(msgSend.Args))
}

if msgSend.Receiver != nil {
t.Error("Expected nil receiver for super send")
}
}

// TestParseSuperKeywordMessage tests parsing super with a keyword message
func TestParseSuperKeywordMessage(t *testing.T) {
input := "super at: 5 put: 10"

p := New(input)
program, err := p.Parse()

if err != nil {
t.Fatalf("Parse returned error: %v", err)
}

if len(program.Statements) != 1 {
t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
}

stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
if !ok {
t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
}

msgSend, ok := stmt.Expression.(*ast.MessageSend)
if !ok {
t.Fatalf("Expected MessageSend, got %T", stmt.Expression)
}

if !msgSend.IsSuper {
t.Error("Expected IsSuper to be true")
}

if msgSend.Selector != "at:put:" {
t.Errorf("Expected selector 'at:put:', got %s", msgSend.Selector)
}

if len(msgSend.Args) != 2 {
t.Fatalf("Expected 2 arguments, got %d", len(msgSend.Args))
}

// Check first argument
arg1, ok := msgSend.Args[0].(*ast.IntegerLiteral)
if !ok {
t.Fatalf("Expected first arg to be IntegerLiteral, got %T", msgSend.Args[0])
}
if arg1.Value != 5 {
t.Errorf("Expected first arg value 5, got %d", arg1.Value)
}

// Check second argument
arg2, ok := msgSend.Args[1].(*ast.IntegerLiteral)
if !ok {
t.Fatalf("Expected second arg to be IntegerLiteral, got %T", msgSend.Args[1])
}
if arg2.Value != 10 {
t.Errorf("Expected second arg value 10, got %d", arg2.Value)
}
}

// TestParseSuperBinaryMessage tests parsing super with a binary message
func TestParseSuperBinaryMessage(t *testing.T) {
input := "super + 5"

p := New(input)
program, err := p.Parse()

if err != nil {
t.Fatalf("Parse returned error: %v", err)
}

if len(program.Statements) != 1 {
t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
}

stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
if !ok {
t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
}

msgSend, ok := stmt.Expression.(*ast.MessageSend)
if !ok {
t.Fatalf("Expected MessageSend, got %T", stmt.Expression)
}

if !msgSend.IsSuper {
t.Error("Expected IsSuper to be true")
}

if msgSend.Selector != "+" {
t.Errorf("Expected selector '+', got %s", msgSend.Selector)
}

if len(msgSend.Args) != 1 {
t.Fatalf("Expected 1 argument, got %d", len(msgSend.Args))
}

arg, ok := msgSend.Args[0].(*ast.IntegerLiteral)
if !ok {
t.Fatalf("Expected arg to be IntegerLiteral, got %T", msgSend.Args[0])
}
if arg.Value != 5 {
t.Errorf("Expected arg value 5, got %d", arg.Value)
}
}

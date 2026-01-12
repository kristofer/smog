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

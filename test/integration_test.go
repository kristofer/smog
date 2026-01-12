// Package test provides integration tests for smog.
package test

import (
	"testing"

	"github.com/kristofer/smog/pkg/ast"
	"github.com/kristofer/smog/pkg/lexer"
	"github.com/kristofer/smog/pkg/parser"
)

// TestVersion0_1_0_Lexer tests lexer functionality for version 0.1.0
func TestVersion0_1_0_Lexer(t *testing.T) {
	t.Run("TokenizesHelloWorld", func(t *testing.T) {
		input := "'Hello, World!' println."

		l := lexer.New(input)
		tokens, err := l.Tokenize()

		if err != nil {
			t.Fatalf("Tokenize failed: %v", err)
		}

		expectedCount := 4 // STRING, IDENTIFIER, PERIOD, EOF
		if len(tokens) != expectedCount {
			t.Errorf("Expected %d tokens, got %d", expectedCount, len(tokens))
		}

		if tokens[0].Type != lexer.TokenString {
			t.Errorf("Expected first token to be STRING, got %v", tokens[0].Type)
		}

		if tokens[1].Type != lexer.TokenIdentifier {
			t.Errorf("Expected second token to be IDENTIFIER, got %v", tokens[1].Type)
		}
	})

	t.Run("TokenizesNumbers", func(t *testing.T) {
		input := "42 3.14 -17"

		l := lexer.New(input)
		tokens, err := l.Tokenize()

		if err != nil {
			t.Fatalf("Tokenize failed: %v", err)
		}

		if tokens[0].Type != lexer.TokenInteger || tokens[0].Literal != "42" {
			t.Errorf("Expected INTEGER token with value 42")
		}

		if tokens[1].Type != lexer.TokenFloat || tokens[1].Literal != "3.14" {
			t.Errorf("Expected FLOAT token with value 3.14")
		}

		if tokens[2].Type != lexer.TokenInteger || tokens[2].Literal != "-17" {
			t.Errorf("Expected INTEGER token with value -17")
		}
	})

	t.Run("SkipsComments", func(t *testing.T) {
		input := `" This is a comment " 42`

		l := lexer.New(input)
		tokens, err := l.Tokenize()

		if err != nil {
			t.Fatalf("Tokenize failed: %v", err)
		}

		// Should only have INTEGER and EOF, no comment token
		if len(tokens) != 2 {
			t.Errorf("Expected 2 tokens (INTEGER, EOF), got %d", len(tokens))
		}

		if tokens[0].Type != lexer.TokenInteger {
			t.Errorf("Expected first token to be INTEGER, got %v", tokens[0].Type)
		}
	})
}

// TestVersion0_1_0_Parser tests parser functionality for version 0.1.0
func TestVersion0_1_0_Parser(t *testing.T) {
	t.Run("ParsesIntegerLiteral", func(t *testing.T) {
		input := "42"

		p := parser.New(input)
		program, err := p.Parse()

		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
		}

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		intLit := stmt.Expression.(*ast.IntegerLiteral)

		if intLit.Value != 42 {
			t.Errorf("Expected integer value 42, got %d", intLit.Value)
		}
	})

	t.Run("ParsesStringLiteral", func(t *testing.T) {
		input := "'Hello, World!'"

		p := parser.New(input)
		program, err := p.Parse()

		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
		}

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		strLit := stmt.Expression.(*ast.StringLiteral)

		if strLit.Value != "Hello, World!" {
			t.Errorf("Expected string value 'Hello, World!', got %s", strLit.Value)
		}
	})

	t.Run("ParsesBooleanAndNil", func(t *testing.T) {
		tests := []struct {
			input    string
			nodeType string
		}{
			{"true", "boolean"},
			{"false", "boolean"},
			{"nil", "nil"},
		}

		for _, tt := range tests {
			p := parser.New(tt.input)
			program, err := p.Parse()

			if err != nil {
				t.Fatalf("Parse failed for '%s': %v", tt.input, err)
			}

			if len(program.Statements) != 1 {
				t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
			}

			stmt := program.Statements[0].(*ast.ExpressionStatement)

			switch tt.nodeType {
			case "boolean":
				if _, ok := stmt.Expression.(*ast.BooleanLiteral); !ok {
					t.Errorf("Expected BooleanLiteral for '%s', got %T", tt.input, stmt.Expression)
				}
			case "nil":
				if _, ok := stmt.Expression.(*ast.NilLiteral); !ok {
					t.Errorf("Expected NilLiteral for '%s', got %T", tt.input, stmt.Expression)
				}
			}
		}
	})

	t.Run("ParsesMultipleStatements", func(t *testing.T) {
		input := `42.
'hello'.
true.`

		p := parser.New(input)
		program, err := p.Parse()

		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		if len(program.Statements) != 3 {
			t.Fatalf("Expected 3 statements, got %d", len(program.Statements))
		}
	})
}

// TestVersion0_1_0_LexerAndParser tests the integration of lexer and parser
func TestVersion0_1_0_LexerAndParser(t *testing.T) {
	t.Run("EndToEndParsing", func(t *testing.T) {
		input := `" Simple program "
42.
'Hello'.
true.`

		// First, tokenize
		l := lexer.New(input)
		tokens, err := l.Tokenize()

		if err != nil {
			t.Fatalf("Tokenize failed: %v", err)
		}

		// Verify we got tokens
		if len(tokens) < 4 {
			t.Fatalf("Expected at least 4 tokens, got %d", len(tokens))
		}

		// Then parse
		p := parser.New(input)
		program, err := p.Parse()

		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		// Verify we got 3 statements (comment is skipped)
		if len(program.Statements) != 3 {
			t.Fatalf("Expected 3 statements, got %d", len(program.Statements))
		}

		// Verify statement types
		stmt0 := program.Statements[0].(*ast.ExpressionStatement)
		if _, ok := stmt0.Expression.(*ast.IntegerLiteral); !ok {
			t.Errorf("Expected first statement to be IntegerLiteral, got %T", stmt0.Expression)
		}

		stmt1 := program.Statements[1].(*ast.ExpressionStatement)
		if _, ok := stmt1.Expression.(*ast.StringLiteral); !ok {
			t.Errorf("Expected second statement to be StringLiteral, got %T", stmt1.Expression)
		}

		stmt2 := program.Statements[2].(*ast.ExpressionStatement)
		if _, ok := stmt2.Expression.(*ast.BooleanLiteral); !ok {
			t.Errorf("Expected third statement to be BooleanLiteral, got %T", stmt2.Expression)
		}
	})
}

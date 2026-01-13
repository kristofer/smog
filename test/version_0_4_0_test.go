// Package test provides integration tests for smog v0.4.0 features.
package test

import (
	"testing"

	"github.com/kristofer/smog/pkg/ast"
	"github.com/kristofer/smog/pkg/compiler"
	"github.com/kristofer/smog/pkg/parser"
)

// TestVersion0_4_0_SuperKeyword tests the super keyword parsing
func TestVersion0_4_0_SuperKeyword(t *testing.T) {
	t.Run("ParseSuperUnaryMessage", func(t *testing.T) {
		input := "super initialize"
		p := parser.New(input)
		program, err := p.Parse()

		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
		}

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		msgSend := stmt.Expression.(*ast.MessageSend)

		if !msgSend.IsSuper {
			t.Error("Expected IsSuper to be true")
		}

		if msgSend.Selector != "initialize" {
			t.Errorf("Expected selector 'initialize', got %s", msgSend.Selector)
		}
	})

	t.Run("ParseSuperKeywordMessage", func(t *testing.T) {
		input := "super at: 5"
		p := parser.New(input)
		program, err := p.Parse()

		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		msgSend := stmt.Expression.(*ast.MessageSend)

		if !msgSend.IsSuper {
			t.Error("Expected IsSuper to be true")
		}

		if msgSend.Selector != "at:" {
			t.Errorf("Expected selector 'at:', got %s", msgSend.Selector)
		}

		if len(msgSend.Args) != 1 {
			t.Errorf("Expected 1 argument, got %d", len(msgSend.Args))
		}
	})

	t.Run("CompileSuperMessage", func(t *testing.T) {
		input := "super initialize"
		p := parser.New(input)
		program, err := p.Parse()

		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		c := compiler.New()
		bytecode, err := c.Compile(program)

		if err != nil {
			t.Fatalf("Compile failed: %v", err)
		}

		// Should have PUSH_SELF and SUPER_SEND instructions
		if len(bytecode.Instructions) < 2 {
			t.Errorf("Expected at least 2 instructions, got %d", len(bytecode.Instructions))
		}
	})
}

// TestVersion0_4_0_DictionaryLiterals tests dictionary literal parsing and compilation
func TestVersion0_4_0_DictionaryLiterals(t *testing.T) {
	t.Run("ParseDictionaryLiteral", func(t *testing.T) {
		input := "#{'name' -> 'Alice'. 'age' -> 30}"
		p := parser.New(input)
		program, err := p.Parse()

		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		dictLit, ok := stmt.Expression.(*ast.DictionaryLiteral)
		if !ok {
			t.Fatalf("Expected DictionaryLiteral, got %T", stmt.Expression)
		}

		if len(dictLit.Pairs) != 2 {
			t.Fatalf("Expected 2 pairs, got %d", len(dictLit.Pairs))
		}
	})

	t.Run("ParseEmptyDictionary", func(t *testing.T) {
		input := "#{}"
		p := parser.New(input)
		program, err := p.Parse()

		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		dictLit := stmt.Expression.(*ast.DictionaryLiteral)

		if len(dictLit.Pairs) != 0 {
			t.Errorf("Expected 0 pairs, got %d", len(dictLit.Pairs))
		}
	})

	t.Run("CompileDictionaryLiteral", func(t *testing.T) {
		input := "#{'x' -> 10}"
		p := parser.New(input)
		program, err := p.Parse()

		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		c := compiler.New()
		bytecode, err := c.Compile(program)

		if err != nil {
			t.Fatalf("Compile failed: %v", err)
		}

		// Should have PUSH instructions for key and value, and MAKE_DICTIONARY
		if len(bytecode.Instructions) < 3 {
			t.Errorf("Expected at least 3 instructions, got %d", len(bytecode.Instructions))
		}
	})
}

// TestVersion0_4_0_CascadingMessages tests cascading message sends
func TestVersion0_4_0_CascadingMessages(t *testing.T) {
	t.Run("ParseSimpleCascade", func(t *testing.T) {
		input := "point x: 10; y: 20"
		p := parser.New(input)
		program, err := p.Parse()

		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		cascade, ok := stmt.Expression.(*ast.CascadeExpression)
		if !ok {
			t.Fatalf("Expected CascadeExpression, got %T", stmt.Expression)
		}

		if len(cascade.Messages) != 2 {
			t.Fatalf("Expected 2 messages, got %d", len(cascade.Messages))
		}

		if cascade.Messages[0].Selector != "x:" {
			t.Errorf("Expected first selector 'x:', got %s", cascade.Messages[0].Selector)
		}

		if cascade.Messages[1].Selector != "y:" {
			t.Errorf("Expected second selector 'y:', got %s", cascade.Messages[1].Selector)
		}
	})

	t.Run("ParseUnaryCascade", func(t *testing.T) {
		input := "obj method1; method2; method3"
		p := parser.New(input)
		program, err := p.Parse()

		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		cascade := stmt.Expression.(*ast.CascadeExpression)

		if len(cascade.Messages) != 3 {
			t.Fatalf("Expected 3 messages, got %d", len(cascade.Messages))
		}
	})

	t.Run("CompileCascade", func(t *testing.T) {
		input := "obj m1; m2"
		p := parser.New(input)
		program, err := p.Parse()

		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		c := compiler.New()
		bytecode, err := c.Compile(program)

		if err != nil {
			t.Fatalf("Compile failed: %v", err)
		}

		// Should have DUP instructions for cascading
		if len(bytecode.Instructions) < 5 {
			t.Errorf("Expected at least 5 instructions, got %d", len(bytecode.Instructions))
		}
	})
}

// TestVersion0_4_0_SelfKeyword tests the self keyword
func TestVersion0_4_0_SelfKeyword(t *testing.T) {
	t.Run("ParseSelf", func(t *testing.T) {
		input := "self"
		p := parser.New(input)
		program, err := p.Parse()

		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		ident, ok := stmt.Expression.(*ast.Identifier)
		if !ok {
			t.Fatalf("Expected Identifier, got %T", stmt.Expression)
		}

		if ident.Name != "self" {
			t.Errorf("Expected identifier 'self', got %s", ident.Name)
		}
	})

	t.Run("ParseSelfMessageSend", func(t *testing.T) {
		input := "self initialize"
		p := parser.New(input)
		program, err := p.Parse()

		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		msgSend := stmt.Expression.(*ast.MessageSend)

		receiver, ok := msgSend.Receiver.(*ast.Identifier)
		if !ok || receiver.Name != "self" {
			t.Error("Expected receiver to be 'self'")
		}

		if msgSend.Selector != "initialize" {
			t.Errorf("Expected selector 'initialize', got %s", msgSend.Selector)
		}
	})
}

// TestVersion0_4_0_Integration tests combinations of v0.4.0 features
func TestVersion0_4_0_Integration(t *testing.T) {
	t.Run("SelfWithCascade", func(t *testing.T) {
		input := "self x: 10; y: 20"
		p := parser.New(input)
		program, err := p.Parse()

		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		cascade, ok := stmt.Expression.(*ast.CascadeExpression)
		if !ok {
			t.Fatalf("Expected CascadeExpression, got %T", stmt.Expression)
		}

		receiver, ok := cascade.Receiver.(*ast.Identifier)
		if !ok || receiver.Name != "self" {
			t.Error("Expected receiver to be 'self'")
		}
	})

	t.Run("DictionaryWithVariousTypes", func(t *testing.T) {
		input := "#{1 -> 'one'. true -> false. 'key' -> 42}"
		p := parser.New(input)
		program, err := p.Parse()

		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		dictLit := stmt.Expression.(*ast.DictionaryLiteral)

		if len(dictLit.Pairs) != 3 {
			t.Errorf("Expected 3 pairs, got %d", len(dictLit.Pairs))
		}
	})
}

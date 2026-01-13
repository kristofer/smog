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

// TestVersion0_4_0_ClassParsing tests parsing of class definitions
func TestVersion0_4_0_ClassParsing(t *testing.T) {
t.Run("ParseSimpleClass", func(t *testing.T) {
input := `Object subclass: #Counter [
| count |

initialize [
count := 0.
]
]`

p := parser.New(input)
program, err := p.Parse()

if err != nil {
t.Fatalf("Parse failed: %v", err)
}

if len(program.Statements) != 1 {
t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
}

class, ok := program.Statements[0].(*ast.Class)
if !ok {
t.Fatalf("Expected Class, got %T", program.Statements[0])
}

if class.Name != "Counter" {
t.Errorf("Expected class name 'Counter', got '%s'", class.Name)
}

if class.SuperClass != "Object" {
t.Errorf("Expected superclass 'Object', got '%s'", class.SuperClass)
}

if len(class.Fields) != 1 || class.Fields[0] != "count" {
t.Errorf("Expected 1 instance variable 'count', got %v", class.Fields)
}

if len(class.Methods) != 1 {
t.Fatalf("Expected 1 method, got %d", len(class.Methods))
}

if class.Methods[0].Name != "initialize" {
t.Errorf("Expected method 'initialize', got '%s'", class.Methods[0].Name)
}
})

t.Run("ParseClassWithClassVariablesAndMethods", func(t *testing.T) {
input := `Object subclass: #Counter [
| count |
<| totalCount |>

initialize [
count := 0.
]

<resetTotal [
totalCount := 0.
]>
]`

p := parser.New(input)
program, err := p.Parse()

if err != nil {
t.Fatalf("Parse failed: %v", err)
}

class := program.Statements[0].(*ast.Class)

if len(class.ClassVariables) != 1 || class.ClassVariables[0] != "totalCount" {
t.Errorf("Expected 1 class variable 'totalCount', got %v", class.ClassVariables)
}

if len(class.ClassMethods) != 1 {
t.Fatalf("Expected 1 class method, got %d", len(class.ClassMethods))
}

if class.ClassMethods[0].Name != "resetTotal" {
t.Errorf("Expected class method 'resetTotal', got '%s'", class.ClassMethods[0].Name)
}
})

t.Run("ParseMultipleClasses", func(t *testing.T) {
input := `Object subclass: #Vehicle [
| speed |

initialize [
speed := 0.
]
]

Vehicle subclass: #Car [
| turboBoost |

initialize [
super initialize.
turboBoost := false.
]
]`

p := parser.New(input)
program, err := p.Parse()

if err != nil {
t.Fatalf("Parse failed: %v", err)
}

if len(program.Statements) != 2 {
t.Fatalf("Expected 2 statements, got %d", len(program.Statements))
}

class1, ok1 := program.Statements[0].(*ast.Class)
class2, ok2 := program.Statements[1].(*ast.Class)

if !ok1 || !ok2 {
t.Fatalf("Expected both statements to be Classes")
}

if class1.Name != "Vehicle" {
t.Errorf("Expected first class 'Vehicle', got '%s'", class1.Name)
}

if class2.Name != "Car" {
t.Errorf("Expected second class 'Car', got '%s'", class2.Name)
}

if class2.SuperClass != "Vehicle" {
t.Errorf("Expected Car's superclass 'Vehicle', got '%s'", class2.SuperClass)
}
})
}

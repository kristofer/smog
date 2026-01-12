package test

import (
	"testing"

	"github.com/kristofer/smog/pkg/compiler"
	"github.com/kristofer/smog/pkg/parser"
	"github.com/kristofer/smog/pkg/vm"
)

// TestVersion0_2_0_Literals tests literal support for version 0.2.0
func TestVersion0_2_0_Literals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{"Integer", "42", int64(42)},
		{"Float", "3.14", 3.14},
		{"String", "'Hello'", "Hello"},
		{"True", "true", true},
		{"False", "false", false},
		{"Nil", "nil", nil},
		{"NegativeInteger", "-17", int64(-17)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New(tt.input)
			program, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			c := compiler.New()
			bc, err := c.Compile(program)
			if err != nil {
				t.Fatalf("Compile error: %v", err)
			}

			vm := vm.New()
			err = vm.Run(bc)
			if err != nil {
				t.Fatalf("VM error: %v", err)
			}

			result := vm.StackTop()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestVersion0_2_0_Variables tests variable declarations and assignments
func TestVersion0_2_0_Variables(t *testing.T) {
	t.Run("SimpleAssignment", func(t *testing.T) {
		input := `| x |
x := 42.
x`

		p := parser.New(input)
		program, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		c := compiler.New()
		bc, err := c.Compile(program)
		if err != nil {
			t.Fatalf("Compile error: %v", err)
		}

		vm := vm.New()
		err = vm.Run(bc)
		if err != nil {
			t.Fatalf("VM error: %v", err)
		}

		result := vm.StackTop()
		if result != int64(42) {
			t.Errorf("Expected 42, got %v", result)
		}
	})

	t.Run("MultipleVariables", func(t *testing.T) {
		input := `| x y z |
x := 10.
y := 20.
z := 30.
z`

		p := parser.New(input)
		program, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		c := compiler.New()
		bc, err := c.Compile(program)
		if err != nil {
			t.Fatalf("Compile error: %v", err)
		}

		vm := vm.New()
		err = vm.Run(bc)
		if err != nil {
			t.Fatalf("VM error: %v", err)
		}

		result := vm.StackTop()
		if result != int64(30) {
			t.Errorf("Expected 30, got %v", result)
		}
	})
}

// TestVersion0_2_0_BinaryMessages tests binary message sends (arithmetic and comparison)
func TestVersion0_2_0_BinaryMessages(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{"Addition", "3 + 4", int64(7)},
		{"Subtraction", "10 - 3", int64(7)},
		{"Multiplication", "3 * 4", int64(12)},
		{"Division", "12 / 3", int64(4)},
		{"LessThan", "3 < 4", true},
		{"GreaterThan", "4 > 3", true},
		{"LessOrEqual", "3 <= 3", true},
		{"GreaterOrEqual", "3 >= 3", true},
		{"Equal", "3 = 3", true},
		{"NotEqual", "3 ~= 4", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New(tt.input)
			program, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			c := compiler.New()
			bc, err := c.Compile(program)
			if err != nil {
				t.Fatalf("Compile error: %v", err)
			}

			vm := vm.New()
			err = vm.Run(bc)
			if err != nil {
				t.Fatalf("VM error: %v", err)
			}

			result := vm.StackTop()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestVersion0_2_0_ComplexExpressions tests combinations of features
func TestVersion0_2_0_ComplexExpressions(t *testing.T) {
	t.Run("ArithmeticWithVariables", func(t *testing.T) {
		input := `| x y |
x := 10.
y := 20.
x + y`

		p := parser.New(input)
		program, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		c := compiler.New()
		bc, err := c.Compile(program)
		if err != nil {
			t.Fatalf("Compile error: %v", err)
		}

		vm := vm.New()
		err = vm.Run(bc)
		if err != nil {
			t.Fatalf("VM error: %v", err)
		}

		result := vm.StackTop()
		if result != int64(30) {
			t.Errorf("Expected 30, got %v", result)
		}
	})

	t.Run("ComparisonWithVariables", func(t *testing.T) {
		input := `| x y |
x := 10.
y := 20.
x < y`

		p := parser.New(input)
		program, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		c := compiler.New()
		bc, err := c.Compile(program)
		if err != nil {
			t.Fatalf("Compile error: %v", err)
		}

		vm := vm.New()
		err = vm.Run(bc)
		if err != nil {
			t.Fatalf("VM error: %v", err)
		}

		result := vm.StackTop()
		if result != true {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("MultipleStatements", func(t *testing.T) {
		input := `| a b c |
a := 5.
b := 10.
c := a + b.
c * 2`

		p := parser.New(input)
		program, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		c := compiler.New()
		bc, err := c.Compile(program)
		if err != nil {
			t.Fatalf("Compile error: %v", err)
		}

		vm := vm.New()
		err = vm.Run(bc)
		if err != nil {
			t.Fatalf("VM error: %v", err)
		}

		result := vm.StackTop()
		if result != int64(30) {
			t.Errorf("Expected 30, got %v", result)
		}
	})
}

// TestVersion0_2_0_EndToEnd tests a complete program flow
func TestVersion0_2_0_EndToEnd(t *testing.T) {
	input := `" Calculate the sum of two numbers "
| x y result |
x := 15.
y := 27.
result := x + y.
result`

	p := parser.New(input)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// Variable declaration + 3 assignments + 1 expression = 5 statements
	if len(program.Statements) != 5 {
		t.Fatalf("Expected 5 statements, got %d", len(program.Statements))
	}

	c := compiler.New()
	bc, err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}

	if len(bc.Instructions) == 0 {
		t.Fatal("Expected bytecode instructions")
	}

	vm := vm.New()
	err = vm.Run(bc)
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	result := vm.StackTop()
	if result != int64(42) {
		t.Errorf("Expected 42, got %v", result)
	}
}

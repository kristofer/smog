package vm

import (
	"testing"

	"github.com/kristofer/smog/pkg/compiler"
	"github.com/kristofer/smog/pkg/parser"
)

func TestVMIntegerLiteral(t *testing.T) {
	input := "42"

	p := parser.New(input)
	program, _ := p.Parse()

	c := compiler.New()
	bc, _ := c.Compile(program)

	vm := New()
	err := vm.Run(bc)
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	result := vm.StackTop()
	if result != int64(42) {
		t.Errorf("Expected 42, got %v", result)
	}
}

func TestVMStringLiteral(t *testing.T) {
	input := "'Hello'"

	p := parser.New(input)
	program, _ := p.Parse()

	c := compiler.New()
	bc, _ := c.Compile(program)

	vm := New()
	err := vm.Run(bc)
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	result := vm.StackTop()
	if result != "Hello" {
		t.Errorf("Expected 'Hello', got %v", result)
	}
}

func TestVMBooleanLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		p := parser.New(tt.input)
		program, _ := p.Parse()

		c := compiler.New()
		bc, _ := c.Compile(program)

		vm := New()
		err := vm.Run(bc)
		if err != nil {
			t.Fatalf("VM error: %v", err)
		}

		result := vm.StackTop()
		if result != tt.expected {
			t.Errorf("Expected %v, got %v", tt.expected, result)
		}
	}
}

func TestVMNilLiteral(t *testing.T) {
	input := "nil"

	p := parser.New(input)
	program, _ := p.Parse()

	c := compiler.New()
	bc, _ := c.Compile(program)

	vm := New()
	err := vm.Run(bc)
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	result := vm.StackTop()
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
}

func TestVMArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"3 + 4", 7},
		{"10 - 3", 7},
		{"3 * 4", 12},
		{"12 / 3", 4},
	}

	for _, tt := range tests {
		p := parser.New(tt.input)
		program, _ := p.Parse()

		c := compiler.New()
		bc, _ := c.Compile(program)

		vm := New()
		err := vm.Run(bc)
		if err != nil {
			t.Fatalf("VM error for %s: %v", tt.input, err)
		}

		result := vm.StackTop()
		if result != tt.expected {
			t.Errorf("For %s: expected %v, got %v", tt.input, tt.expected, result)
		}
	}
}

func TestVMComparison(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"3 < 4", true},
		{"4 < 3", false},
		{"3 > 4", false},
		{"4 > 3", true},
		{"3 <= 3", true},
		{"3 >= 3", true},
		{"3 = 3", true},
		{"3 ~= 4", true},
	}

	for _, tt := range tests {
		p := parser.New(tt.input)
		program, _ := p.Parse()

		c := compiler.New()
		bc, _ := c.Compile(program)

		vm := New()
		err := vm.Run(bc)
		if err != nil {
			t.Fatalf("VM error for %s: %v", tt.input, err)
		}

		result := vm.StackTop()
		if result != tt.expected {
			t.Errorf("For %s: expected %v, got %v", tt.input, tt.expected, result)
		}
	}
}

func TestVMVariableDeclarationAndAssignment(t *testing.T) {
	input := `| x |
x := 42.
x`

	p := parser.New(input)
	program, _ := p.Parse()

	c := compiler.New()
	bc, _ := c.Compile(program)

	vm := New()
	err := vm.Run(bc)
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	result := vm.StackTop()
	if result != int64(42) {
		t.Errorf("Expected 42, got %v", result)
	}
}

func TestVMMultipleStatements(t *testing.T) {
	input := `| x y |
x := 10.
y := 20.
x + y`

	p := parser.New(input)
	program, _ := p.Parse()

	c := compiler.New()
	bc, _ := c.Compile(program)

	vm := New()
	err := vm.Run(bc)
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	result := vm.StackTop()
	if result != int64(30) {
		t.Errorf("Expected 30, got %v", result)
	}
}

func TestVMSimpleBlock(t *testing.T) {
input := "[ 42 ] value"

p := parser.New(input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error: %v", err)
}

result := vm.StackTop()
if result != int64(42) {
t.Errorf("Expected 42, got %v", result)
}
}

func TestVMBlockWithOneParameter(t *testing.T) {
input := "[ :x | x * 2 ] value: 5"

p := parser.New(input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error: %v", err)
}

result := vm.StackTop()
if result != int64(10) {
t.Errorf("Expected 10, got %v", result)
}
}

func TestVMBlockWithTwoParameters(t *testing.T) {
input := "[ :x :y | x + y ] value: 3 value: 7"

p := parser.New(input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error: %v", err)
}

result := vm.StackTop()
if result != int64(10) {
t.Errorf("Expected 10, got %v", result)
}
}

func TestVMArrayLiteral(t *testing.T) {
input := "#(1 2 3) size"

p := parser.New(input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error: %v", err)
}

result := vm.StackTop()
if result != int64(3) {
t.Errorf("Expected 3, got %v", result)
}
}

func TestVMArrayAt(t *testing.T) {
input := "#(10 20 30) at: 2"

p := parser.New(input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error: %v", err)
}

result := vm.StackTop()
if result != int64(20) {
t.Errorf("Expected 20, got %v", result)
}
}

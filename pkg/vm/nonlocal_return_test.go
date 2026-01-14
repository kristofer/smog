package vm

import (
	"testing"

	"github.com/kristofer/smog/pkg/compiler"
	"github.com/kristofer/smog/pkg/parser"
)

// TestNonLocalReturnInBlock tests that a return statement inside a block
// returns from the enclosing method, not just from the block.
func TestNonLocalReturnInBlock(t *testing.T) {
	source := `
Object subclass: #TestClass [
    testMethod [
        (true) ifTrue: [
            ^42
        ].
        'This should not execute' println.
        ^99
    ]
]

| obj result |
obj := TestClass new.
result := obj testMethod.
result
`

	p := parser.New(source)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	c := compiler.New()
	bc, err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}

	vm := New()
	err = vm.Run(bc)
	if err != nil {
		t.Fatalf("Runtime error: %v", err)
	}

	result := vm.StackTop()
	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// The result should be 42 from the non-local return,
	// not 99 which would be returned if execution continued
	resultInt, ok := result.(int64)
	if !ok {
		t.Fatalf("Expected int64 result, got %T", result)
	}

	if resultInt != 42 {
		t.Errorf("Expected 42, got %d", resultInt)
	}
}

// TestNonLocalReturnInNestedBlocks tests that non-local return works
// through multiple levels of block nesting.
func TestNonLocalReturnInNestedBlocks(t *testing.T) {
	source := `
Object subclass: #TestClass [
    testMethod [
        (true) ifTrue: [
            (true) ifTrue: [
                ^123
            ]
        ].
        ^456
    ]
]

| obj result |
obj := TestClass new.
result := obj testMethod.
result
`

	p := parser.New(source)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	c := compiler.New()
	bc, err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}

	vm := New()
	err = vm.Run(bc)
	if err != nil {
		t.Fatalf("Runtime error: %v", err)
	}

	result := vm.StackTop()
	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	resultInt, ok := result.(int64)
	if !ok {
		t.Fatalf("Expected int64 result, got %T", result)
	}

	if resultInt != 123 {
		t.Errorf("Expected 123, got %d", resultInt)
	}
}

// TestLocalReturnInMethod tests that a return statement in a method
// (not in a block) still works as expected.
func TestLocalReturnInMethod(t *testing.T) {
	source := `
Object subclass: #TestClass [
    testMethod [
        ^77
        'This should not execute' println.
        ^88
    ]
]

| obj result |
obj := TestClass new.
result := obj testMethod.
result
`

	p := parser.New(source)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	c := compiler.New()
	bc, err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}

	vm := New()
	err = vm.Run(bc)
	if err != nil {
		t.Fatalf("Runtime error: %v", err)
	}

	result := vm.StackTop()
	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	resultInt, ok := result.(int64)
	if !ok {
		t.Fatalf("Expected int64 result, got %T", result)
	}

	if resultInt != 77 {
		t.Errorf("Expected 77, got %d", resultInt)
	}
}

// TestNonLocalReturnInIfFalse tests non-local return in ifFalse: block.
func TestNonLocalReturnInIfFalse(t *testing.T) {
	source := `
Object subclass: #TestClass [
    testMethod [
        (false) ifTrue: [
            ^11
        ] ifFalse: [
            ^22
        ].
        ^33
    ]
]

| obj result |
obj := TestClass new.
result := obj testMethod.
result
`

	p := parser.New(source)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	c := compiler.New()
	bc, err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}

	vm := New()
	err = vm.Run(bc)
	if err != nil {
		t.Fatalf("Runtime error: %v", err)
	}

	result := vm.StackTop()
	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	resultInt, ok := result.(int64)
	if !ok {
		t.Fatalf("Expected int64 result, got %T", result)
	}

	if resultInt != 22 {
		t.Errorf("Expected 22, got %d", resultInt)
	}
}

// TestNonLocalReturnDoesNotAffectOtherMethods tests that a non-local
// return in one method doesn't affect execution in other methods.
func TestNonLocalReturnDoesNotAffectOtherMethods(t *testing.T) {
	source := `
Object subclass: #TestClass [
    method1 [
        (true) ifTrue: [ ^10 ].
        ^20
    ]
    
    method2 [
        (false) ifTrue: [ ^30 ].
        ^40
    ]
]

| obj r1 r2 |
obj := TestClass new.
r1 := obj method1.
r2 := obj method2.
r1 + r2
`

	p := parser.New(source)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	c := compiler.New()
	bc, err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}

	vm := New()
	err = vm.Run(bc)
	if err != nil {
		t.Fatalf("Runtime error: %v", err)
	}

	result := vm.StackTop()
	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	resultInt, ok := result.(int64)
	if !ok {
		t.Fatalf("Expected int64 result, got %T", result)
	}

	// method1 returns 10, method2 returns 40, so 10 + 40 = 50
	if resultInt != 50 {
		t.Errorf("Expected 50, got %d", resultInt)
	}
}

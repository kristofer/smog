package vm

import (
	"strings"
	"testing"

	"github.com/kristofer/smog/pkg/compiler"
	"github.com/kristofer/smog/pkg/parser"
)

// TestStackTraceOnError tests that runtime errors include stack trace information
func TestStackTraceOnError(t *testing.T) {
	source := `
| x y |
x := 10.
y := 0.
x / y
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
	if err == nil {
		t.Fatal("Expected division by zero error, got nil")
	}

	t.Logf("Error type: %T", err)
	t.Logf("Error value: %#v", err)

	// Check if error is a RuntimeError with stack trace
	runtimeErr, ok := err.(*RuntimeError)
	if !ok {
		t.Fatalf("Expected RuntimeError, got %T: %v", err, err)
	}

	// Verify the error message contains the expected text
	errMsg := runtimeErr.Error()
	if !strings.Contains(errMsg, "division by zero") {
		t.Errorf("Expected error message to contain 'division by zero', got: %v", errMsg)
	}

	// Verify stack trace is present
	if !strings.Contains(errMsg, "Stack trace:") {
		t.Errorf("Expected stack trace in error message, got: %v", errMsg)
	}
}

// TestStackTraceWithNestedCalls tests stack traces with nested message sends
func TestStackTraceWithNestedCalls(t *testing.T) {
	source := `
Object subclass: #TestClass [
    method1 [
        ^self method2
    ]
    
    method2 [
        ^self method3
    ]
    
    method3 [
        ^1 / 0
    ]
]

| obj |
obj := TestClass new.
obj method1
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
	if err == nil {
		t.Fatal("Expected division by zero error, got nil")
	}

	// Check if error is a RuntimeError with stack trace
	runtimeErr, ok := err.(*RuntimeError)
	if !ok {
		t.Fatalf("Expected RuntimeError, got %T: %v", err, err)
	}

	// Verify the error message contains expected information
	errMsg := runtimeErr.Error()
	if !strings.Contains(errMsg, "division by zero") {
		t.Errorf("Expected error message to contain 'division by zero', got: %v", errMsg)
	}

	// Verify stack trace is present
	if !strings.Contains(errMsg, "Stack trace:") {
		t.Errorf("Expected stack trace in error message, got: %v", errMsg)
	}

	// The stack trace should show nested calls
	// (exact format may vary, but it should have multiple frames)
	if len(runtimeErr.StackTrace) == 0 {
		t.Error("Expected non-empty stack trace")
	}
}

// TestNoStackTraceOnSuccess tests that successful execution doesn't create stack traces
func TestNoStackTraceOnSuccess(t *testing.T) {
	source := `
| x y |
x := 10.
y := 2.
x / y
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
		t.Fatalf("Expected successful execution, got error: %v", err)
	}

	result := vm.StackTop()
	if result != int64(5) {
		t.Errorf("Expected result 5, got %v", result)
	}
}

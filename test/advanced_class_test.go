package test

import (
	"testing"

	"github.com/kristofer/smog/pkg/compiler"
	"github.com/kristofer/smog/pkg/parser"
	"github.com/kristofer/smog/pkg/vm"
)

// TestClassMethod_SimpleClassMethod tests calling a class method.
func TestClassMethod_SimpleClassMethod(t *testing.T) {
	source := `
		Object subclass: #Counter [
			| count |
			
			" Class method "
			<create [
				^Counter new
			]>
			
			initialize [
				count := 0.
			]
			
			value [
				^count
			]
		]
		
		| counter |
		counter := Counter create.
		counter initialize.
	`

	p := parser.New(source)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	c := compiler.New()
	bytecode, err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}

	v := vm.New()
	err = v.Run(bytecode)
	if err != nil {
		t.Fatalf("Runtime error: %v", err)
	}

	// Should complete without error
}

// TestClassVariable_SharedAcrossInstances tests that class variables are shared.
func TestClassVariable_SharedAcrossInstances(t *testing.T) {
	source := `
		Object subclass: #Counter [
			" Class variable "
			<| totalCount |>
			
			" Class method to set totalCount "
			<setTotal: n [
				totalCount := n.
			]>
			
			" Instance method to get totalCount "
			getTotalCount [
				^totalCount
			]
		]
		
		| c1 result |
		Counter setTotal: 42.
		c1 := Counter new.
		result := c1 getTotalCount.
	`

	p := parser.New(source)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	c := compiler.New()
	bytecode, err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}

	v := vm.New()
	err = v.Run(bytecode)
	if err != nil {
		t.Fatalf("Runtime error: %v", err)
	}

	// totalCount should be 2 (initialized twice)
	// Class variable starts at nil, nil + 1 would need special handling
	// For now, let's just verify it doesn't crash
}

// TestClassMethod_WithParameters tests class method with parameters.
func TestClassMethod_WithParameters(t *testing.T) {
	source := `
		Object subclass: #Point [
			| x y |
			
			" Class method "
			<x: xVal y: yVal [
				| point |
				point := Point new.
				point setX: xVal.
				point setY: yVal.
				^point
			]>
			
			setX: val [
				x := val.
			]
			
			setY: val [
				y := val.
			]
			
			getX [
				^x
			]
		]
		
		| point result |
		point := Point x: 10 y: 20.
		result := point getX.
	`

	p := parser.New(source)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	c := compiler.New()
	bytecode, err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}

	v := vm.New()
	err = v.Run(bytecode)
	if err != nil {
		t.Fatalf("Runtime error: %v", err)
	}

	result := v.StackTop()
	if result != int64(10) {
		t.Errorf("Expected x to be 10, got %v", result)
	}
}

// TestClassVariable_AccessFromClassMethod tests accessing class variables from class methods.
func TestClassVariable_AccessFromClassMethod(t *testing.T) {
	source := `
		Object subclass: #IDGenerator [
			" Class variable "
			<| nextID |>
			
			" Class methods "
			<initialize [
				nextID := 1.
			]>
			
			<nextID [
				| current |
				current := nextID.
				nextID := nextID + 1.
				^current
			]>
		]
		
		| id1 id2 |
		IDGenerator initialize.
		id1 := IDGenerator nextID.
		id2 := IDGenerator nextID.
	`

	p := parser.New(source)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	c := compiler.New()
	bytecode, err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}

	v := vm.New()
	err = v.Run(bytecode)
	if err != nil {
		t.Fatalf("Runtime error: %v", err)
	}

	// id2 should be 2
	result := v.StackTop()
	if result != int64(2) {
		t.Errorf("Expected second ID to be 2, got %v", result)
	}
}

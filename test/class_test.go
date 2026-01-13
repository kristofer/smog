package test

import (
	"testing"

	"github.com/kristofer/smog/pkg/compiler"
	"github.com/kristofer/smog/pkg/parser"
	"github.com/kristofer/smog/pkg/vm"
)

// TestSimpleClassDefinition tests defining a simple class.
func TestSimpleClassDefinition(t *testing.T) {
	source := `
		Object subclass: #Counter [
			| count |
			
			initialize [
				count := 0.
			]
			
			value [
				^count
			]
		]
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

	// Check that Counter class was registered
	if v.GetGlobal("Counter") == nil {
		t.Fatal("Counter class not registered as global")
	}
}

// TestClassInstantiation tests creating an instance of a class.
func TestClassInstantiation(t *testing.T) {
	source := `
		Object subclass: #Counter [
			| count |
			
			initialize [
				count := 0.
			]
		]
		
		| counter |
		counter := Counter new.
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

	// The instance should be on the stack or in a local variable
	// Just verify it doesn't crash
}

// TestMethodCall tests calling a method on an instance.
func TestMethodCall(t *testing.T) {
	source := `
		Object subclass: #Counter [
			| count |
			
			initialize [
				count := 0.
			]
			
			value [
				^count
			]
		]
		
		| counter result |
		counter := Counter new.
		counter initialize.
		result := counter value.
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

	// The result should be 0 (from counter.value)
	result := v.StackTop()
	if result != int64(0) {
		t.Errorf("Expected counter value to be 0, got %v", result)
	}
}

// TestMethodWithModification tests modifying instance variables.
func TestMethodWithModification(t *testing.T) {
	source := `
		Object subclass: #Counter [
			| count |
			
			initialize [
				count := 0.
			]
			
			increment [
				count := count + 1.
			]
			
			value [
				^count
			]
		]
		
		| counter result |
		counter := Counter new.
		counter initialize.
		counter increment.
		counter increment.
		counter increment.
		result := counter value.
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

	// The result should be 3 (incremented 3 times)
	result := v.StackTop()
	if result != int64(3) {
		t.Errorf("Expected counter value to be 3, got %v", result)
	}
}

// TestMethodWithParameters tests methods with parameters.
func TestMethodWithParameters(t *testing.T) {
	source := `
		Object subclass: #Adder [
			| value |
			
			setValue: v [
				value := v.
			]
			
			add: n [
				value := value + n.
			]
			
			getValue [
				^value
			]
		]
		
		| adder result |
		adder := Adder new.
		adder setValue: 10.
		adder add: 5.
		result := adder getValue.
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

	// The result should be 15 (10 + 5)
	result := v.StackTop()
	if result != int64(15) {
		t.Errorf("Expected adder value to be 15, got %v", result)
	}
}

// TestMultipleInstances tests creating multiple instances of a class.
func TestMultipleInstances(t *testing.T) {
	source := `
		Object subclass: #Counter [
			| count |
			
			initialize [
				count := 0.
			]
			
			increment [
				count := count + 1.
			]
			
			value [
				^count
			]
		]
		
		| counter1 counter2 result |
		counter1 := Counter new.
		counter1 initialize.
		counter1 increment.
		counter1 increment.
		
		counter2 := Counter new.
		counter2 initialize.
		counter2 increment.
		
		result := counter2 value.
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

	// counter2 value should be 1 (not affected by counter1's increments)
	result := v.StackTop()
	if result != int64(1) {
		t.Errorf("Expected counter2 value to be 1, got %v", result)
	}
}

// TestMultipleFields tests a class with multiple instance variables.
func TestMultipleFields(t *testing.T) {
	source := `
		Object subclass: #Point [
			| x y |
			
			x: xValue y: yValue [
				x := xValue.
				y := yValue.
			]
			
			getX [
				^x
			]
			
			getY [
				^y
			]
		]
		
		| point result |
		point := Point new.
		point x: 10 y: 20.
		result := point getY.
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

	// The result should be 20 (the y value)
	result := v.StackTop()
	if result != int64(20) {
		t.Errorf("Expected point y to be 20, got %v", result)
	}
}

// TestCompleteCounterWorkflow tests a complete Counter workflow.
func TestCompleteCounterWorkflow(t *testing.T) {
	source := `
		Object subclass: #Counter [
			| count |
			
			initialize [
				count := 0.
			]
			
			increment [
				count := count + 1.
			]
			
			decrement [
				count := count - 1.
			]
			
			value [
				^count
			]
			
			reset [
				count := 0.
			]
		]
		
		| counter val1 val2 val3 |
		counter := Counter new.
		counter initialize.
		
		" Test initial value "
		val1 := counter value.
		
		" Test increment "
		counter increment.
		counter increment.
		counter increment.
		val2 := counter value.
		
		" Test decrement "
		counter decrement.
		val3 := counter value.
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

	// val3 should be 2 (incremented 3 times, decremented once)
	result := v.StackTop()
	if result != int64(2) {
		t.Errorf("Expected counter value to be 2, got %v", result)
	}
}

package test

import (
	"testing"

	"github.com/kristofer/smog/pkg/compiler"
	"github.com/kristofer/smog/pkg/parser"
	"github.com/kristofer/smog/pkg/vm"
)

// TestInheritance_MethodOverride tests method override in subclass.
func TestInheritance_MethodOverride(t *testing.T) {
	source := `
		Object subclass: #Animal [
			| name |
			
			setName: n [
				name := n.
			]
			
			speak [
				^'Some sound'
			]
		]
		
		Animal subclass: #Dog [
			speak [
				^'Woof!'
			]
		]
		
		| dog result |
		dog := Dog new.
		result := dog speak.
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
	if result != "Woof!" {
		t.Errorf("Expected dog to say 'Woof!', got %v", result)
	}
}

// TestInheritance_InheritedMethod tests calling a method from parent class.
func TestInheritance_InheritedMethod(t *testing.T) {
	source := `
		Object subclass: #Animal [
			| name |
			
			setName: n [
				name := n.
			]
			
			getName [
				^name
			]
		]
		
		Animal subclass: #Dog [
			bark [
				^'Woof!'
			]
		]
		
		| dog result |
		dog := Dog new.
		dog setName: 'Buddy'.
		result := dog getName.
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
	if result != "Buddy" {
		t.Errorf("Expected dog name to be 'Buddy', got %v", result)
	}
}

// TestInheritance_SuperSend tests super message sends.
func TestInheritance_SuperSend(t *testing.T) {
	source := `
		Object subclass: #Vehicle [
			| speed |
			
			initialize [
				speed := 0.
			]
			
			accelerate [
				speed := speed + 10.
				^speed
			]
		]
		
		Vehicle subclass: #Car [
			| turbo |
			
			initialize [
				super initialize.
				turbo := 5.
			]
			
			accelerate [
				| baseSpeed |
				baseSpeed := super accelerate.
				^baseSpeed + turbo
			]
		]
		
		| car result |
		car := Car new.
		car initialize.
		result := car accelerate.
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
	// Vehicle accelerate returns 10, Car adds 5 = 15
	if result != int64(15) {
		t.Errorf("Expected car speed to be 15, got %v", result)
	}
}

// TestInheritance_ThreeLevelHierarchy tests inheritance through three levels.
func TestInheritance_ThreeLevelHierarchy(t *testing.T) {
	source := `
		Object subclass: #Animal [
			| name |
			
			setName: n [
				name := n.
			]
			
			getName [
				^name
			]
		]
		
		Animal subclass: #Mammal [
			| furColor |
			
			setFurColor: c [
				furColor := c.
			]
			
			getFurColor [
				^furColor
			]
		]
		
		Mammal subclass: #Dog [
			bark [
				^'Woof!'
			]
		]
		
		| dog name color |
		dog := Dog new.
		dog setName: 'Buddy'.
		dog setFurColor: 'brown'.
		name := dog getName.
		color := dog getFurColor.
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

	// Check the color on stack
	result := v.StackTop()
	if result != "brown" {
		t.Errorf("Expected furColor to be 'brown', got %v", result)
	}
}

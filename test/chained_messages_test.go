package test

import (
	"testing"

	"github.com/kristofer/smog/pkg/compiler"
	"github.com/kristofer/smog/pkg/parser"
	"github.com/kristofer/smog/pkg/vm"
)

// TestChainedMessageSends tests chaining method calls.
func TestChainedMessageSends(t *testing.T) {
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
		result := counter value.
		result println.
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

	// Should print 1 and complete without error
}

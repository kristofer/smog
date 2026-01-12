package compiler

import (
	"testing"

	"github.com/kristofer/smog/pkg/bytecode"
	"github.com/kristofer/smog/pkg/parser"
)

func TestCompileIntegerLiteral(t *testing.T) {
	input := "42"

	p := parser.New(input)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	c := New()
	bc, err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// Should have: PUSH constant, RETURN
	if len(bc.Instructions) != 2 {
		t.Fatalf("Expected 2 instructions, got %d", len(bc.Instructions))
	}

	if bc.Instructions[0].Op != bytecode.OpPush {
		t.Errorf("Expected PUSH instruction, got %v", bc.Instructions[0].Op)
	}

	if bc.Instructions[1].Op != bytecode.OpReturn {
		t.Errorf("Expected RETURN instruction, got %v", bc.Instructions[1].Op)
	}

	// Check constant pool
	if len(bc.Constants) != 1 {
		t.Fatalf("Expected 1 constant, got %d", len(bc.Constants))
	}

	if bc.Constants[0] != int64(42) {
		t.Errorf("Expected constant 42, got %v", bc.Constants[0])
	}
}

func TestCompileStringLiteral(t *testing.T) {
	input := "'Hello'"

	p := parser.New(input)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	c := New()
	bc, err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	if len(bc.Instructions) != 2 {
		t.Fatalf("Expected 2 instructions, got %d", len(bc.Instructions))
	}

	if bc.Instructions[0].Op != bytecode.OpPush {
		t.Errorf("Expected PUSH instruction, got %v", bc.Instructions[0].Op)
	}

	if bc.Constants[0] != "Hello" {
		t.Errorf("Expected constant 'Hello', got %v", bc.Constants[0])
	}
}

func TestCompileBooleanLiterals(t *testing.T) {
	tests := []struct {
		input      string
		expectedOp bytecode.Opcode
	}{
		{"true", bytecode.OpPushTrue},
		{"false", bytecode.OpPushFalse},
	}

	for _, tt := range tests {
		p := parser.New(tt.input)
		program, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		c := New()
		bc, err := c.Compile(program)
		if err != nil {
			t.Fatalf("Compile failed: %v", err)
		}

		if len(bc.Instructions) != 2 {
			t.Fatalf("Expected 2 instructions, got %d", len(bc.Instructions))
		}

		if bc.Instructions[0].Op != tt.expectedOp {
			t.Errorf("Expected %v instruction, got %v", tt.expectedOp, bc.Instructions[0].Op)
		}
	}
}

func TestCompileNilLiteral(t *testing.T) {
	input := "nil"

	p := parser.New(input)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	c := New()
	bc, err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	if bc.Instructions[0].Op != bytecode.OpPushNil {
		t.Errorf("Expected PUSH_NIL instruction, got %v", bc.Instructions[0].Op)
	}
}

func TestCompileVariableDeclarationAndAssignment(t *testing.T) {
	input := `| x |
x := 42`

	p := parser.New(input)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	c := New()
	bc, err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// Should have: PUSH 42, STORE_LOCAL 0, RETURN
	if len(bc.Instructions) != 3 {
		t.Fatalf("Expected 3 instructions, got %d", len(bc.Instructions))
	}

	if bc.Instructions[0].Op != bytecode.OpPush {
		t.Errorf("Expected PUSH instruction, got %v", bc.Instructions[0].Op)
	}

	if bc.Instructions[1].Op != bytecode.OpStoreLocal {
		t.Errorf("Expected STORE_LOCAL instruction, got %v", bc.Instructions[1].Op)
	}

	if bc.Instructions[1].Operand != 0 {
		t.Errorf("Expected operand 0, got %d", bc.Instructions[1].Operand)
	}
}

func TestCompileUnaryMessageSend(t *testing.T) {
	input := "'Hello' println"

	p := parser.New(input)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	c := New()
	bc, err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// Should have: PUSH "Hello", SEND println, RETURN
	if len(bc.Instructions) != 3 {
		t.Fatalf("Expected 3 instructions, got %d", len(bc.Instructions))
	}

	if bc.Instructions[0].Op != bytecode.OpPush {
		t.Errorf("Expected PUSH instruction, got %v", bc.Instructions[0].Op)
	}

	if bc.Instructions[1].Op != bytecode.OpSend {
		t.Errorf("Expected SEND instruction, got %v", bc.Instructions[1].Op)
	}

	// Check that "println" is in the constants
	found := false
	for _, c := range bc.Constants {
		if c == "println" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected 'println' in constants")
	}
}

func TestCompileBinaryMessageSend(t *testing.T) {
	input := "3 + 4"

	p := parser.New(input)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	c := New()
	bc, err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// Should have: PUSH 3, PUSH 4, SEND +, RETURN
	if len(bc.Instructions) != 4 {
		t.Fatalf("Expected 4 instructions, got %d", len(bc.Instructions))
	}

	if bc.Instructions[0].Op != bytecode.OpPush {
		t.Errorf("Expected first PUSH instruction, got %v", bc.Instructions[0].Op)
	}

	if bc.Instructions[1].Op != bytecode.OpPush {
		t.Errorf("Expected second PUSH instruction, got %v", bc.Instructions[1].Op)
	}

	if bc.Instructions[2].Op != bytecode.OpSend {
		t.Errorf("Expected SEND instruction, got %v", bc.Instructions[2].Op)
	}

	// Check constants
	if bc.Constants[0] != int64(3) {
		t.Errorf("Expected constant 3, got %v", bc.Constants[0])
	}

	if bc.Constants[1] != int64(4) {
		t.Errorf("Expected constant 4, got %v", bc.Constants[1])
	}
}

func TestCompileKeywordMessageSend(t *testing.T) {
	input := "point x: 10 y: 20"

	p := parser.New(input)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	c := New()
	bc, err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// Should have: LOAD_GLOBAL point, PUSH 10, PUSH 20, SEND x:y:, RETURN
	if len(bc.Instructions) != 5 {
		t.Fatalf("Expected 5 instructions, got %d", len(bc.Instructions))
	}

	if bc.Instructions[0].Op != bytecode.OpLoadGlobal {
		t.Errorf("Expected LOAD_GLOBAL instruction, got %v", bc.Instructions[0].Op)
	}

	if bc.Instructions[1].Op != bytecode.OpPush {
		t.Errorf("Expected first PUSH instruction, got %v", bc.Instructions[1].Op)
	}

	if bc.Instructions[2].Op != bytecode.OpPush {
		t.Errorf("Expected second PUSH instruction, got %v", bc.Instructions[2].Op)
	}

	if bc.Instructions[3].Op != bytecode.OpSend {
		t.Errorf("Expected SEND instruction, got %v", bc.Instructions[3].Op)
	}
}

func TestCompileMultipleStatements(t *testing.T) {
	input := `42.
'hello'.
true.`

	p := parser.New(input)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	c := New()
	bc, err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// Should have: PUSH 42, PUSH "hello", PUSH_TRUE, RETURN
	if len(bc.Instructions) != 4 {
		t.Fatalf("Expected 4 instructions, got %d", len(bc.Instructions))
	}

	if bc.Instructions[0].Op != bytecode.OpPush {
		t.Errorf("Expected first PUSH instruction, got %v", bc.Instructions[0].Op)
	}

	if bc.Instructions[1].Op != bytecode.OpPush {
		t.Errorf("Expected second PUSH instruction, got %v", bc.Instructions[1].Op)
	}

	if bc.Instructions[2].Op != bytecode.OpPushTrue {
		t.Errorf("Expected PUSH_TRUE instruction, got %v", bc.Instructions[2].Op)
	}

	if bc.Instructions[3].Op != bytecode.OpReturn {
		t.Errorf("Expected RETURN instruction, got %v", bc.Instructions[3].Op)
	}
}

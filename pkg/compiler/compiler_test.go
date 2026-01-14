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

	// Should have: PUSH 42, POP, PUSH "hello", POP, PUSH_TRUE, RETURN
	// (POP after each non-last expression statement to clean stack)
	if len(bc.Instructions) != 6 {
		t.Fatalf("Expected 6 instructions, got %d", len(bc.Instructions))
	}

	if bc.Instructions[0].Op != bytecode.OpPush {
		t.Errorf("Expected first PUSH instruction, got %v", bc.Instructions[0].Op)
	}

	if bc.Instructions[1].Op != bytecode.OpPop {
		t.Errorf("Expected first POP instruction, got %v", bc.Instructions[1].Op)
	}

	if bc.Instructions[2].Op != bytecode.OpPush {
		t.Errorf("Expected second PUSH instruction, got %v", bc.Instructions[2].Op)
	}

	if bc.Instructions[3].Op != bytecode.OpPop {
		t.Errorf("Expected second POP instruction, got %v", bc.Instructions[3].Op)
	}

	if bc.Instructions[4].Op != bytecode.OpPushTrue {
		t.Errorf("Expected PUSH_TRUE instruction, got %v", bc.Instructions[4].Op)
	}

	if bc.Instructions[5].Op != bytecode.OpReturn {
		t.Errorf("Expected RETURN instruction, got %v", bc.Instructions[5].Op)
	}
}

func TestCompileSimpleBlock(t *testing.T) {
input := "[ 42 ]"

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

// Should have: MAKE_CLOSURE, RETURN
if len(bc.Instructions) != 2 {
t.Fatalf("Expected 2 instructions, got %d", len(bc.Instructions))
}

if bc.Instructions[0].Op != bytecode.OpMakeClosure {
t.Errorf("Expected MAKE_CLOSURE instruction, got %v", bc.Instructions[0].Op)
}

// Check that block bytecode is in constants
if len(bc.Constants) < 1 {
t.Fatalf("Expected at least 1 constant (block bytecode), got %d", len(bc.Constants))
}

blockBC, ok := bc.Constants[0].(*bytecode.Bytecode)
if !ok {
t.Fatalf("Expected first constant to be Bytecode, got %T", bc.Constants[0])
}

// Block should have: PUSH 42, RETURN
if len(blockBC.Instructions) != 2 {
t.Errorf("Expected 2 instructions in block, got %d", len(blockBC.Instructions))
}
}

func TestCompileBlockWithParameter(t *testing.T) {
input := "[ :x | x + 1 ]"

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

// Should have: MAKE_CLOSURE, RETURN
if len(bc.Instructions) < 1 {
t.Fatalf("Expected at least 1 instruction, got %d", len(bc.Instructions))
}

if bc.Instructions[0].Op != bytecode.OpMakeClosure {
t.Errorf("Expected MAKE_CLOSURE instruction, got %v", bc.Instructions[0].Op)
}

// Check parameter count is encoded in operand (low 8 bits)
paramCount := bc.Instructions[0].Operand & bytecode.ArgCountMask
if paramCount != 1 {
t.Errorf("Expected 1 parameter, got %d", paramCount)
}
}

func TestCompileArrayLiteral(t *testing.T) {
input := "#(1 2 3)"

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

// Should have: PUSH 1, PUSH 2, PUSH 3, MAKE_ARRAY 3, RETURN
if len(bc.Instructions) != 5 {
t.Fatalf("Expected 5 instructions, got %d", len(bc.Instructions))
}

// Check for three PUSH instructions
for i := 0; i < 3; i++ {
if bc.Instructions[i].Op != bytecode.OpPush {
t.Errorf("Expected PUSH instruction at index %d, got %v", i, bc.Instructions[i].Op)
}
}

// Check MAKE_ARRAY instruction
if bc.Instructions[3].Op != bytecode.OpMakeArray {
t.Errorf("Expected MAKE_ARRAY instruction, got %v", bc.Instructions[3].Op)
}

// Check element count
if bc.Instructions[3].Operand != 3 {
t.Errorf("Expected MAKE_ARRAY operand 3, got %d", bc.Instructions[3].Operand)
}
}

// TestCompileIncremental tests that CompileIncremental preserves the symbol table
// across multiple compilations, which is needed for REPL functionality.
func TestCompileIncremental(t *testing.T) {
c := New()

// First compilation: declare a variable
input1 := "| x |"
p1 := parser.New(input1)
program1, err := p1.Parse()
if err != nil {
t.Fatalf("Parse failed for input1: %v", err)
}

bc1, err := c.CompileIncremental(program1)
if err != nil {
t.Fatalf("CompileIncremental failed for input1: %v", err)
}

// First bytecode should just have RETURN (variable declaration doesn't generate code)
if len(bc1.Instructions) != 1 || bc1.Instructions[0].Op != bytecode.OpReturn {
t.Logf("First compilation generated %d instructions", len(bc1.Instructions))
}

// Check that x is in the symbol table (slot 0)
// We can't directly check the symbol table since it's unexported,
// but we can verify behavior by compiling code that uses x

// Second compilation: assign to x
input2 := "x := 42."
p2 := parser.New(input2)
program2, err := p2.Parse()
if err != nil {
t.Fatalf("Parse failed for input2: %v", err)
}

bc2, err := c.CompileIncremental(program2)
if err != nil {
t.Fatalf("CompileIncremental failed for input2: %v", err)
}

// Should have: PUSH 42, STORE_LOCAL 0, RETURN
if len(bc2.Instructions) != 3 {
t.Fatalf("Expected 3 instructions in bc2, got %d", len(bc2.Instructions))
}

if bc2.Instructions[0].Op != bytecode.OpPush {
t.Errorf("Expected PUSH instruction, got %v", bc2.Instructions[0].Op)
}

if bc2.Instructions[1].Op != bytecode.OpStoreLocal {
t.Errorf("Expected STORE_LOCAL instruction, got %v", bc2.Instructions[1].Op)
}

// The operand should be 0 (slot for x)
if bc2.Instructions[1].Operand != 0 {
t.Errorf("Expected STORE_LOCAL to use slot 0, got %d", bc2.Instructions[1].Operand)
}

// Third compilation: read x
input3 := "x println."
p3 := parser.New(input3)
program3, err := p3.Parse()
if err != nil {
t.Fatalf("Parse failed for input3: %v", err)
}

bc3, err := c.CompileIncremental(program3)
if err != nil {
t.Fatalf("CompileIncremental failed for input3: %v", err)
}

// Should have: LOAD_LOCAL 0, PUSH "println", SEND, RETURN
if bc3.Instructions[0].Op != bytecode.OpLoadLocal {
t.Errorf("Expected LOAD_LOCAL instruction, got %v", bc3.Instructions[0].Op)
}

// The operand should be 0 (slot for x)
if bc3.Instructions[0].Operand != 0 {
t.Errorf("Expected LOAD_LOCAL to use slot 0, got %d", bc3.Instructions[0].Operand)
}
}

// TestCompileIncrementalMultipleVars tests that multiple local variables
// are tracked correctly across incremental compilations.
func TestCompileIncrementalMultipleVars(t *testing.T) {
c := New()

// Declare two variables
input1 := "| x y |"
p1 := parser.New(input1)
program1, err := p1.Parse()
if err != nil {
t.Fatalf("Parse failed for input1: %v", err)
}

_, err = c.CompileIncremental(program1)
if err != nil {
t.Fatalf("CompileIncremental failed for input1: %v", err)
}

// Assign to both variables
input2 := "x := 10. y := 20."
p2 := parser.New(input2)
program2, err := p2.Parse()
if err != nil {
t.Fatalf("Parse failed for input2: %v", err)
}

bc2, err := c.CompileIncremental(program2)
if err != nil {
t.Fatalf("CompileIncremental failed for input2: %v", err)
}

// Should use STORE_LOCAL with slots 0 and 1
storeCount := 0
for _, inst := range bc2.Instructions {
if inst.Op == bytecode.OpStoreLocal {
storeCount++
if inst.Operand != 0 && inst.Operand != 1 {
t.Errorf("Expected STORE_LOCAL to use slot 0 or 1, got %d", inst.Operand)
}
}
}

if storeCount != 2 {
t.Errorf("Expected 2 STORE_LOCAL instructions, got %d", storeCount)
}

// Use both variables
input3 := "x + y."
p3 := parser.New(input3)
program3, err := p3.Parse()
if err != nil {
t.Fatalf("Parse failed for input3: %v", err)
}

bc3, err := c.CompileIncremental(program3)
if err != nil {
t.Fatalf("CompileIncremental failed for input3: %v", err)
}

// Should have two LOAD_LOCAL instructions
loadCount := 0
for _, inst := range bc3.Instructions {
if inst.Op == bytecode.OpLoadLocal {
loadCount++
if inst.Operand != 0 && inst.Operand != 1 {
t.Errorf("Expected LOAD_LOCAL to use slot 0 or 1, got %d", inst.Operand)
}
}
}

if loadCount != 2 {
t.Errorf("Expected 2 LOAD_LOCAL instructions, got %d", loadCount)
}
}

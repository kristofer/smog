package bytecode

import (
	"bytes"
	"testing"
)

// TestEncodeDecodeSimpleBytecode tests round-trip encoding and decoding
// of basic bytecode with simple instructions and constants.
func TestEncodeDecodeSimpleBytecode(t *testing.T) {
	// Create a simple bytecode: PUSH 42, RETURN
	original := &Bytecode{
		Instructions: []Instruction{
			{Op: OpPush, Operand: 0},
			{Op: OpReturn, Operand: 0},
		},
		Constants: []interface{}{
			int64(42),
		},
	}

	// Encode to bytes
	var buf bytes.Buffer
	if err := Encode(original, &buf); err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Verify something was written
	if buf.Len() == 0 {
		t.Fatal("No data was encoded")
	}

	// Decode back
	decoded, err := Decode(&buf)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// Verify instructions match
	if len(decoded.Instructions) != len(original.Instructions) {
		t.Fatalf("Instruction count mismatch: got %d, want %d",
			len(decoded.Instructions), len(original.Instructions))
	}

	for i, instr := range decoded.Instructions {
		if instr.Op != original.Instructions[i].Op {
			t.Errorf("Instruction %d opcode mismatch: got %v, want %v",
				i, instr.Op, original.Instructions[i].Op)
		}
		if instr.Operand != original.Instructions[i].Operand {
			t.Errorf("Instruction %d operand mismatch: got %d, want %d",
				i, instr.Operand, original.Instructions[i].Operand)
		}
	}

	// Verify constants match
	if len(decoded.Constants) != len(original.Constants) {
		t.Fatalf("Constant count mismatch: got %d, want %d",
			len(decoded.Constants), len(original.Constants))
	}

	if decoded.Constants[0] != int64(42) {
		t.Errorf("Constant value mismatch: got %v, want 42", decoded.Constants[0])
	}
}

// TestEncodeDecodeAllConstantTypes tests encoding and decoding of all
// supported constant types.
func TestEncodeDecodeAllConstantTypes(t *testing.T) {
	// Create bytecode with various constant types
	original := &Bytecode{
		Instructions: []Instruction{
			{Op: OpReturn, Operand: 0},
		},
		Constants: []interface{}{
			int64(123),          // Integer
			float64(3.14),       // Float
			"Hello, World!",     // String
			true,                // Boolean true
			false,               // Boolean false
			nil,                 // Nil
		},
	}

	// Encode and decode
	var buf bytes.Buffer
	if err := Encode(original, &buf); err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	decoded, err := Decode(&buf)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// Verify all constants
	if len(decoded.Constants) != len(original.Constants) {
		t.Fatalf("Constant count mismatch: got %d, want %d",
			len(decoded.Constants), len(original.Constants))
	}

	// Integer
	if decoded.Constants[0] != int64(123) {
		t.Errorf("Integer constant mismatch: got %v, want 123", decoded.Constants[0])
	}

	// Float
	if decoded.Constants[1] != float64(3.14) {
		t.Errorf("Float constant mismatch: got %v, want 3.14", decoded.Constants[1])
	}

	// String
	if decoded.Constants[2] != "Hello, World!" {
		t.Errorf("String constant mismatch: got %v, want 'Hello, World!'", decoded.Constants[2])
	}

	// Boolean true
	if decoded.Constants[3] != true {
		t.Errorf("Boolean true constant mismatch: got %v, want true", decoded.Constants[3])
	}

	// Boolean false
	if decoded.Constants[4] != false {
		t.Errorf("Boolean false constant mismatch: got %v, want false", decoded.Constants[4])
	}

	// Nil
	if decoded.Constants[5] != nil {
		t.Errorf("Nil constant mismatch: got %v, want nil", decoded.Constants[5])
	}
}

// TestEncodeDecodeAllOpcodes tests encoding and decoding of all opcodes.
func TestEncodeDecodeAllOpcodes(t *testing.T) {
	// Create bytecode with various opcodes
	original := &Bytecode{
		Instructions: []Instruction{
			{Op: OpPush, Operand: 0},
			{Op: OpPop, Operand: 0},
			{Op: OpDup, Operand: 0},
			{Op: OpSend, Operand: (1 << 8) | 2},
			{Op: OpSuperSend, Operand: (2 << 8) | 1},
			{Op: OpLoadLocal, Operand: 0},
			{Op: OpStoreLocal, Operand: 1},
			{Op: OpLoadField, Operand: 2},
			{Op: OpStoreField, Operand: 3},
			{Op: OpLoadClassVar, Operand: 0},
			{Op: OpStoreClassVar, Operand: 1},
			{Op: OpLoadGlobal, Operand: 4},
			{Op: OpStoreGlobal, Operand: 5},
			{Op: OpJump, Operand: 10},
			{Op: OpJumpIfFalse, Operand: 20},
			{Op: OpReturn, Operand: 0},
			{Op: OpPushSelf, Operand: 0},
			{Op: OpPushNil, Operand: 0},
			{Op: OpPushTrue, Operand: 0},
			{Op: OpPushFalse, Operand: 0},
			{Op: OpDefineClass, Operand: 0},
			{Op: OpNewObject, Operand: 0},
			{Op: OpMakeClosure, Operand: (6 << 8) | 2},
			{Op: OpCallBlock, Operand: 3},
			{Op: OpMakeArray, Operand: 5},
			{Op: OpMakeDictionary, Operand: 3},
		},
		Constants: []interface{}{
			int64(0), "selector1", "selector2", "global1", "global2", "global3", int64(1),
		},
	}

	// Encode and decode
	var buf bytes.Buffer
	if err := Encode(original, &buf); err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	decoded, err := Decode(&buf)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// Verify all instructions
	if len(decoded.Instructions) != len(original.Instructions) {
		t.Fatalf("Instruction count mismatch: got %d, want %d",
			len(decoded.Instructions), len(original.Instructions))
	}

	for i, instr := range decoded.Instructions {
		if instr.Op != original.Instructions[i].Op {
			t.Errorf("Instruction %d opcode mismatch: got %v, want %v",
				i, instr.Op, original.Instructions[i].Op)
		}
		if instr.Operand != original.Instructions[i].Operand {
			t.Errorf("Instruction %d operand mismatch: got %d, want %d",
				i, instr.Operand, original.Instructions[i].Operand)
		}
	}
}

// TestEncodeDecodeNestedBytecode tests encoding and decoding of bytecode
// containing nested bytecode (for blocks/closures).
func TestEncodeDecodeNestedBytecode(t *testing.T) {
	// Create nested bytecode (block inside main code)
	blockCode := &Bytecode{
		Instructions: []Instruction{
			{Op: OpLoadLocal, Operand: 0},
			{Op: OpPush, Operand: 0},
			{Op: OpSend, Operand: (1 << 8) | 1}, // + message
			{Op: OpReturn, Operand: 0},
		},
		Constants: []interface{}{
			int64(1),
			"+",
		},
	}

	original := &Bytecode{
		Instructions: []Instruction{
			{Op: OpMakeClosure, Operand: (0 << 8) | 1}, // Block with 1 param
			{Op: OpPush, Operand: 1},
			{Op: OpSend, Operand: (2 << 8) | 1}, // value: message
			{Op: OpReturn, Operand: 0},
		},
		Constants: []interface{}{
			blockCode,     // Nested bytecode
			int64(5),
			"value:",
		},
	}

	// Encode and decode
	var buf bytes.Buffer
	if err := Encode(original, &buf); err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	decoded, err := Decode(&buf)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// Verify nested bytecode
	if len(decoded.Constants) != 3 {
		t.Fatalf("Constant count mismatch: got %d, want 3", len(decoded.Constants))
	}

	nestedBC, ok := decoded.Constants[0].(*Bytecode)
	if !ok {
		t.Fatalf("First constant is not *Bytecode: got %T", decoded.Constants[0])
	}

	if len(nestedBC.Instructions) != 4 {
		t.Errorf("Nested bytecode instruction count mismatch: got %d, want 4",
			len(nestedBC.Instructions))
	}

	if len(nestedBC.Constants) != 2 {
		t.Errorf("Nested bytecode constant count mismatch: got %d, want 2",
			len(nestedBC.Constants))
	}
}

// TestEncodeDecodeClassDefinition tests encoding and decoding of
// ClassDefinition constants.
func TestEncodeDecodeClassDefinition(t *testing.T) {
	// Create a method for the class
	methodCode := &Bytecode{
		Instructions: []Instruction{
			{Op: OpLoadField, Operand: 0},
			{Op: OpReturn, Operand: 0},
		},
		Constants: []interface{}{},
	}

	classDef := &ClassDefinition{
		Name:       "Counter",
		SuperClass: "Object",
		Fields:     []string{"count"},
		ClassVariables: []string{"instanceCount"},
		ClassVarValues: make(map[string]interface{}),
		Methods: []*MethodDefinition{
			{
				Selector:   "value",
				Parameters: []string{},
				Code:       methodCode,
			},
		},
		ClassMethods: []*MethodDefinition{},
	}

	original := &Bytecode{
		Instructions: []Instruction{
			{Op: OpDefineClass, Operand: 0},
			{Op: OpReturn, Operand: 0},
		},
		Constants: []interface{}{
			classDef,
		},
	}

	// Encode and decode
	var buf bytes.Buffer
	if err := Encode(original, &buf); err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	decoded, err := Decode(&buf)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// Verify class definition
	if len(decoded.Constants) != 1 {
		t.Fatalf("Constant count mismatch: got %d, want 1", len(decoded.Constants))
	}

	decodedClass, ok := decoded.Constants[0].(*ClassDefinition)
	if !ok {
		t.Fatalf("First constant is not *ClassDefinition: got %T", decoded.Constants[0])
	}

	if decodedClass.Name != "Counter" {
		t.Errorf("Class name mismatch: got %s, want Counter", decodedClass.Name)
	}

	if decodedClass.SuperClass != "Object" {
		t.Errorf("Superclass name mismatch: got %s, want Object", decodedClass.SuperClass)
	}

	if len(decodedClass.Fields) != 1 || decodedClass.Fields[0] != "count" {
		t.Errorf("Fields mismatch: got %v, want [count]", decodedClass.Fields)
	}

	if len(decodedClass.ClassVariables) != 1 || decodedClass.ClassVariables[0] != "instanceCount" {
		t.Errorf("ClassVariables mismatch: got %v, want [instanceCount]", decodedClass.ClassVariables)
	}

	if len(decodedClass.Methods) != 1 {
		t.Fatalf("Method count mismatch: got %d, want 1", len(decodedClass.Methods))
	}

	if decodedClass.Methods[0].Selector != "value" {
		t.Errorf("Method selector mismatch: got %s, want value", decodedClass.Methods[0].Selector)
	}
}

// TestInvalidMagicNumber tests that decoding fails with wrong magic number.
func TestInvalidMagicNumber(t *testing.T) {
	// Create buffer with wrong magic number
	var buf bytes.Buffer
	wrongMagic := uint32(0x12345678)
	
	// Write wrong header manually
	buf.Write([]byte{
		byte(wrongMagic), byte(wrongMagic >> 8), byte(wrongMagic >> 16), byte(wrongMagic >> 24),
		0, 0, 0, 0, // version
		0, 0, 0, 0, // flags
	})

	// Try to decode
	_, err := Decode(&buf)
	if err == nil {
		t.Fatal("Expected error for invalid magic number, got nil")
	}
}

// TestUnsupportedVersion tests that decoding fails with unsupported version.
func TestUnsupportedVersion(t *testing.T) {
	// Create buffer with unsupported version
	var buf bytes.Buffer
	
	// Write header with unsupported version
	buf.Write([]byte{
		0x47, 0x4F, 0x4D, 0x53, // SMOG magic number
		99, 0, 0, 0,            // version 99
		0, 0, 0, 0,             // flags
	})

	// Try to decode
	_, err := Decode(&buf)
	if err == nil {
		t.Fatal("Expected error for unsupported version, got nil")
	}
}

// TestEmptyBytecode tests encoding and decoding of empty bytecode.
func TestEmptyBytecode(t *testing.T) {
	original := &Bytecode{
		Instructions: []Instruction{},
		Constants:    []interface{}{},
	}

	var buf bytes.Buffer
	if err := Encode(original, &buf); err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	decoded, err := Decode(&buf)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if len(decoded.Instructions) != 0 {
		t.Errorf("Expected 0 instructions, got %d", len(decoded.Instructions))
	}

	if len(decoded.Constants) != 0 {
		t.Errorf("Expected 0 constants, got %d", len(decoded.Constants))
	}
}

// TestLargeOperands tests encoding and decoding of instructions with
// large operand values (both positive and negative).
func TestLargeOperands(t *testing.T) {
	original := &Bytecode{
		Instructions: []Instruction{
			{Op: OpJump, Operand: 100000},
			{Op: OpJump, Operand: -100000},
			{Op: OpSend, Operand: (50000 << 8) | 255},
			{Op: OpReturn, Operand: 0},
		},
		Constants: []interface{}{},
	}

	var buf bytes.Buffer
	if err := Encode(original, &buf); err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	decoded, err := Decode(&buf)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if len(decoded.Instructions) != 4 {
		t.Fatalf("Instruction count mismatch: got %d, want 4", len(decoded.Instructions))
	}

	// Verify large positive operand
	if decoded.Instructions[0].Operand != 100000 {
		t.Errorf("Large positive operand mismatch: got %d, want 100000",
			decoded.Instructions[0].Operand)
	}

	// Verify large negative operand
	if decoded.Instructions[1].Operand != -100000 {
		t.Errorf("Large negative operand mismatch: got %d, want -100000",
			decoded.Instructions[1].Operand)
	}

	// Verify packed operand
	expectedPacked := (50000 << 8) | 255
	if decoded.Instructions[2].Operand != expectedPacked {
		t.Errorf("Packed operand mismatch: got %d, want %d",
			decoded.Instructions[2].Operand, expectedPacked)
	}
}

// TestUnicodeStrings tests encoding and decoding of Unicode strings.
func TestUnicodeStrings(t *testing.T) {
	original := &Bytecode{
		Instructions: []Instruction{
			{Op: OpReturn, Operand: 0},
		},
		Constants: []interface{}{
			"Hello, 世界",           // Chinese
			"Привет, мир",         // Russian
			"مرحبا بالعالم",       // Arabic
			"🎉🎊✨",               // Emojis
		},
	}

	var buf bytes.Buffer
	if err := Encode(original, &buf); err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	decoded, err := Decode(&buf)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if len(decoded.Constants) != 4 {
		t.Fatalf("Constant count mismatch: got %d, want 4", len(decoded.Constants))
	}

	expected := []string{
		"Hello, 世界",
		"Привет, мир",
		"مرحبا بالعالم",
		"🎉🎊✨",
	}

	for i, exp := range expected {
		if decoded.Constants[i] != exp {
			t.Errorf("Unicode string %d mismatch: got %s, want %s",
				i, decoded.Constants[i], exp)
		}
	}
}

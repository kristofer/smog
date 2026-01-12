// Package bytecode defines the bytecode format and opcodes for smog.
package bytecode

// Opcode represents a bytecode instruction
type Opcode byte

const (
	// Stack operations
	OpPush Opcode = iota
	OpPop
	OpDup

	// Message sends
	OpSend
	OpSuperSend

	// Variables
	OpLoadLocal
	OpStoreLocal
	OpLoadField
	OpStoreField
	OpLoadGlobal
	OpStoreGlobal

	// Control flow
	OpJump
	OpJumpIfFalse
	OpReturn

	// Literals
	OpPushSelf
	OpPushNil
	OpPushTrue
	OpPushFalse

	// Object creation
	OpNewObject
)

// Instruction represents a single bytecode instruction
type Instruction struct {
	Op      Opcode
	Operand int
}

// Bytecode represents compiled bytecode
type Bytecode struct {
	Instructions []Instruction
	Constants    []interface{}
}

// Constants for encoding/decoding message send operands
const (
	// SelectorIndexShift is the number of bits to shift the selector index
	SelectorIndexShift = 8
	// ArgCountMask is the mask for extracting the argument count
	ArgCountMask = 0xFF
)

// String returns a string representation of an opcode
func (op Opcode) String() string {
	switch op {
	case OpPush:
		return "PUSH"
	case OpPop:
		return "POP"
	case OpDup:
		return "DUP"
	case OpSend:
		return "SEND"
	case OpSuperSend:
		return "SUPER_SEND"
	case OpLoadLocal:
		return "LOAD_LOCAL"
	case OpStoreLocal:
		return "STORE_LOCAL"
	case OpLoadField:
		return "LOAD_FIELD"
	case OpStoreField:
		return "STORE_FIELD"
	case OpLoadGlobal:
		return "LOAD_GLOBAL"
	case OpStoreGlobal:
		return "STORE_GLOBAL"
	case OpJump:
		return "JUMP"
	case OpJumpIfFalse:
		return "JUMP_IF_FALSE"
	case OpReturn:
		return "RETURN"
	case OpPushSelf:
		return "PUSH_SELF"
	case OpPushNil:
		return "PUSH_NIL"
	case OpPushTrue:
		return "PUSH_TRUE"
	case OpPushFalse:
		return "PUSH_FALSE"
	case OpNewObject:
		return "NEW_OBJECT"
	default:
		return "UNKNOWN"
	}
}

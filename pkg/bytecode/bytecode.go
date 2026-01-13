// Package bytecode defines the bytecode format and opcodes for smog.
//
// The bytecode is the low-level intermediate representation that the smog
// virtual machine (VM) executes. It consists of a sequence of instructions,
// each with an opcode and an optional operand, plus a constant pool for
// literal values.
//
// Architecture:
//
// The bytecode system follows a stack-based architecture where:
//   1. Values are pushed onto and popped from a runtime stack
//   2. Operations consume values from the stack and push results back
//   3. Variables are stored in separate local and global storage
//   4. Message sends use dynamic dispatch to find and execute methods
//
// Example compilation:
//
//   Source:  x := 10. x + 5.
//
//   Bytecode:
//     PUSH 10         ; Load constant 10 onto stack
//     STORE_LOCAL 0   ; Store to local variable x (slot 0)
//     LOAD_LOCAL 0    ; Load x back onto stack
//     PUSH 5          ; Load constant 5 onto stack
//     SEND +, 1       ; Send + message with 1 argument
//     RETURN          ; End of program
//
//   Constants: [10, 5, "+"]
//
// Instruction Format:
//
// Each instruction consists of:
//   - Opcode (byte): The operation to perform
//   - Operand (int): Additional data for the instruction
//
// The operand's meaning depends on the opcode:
//   - PUSH: index into constant pool
//   - LOAD_LOCAL/STORE_LOCAL: local variable slot number
//   - SEND: packed value containing selector index and arg count
//
// Design Philosophy:
//
// The bytecode design balances simplicity with efficiency:
//   - Simple opcodes are easy to implement and debug
//   - Stack-based design minimizes temporary storage complexity
//   - Constant pool reduces bytecode size (literals referenced by index)
//   - Separation of concerns: bytecode describes "what to do", VM decides "how"
package bytecode

// Opcode represents a bytecode instruction operation.
//
// Each opcode tells the VM what operation to perform. Opcodes are
// single bytes (0-255), making them compact and fast to decode.
type Opcode byte

// Bytecode instruction opcodes.
//
// These are organized by category for clarity:
const (
	// === Stack Operations ===
	//
	// These opcodes manipulate the value stack directly.

	// OpPush loads a constant from the constant pool onto the stack.
	// Operand: index into the constant pool
	//
	// Example: PUSH 0  ; loads constant[0] onto stack
	OpPush Opcode = iota

	// OpPop removes the top value from the stack.
	// Operand: unused (typically 0)
	//
	// Used to discard values that aren't needed.
	OpPop

	// OpDup duplicates the top value on the stack.
	// Operand: unused (typically 0)
	//
	// Useful when the same value is needed multiple times.
	OpDup

	// === Message Operations ===
	//
	// These opcodes handle message sending - the core operation in smog.

	// OpSend sends a message to an object.
	// Operand: packed value containing:
	//   - High bits: selector index in constant pool
	//   - Low 8 bits: number of arguments
	//
	// Stack before: [receiver, arg1, arg2, ..., argN]
	// Stack after:  [result]
	//
	// This is the most important instruction - it implements the message
	// passing semantics that make smog object-oriented.
	OpSend

	// OpSuperSend sends a message to the superclass.
	// Operand: same format as OpSend
	//
	// Like OpSend but starts method lookup in the superclass rather
	// than the receiver's class. Used for super calls.
	OpSuperSend

	// === Variable Operations ===
	//
	// These opcodes handle variable access and assignment.

	// OpLoadLocal loads a local variable onto the stack.
	// Operand: index of the local variable slot
	//
	// Local variables are function/method/block scoped and stored
	// in a fixed-size array indexed by their declaration order.
	OpLoadLocal

	// OpStoreLocal stores a value to a local variable.
	// Operand: index of the local variable slot
	//
	// Pops the top value from the stack and stores it in the local
	// variable, then pushes the value back (assignments return values).
	OpStoreLocal

	// OpLoadField loads an instance variable from the current object.
	// Operand: index of the field
	//
	// Used within methods to access instance variables (fields)
	// of the receiver (self).
	OpLoadField

	// OpStoreField stores a value to an instance variable.
	// Operand: index of the field
	//
	// Used within methods to modify instance variables of self.
	OpStoreField

	// OpLoadGlobal loads a global variable onto the stack.
	// Operand: index of the variable name in constant pool
	//
	// Global variables are stored in a hash map keyed by name.
	OpLoadGlobal

	// OpStoreGlobal stores a value to a global variable.
	// Operand: index of the variable name in constant pool
	OpStoreGlobal

	// === Control Flow Operations ===
	//
	// These opcodes control program flow.

	// OpJump unconditionally jumps to a new instruction.
	// Operand: target instruction index
	//
	// Used to implement loops and skip over code sections.
	OpJump

	// OpJumpIfFalse conditionally jumps if the top stack value is false.
	// Operand: target instruction index
	//
	// Pops a boolean from the stack and jumps if it's false.
	// Used to implement if-statements and short-circuit logic.
	OpJumpIfFalse

	// OpReturn returns from the current method/block/program.
	// Operand: unused
	//
	// Ends execution of the current code context. If there's a value
	// on the stack, it becomes the return value.
	OpReturn

	// === Literal Operations ===
	//
	// These opcodes push special constant values onto the stack.
	// They're more efficient than using OpPush with a constant pool entry.

	// OpPushSelf pushes the current receiver (self) onto the stack.
	// Operand: unused
	//
	// Used within methods to refer to the object receiving the message.
	OpPushSelf

	// OpPushNil pushes the nil value onto the stack.
	// Operand: unused
	OpPushNil

	// OpPushTrue pushes the boolean true value onto the stack.
	// Operand: unused
	OpPushTrue

	// OpPushFalse pushes the boolean false value onto the stack.
	// Operand: unused
	OpPushFalse

	// === Object Operations ===

	// OpNewObject creates a new instance of a class.
	// Operand: class identifier
	//
	// Allocates a new object with space for instance variables.
	OpNewObject
)

// Instruction represents a single bytecode instruction.
//
// Each instruction consists of an operation (opcode) and an operand.
// The operand's meaning depends on the opcode - it might be an index,
// a count, an offset, or unused.
//
// Example:
//   Instruction{Op: OpPush, Operand: 3}
//     -> Push constant[3] onto the stack
//
//   Instruction{Op: OpLoadLocal, Operand: 0}
//     -> Load local variable at index 0 onto the stack
//
//   Instruction{Op: OpSend, Operand: (2 << 8) | 1}
//     -> Send message with selector at constant[2] with 1 argument
type Instruction struct {
	Op      Opcode // The operation to perform
	Operand int    // Additional data for the instruction
}

// Bytecode represents a complete compiled program or method.
//
// A Bytecode contains everything needed to execute a piece of smog code:
//   - Instructions: The sequence of operations to perform
//   - Constants: The pool of literal values referenced by instructions
//
// The constant pool stores:
//   - Numbers (int64, float64)
//   - Strings (for string literals and selectors)
//   - Variable names (for global access)
//
// Why use a constant pool?
//   - Reduces bytecode size (reference by index instead of embedding)
//   - Allows sharing of common values
//   - Simplifies instruction format (fixed-size operands)
//
// Example:
//
//   Source: 'Hello' println. 42.
//
//   Bytecode{
//     Instructions: [
//       {OpPush, 0},       ; Push constant[0] ("Hello")
//       {OpSend, (1<<8)|0},; Send constant[1] ("println") with 0 args
//       {OpPop, 0},        ; Discard result
//       {OpPush, 2},       ; Push constant[2] (42)
//       {OpReturn, 0},     ; End
//     ],
//     Constants: ["Hello", "println", 42]
//   }
type Bytecode struct {
	Instructions []Instruction // Sequence of bytecode instructions
	Constants    []interface{} // Pool of constant values
}

// Constants for encoding/decoding message send operands.
//
// For OpSend and OpSuperSend instructions, we need to encode two pieces
// of information in a single operand:
//   1. The selector (message name) - index into constant pool
//   2. The number of arguments
//
// We pack these together using bit manipulation:
//   - High bits (8 and above): selector index
//   - Low 8 bits: argument count (0-255)
//
// Example:
//   Selector index: 5
//   Arg count: 2
//   Packed operand: (5 << 8) | 2 = 1282
//
// To unpack:
//   selectorIndex := operand >> 8        // Right shift 8 bits -> 5
//   argCount := operand & 0xFF           // Mask low 8 bits -> 2
//
// This approach allows us to keep the Instruction format simple with
// a single operand field while still encoding the necessary information.
const (
	// SelectorIndexShift is the number of bits to shift left when encoding
	// the selector index, or shift right when decoding it.
	SelectorIndexShift = 8

	// ArgCountMask is the bitmask for extracting the argument count from
	// the low 8 bits of the operand.
	ArgCountMask = 0xFF
)

// String returns a human-readable name for an opcode.
//
// This is primarily used for debugging, logging, and disassembling bytecode.
// It allows us to print instructions in a readable format like:
//   PUSH 0
//   LOAD_LOCAL 1
//   SEND 2
// instead of opaque numbers.
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

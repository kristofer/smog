// Package vm implements the bytecode virtual machine for smog.
//
// The VM is a stack-based interpreter that executes bytecode instructions.
// It's the final stage in the execution pipeline:
//
//   Source Code -> Lexer -> Parser -> AST -> Compiler -> Bytecode -> VM -> Execution
//
// Virtual Machine Architecture:
//
// The VM uses a stack-based architecture with the following components:
//
//   1. Value Stack: Holds intermediate values during computation
//   2. Stack Pointer (sp): Tracks the top of the value stack
//   3. Local Variables: Array of local variable values
//   4. Global Variables: Hash map of global variable values
//   5. Constants: Pool of literal values from the bytecode
//
// Execution Model:
//
// The VM executes instructions sequentially using an instruction pointer (ip).
// Each instruction manipulates the stack, variables, or control flow.
//
// Example Execution:
//
//   Source: x := 5. x + 3.
//
//   Bytecode:
//     0: PUSH 0          ; constant[0] = 5
//     1: STORE_LOCAL 0   ; x is slot 0
//     2: LOAD_LOCAL 0    ; load x
//     3: PUSH 1          ; constant[1] = 3
//     4: SEND 2, 1       ; constant[2] = "+", 1 argument
//     5: RETURN
//
//   Execution trace:
//     IP=0: PUSH 0        -> stack=[5]
//     IP=1: STORE_LOCAL 0 -> stack=[5], locals[0]=5
//     IP=2: LOAD_LOCAL 0  -> stack=[5,5]
//     IP=3: PUSH 1        -> stack=[5,5,3]
//     IP=4: SEND +, 1     -> stack=[5,8]  (5+3=8)
//     IP=5: RETURN        -> done
//
// Stack Operations:
//
// Most operations follow a pattern:
//   1. Pop operands from stack
//   2. Perform operation
//   3. Push result back onto stack
//
// This keeps the VM simple and uniform. For example, binary operations
// like + always pop two values and push one result.
//
// Message Dispatch:
//
// The send() method implements message dispatch. Currently, it handles
// primitive operations (arithmetic, comparison, I/O) directly. In a full
// implementation, it would look up methods in the receiver's class.
//
// Error Handling:
//
// The VM returns errors for runtime problems:
//   - Stack overflow/underflow
//   - Invalid operands (e.g., adding string to number)
//   - Division by zero
//   - Unknown messages
//
// Design Philosophy:
//
// The VM is designed to be:
//   - Simple: Easy to understand and debug
//   - Efficient: Minimal overhead for common operations
//   - Safe: Checks bounds and types to prevent crashes
//   - Extensible: Easy to add new operations and types
package vm

import (
	"fmt"

	"github.com/kristofer/smog/pkg/bytecode"
)

// VM represents the virtual machine that executes bytecode.
//
// State Components:
//
//   stack: The value stack for intermediate computations
//     - Fixed size (1024 entries)
//     - Grows upward as values are pushed
//     - Values can be any Go type (int64, float64, string, bool, nil, objects)
//
//   sp: Stack pointer - index of the next free slot
//     - Points one past the top element
//     - sp=0 means stack is empty
//     - sp=N means there are N elements, top is at stack[N-1]
//
//   locals: Local variable storage
//     - Fixed size array (256 slots)
//     - Indexed by variable slot number from compiler
//     - Initialized to nil
//
//   globals: Global variable storage
//     - Hash map keyed by variable name
//     - Created on first assignment
//     - Persists across multiple Run() calls
//
//   constants: Constant pool from bytecode
//     - Set at the start of Run()
//     - Contains literals and identifiers
//     - Referenced by index in instructions
type VM struct {
	stack     []interface{}         // Value stack for computation
	sp        int                   // Stack pointer (index of next free slot)
	locals    []interface{}         // Local variable storage
	globals   map[string]interface{} // Global variable storage
	constants []interface{}         // Constant pool from bytecode
}

// New creates a new virtual machine instance.
//
// Initializes:
//   - Empty value stack with 1024 slots
//   - Stack pointer at 0 (empty)
//   - Local variable array with 256 slots
//   - Empty global variable map
//
// The VM is reusable - you can call Run() multiple times on the same VM.
// Global variables persist across runs, but the stack and locals are reset.
func New() *VM {
	return &VM{
		stack:   make([]interface{}, 1024),
		sp:      0,
		locals:  make([]interface{}, 256),
		globals: make(map[string]interface{}),
	}
}

// Run executes bytecode on the virtual machine.
//
// This is the main execution loop of the VM. It processes instructions
// sequentially from the bytecode until hitting a RETURN or an error.
//
// Execution Process:
//   1. Reset VM state (stack and locals cleared)
//   2. Load the constant pool from bytecode
//   3. Execute instructions from IP=0 until RETURN or error
//   4. Each instruction updates stack, variables, or control flow
//
// Parameters:
//   - bc: The bytecode to execute (instructions + constants)
//
// Returns:
//   - nil if execution completed successfully
//   - error if a runtime error occurred
//
// State Management:
//   The VM resets its stack and locals before each run to ensure clean
//   execution. However, global variables persist across runs, allowing
//   state to be maintained between executions.
//
// Example:
//
//   vm := vm.New()
//   bytecode, _ := compiler.Compile(program)
//   err := vm.Run(bytecode)
//   if err != nil {
//     fmt.Println("Runtime error:", err)
//   }
//   result := vm.StackTop() // Get the final result
func (vm *VM) Run(bc *bytecode.Bytecode) error {
	// Reset state for clean execution
	// Stack pointer back to 0 (empty stack)
	vm.sp = 0
	
	// Clear all local variables to nil
	// This ensures no leftover state from previous runs
	for i := range vm.locals {
		vm.locals[i] = nil
	}
	
	// Load the constant pool from the bytecode
	vm.constants = bc.Constants

	// Main execution loop
	// Process instructions sequentially using instruction pointer (ip)
	for ip := 0; ip < len(bc.Instructions); ip++ {
		inst := bc.Instructions[ip]

		// Dispatch to instruction handler based on opcode
		switch inst.Op {
		case bytecode.OpPush:
			// PUSH: Load a constant onto the stack
			// Operand: index into constant pool
			//
			// Example: PUSH 2 loads constant[2] onto stack
			if inst.Operand < 0 || inst.Operand >= len(vm.constants) {
				return fmt.Errorf("constant index out of bounds: %d", inst.Operand)
			}
			if err := vm.push(vm.constants[inst.Operand]); err != nil {
				return err
			}

		case bytecode.OpPop:
			// POP: Discard the top value from the stack
			// Operand: unused
			//
			// Used to clean up unwanted values
			if _, err := vm.pop(); err != nil {
				return err
			}

		case bytecode.OpPushTrue:
			// PUSH_TRUE: Push boolean true onto the stack
			// Operand: unused
			//
			// More efficient than using OpPush with a constant
			if err := vm.push(true); err != nil {
				return err
			}

		case bytecode.OpPushFalse:
			// PUSH_FALSE: Push boolean false onto the stack
			// Operand: unused
			if err := vm.push(false); err != nil {
				return err
			}

		case bytecode.OpPushNil:
			// PUSH_NIL: Push nil onto the stack
			// Operand: unused
			if err := vm.push(nil); err != nil {
				return err
			}

		case bytecode.OpLoadLocal:
			// LOAD_LOCAL: Load a local variable onto the stack
			// Operand: local variable slot index
			//
			// Example: LOAD_LOCAL 0 loads locals[0]
			if inst.Operand < 0 || inst.Operand >= len(vm.locals) {
				return fmt.Errorf("local variable index out of bounds: %d", inst.Operand)
			}
			if err := vm.push(vm.locals[inst.Operand]); err != nil {
				return err
			}

		case bytecode.OpStoreLocal:
			// STORE_LOCAL: Store the top stack value to a local variable
			// Operand: local variable slot index
			//
			// The value is popped from the stack, stored, then pushed back
			// because assignments return their value in smog.
			//
			// Example: STORE_LOCAL 0 stores to locals[0]
			if inst.Operand < 0 || inst.Operand >= len(vm.locals) {
				return fmt.Errorf("local variable index out of bounds: %d", inst.Operand)
			}
			val, err := vm.pop()
			if err != nil {
				return err
			}
			vm.locals[inst.Operand] = val
			// Push the value back (assignment returns the value)
			if err := vm.push(val); err != nil {
				return err
			}

		case bytecode.OpLoadGlobal:
			// LOAD_GLOBAL: Load a global variable onto the stack
			// Operand: index of variable name in constant pool
			//
			// Global variables are stored in a map by name.
			// The name is retrieved from the constant pool.
			//
			// Example: LOAD_GLOBAL 5 where constant[5]="MyClass"
			//   -> loads globals["MyClass"]
			if inst.Operand < 0 || inst.Operand >= len(vm.constants) {
				return fmt.Errorf("constant index out of bounds: %d", inst.Operand)
			}
			name, ok := vm.constants[inst.Operand].(string)
			if !ok {
				return fmt.Errorf("expected string constant for global name")
			}
			val, ok := vm.globals[name]
			if !ok {
				return fmt.Errorf("undefined global variable: %s", name)
			}
			if err := vm.push(val); err != nil {
				return err
			}

		case bytecode.OpStoreGlobal:
			// STORE_GLOBAL: Store the top stack value to a global variable
			// Operand: index of variable name in constant pool
			//
			// Creates the global if it doesn't exist.
			// Like local stores, the value is pushed back.
			if inst.Operand < 0 || inst.Operand >= len(vm.constants) {
				return fmt.Errorf("constant index out of bounds: %d", inst.Operand)
			}
			name, ok := vm.constants[inst.Operand].(string)
			if !ok {
				return fmt.Errorf("expected string constant for global name")
			}
			val, err := vm.pop()
			if err != nil {
				return err
			}
			vm.globals[name] = val
			// Push the value back
			if err := vm.push(val); err != nil {
				return err
			}

		case bytecode.OpSend:
			// SEND: Send a message to an object
			// Operand: packed value with selector index and arg count
			//
			// This is the core operation that implements message passing.
			//
			// Process:
			//   1. Decode selector index and arg count from operand
			//   2. Pop arguments from stack (in reverse order)
			//   3. Pop receiver from stack
			//   4. Execute the message send (via send() method)
			//   5. Push result onto stack
			//
			// Stack before: [receiver, arg1, arg2, ..., argN]
			// Stack after:  [result]

			// Decode operand using bit manipulation
			// High bits: selector index in constant pool
			// Low 8 bits: argument count
			selectorIdx := inst.Operand >> bytecode.SelectorIndexShift
			argCount := inst.Operand & bytecode.ArgCountMask

			// Get the selector string from constants
			if selectorIdx < 0 || selectorIdx >= len(vm.constants) {
				return fmt.Errorf("selector index out of bounds: %d", selectorIdx)
			}
			selector, ok := vm.constants[selectorIdx].(string)
			if !ok {
				return fmt.Errorf("expected string constant for selector")
			}

			// Pop arguments in reverse order
			// They were pushed left-to-right, so we pop right-to-left
			// to get them back in the correct order
			args := make([]interface{}, argCount)
			for i := argCount - 1; i >= 0; i-- {
				arg, err := vm.pop()
				if err != nil {
					return err
				}
				args[i] = arg
			}

			// Pop receiver
			receiver, err := vm.pop()
			if err != nil {
				return err
			}

			// Execute the message send
			result, err := vm.send(receiver, selector, args)
			if err != nil {
				return err
			}

			// Push result onto stack
			if err := vm.push(result); err != nil {
				return err
			}

		case bytecode.OpReturn:
			// RETURN: End execution
			// Operand: unused
			//
			// Exits the execution loop. The final value (if any) remains
			// on the stack and can be retrieved with StackTop().
			return nil

		default:
			return fmt.Errorf("unknown opcode: %v", inst.Op)
		}
	}

	return nil
}

// send executes a message send operation.
//
// This method implements the message dispatch mechanism - the core of
// object-oriented programming in smog. When a message is sent to an object,
// this method determines what action to take.
//
// Current Implementation:
//   This is a simplified implementation that handles only primitive operations.
//   In a full Smalltalk-style implementation, this would:
//     1. Look up the receiver's class
//     2. Search for a method matching the selector
//     3. Execute the method in a new context
//     4. Return the result
//
// Primitive Operations:
//   For now, we handle these selectors as built-in primitives:
//     - Arithmetic: +, -, *, /
//     - Comparison: <, >, <=, >=, =, ~=
//     - I/O: print, println
//
// Parameters:
//   - receiver: The object receiving the message
//   - selector: The message name (method selector)
//   - args: Arguments to the message
//
// Returns:
//   - The result of the operation
//   - Error if the message is unknown or arguments are invalid
//
// Example:
//   send(5, "+", [3]) -> 8
//   send("Hello", "println", []) -> "Hello" (and prints it)
func (vm *VM) send(receiver interface{}, selector string, args []interface{}) (interface{}, error) {
	// Handle primitive operations
	// These are built directly into the VM for efficiency
	switch selector {
	case "+":
		return vm.add(receiver, args[0])
	case "-":
		return vm.subtract(receiver, args[0])
	case "*":
		return vm.multiply(receiver, args[0])
	case "/":
		return vm.divide(receiver, args[0])
	case "<":
		return vm.lessThan(receiver, args[0])
	case ">":
		return vm.greaterThan(receiver, args[0])
	case "<=":
		return vm.lessOrEqual(receiver, args[0])
	case ">=":
		return vm.greaterOrEqual(receiver, args[0])
	case "=":
		return vm.equal(receiver, args[0])
	case "~=":
		return vm.notEqual(receiver, args[0])
	case "println":
		// Print the receiver followed by a newline
		fmt.Println(receiver)
		// Return the receiver (allows method chaining)
		return receiver, nil
	case "print":
		// Print the receiver without a newline
		fmt.Print(receiver)
		return receiver, nil
	default:
		return nil, fmt.Errorf("unknown message: %s", selector)
	}
}

// Primitive operations for arithmetic and comparison.
//
// These implement the basic mathematical and logical operations that form
// the foundation of computation. Each operation:
//   1. Type-checks the operands
//   2. Performs the operation
//   3. Returns the result or an error
//
// Type Support:
//   Currently supports int64 and float64 for numeric operations.
//   A full implementation would use polymorphic method dispatch instead.

// add implements the + binary message.
//
// Supported types:
//   - int64 + int64 -> int64
//   - float64 + float64 -> float64
//
// Examples:
//   add(5, 3) -> 8
//   add(2.5, 1.5) -> 4.0
//
// Errors:
//   - Type mismatch (e.g., int + float)
//   - Unsupported types
func (vm *VM) add(a, b interface{}) (interface{}, error) {
	switch aVal := a.(type) {
	case int64:
		if bVal, ok := b.(int64); ok {
			return aVal + bVal, nil
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			return aVal + bVal, nil
		}
	}
	return nil, fmt.Errorf("cannot add %T and %T", a, b)
}

// subtract implements the - binary message.
//
// Supported types:
//   - int64 - int64 -> int64
//   - float64 - float64 -> float64
func (vm *VM) subtract(a, b interface{}) (interface{}, error) {
	switch aVal := a.(type) {
	case int64:
		if bVal, ok := b.(int64); ok {
			return aVal - bVal, nil
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			return aVal - bVal, nil
		}
	}
	return nil, fmt.Errorf("cannot subtract %T and %T", a, b)
}

// multiply implements the * binary message.
//
// Supported types:
//   - int64 * int64 -> int64
//   - float64 * float64 -> float64
func (vm *VM) multiply(a, b interface{}) (interface{}, error) {
	switch aVal := a.(type) {
	case int64:
		if bVal, ok := b.(int64); ok {
			return aVal * bVal, nil
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			return aVal * bVal, nil
		}
	}
	return nil, fmt.Errorf("cannot multiply %T and %T", a, b)
}

// divide implements the / binary message.
//
// Supported types:
//   - int64 / int64 -> int64 (integer division)
//   - float64 / float64 -> float64
//
// Errors:
//   - Division by zero
//   - Type mismatch
func (vm *VM) divide(a, b interface{}) (interface{}, error) {
	switch aVal := a.(type) {
	case int64:
		if bVal, ok := b.(int64); ok {
			if bVal == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return aVal / bVal, nil
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			if bVal == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return aVal / bVal, nil
		}
	}
	return nil, fmt.Errorf("cannot divide %T and %T", a, b)
}

// Comparison operations return boolean values.
//
// These implement the relational operators that allow comparing values.
// All return true or false.

// lessThan implements the < binary message.
func (vm *VM) lessThan(a, b interface{}) (interface{}, error) {
	switch aVal := a.(type) {
	case int64:
		if bVal, ok := b.(int64); ok {
			return aVal < bVal, nil
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			return aVal < bVal, nil
		}
	}
	return nil, fmt.Errorf("cannot compare %T and %T", a, b)
}

// greaterThan implements the > binary message.
func (vm *VM) greaterThan(a, b interface{}) (interface{}, error) {
	switch aVal := a.(type) {
	case int64:
		if bVal, ok := b.(int64); ok {
			return aVal > bVal, nil
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			return aVal > bVal, nil
		}
	}
	return nil, fmt.Errorf("cannot compare %T and %T", a, b)
}

// lessOrEqual implements the <= binary message.
func (vm *VM) lessOrEqual(a, b interface{}) (interface{}, error) {
	switch aVal := a.(type) {
	case int64:
		if bVal, ok := b.(int64); ok {
			return aVal <= bVal, nil
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			return aVal <= bVal, nil
		}
	}
	return nil, fmt.Errorf("cannot compare %T and %T", a, b)
}

// greaterOrEqual implements the >= binary message.
func (vm *VM) greaterOrEqual(a, b interface{}) (interface{}, error) {
	switch aVal := a.(type) {
	case int64:
		if bVal, ok := b.(int64); ok {
			return aVal >= bVal, nil
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			return aVal >= bVal, nil
		}
	}
	return nil, fmt.Errorf("cannot compare %T and %T", a, b)
}

// equal implements the = binary message.
//
// Uses Go's == operator, which handles most types correctly.
// Returns true if the values are equal, false otherwise.
func (vm *VM) equal(a, b interface{}) (interface{}, error) {
	return a == b, nil
}

// notEqual implements the ~= binary message.
//
// Complement of equal - returns true if values are different.
func (vm *VM) notEqual(a, b interface{}) (interface{}, error) {
	return a != b, nil
}

// Stack manipulation methods.
//
// These implement the basic stack operations used throughout the VM.
// The stack is a fundamental data structure for expression evaluation.

// push adds a value to the top of the stack.
//
// The stack grows upward. Each push:
//   1. Checks for stack overflow
//   2. Stores the value at stack[sp]
//   3. Increments the stack pointer
//
// Parameters:
//   - obj: The value to push (can be any type)
//
// Returns:
//   - nil if successful
//   - error if stack overflow
//
// Example:
//   Initial: stack=[], sp=0
//   push(5): stack=[5], sp=1
//   push(3): stack=[5,3], sp=2
func (vm *VM) push(obj interface{}) error {
	if vm.sp >= len(vm.stack) {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = obj
	vm.sp++
	return nil
}

// pop removes and returns the value from the top of the stack.
//
// The stack shrinks downward. Each pop:
//   1. Checks for stack underflow
//   2. Decrements the stack pointer
//   3. Returns the value at the new top
//
// Returns:
//   - The popped value
//   - error if stack underflow
//
// Example:
//   Initial: stack=[5,3], sp=2
//   pop(): returns 3, stack=[5], sp=1
//   pop(): returns 5, stack=[], sp=0
func (vm *VM) pop() (interface{}, error) {
	if vm.sp <= 0 {
		return nil, fmt.Errorf("stack underflow")
	}
	vm.sp--
	return vm.stack[vm.sp], nil
}

// StackTop returns the value at the top of the stack without removing it.
//
// This is useful for inspecting the final result after execution without
// modifying the stack state.
//
// Returns:
//   - The top stack value, or nil if stack is empty
//
// Example:
//   After executing "3 + 4", StackTop() returns 7
func (vm *VM) StackTop() interface{} {
	if vm.sp <= 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

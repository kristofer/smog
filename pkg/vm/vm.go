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
	stack        []interface{}                        // Value stack for computation
	sp           int                                  // Stack pointer (index of next free slot)
	locals       []interface{}                        // Local variable storage
	globals      map[string]interface{}               // Global variable storage
	constants    []interface{}                        // Constant pool from bytecode
	self         interface{}                          // Current receiver (self) for method execution
	currentClass *bytecode.ClassDefinition            // Current class context (for super sends)
	fieldOffset  int                                  // Offset for field indices (for inheritance)
	classes      map[string]*bytecode.ClassDefinition // Registered classes by name
	homeContext  *VM                                  // Home context for non-local returns (nil for methods, set for blocks)
	callStack    []StackFrame                         // Call stack for debugging and error reporting
	ip           int                                  // Current instruction pointer (for error reporting)
	debugger     *Debugger                            // Optional debugger for interactive debugging
}

// New creates a new virtual machine instance.
//
// Initializes:
//   - Empty value stack with 1024 slots
//   - Stack pointer at 0 (empty)
//   - Local variable array with 256 slots
//   - Empty global variable map
//   - Empty class registry
//
// The VM is reusable - you can call Run() multiple times on the same VM.
// Global variables and registered classes persist across runs, but the 
// stack and locals are reset.
func New() *VM {
	return &VM{
		stack:     make([]interface{}, 1024),
		sp:        0,
		locals:    make([]interface{}, 256),
		globals:   make(map[string]interface{}),
		classes:   make(map[string]*bytecode.ClassDefinition),
		callStack: make([]StackFrame, 0, 64), // Preallocate space for 64 frames
	}
}

// Run executes bytecode on the virtual machine.
//
// This is the main execution loop of the VM. It processes instructions
// sequentially from the bytecode until hitting a RETURN or an error.
//
// Execution Process:
//   1. Reset VM state (stack cleared; locals cleared only if all are nil)
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
//   The VM resets its stack before each run. Locals are only cleared
//   if they appear to be uninitialized (all nil). This allows blocks
//   to pre-load parameter values before calling Run().
//   Global variables persist across runs, allowing state to be maintained.
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
	// Reset stack pointer to 0 (empty stack)
	vm.sp = 0
	
	// Check if locals need to be cleared
	// If any local is non-nil, we assume they've been pre-initialized
	// (e.g., for block parameters) and don't clear them
	hasInitializedLocals := false
	for i := range vm.locals {
		if vm.locals[i] != nil {
			hasInitializedLocals = true
			break
		}
	}
	
	// Only clear locals if none are initialized
	if !hasInitializedLocals {
		for i := range vm.locals {
			vm.locals[i] = nil
		}
	}
	
	// Load the constant pool from the bytecode
	vm.constants = bc.Constants

	// Push a frame for the main program execution
	vm.pushFrame("main program", "")
	// Use defer to ensure frame is popped even on error
	defer vm.popFrame()

	// Main execution loop
	// Process instructions sequentially using instruction pointer (ip)
	for vm.ip = 0; vm.ip < len(bc.Instructions); vm.ip++ {
		inst := bc.Instructions[vm.ip]

		// Check for debugger breakpoints
		if vm.debugger != nil && vm.debugger.ShouldPause() {
			if !vm.debugger.InteractivePrompt(bc) {
				// User chose to quit
				return fmt.Errorf("debugging session terminated")
			}
		}

		// Dispatch to instruction handler based on opcode
		switch inst.Op {
		case bytecode.OpPush:
			// PUSH: Load a constant onto the stack
			// Operand: index into constant pool
			//
			// Example: PUSH 2 loads constant[2] onto stack
			if inst.Operand < 0 || inst.Operand >= len(vm.constants) {
				return vm.runtimeError(fmt.Sprintf("constant index out of bounds: %d", inst.Operand))
			}
			if err := vm.push(vm.constants[inst.Operand]); err != nil {
				return vm.runtimeError(err.Error())
			}

		case bytecode.OpPop:
			// POP: Discard the top value from the stack
			// Operand: unused
			//
			// Used to clean up unwanted values
			if _, err := vm.pop(); err != nil {
				return err
			}

		case bytecode.OpDup:
			// DUP: Duplicate the top value on the stack
			// Operand: unused
			//
			// Creates a copy of the top stack value and pushes it.
			// Stack before: [..., value]
			// Stack after:  [..., value, value]
			//
			// This is used in cascading messages to keep the receiver
			// available for multiple message sends.
			if vm.sp == 0 {
				return vm.runtimeError("stack underflow: cannot duplicate empty stack")
			}
			topValue := vm.stack[vm.sp-1]
			if err := vm.push(topValue); err != nil {
				return vm.runtimeError(err.Error())
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

		case bytecode.OpPushSelf:
			// PUSH_SELF: Push the current receiver (self) onto the stack
			// Operand: unused
			//
			// In methods, self refers to the current object instance.
			// Outside of methods, self is nil.
			if err := vm.push(vm.self); err != nil {
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
				return vm.runtimeError(fmt.Sprintf("selector index out of bounds: %d", selectorIdx))
			}
			selector, ok := vm.constants[selectorIdx].(string)
			if !ok {
				return vm.runtimeError("expected string constant for selector")
			}

			// Pop arguments in reverse order
			// They were pushed left-to-right, so we pop right-to-left
			// to get them back in the correct order
			args := make([]interface{}, argCount)
			for i := argCount - 1; i >= 0; i-- {
				arg, err := vm.pop()
				if err != nil {
					return vm.runtimeError(err.Error())
				}
				args[i] = arg
			}

			// Pop receiver
			receiver, err := vm.pop()
			if err != nil {
				return vm.runtimeError(err.Error())
			}

			// Push call frame for stack trace
			vm.pushFrame("message send", selector)

			// Execute the message send
			result, err := vm.send(receiver, selector, args)
			
			// Pop call frame
			vm.popFrame()
			
			if err != nil {
				// Preserve NonLocalReturn errors without wrapping
				if _, isNonLocal := err.(*NonLocalReturn); isNonLocal {
					return err
				}
				return vm.runtimeError(err.Error())
			}

			// Push result onto stack
			if err := vm.push(result); err != nil {
				return vm.runtimeError(err.Error())
			}

		case bytecode.OpSuperSend:
			// SUPER_SEND: Send a message to the superclass
			// Operand: packed value with selector index and arg count
			//
			// This looks up the method starting from the superclass of the
			// current class context, allowing proper super message sends.

			// Decode operand (same as OpSend)
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
			args := make([]interface{}, argCount)
			for i := argCount - 1; i >= 0; i-- {
				arg, err := vm.pop()
				if err != nil {
					return err
				}
				args[i] = arg
			}

			// Pop receiver (should be self)
			receiver, err := vm.pop()
			if err != nil {
				return err
			}

			// Super sends only work on instances with a current class context
			instance, ok := receiver.(*Instance)
			if !ok {
				return fmt.Errorf("super can only be used within instance methods")
			}

			if vm.currentClass == nil {
				return fmt.Errorf("super used without class context")
			}

			// Dispatch to superclass method
			result, err := vm.superSend(instance, selector, args)
			if err != nil {
				return err
			}

			// Push result onto stack
			if err := vm.push(result); err != nil {
				return err
			}

		case bytecode.OpMakeClosure:
			// MAKE_CLOSURE: Create a block (closure) object
			// Operand: packed value with bytecode index, parent local count, and parameter count
			// Format: [blockIdx (high 16 bits)] [parentLocalCount (bits 8-15)] [paramCount (bits 0-7)]
			//
			// Process:
			//   1. Decode bytecode index, parent local count, and parameter count
			//   2. Get the block bytecode from constants
			//   3. Create a Block object with closure information
			//   4. Push the block onto the stack
			//
			// The block can later be executed with the 'value' message.

			// Decode operand
			bytecodeIdx := inst.Operand >> 16
			parentLocalCount := (inst.Operand >> 8) & 0xFF
			paramCount := inst.Operand & 0xFF

			// Get block bytecode from constants
			if bytecodeIdx < 0 || bytecodeIdx >= len(vm.constants) {
				return fmt.Errorf("bytecode index out of bounds: %d", bytecodeIdx)
			}
			blockBC, ok := vm.constants[bytecodeIdx].(*bytecode.Bytecode)
			if !ok {
				return fmt.Errorf("expected Bytecode in constant pool for block")
			}
			
			block := &Block{
				Bytecode:         blockBC,
				ParamCount:       paramCount,
				ParentLocalCount: parentLocalCount,
				// Capture the home context for non-local returns
				// If we're in a block (vm.homeContext is set), use that
				// Otherwise, use the current VM (we're in a method)
				HomeContext:      vm.homeContext,
			}
			
			// If homeContext is nil, we're in a method or top-level, so set it to current VM
			if block.HomeContext == nil {
				block.HomeContext = vm
			}

			// Push block onto stack
			if err := vm.push(block); err != nil {
				return err
			}

		case bytecode.OpMakeArray:
			// MAKE_ARRAY: Create an array from stack elements
			// Operand: number of elements
			//
			// Process:
			//   1. Pop N elements from stack
			//   2. Create an Array object containing them
			//   3. Push the array onto the stack
			//
			// Stack before: [elem1, elem2, ..., elemN]
			// Stack after:  [array]

			elemCount := inst.Operand

			// Pop elements (in reverse order to maintain order)
			elements := make([]interface{}, elemCount)
			for i := elemCount - 1; i >= 0; i-- {
				elem, err := vm.pop()
				if err != nil {
					return err
				}
				elements[i] = elem
			}

			// Create array object
			array := &Array{Elements: elements}

			// Push array onto stack
			if err := vm.push(array); err != nil {
				return err
			}

		case bytecode.OpMakeDictionary:
			// MAKE_DICTIONARY: Create a dictionary from stack elements
			// Operand: number of key-value pairs
			//
			// Process:
			//   1. Pop 2N elements from stack (N key-value pairs)
			//   2. Create a map/dictionary object containing them
			//   3. Push the dictionary onto the stack
			//
			// Stack before: [key1, value1, key2, value2, ..., keyN, valueN]
			// Stack after:  [dictionary]
			//
			// Note: In Go, map keys must be comparable types (no slices, maps, or functions).
			// Using non-comparable types as dictionary keys will cause a runtime panic.
			// This is a known limitation of the current implementation.

			pairCount := inst.Operand

			// Create the dictionary map
			dict := make(map[interface{}]interface{})

			// Pop key-value pairs (in reverse order)
			for i := pairCount - 1; i >= 0; i-- {
				// Pop value first, then key (they're pushed in key, value order)
				value, err := vm.pop()
				if err != nil {
					return err
				}
				key, err := vm.pop()
				if err != nil {
					return err
				}
				
				// Note: No validation of key type here. Using non-comparable types
				// (slices, maps, functions) will cause a panic.
				// TODO: Add key type validation or use a custom map implementation
				dict[key] = value
			}

			// Push dictionary onto stack
			if err := vm.push(dict); err != nil {
				return err
			}

		case bytecode.OpDefineClass:
			// DEFINE_CLASS: Register a class definition
			// Operand: index into constant pool for ClassDefinition
			//
			// Retrieves the ClassDefinition from constants and registers
			// it in the VM's class registry, making it available for
			// instantiation via the 'new' message.
			if inst.Operand < 0 || inst.Operand >= len(vm.constants) {
				return fmt.Errorf("constant index out of bounds: %d", inst.Operand)
			}

			classDef, ok := vm.constants[inst.Operand].(*bytecode.ClassDefinition)
			if !ok {
				return fmt.Errorf("expected ClassDefinition at constant[%d], got %T", 
					inst.Operand, vm.constants[inst.Operand])
			}

			// Register the class in the global class registry
			vm.classes[classDef.Name] = classDef

			// Also register the class as a global variable so it can be referenced
			vm.globals[classDef.Name] = classDef

		case bytecode.OpLoadField:
			// LOAD_FIELD: Load an instance variable onto the stack
			// Operand: field index
			//
			// Loads a field from the current object (self).
			// Only valid within method context where self is an Instance.
			// Field indices are absolute (methods are compiled with all inherited fields).
			instance, ok := vm.self.(*Instance)
			if !ok {
				return fmt.Errorf("LOAD_FIELD requires self to be an Instance, got %T", vm.self)
			}

			if inst.Operand < 0 || inst.Operand >= len(instance.Fields) {
				return fmt.Errorf("field index out of bounds: %d", inst.Operand)
			}

			if err := vm.push(instance.Fields[inst.Operand]); err != nil {
				return err
			}

		case bytecode.OpStoreField:
			// STORE_FIELD: Store a value to an instance variable
			// Operand: field index
			//
			// Stores the top stack value to a field of the current object (self).
			// The value is popped, stored, then pushed back (assignments return values).
			// Field indices are absolute (methods are compiled with all inherited fields).
			instance, ok := vm.self.(*Instance)
			if !ok {
				return fmt.Errorf("STORE_FIELD requires self to be an Instance, got %T", vm.self)
			}

			if inst.Operand < 0 || inst.Operand >= len(instance.Fields) {
				return fmt.Errorf("field index out of bounds: %d", inst.Operand)
			}

			val, err := vm.pop()
			if err != nil {
				return err
			}

			instance.Fields[inst.Operand] = val

			// Push the value back (assignment returns the value)
			if err := vm.push(val); err != nil {
				return err
			}

		case bytecode.OpLoadClassVar:
			// LOAD_CLASS_VAR: Load a class variable onto the stack
			// Operand: class variable index
			//
			// Loads a class variable from the current class.
			// Class variables are shared across all instances of a class.
			if vm.currentClass == nil {
				return fmt.Errorf("LOAD_CLASS_VAR requires a class context")
			}

			if inst.Operand < 0 || inst.Operand >= len(vm.currentClass.ClassVariables) {
				return fmt.Errorf("class variable index out of bounds: %d", inst.Operand)
			}

			varName := vm.currentClass.ClassVariables[inst.Operand]
			val, exists := vm.currentClass.ClassVarValues[varName]
			if !exists {
				// Class variable not yet initialized - push nil
				val = nil
			}

			if err := vm.push(val); err != nil {
				return err
			}

		case bytecode.OpStoreClassVar:
			// STORE_CLASS_VAR: Store a value to a class variable
			// Operand: class variable index
			//
			// Stores the top stack value to a class variable.
			// The value is popped, stored, then pushed back (assignments return values).
			if vm.currentClass == nil {
				return fmt.Errorf("STORE_CLASS_VAR requires a class context")
			}

			if inst.Operand < 0 || inst.Operand >= len(vm.currentClass.ClassVariables) {
				return fmt.Errorf("class variable index out of bounds: %d", inst.Operand)
			}

			val, err := vm.pop()
			if err != nil {
				return err
			}

			varName := vm.currentClass.ClassVariables[inst.Operand]
			vm.currentClass.ClassVarValues[varName] = val

			// Push the value back (assignment returns the value)
			if err := vm.push(val); err != nil {
				return err
			}

		case bytecode.OpReturn:
			// RETURN: End execution (local return)
			// Operand: unused
			//
			// Exits the execution loop. The final value (if any) remains
			// on the stack and can be retrieved with StackTop().
			// This is a local return - it only exits the current context.
			return nil

		case bytecode.OpNonLocalReturn:
			// NON_LOCAL_RETURN: Perform a non-local return
			// Operand: unused
			//
			// Returns from the method that created the currently executing block,
			// not just from the block itself. This implements Smalltalk-style
			// non-local return semantics.
			//
			// The return value is the top of the stack. We create a NonLocalReturn
			// error with the value and the home context.
			//
			// If vm.homeContext is set (we're in a block), the home context is
			// the VM that created the block. Otherwise (we're in a method or top-level),
			// homeContext is nil and this behaves like a normal return.
			var returnValue interface{}
			if vm.sp > 0 {
				returnValue = vm.stack[vm.sp-1]
			}
			
			if vm.homeContext != nil {
				// We're in a block - return to the home context
				return &NonLocalReturn{
					Value:       returnValue,
					HomeContext: vm.homeContext,
				}
			}
			// We're at method/top level - treat as normal return
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
	// Check if receiver is a Block and selector is 'value' or starts with 'value:'
	if block, ok := receiver.(*Block); ok {
		// Match 'value' (no args) or 'value:' with varying arg counts
		if selector == "value" || (len(selector) >= 6 && selector[:6] == "value:") {
			return vm.executeBlock(block, args)
		}

		// Handle whileTrue: and whileFalse:
		switch selector {
		case "whileTrue:":
			if len(args) != 1 {
				return nil, fmt.Errorf("whileTrue: expects 1 argument (block), got %d", len(args))
			}
			bodyBlock, ok := args[0].(*Block)
			if !ok {
				return nil, fmt.Errorf("whileTrue: argument must be a block")
			}

			// Execute the condition block, and while it returns true, execute the body
			for {
				result, err := vm.executeBlock(block, []interface{}{})
				if err != nil {
					return nil, err
				}

				// Check if result is a boolean true
				conditionTrue, ok := result.(bool)
				if !ok {
					return nil, fmt.Errorf("whileTrue: condition block must return a boolean")
				}

				if !conditionTrue {
					break
				}

				// Execute the body block
				_, err = vm.executeBlock(bodyBlock, []interface{}{})
				if err != nil {
					return nil, err
				}
			}
			return nil, nil

		case "whileFalse:":
			if len(args) != 1 {
				return nil, fmt.Errorf("whileFalse: expects 1 argument (block), got %d", len(args))
			}
			bodyBlock, ok := args[0].(*Block)
			if !ok {
				return nil, fmt.Errorf("whileFalse: argument must be a block")
			}

			// Execute the condition block, and while it returns false, execute the body
			for {
				result, err := vm.executeBlock(block, []interface{}{})
				if err != nil {
					return nil, err
				}

				// Check if result is a boolean false
				conditionFalse, ok := result.(bool)
				if !ok {
					return nil, fmt.Errorf("whileFalse: condition block must return a boolean")
				}

				if conditionFalse {
					break
				}

				// Execute the body block
				_, err = vm.executeBlock(bodyBlock, []interface{}{})
				if err != nil {
					return nil, err
				}
			}
			return nil, nil
		}
	}

	// Check if receiver is a Boolean and handle boolean control flow
	if b, ok := receiver.(bool); ok {
		switch selector {
		case "ifTrue:":
			if len(args) != 1 {
				return nil, fmt.Errorf("ifTrue: expects 1 argument (block), got %d", len(args))
			}
			block, ok := args[0].(*Block)
			if !ok {
				return nil, fmt.Errorf("ifTrue: argument must be a block")
			}
			if b {
				return vm.executeBlock(block, []interface{}{})
			}
			return nil, nil
		case "ifFalse:":
			if len(args) != 1 {
				return nil, fmt.Errorf("ifFalse: expects 1 argument (block), got %d", len(args))
			}
			block, ok := args[0].(*Block)
			if !ok {
				return nil, fmt.Errorf("ifFalse: argument must be a block")
			}
			if !b {
				return vm.executeBlock(block, []interface{}{})
			}
			return nil, nil
		case "ifTrue:ifFalse:":
			if len(args) != 2 {
				return nil, fmt.Errorf("ifTrue:ifFalse: expects 2 arguments (blocks), got %d", len(args))
			}
			trueBlock, ok1 := args[0].(*Block)
			falseBlock, ok2 := args[1].(*Block)
			if !ok1 || !ok2 {
				return nil, fmt.Errorf("ifTrue:ifFalse: arguments must be blocks")
			}
			if b {
				return vm.executeBlock(trueBlock, []interface{}{})
			}
			return vm.executeBlock(falseBlock, []interface{}{})
		}
	}

	// Check if receiver is an Integer and handle integer messages
	if num, ok := receiver.(int64); ok {
		switch selector {
		case "timesRepeat:":
			if len(args) != 1 {
				return nil, fmt.Errorf("timesRepeat: expects 1 argument (block), got %d", len(args))
			}
			block, ok := args[0].(*Block)
			if !ok {
				return nil, fmt.Errorf("timesRepeat: argument must be a block")
			}
			for i := int64(0); i < num; i++ {
				_, err := vm.executeBlock(block, []interface{}{})
				if err != nil {
					return nil, err
				}
			}
			return nil, nil
		}
	}

	// Check if receiver is an Array and handle array messages
	if array, ok := receiver.(*Array); ok {
		switch selector {
		case "size":
			return int64(len(array.Elements)), nil
		case "at:":
			// Array indexing (1-based like Smalltalk)
			if len(args) != 1 {
				return nil, fmt.Errorf("at: expects 1 argument, got %d", len(args))
			}
			idx, ok := args[0].(int64)
			if !ok {
				return nil, fmt.Errorf("array index must be integer")
			}
			if idx < 1 || idx > int64(len(array.Elements)) {
				return nil, fmt.Errorf("array index out of bounds: %d", idx)
			}
			return array.Elements[idx-1], nil
		case "at:put:":
			// Array element assignment (1-based like Smalltalk)
			if len(args) != 2 {
				return nil, fmt.Errorf("at:put: expects 2 arguments, got %d", len(args))
			}
			idx, ok := args[0].(int64)
			if !ok {
				return nil, fmt.Errorf("array index must be integer")
			}
			if idx < 1 || idx > int64(len(array.Elements)) {
				return nil, fmt.Errorf("array index out of bounds: %d", idx)
			}
			value := args[1]
			array.Elements[idx-1] = value
			return value, nil
		case "do:":
			// Iterate over array elements with a block
			if len(args) != 1 {
				return nil, fmt.Errorf("do: expects 1 argument (block), got %d", len(args))
			}
			block, ok := args[0].(*Block)
			if !ok {
				return nil, fmt.Errorf("do: argument must be a block")
			}
			for _, elem := range array.Elements {
				_, err := vm.executeBlock(block, []interface{}{elem})
				if err != nil {
					return nil, err
				}
			}
			return array, nil
		}
	}

	// Check if receiver is a ClassDefinition (class object)
	if classDef, ok := receiver.(*bytecode.ClassDefinition); ok {
		switch selector {
		case "new":
			// Create a new instance of the class
			// Allocate fields for this class and all superclasses
			totalFields := vm.countAllFields(classDef)
			instance := &Instance{
				Class:  classDef,
				Fields: make([]interface{}, totalFields),
			}
			return instance, nil
		default:
			// Look up class method
			return vm.executeClassMethod(classDef, selector, args)
		}
	}

	// Check if receiver is an Instance (object instance)
	if instance, ok := receiver.(*Instance); ok {
		// Look up method in the instance's class
		return vm.executeMethod(instance, selector, args)
	}

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

	// HTTP primitives
	case "httpGet:":
		if len(args) != 1 {
			return nil, fmt.Errorf("httpGet: expects 1 argument")
		}
		url, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("httpGet: URL must be a string")
		}
		return vm.httpGet(url)

	case "httpPost:body:":
		if len(args) != 2 {
			return nil, fmt.Errorf("httpPost:body: expects 2 arguments")
		}
		url, ok1 := args[0].(string)
		body, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("httpPost:body: arguments must be strings")
		}
		return vm.httpPost(url, body)

	// Crypto primitives
	case "aesEncrypt:key:":
		if len(args) != 2 {
			return nil, fmt.Errorf("aesEncrypt:key: expects 2 arguments")
		}
		data, ok1 := args[0].(string)
		key, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("aesEncrypt:key: arguments must be strings")
		}
		return vm.aesEncrypt(data, key)

	case "aesDecrypt:key:":
		if len(args) != 2 {
			return nil, fmt.Errorf("aesDecrypt:key: expects 2 arguments")
		}
		data, ok1 := args[0].(string)
		key, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("aesDecrypt:key: arguments must be strings")
		}
		return vm.aesDecrypt(data, key)

	case "aesGenerateKey":
		return vm.aesGenerateKey()

	case "sha256:":
		if len(args) != 1 {
			return nil, fmt.Errorf("sha256: expects 1 argument")
		}
		data, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("sha256: argument must be a string")
		}
		return vm.sha256Hash(data), nil

	case "sha512:":
		if len(args) != 1 {
			return nil, fmt.Errorf("sha512: expects 1 argument")
		}
		data, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("sha512: argument must be a string")
		}
		return vm.sha512Hash(data), nil

	case "md5:":
		if len(args) != 1 {
			return nil, fmt.Errorf("md5: expects 1 argument")
		}
		data, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("md5: argument must be a string")
		}
		return vm.md5Hash(data), nil

	case "base64Encode:":
		if len(args) != 1 {
			return nil, fmt.Errorf("base64Encode: expects 1 argument")
		}
		data, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("base64Encode: argument must be a string")
		}
		return vm.base64Encode(data), nil

	case "base64Decode:":
		if len(args) != 1 {
			return nil, fmt.Errorf("base64Decode: expects 1 argument")
		}
		data, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("base64Decode: argument must be a string")
		}
		return vm.base64Decode(data)

	// Compression primitives
	case "zipCompress:":
		if len(args) != 1 {
			return nil, fmt.Errorf("zipCompress: expects 1 argument")
		}
		data, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("zipCompress: argument must be a string")
		}
		return vm.zipCompress(data)

	case "zipDecompress:":
		if len(args) != 1 {
			return nil, fmt.Errorf("zipDecompress: expects 1 argument")
		}
		data, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("zipDecompress: argument must be a string")
		}
		return vm.zipDecompress(data)

	case "gzipCompress:":
		if len(args) != 1 {
			return nil, fmt.Errorf("gzipCompress: expects 1 argument")
		}
		data, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("gzipCompress: argument must be a string")
		}
		return vm.gzipCompress(data)

	case "gzipDecompress:":
		if len(args) != 1 {
			return nil, fmt.Errorf("gzipDecompress: expects 1 argument")
		}
		data, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("gzipDecompress: argument must be a string")
		}
		return vm.gzipDecompress(data)

	// File I/O primitives
	case "fileRead:":
		if len(args) != 1 {
			return nil, fmt.Errorf("fileRead: expects 1 argument")
		}
		path, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("fileRead: path must be a string")
		}
		return vm.fileRead(path)

	case "fileWrite:content:":
		if len(args) != 2 {
			return nil, fmt.Errorf("fileWrite:content: expects 2 arguments")
		}
		path, ok1 := args[0].(string)
		content, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("fileWrite:content: arguments must be strings")
		}
		err := vm.fileWrite(path, content)
		if err != nil {
			return nil, err
		}
		return nil, nil

	case "fileExists:":
		if len(args) != 1 {
			return nil, fmt.Errorf("fileExists: expects 1 argument")
		}
		path, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("fileExists: path must be a string")
		}
		return vm.fileExists(path), nil

	case "fileDelete:":
		if len(args) != 1 {
			return nil, fmt.Errorf("fileDelete: expects 1 argument")
		}
		path, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("fileDelete: path must be a string")
		}
		err := vm.fileDelete(path)
		if err != nil {
			return nil, err
		}
		return nil, nil

	// JSON primitives
	case "jsonParse:":
		if len(args) != 1 {
			return nil, fmt.Errorf("jsonParse: expects 1 argument")
		}
		data, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("jsonParse: argument must be a string")
		}
		return vm.jsonParse(data)

	case "jsonGenerate:":
		if len(args) != 1 {
			return nil, fmt.Errorf("jsonGenerate: expects 1 argument")
		}
		return vm.jsonGenerate(args[0])

	// Regex primitives
	case "regexMatch:text:":
		if len(args) != 2 {
			return nil, fmt.Errorf("regexMatch:text: expects 2 arguments")
		}
		pattern, ok1 := args[0].(string)
		text, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("regexMatch:text: arguments must be strings")
		}
		return vm.regexMatch(pattern, text)

	case "regexFindAll:text:":
		if len(args) != 2 {
			return nil, fmt.Errorf("regexFindAll:text: expects 2 arguments")
		}
		pattern, ok1 := args[0].(string)
		text, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("regexFindAll:text: arguments must be strings")
		}
		return vm.regexFindAll(pattern, text)

	case "regexReplace:text:with:":
		if len(args) != 3 {
			return nil, fmt.Errorf("regexReplace:text:with: expects 3 arguments")
		}
		pattern, ok1 := args[0].(string)
		text, ok2 := args[1].(string)
		replacement, ok3 := args[2].(string)
		if !ok1 || !ok2 || !ok3 {
			return nil, fmt.Errorf("regexReplace:text:with: arguments must be strings")
		}
		return vm.regexReplace(pattern, text, replacement)

	// Random number generation primitives
	case "randomInt:max:":
		if len(args) != 2 {
			return nil, fmt.Errorf("randomInt:max: expects 2 arguments")
		}
		min, ok1 := args[0].(int64)
		max, ok2 := args[1].(int64)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("randomInt:max: arguments must be integers")
		}
		return vm.randomInt(min, max)

	case "randomFloat":
		return vm.randomFloat()

	case "randomBytes:":
		if len(args) != 1 {
			return nil, fmt.Errorf("randomBytes: expects 1 argument")
		}
		length, ok := args[0].(int64)
		if !ok {
			return nil, fmt.Errorf("randomBytes: argument must be an integer")
		}
		return vm.randomBytes(length)

	// Date/Time primitives
	case "dateNow":
		return vm.dateNow(), nil

	case "dateFormat:format:":
		if len(args) != 2 {
			return nil, fmt.Errorf("dateFormat:format: expects 2 arguments")
		}
		timestamp, ok1 := args[0].(int64)
		format, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("dateFormat:format: arguments must be integer and string")
		}
		return vm.dateFormat(timestamp, format), nil

	case "dateParse:format:":
		if len(args) != 2 {
			return nil, fmt.Errorf("dateParse:format: expects 2 arguments")
		}
		dateStr, ok1 := args[0].(string)
		format, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("dateParse:format: arguments must be strings")
		}
		return vm.dateParse(dateStr, format)

	case "timeYear:":
		if len(args) != 1 {
			return nil, fmt.Errorf("timeYear: expects 1 argument")
		}
		timestamp, ok := args[0].(int64)
		if !ok {
			return nil, fmt.Errorf("timeYear: argument must be an integer")
		}
		return vm.timeYear(timestamp), nil

	case "timeMonth:":
		if len(args) != 1 {
			return nil, fmt.Errorf("timeMonth: expects 1 argument")
		}
		timestamp, ok := args[0].(int64)
		if !ok {
			return nil, fmt.Errorf("timeMonth: argument must be an integer")
		}
		return vm.timeMonth(timestamp), nil

	case "timeDay:":
		if len(args) != 1 {
			return nil, fmt.Errorf("timeDay: expects 1 argument")
		}
		timestamp, ok := args[0].(int64)
		if !ok {
			return nil, fmt.Errorf("timeDay: argument must be an integer")
		}
		return vm.timeDay(timestamp), nil

	case "timeHour:":
		if len(args) != 1 {
			return nil, fmt.Errorf("timeHour: expects 1 argument")
		}
		timestamp, ok := args[0].(int64)
		if !ok {
			return nil, fmt.Errorf("timeHour: argument must be an integer")
		}
		return vm.timeHour(timestamp), nil

	case "timeMinute:":
		if len(args) != 1 {
			return nil, fmt.Errorf("timeMinute: expects 1 argument")
		}
		timestamp, ok := args[0].(int64)
		if !ok {
			return nil, fmt.Errorf("timeMinute: argument must be an integer")
		}
		return vm.timeMinute(timestamp), nil

	case "timeSecond:":
		if len(args) != 1 {
			return nil, fmt.Errorf("timeSecond: expects 1 argument")
		}
		timestamp, ok := args[0].(int64)
		if !ok {
			return nil, fmt.Errorf("timeSecond: argument must be an integer")
		}
		return vm.timeSecond(timestamp), nil

	default:
		return nil, fmt.Errorf("unknown message: %s", selector)
	}
}

// tryPrimitive attempts to execute a primitive operation.
// Returns (result, nil) if the primitive was handled, or (nil, error) if not a primitive.
// This allows falling back to method lookup when primitives don't apply.
func (vm *VM) tryPrimitive(receiver interface{}, selector string, args []interface{}) (interface{}, error) {
	// Handle primitive operations
	// These are built directly into the VM for efficiency
	switch selector {
	case "+":
		if len(args) != 1 {
			return nil, fmt.Errorf("not a primitive")
		}
		return vm.add(receiver, args[0])
	case "-":
		if len(args) != 1 {
			return nil, fmt.Errorf("not a primitive")
		}
		return vm.subtract(receiver, args[0])
	case "*":
		if len(args) != 1 {
			return nil, fmt.Errorf("not a primitive")
		}
		return vm.multiply(receiver, args[0])
	case "/":
		if len(args) != 1 {
			return nil, fmt.Errorf("not a primitive")
		}
		return vm.divide(receiver, args[0])
	case "<":
		if len(args) != 1 {
			return nil, fmt.Errorf("not a primitive")
		}
		return vm.lessThan(receiver, args[0])
	case ">":
		if len(args) != 1 {
			return nil, fmt.Errorf("not a primitive")
		}
		return vm.greaterThan(receiver, args[0])
	case "<=":
		if len(args) != 1 {
			return nil, fmt.Errorf("not a primitive")
		}
		return vm.lessOrEqual(receiver, args[0])
	case ">=":
		if len(args) != 1 {
			return nil, fmt.Errorf("not a primitive")
		}
		return vm.greaterOrEqual(receiver, args[0])
	case "=":
		if len(args) != 1 {
			return nil, fmt.Errorf("not a primitive")
		}
		return vm.equal(receiver, args[0])
	case "~=":
		if len(args) != 1 {
			return nil, fmt.Errorf("not a primitive")
		}
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
		// Not a basic primitive
		return nil, fmt.Errorf("not a primitive")
	}
}

// executeBlock executes a block with the given arguments.
//
// Process:
//   1. Check argument count matches parameter count
//   2. Create a new VM instance for the block execution
//   3. Set up parameters as local variables BEFORE calling Run()
//   4. Run the block's bytecode
//   5. Return the result
//
// Parameters:
//   - block: The Block object to execute
//   - args: Arguments to pass to the block
//
// Returns:
//   - The result of executing the block
//   - Error if execution fails or argument count doesn't match
func (vm *VM) executeBlock(block *Block, args []interface{}) (interface{}, error) {
	// Check argument count
	if len(args) != block.ParamCount {
		return nil, fmt.Errorf("block expects %d arguments, got %d", block.ParamCount, len(args))
	}

	// Create a new VM for block execution
	// Blocks share the parent's locals array to support closures
	// This allows blocks to access and modify variables from the enclosing scope
	blockVM := &VM{
		stack:       make([]interface{}, 1024),
		sp:          0,
		locals:      vm.locals,  // Share locals with parent for closure support
		globals:     vm.globals, // Share globals with parent VM
		constants:   block.Bytecode.Constants, // Will be overwritten by Run() anyway
		classes:     vm.classes, // Share class registry
		self:        vm.self,    // Share self reference
		homeContext: block.HomeContext, // Set the home context for non-local returns
	}

	// Block parameters are stored starting at the parent's local count
	// The compiler allocated them at slots starting from parent's localCount
	// We use the ParentLocalCount stored in the block
	parentLocalCount := block.ParentLocalCount
	requiredSize := parentLocalCount + block.ParamCount
	
	if cap(vm.locals) < requiredSize {
		// Need to expand capacity
		newLocals := make([]interface{}, requiredSize)
		copy(newLocals, vm.locals)
		vm.locals = newLocals
		blockVM.locals = newLocals  // Share the new array with blockVM
	} else if len(vm.locals) < requiredSize {
		// Just extend the slice
		vm.locals = vm.locals[:requiredSize]
		blockVM.locals = vm.locals  // Ensure blockVM has the extended slice
	}

	// Set block parameters in the locals array
	// They start at parentLocalCount
	for i, arg := range args {
		blockVM.locals[parentLocalCount+i] = arg
	}

	// Execute the block bytecode
	if err := blockVM.Run(block.Bytecode); err != nil {
		// Check if this is a non-local return
		if nlr, ok := err.(*NonLocalReturn); ok {
			// Non-local returns always propagate up through blocks.
			// The method execution (executeMethod) will catch it and convert
			// to a normal return when nlr.HomeContext matches the method's VM.
			return nil, nlr
		}
		// Other errors propagate normally
		return nil, err
	}

	// Restore locals length to what it was before (cleanup block parameters)
	vm.locals = vm.locals[:parentLocalCount]

	// Return the top value from the block's stack
	result := blockVM.StackTop()
	if result == nil {
		// Blocks return nil if they don't have an explicit result
		return nil, nil
	}
	return result, nil
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

// Block represents a runtime block (closure) object.
//
// Blocks are first-class objects that encapsulate code and can be
// executed later. They capture their environment (closure semantics).
//
// A block contains:
//   - Bytecode: The compiled code to execute
//   - ParamCount: Number of parameters the block expects
//   - ParentLocalCount: Number of locals in the parent context (for closure support)
//   - HomeContext: The VM context where the block was created (for non-local returns)
type Block struct {
	Bytecode         *bytecode.Bytecode // The block's compiled code
	ParamCount       int                // Number of parameters
	ParentLocalCount int                // Number of locals in parent context
	HomeContext      *VM                // The VM context that created this block (for non-local returns)
}

// NonLocalReturn is a special error type used to implement non-local returns.
//
// In Smalltalk-style languages, a return statement (^) inside a block doesn't
// just return from the block - it returns from the method that created the block.
// This is called a "non-local return" because it exits from a context other than
// the immediately executing one.
//
// When a block executes OpNonLocalReturn, it creates a NonLocalReturn error
// containing the return value. This error propagates up through executeBlock()
// calls until it reaches the method execution context, where it's caught and
// converted into a normal return.
//
// Example flow:
//   1. Method M creates a block B and passes it to ifTrue:
//   2. ifTrue: calls executeBlock(B)
//   3. Block B executes OpNonLocalReturn with value 42
//   4. NonLocalReturn{Value: 42, HomeContext: M's VM} is created
//   5. executeBlock returns this as an error
//   6. ifTrue: propagates the error up
//   7. Method M catches it and returns 42
type NonLocalReturn struct {
	Value       interface{} // The value to return
	HomeContext *VM         // The target context to return to (the method's VM)
}

// Error implements the error interface for NonLocalReturn.
func (nlr *NonLocalReturn) Error() string {
	return "non-local return"
}


// Array represents a runtime array object.
//
// Arrays are ordered collections of values.
type Array struct {
	Elements []interface{} // The array elements
}

// Instance represents a runtime object instance.
//
// An Instance is created from a ClassDefinition and contains:
//   - Class: Reference to the class definition
//   - Fields: Values of the instance variables
//
// Example:
//   For a Counter class with one field 'count':
//     Instance{Class: CounterClassDef, Fields: [0]}
type Instance struct {
	Class  *bytecode.ClassDefinition // The class this is an instance of
	Fields []interface{}              // Instance variable values
}

// count AllFields counts total fields in class hierarchy.
//
// This counts all instance variables from this class and all superclasses.
// Fields are ordered from superclass to subclass.
func (vm *VM) countAllFields(class *bytecode.ClassDefinition) int {
	total := len(class.Fields)
	currentClass := class
	
	// Walk up the hierarchy counting fields
	for currentClass.SuperClass != "" && currentClass.SuperClass != "Object" {
		superClass, exists := vm.classes[currentClass.SuperClass]
		if !exists {
			break
		}
		total += len(superClass.Fields)
		currentClass = superClass
	}
	
	return total
}

// getFieldOffset calculates the field offset for a class in the inheritance hierarchy.
//
// This returns the starting index for this class's fields in the instance field array.
// Superclass fields come first, so the offset is the sum of all superclass field counts.
func (vm *VM) getFieldOffset(class *bytecode.ClassDefinition) int {
	offset := 0
	currentClass := class
	
	// Walk up the hierarchy counting superclass fields
	for currentClass.SuperClass != "" && currentClass.SuperClass != "Object" {
		superClass, exists := vm.classes[currentClass.SuperClass]
		if !exists {
			break
		}
		offset += len(superClass.Fields)
		currentClass = superClass
	}
	
	return offset
}

// lookupMethod searches for a method in a class and its superclass chain.
//
// This implements the method lookup algorithm for inheritance:
//   1. Search for the method in the given class
//   2. If not found and class has a superclass, search in superclass
//   3. Continue up the hierarchy until method is found or chain ends
//
// Parameters:
//   - class: The class to start searching from
//   - selector: The method name to find
//
// Returns:
//   - The method definition if found, nil otherwise
//   - The class where the method was found (for super sends)
func (vm *VM) lookupMethod(class *bytecode.ClassDefinition, selector string) (*bytecode.MethodDefinition, *bytecode.ClassDefinition) {
	currentClass := class
	
	// Walk up the class hierarchy
	for currentClass != nil {
		// Search for method in current class
		for _, m := range currentClass.Methods {
			if m.Selector == selector {
				return m, currentClass
			}
		}
		
		// Method not found in this class, try superclass
		if currentClass.SuperClass == "" || currentClass.SuperClass == "Object" {
			// No superclass or reached Object (root of hierarchy)
			break
		}
		
		// Get the superclass definition
		superClass, exists := vm.classes[currentClass.SuperClass]
		if !exists {
			// Superclass not found - stop searching
			break
		}
		
		currentClass = superClass
	}
	
	// Method not found in hierarchy
	return nil, nil
}

// superSend executes a method from the superclass.
//
// This implements super message sends by starting the method lookup
// from the superclass of the current class context.
//
// Parameters:
//   - instance: The object instance (self)
//   - selector: The method name
//   - args: Arguments to the method
//
// Returns:
//   - The method's return value
//   - Error if method not found or execution fails
func (vm *VM) superSend(instance *Instance, selector string, args []interface{}) (interface{}, error) {
	// Get the superclass of the current class context
	if vm.currentClass.SuperClass == "" || vm.currentClass.SuperClass == "Object" {
		return nil, fmt.Errorf("class %s has no superclass to send '%s' to", 
			vm.currentClass.Name, selector)
	}

	superClass, exists := vm.classes[vm.currentClass.SuperClass]
	if !exists {
		return nil, fmt.Errorf("superclass %s not found for class %s", 
			vm.currentClass.SuperClass, vm.currentClass.Name)
	}

	// Look up the method starting from superclass
	method, class := vm.lookupMethod(superClass, selector)

	if method == nil {
		return nil, fmt.Errorf("superclass of %s does not understand message '%s'", 
			vm.currentClass.Name, selector)
	}

	// Check argument count
	if len(args) != len(method.Parameters) {
		return nil, fmt.Errorf("method %s expects %d arguments, got %d", 
			selector, len(method.Parameters), len(args))
	}

	// Create a new VM for method execution
	methodVM := New()
	methodVM.globals = vm.globals       // Share global variables
	methodVM.classes = vm.classes       // Share class registry
	methodVM.self = instance            // Set self to the instance
	methodVM.currentClass = class       // Set class context to where method was found
	// No field offset needed - methods are compiled with all fields

	// Set up method parameters as local variables
	for i, arg := range args {
		methodVM.locals[i] = arg
	}

	// Execute the method bytecode
	if err := methodVM.Run(method.Code); err != nil {
		// Check if this is a non-local return targeting this method
		if nlr, ok := err.(*NonLocalReturn); ok {
			// If the non-local return's home context is this method's VM,
			// then this is where the return should stop - convert it to a normal return
			if nlr.HomeContext == methodVM {
				// This non-local return is for us - use its value as the return value
				return nlr.Value, nil
			}
			// Otherwise, propagate it further up
			return nil, nlr
		}
		return nil, fmt.Errorf("error in super method %s: %w", selector, err)
	}

	// Return the result (top of stack)
	if methodVM.sp > 0 {
		return methodVM.stack[methodVM.sp-1], nil
	}

	// No value on stack - return nil
	return nil, nil
}

// executeMethod executes a user-defined method on an instance.
//
// This implements the method lookup and dispatch for user-defined classes:
//   1. Find the method by selector in the instance's class
//   2. Check argument count matches parameter count
//   3. Create a new VM context for method execution
//   4. Set self to the instance
//   5. Pass arguments as local variables
//   6. Execute the method bytecode
//   7. Return the result
//
// Parameters:
//   - instance: The object instance receiving the message
//   - selector: The method name
//   - args: Arguments to the method
//
// Returns:
//   - The method's return value
//   - Error if method not found or execution fails
func (vm *VM) executeMethod(instance *Instance, selector string, args []interface{}) (interface{}, error) {
	// Look up the method in the instance's class hierarchy
	method, class := vm.lookupMethod(instance.Class, selector)

	if method == nil {
		// Method not found in class hierarchy - try primitives
		result, err := vm.tryPrimitive(instance, selector, args)
		if err == nil {
			// Primitive handled it
			return result, nil
		}
		// Not a primitive - report error
		return nil, fmt.Errorf("instance of %s does not understand message '%s'", 
			instance.Class.Name, selector)
	}

	// Check argument count
	if len(args) != len(method.Parameters) {
		return nil, fmt.Errorf("method %s expects %d arguments, got %d", 
			selector, len(method.Parameters), len(args))
	}

	// Create a new VM for method execution to isolate its stack and locals
	methodVM := New()
	methodVM.globals = vm.globals       // Share global variables
	methodVM.classes = vm.classes       // Share class registry
	methodVM.self = instance            // Set self to the instance
	methodVM.currentClass = class       // Set current class context for super sends
	// No field offset needed - methods are compiled with all fields

	// Set up method parameters as local variables
	for i, arg := range args {
		methodVM.locals[i] = arg
	}

	// Execute the method bytecode
	if err := methodVM.Run(method.Code); err != nil {
		// Check if this is a non-local return targeting this method
		if nlr, ok := err.(*NonLocalReturn); ok {
			// Check if this non-local return targets this method's execution.
			// When a block is created during this method's execution, it captures
			// methodVM as its HomeContext. The pointer comparison works because:
			// 1. methodVM is created once at the start of this method execution
			// 2. All blocks created during this execution capture the same methodVM pointer
			// 3. When OpNonLocalReturn executes, it uses the captured HomeContext pointer
			// Thus, nlr.HomeContext == methodVM means "return from this method execution"
			if nlr.HomeContext == methodVM {
				// This non-local return is for us - use its value as the return value
				return nlr.Value, nil
			}
			// Otherwise, propagate it further up (shouldn't normally happen in well-formed code)
			return nil, nlr
		}
		return nil, fmt.Errorf("error in method %s: %w", selector, err)
	}

	// Return the result (top of stack)
	if methodVM.sp > 0 {
		return methodVM.stack[methodVM.sp-1], nil
	}

	// No value on stack - return nil
	return nil, nil
}

// executeClassMethod executes a class method.
//
// Class methods are defined on the class itself rather than instances.
// They have access to class variables but not instance variables.
//
// Parameters:
//   - classDef: The class definition
//   - selector: The method name
//   - args: Arguments to the method
//
// Returns:
//   - The method's return value
//   - Error if method not found or execution fails
func (vm *VM) executeClassMethod(classDef *bytecode.ClassDefinition, selector string, args []interface{}) (interface{}, error) {
	// Look up the class method
	var method *bytecode.MethodDefinition
	for _, m := range classDef.ClassMethods {
		if m.Selector == selector {
			method = m
			break
		}
	}

	if method == nil {
		// Class method not found
		return nil, fmt.Errorf("class %s does not understand class message '%s'", 
			classDef.Name, selector)
	}

	// Check argument count
	if len(args) != len(method.Parameters) {
		return nil, fmt.Errorf("class method %s expects %d arguments, got %d", 
			selector, len(method.Parameters), len(args))
	}

	// Create a new VM for method execution
	methodVM := New()
	methodVM.globals = vm.globals       // Share global variables
	methodVM.classes = vm.classes       // Share class registry
	methodVM.self = classDef            // Set self to the class
	methodVM.currentClass = classDef    // Set class context

	// Set up method parameters as local variables
	for i, arg := range args {
		methodVM.locals[i] = arg
	}

	// Execute the method bytecode
	if err := methodVM.Run(method.Code); err != nil {
		// Check if this is a non-local return targeting this method
		if nlr, ok := err.(*NonLocalReturn); ok {
			// If the non-local return's home context is this method's VM,
			// then this is where the return should stop - convert it to a normal return
			if nlr.HomeContext == methodVM {
				// This non-local return is for us - use its value as the return value
				return nlr.Value, nil
			}
			// Otherwise, propagate it further up
			return nil, nlr
		}
		return nil, fmt.Errorf("error in class method %s: %w", selector, err)
	}

	// Return the result (top of stack)
	if methodVM.sp > 0 {
		return methodVM.stack[methodVM.sp-1], nil
	}

	// No value on stack - return nil
	return nil, nil
}

// GetGlobal retrieves a global variable by name.
//
// This is primarily for testing purposes.
//
// Parameters:
//   - name: The global variable name
//
// Returns:
//   - The value of the global, or nil if not found
func (vm *VM) GetGlobal(name string) interface{} {
	return vm.globals[name]
}

// pushFrame adds a new call frame to the call stack.
// This is used for stack trace generation.
func (vm *VM) pushFrame(name, selector string) {
	frame := StackFrame{
		Name:     name,
		Selector: selector,
		IP:       vm.ip,
	}
	vm.callStack = append(vm.callStack, frame)
}

// popFrame removes the top call frame from the call stack.
func (vm *VM) popFrame() {
	if len(vm.callStack) > 0 {
		vm.callStack = vm.callStack[:len(vm.callStack)-1]
	}
}

// runtimeError creates a RuntimeError with the current call stack.
func (vm *VM) runtimeError(message string) error {
	// Make a copy of the call stack
	stack := make([]StackFrame, len(vm.callStack))
	copy(stack, vm.callStack)
	
	// Add current instruction pointer to the last frame if there is one
	if len(stack) > 0 {
		stack[len(stack)-1].IP = vm.ip
	}
	
	return newRuntimeError(message, stack)
}

// EnableDebugger creates and enables a debugger for this VM.
func (vm *VM) EnableDebugger() *Debugger {
	if vm.debugger == nil {
		vm.debugger = NewDebugger(vm)
	}
	vm.debugger.Enable()
	return vm.debugger
}

// GetDebugger returns the debugger instance if debugging is enabled.
func (vm *VM) GetDebugger() *Debugger {
	return vm.debugger
}

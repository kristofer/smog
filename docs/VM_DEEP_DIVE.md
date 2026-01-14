# Smog Virtual Machine Deep Dive

## Overview

The Smog Virtual Machine (VM) is a stack-based bytecode interpreter that executes compiled Smog programs. This guide provides an in-depth look at how the VM works, from architecture to execution mechanics.

## Purpose and Role

The VM is the final stage in the Smog execution pipeline:

```
Source → Lexer → Parser → AST → Compiler → Bytecode → [.sg file] → **VM** → Results
```

The VM can execute bytecode from two sources:
1. **Freshly compiled** - Bytecode generated from source code
2. **Pre-compiled** - Bytecode loaded from .sg files

Its primary responsibilities are:
1. **Bytecode Execution**: Interpret and execute bytecode instructions
2. **Stack Management**: Maintain operand stack for computation
3. **Memory Management**: Handle variables and objects in memory
4. **Message Dispatch**: Route messages to appropriate method implementations
5. **Runtime Services**: Provide primitives (I/O, arithmetic, etc.)

## Architecture

### Core Components

```
┌─────────────────────────────────────┐
│         Virtual Machine             │
├─────────────────────────────────────┤
│  ┌───────────────────────────────┐  │
│  │    Execution Stack (1024)    │  │  ← Operand storage
│  └───────────────────────────────┘  │
│  ┌───────────────────────────────┐  │
│  │   Local Variables (256)      │  │  ← Method locals
│  └───────────────────────────────┘  │
│  ┌───────────────────────────────┐  │
│  │   Global Variables (map)     │  │  ← Globals
│  └───────────────────────────────┘  │
│  ┌───────────────────────────────┐  │
│  │   Instruction Pointer (IP)   │  │  ← Program counter
│  └───────────────────────────────┘  │
│  ┌───────────────────────────────┐  │
│  │   Constant Pool (read-only)  │  │  ← Literals
│  └───────────────────────────────┘  │
└─────────────────────────────────────┘
```

### Memory Model

**1. Execution Stack**
- **Type**: Fixed-size array (1024 slots)
- **Purpose**: Store intermediate values during computation
- **Access Pattern**: LIFO (Last In, First Out)
- **Operations**: push, pop, peek

**2. Local Variables**
- **Type**: Fixed-size array (256 slots)
- **Purpose**: Store method/block local variables
- **Scope**: Reset on each method entry
- **Access**: Direct indexing (fast)

**3. Global Variables**
- **Type**: Hash map (string → value)
- **Purpose**: Store program-wide variables
- **Scope**: Persistent across method calls
- **Access**: By name lookup (slower than locals)

**4. Constant Pool**
- **Type**: Read-only array
- **Purpose**: Store literal values from bytecode
- **Content**: Numbers, strings, symbols, selectors
- **Shared**: Same pool used by all code

## Execution Model

### Fetch-Decode-Execute Cycle

The VM operates in a continuous loop:

```go
func (vm *VM) Run(bytecode *Bytecode) error {
    vm.ip = 0
    vm.constants = bytecode.Constants
    instructions := bytecode.Instructions
    
    for vm.ip < len(instructions) {
        // 1. FETCH
        instruction := instructions[vm.ip]
        opcode := instruction.Opcode
        operand := instruction.Operand
        
        // 2. DECODE & EXECUTE
        switch opcode {
        case OpPush:
            value := vm.constants[operand]
            vm.push(value)
            
        case OpAdd:
            b := vm.pop()
            a := vm.pop()
            vm.push(a + b)
            
        // ... more opcodes
        }
        
        // 3. ADVANCE
        vm.ip++
    }
    
    return nil
}
```

### Instruction Execution Examples

**Example 1: PUSH (Load Constant)**
```
Bytecode: PUSH 0
Constants: [42]

Before:
  Stack: []
  IP: 0

Execute:
  1. Read constants[0] → 42
  2. Push 42 onto stack
  3. IP++

After:
  Stack: [42]
  IP: 1
```

**Example 2: SEND (Message Send)**
```
Bytecode: SEND 0x0201  (selector index=2, argc=1)
Constants: [3, 4, "+"]

Before:
  Stack: [3, 4]  (receiver=3, arg=4)
  IP: 5

Execute:
  1. Decode: selector index=2 ("+"), argc=1
  2. Pop argument: 4
  3. Pop receiver: 3
  4. Lookup primitive "+" for integer
  5. Execute: 3 + 4 = 7
  6. Push result: 7
  7. IP++

After:
  Stack: [7]
  IP: 6
```

**Example 3: Variable Operations**
```
Bytecode: 
  PUSH 0           ; Load 42
  STORE_LOCAL 0    ; x := 42
  LOAD_LOCAL 0     ; Load x
Constants: [42]

Execution:
  PUSH 0
    Stack: [42]
    Locals: [nil, nil, ...]
  
  STORE_LOCAL 0
    Stack: []          ; Pop consumed value
    Locals: [42, nil, ...]
  
  LOAD_LOCAL 0
    Stack: [42]        ; Pushed local[0]
    Locals: [42, nil, ...]
```

## Stack Operations in Detail

### Push Operation

```go
func (vm *VM) push(value interface{}) error {
    if vm.sp >= StackSize {
        return errors.New("stack overflow")
    }
    vm.stack[vm.sp] = value
    vm.sp++
    return nil
}
```

**Example:**
```
Initial:  Stack: [10, 20, __], sp=2
push(30): Stack: [10, 20, 30], sp=3
```

### Pop Operation

```go
func (vm *VM) pop() (interface{}, error) {
    if vm.sp <= 0 {
        return nil, errors.New("stack underflow")
    }
    vm.sp--
    return vm.stack[vm.sp], nil
}
```

**Example:**
```
Initial: Stack: [10, 20, 30], sp=3
pop():   Stack: [10, 20, __], sp=2, returns 30
```

### Peek Operation

```go
func (vm *VM) peek() (interface{}, error) {
    if vm.sp <= 0 {
        return nil, errors.New("empty stack")
    }
    return vm.stack[vm.sp-1], nil
}
```

## Message Dispatch

Message sending is the heart of Smog's object system:

### Dispatch Process

1. **Receiver and Arguments on Stack**
```
Stack: [receiver, arg1, arg2, ..., argN]
```

2. **Extract Message Information**
```go
selectorIndex := operand >> 8
argCount := operand & 0xFF
selector := vm.constants[selectorIndex]
```

3. **Pop Arguments and Receiver**
```go
args := make([]interface{}, argCount)
for i := argCount - 1; i >= 0; i-- {
    args[i] = vm.pop()
}
receiver := vm.pop()
```

4. **Lookup Method**
```go
method := vm.lookupMethod(receiver, selector)
```

5. **Execute Method**
- For primitives: Execute native code
- For user methods: Execute method bytecode
- For blocks: Execute block bytecode

6. **Push Result**
```go
vm.push(result)
```

### Primitive Methods

Primitives are built-in operations implemented in Go:

```go
func (vm *VM) executePrimitive(receiver interface{}, selector string, args []interface{}) (interface{}, error) {
    switch selector {
    case "+":
        a, ok1 := receiver.(int)
        b, ok2 := args[0].(int)
        if ok1 && ok2 {
            return a + b, nil
        }
        return nil, errors.New("type error")
    
    case "println":
        fmt.Println(receiver)
        return receiver, nil
    
    // ... more primitives
    }
}
```

**Common Primitives:**
- Arithmetic: `+`, `-`, `*`, `/`
- Comparison: `<`, `>`, `<=`, `>=`, `=`, `~=`
- I/O: `print`, `println`
- Collections: `at:`, `at:put:`, `size`, `do:`
- Control flow: `ifTrue:`, `ifFalse:`, `timesRepeat:`

## Block/Closure Execution

Blocks are first-class objects that can capture variables:

### Block Creation

```
Bytecode: PUSH_BLOCK 0

Execute:
  1. Create closure object
  2. Capture current environment (local variables)
  3. Store block bytecode reference
  4. Push closure onto stack
```

### Block Evaluation

When a block receives `value:` message:

```
1. Save current VM state (IP, stack, locals)
2. Set up new stack frame for block
3. Bind block parameters to arguments
4. Execute block bytecode
5. Get return value
6. Restore VM state
7. Push return value
```

**Example:**
```smog
| square |
square := [ :x | x * x ].
square value: 5.
```

**Execution:**
```
Create block:
  PUSH_BLOCK 0
  Stack: [Block{bytecode: [...], captures: []}]
  STORE_LOCAL 0  ; square := block

Invoke block:
  LOAD_LOCAL 0   ; Load square
  PUSH 5         ; Argument
  SEND value:, 1
  
  Inside block execution:
    - Bind x = 5
    - Execute: x * x
    - Result: 25
  
  Stack: [25]
```

### Non-Local Returns

Non-local returns are a fundamental feature of Smalltalk-style blocks. When a return statement (`^`) executes in a block, it doesn't just exit the block - it exits the method that created the block.

**Implementation:**

1. **Block Creation**: Each block captures a reference to its "home context" - the VM executing the method that created it.

2. **OpNonLocalReturn**: When compiled, return statements in blocks generate `OpNonLocalReturn` instead of `OpReturn`.

3. **Exception Propagation**: When `OpNonLocalReturn` executes, it creates a `NonLocalReturn` error containing:
   - The return value
   - The home context (target method's VM)

4. **Unwinding**: The error propagates up through nested block calls until it reaches the target method's `executeMethod`, which catches it and converts it to a normal return.

**Example:**
```smog
findFirst: predicate [
    self do: [ :each |
        (predicate value: each) ifTrue: [
            ^each    " Non-local return from findFirst: "
        ].
    ].
    ^nil
]
```

**Execution Flow:**
```
Method VM: findFirst:
  ↓ creates block1 for do:
  ↓ Block1 VM executes
    ↓ creates block2 for ifTrue:
    ↓ Block2 VM executes
      ↓ OpNonLocalReturn with value=each, homeContext=Method VM
      ↑ NonLocalReturn error propagates
    ↑ Block2 returns error
  ↑ Block1 returns error
↑ Method catches error, returns value=each
```

This mechanism allows natural control flow:
- Early returns from search methods
- Breaking out of loops
- Conditional method termination

**Key Points:**
- Blocks created in methods have homeContext = method's VM
- Blocks created in blocks inherit the parent block's homeContext
- Non-local returns only work within the creating method's execution
- After the method returns, blocks with non-local returns become invalid

## Control Flow Implementation

### Conditional: ifTrue:

```smog
x > 0 ifTrue: [ 'positive' println ].
```

**VM Implementation:**

The `ifTrue:` primitive on boolean objects:

```go
// On True object
func (t *True) ifTrue(block Block) interface{} {
    return block.value()  // Execute block
}

// On False object
func (f *False) ifTrue(block Block) interface{} {
    return nil  // Don't execute block
}
```

**Execution:**
```
Stack: [true, Block{...}]
SEND ifTrue:, 1

Execute:
  1. Pop block argument
  2. Pop receiver (true)
  3. Call true.ifTrue(block)
  4. Block executes
  5. Push result
```

### Loop: timesRepeat:

```smog
5 timesRepeat: [ 'hello' println ].
```

**VM Implementation:**

```go
func (i *Integer) timesRepeat(block Block) interface{} {
    for n := 0; n < i.value; n++ {
        block.value()
    }
    return nil
}
```

**Execution:**
```
Stack: [5, Block{...}]
SEND timesRepeat:, 1

Execute:
  1. Pop block
  2. Pop receiver (5)
  3. Call 5.timesRepeat(block)
  4. Loop 5 times, each calling block.value()
  5. Push result (nil)
```

## Error Handling

The VM handles several types of runtime errors:

### Stack Errors

**Stack Overflow:**
```
Error: stack overflow (SP=1024, max=1024)
At: instruction 42
```

**Stack Underflow:**
```
Error: stack underflow (attempted pop with SP=0)
At: instruction 17
```

### Memory Errors

**Invalid Constant Index:**
```
Error: constant index 100 out of bounds (pool size: 50)
At: instruction 5
```

**Undefined Variable:**
```
Error: undefined global variable 'count'
At: instruction 23
```

### Type Errors

**Invalid Operation:**
```
Error: cannot apply '+' to string and integer
At: instruction 30
```

**Message Not Understood:**
```
Error: object of type String does not understand message 'unknownMethod'
At: instruction 45
```

## Performance Characteristics

### Execution Speed

**Instruction Throughput:**
- Simple ops (PUSH, POP): ~10-50 million/sec
- Arithmetic ops: ~5-20 million/sec
- Message sends: ~1-5 million/sec

**Factors Affecting Speed:**
- Stack operations are fast (array access)
- Message dispatch has overhead (lookup + call)
- Primitive methods are faster than user methods
- Block evaluation requires state save/restore

### Memory Usage

**Per VM Instance:**
- Execution stack: ~8 KB (1024 × 8 bytes)
- Local variables: ~2 KB (256 × 8 bytes)
- Global variables: Varies (hash map)
- Total baseline: ~10-15 KB

**Per Bytecode:**
- Instructions: ~5 bytes each
- Constants: Depends on values
- Typical program: 10-100 KB

## Debugging and Introspection

### Stack Inspection

```go
func (vm *VM) StackSnapshot() []interface{} {
    snapshot := make([]interface{}, vm.sp)
    copy(snapshot, vm.stack[:vm.sp])
    return snapshot
}
```

### Execution Trace

Enable tracing to see each instruction:

```go
vm.EnableTrace(true)
vm.Run(bytecode)
```

**Output:**
```
[0000] PUSH 0              Stack: []
[0001] PUSH 1              Stack: [10]
[0002] SEND 2, 1           Stack: [10, 5]
[0003] STORE_LOCAL 0       Stack: [15]
```

### Breakpoints

Set breakpoints on instruction addresses:

```go
vm.SetBreakpoint(10)  // Break at instruction 10
vm.Run(bytecode)      // Pauses at instruction 10
```

## Optimization Strategies

### 1. Inline Caching

Cache method lookup results:

```go
type InlineCache struct {
    receiverType reflect.Type
    method       *Method
}

var cache map[string]*InlineCache

func (vm *VM) sendCached(selector string, receiver interface{}) {
    recvType := reflect.TypeOf(receiver)
    
    if cache[selector] != nil && cache[selector].receiverType == recvType {
        // Cache hit - use cached method
        return cache[selector].method.execute(receiver)
    }
    
    // Cache miss - lookup and cache
    method := vm.lookupMethod(receiver, selector)
    cache[selector] = &InlineCache{recvType, method}
    return method.execute(receiver)
}
```

### 2. Bytecode Verification

Verify bytecode before execution to enable optimizations:

```go
func (vm *VM) Verify(bytecode *Bytecode) error {
    // Check:
    // - All constant references valid
    // - Stack operations balanced
    // - No unreachable code
    // - All paths return
    
    // After verification, can skip runtime checks
    vm.verified = true
}
```

### 3. Direct Threading

Replace switch statement with computed gotos (C/Assembly):

```c
// Instead of switch
dispatch_table[OpPush]  = &&op_push;
dispatch_table[OpPop]   = &&op_pop;
// ...

goto *dispatch_table[opcode];

op_push:
    // Execute PUSH
    goto *dispatch_table[next_opcode];

op_pop:
    // Execute POP
    goto *dispatch_table[next_opcode];
```

## Best Practices

1. **Verify bytecode**: Check validity before execution
2. **Handle errors gracefully**: Provide clear error messages with context
3. **Limit stack depth**: Prevent stack overflow with checks
4. **Optimize hot paths**: Inline caching for frequent operations
5. **Profile performance**: Identify bottlenecks before optimizing
6. **Test edge cases**: Empty stacks, maximum values, etc.

## VM API

### Creating and Running

```go
import "github.com/kristofer/smog/pkg/vm"

// Create VM
vm := vm.New()

// Run bytecode
err := vm.Run(bytecode)
if err != nil {
    log.Fatal(err)
}

// Get result
result := vm.StackTop()
```

### Configuration

```go
// Set stack size (if configurable)
vm.SetStackSize(2048)

// Enable tracing
vm.SetTrace(true)

// Set global variable
vm.SetGlobal("pi", 3.14159)
```

## Testing the VM

Example test:

```go
func TestVMArithmetic(t *testing.T) {
    // Create bytecode: 3 + 4
    bc := &Bytecode{
        Constants: []interface{}{3, 4, "+"},
        Instructions: []Instruction{
            {Opcode: OpPush, Operand: 0},      // Push 3
            {Opcode: OpPush, Operand: 1},      // Push 4
            {Opcode: OpSend, Operand: 0x0201}, // Send + (selector=2, argc=1)
        },
    }
    
    vm := New()
    err := vm.Run(bc)
    
    if err != nil {
        t.Fatalf("VM error: %v", err)
    }
    
    result := vm.StackTop()
    if result != 7 {
        t.Errorf("Expected 7, got %v", result)
    }
}
```

## Related Documentation

- [Bytecode Documentation](BYTECODE_GENERATION.md) - How bytecode is created
- [Bytecode Format Guide](BYTECODE_FORMAT.md) - .sg file format and loading
- [Compiler Documentation](COMPILER.md) - AST to bytecode compilation
- [VM Specification](../pkg/vm/SPECIFICATION.md) - Formal VM specification
- [Language Specification](spec/LANGUAGE_SPEC.md) - Language semantics

## Loading and Executing Bytecode

The VM can execute bytecode from two sources:

### 1. From Memory (Freshly Compiled)

```go
// Parse source
parser := parser.New(source)
program, _ := parser.Parse()

// Compile to bytecode
compiler := compiler.New()
bytecode, _ := compiler.Compile(program)

// Execute immediately
vm := vm.New()
vm.Run(bytecode)
```

### 2. From .sg Files (Pre-compiled)

```go
// Load bytecode from .sg file
file, _ := os.Open("program.sg")
defer file.Close()
bytecode, _ := bytecode.Decode(file)

// Execute directly (no parsing/compilation)
vm := vm.New()
vm.Run(bytecode)
```

**Performance difference:**
- From source: ~5-500ms (parse + compile + execute)
- From .sg: ~1-10ms (load + execute)
- Speedup: 5-50x for typical programs

### Bytecode Validation

The VM should validate bytecode before execution:

```go
func (vm *VM) Run(bc *Bytecode) error {
    // Validate bytecode structure
    if err := vm.validateBytecode(bc); err != nil {
        return fmt.Errorf("invalid bytecode: %w", err)
    }
    
    // Execute validated bytecode
    return vm.execute(bc)
}
```

**Validation checks:**
- Opcode values are valid
- Constant pool indices are in range
- Stack doesn't overflow/underflow
- Jump targets are valid

## Summary

The Smog Virtual Machine executes bytecode through a stack-based architecture, managing memory, dispatching messages, and providing runtime services. It can execute bytecode from memory (freshly compiled) or from .sg files (pre-compiled), with significant performance benefits from pre-compilation. Understanding the VM is crucial for performance tuning, debugging, and extending the language with new primitives. The clean separation between bytecode and execution allows for future optimizations like JIT compilation while maintaining compatibility.

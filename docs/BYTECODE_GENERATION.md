# Smog Bytecode Generation Guide

## Overview

Bytecode is the low-level intermediate representation that bridges high-level Smog source code and the virtual machine. This guide explains how the compiler generates bytecode from an Abstract Syntax Tree (AST).

## What is Bytecode?

Bytecode is a sequence of instructions designed for a stack-based virtual machine. Unlike native machine code (for x86, ARM, etc.), bytecode is:
- **Platform-independent**: Runs on any system with a Smog VM
- **Higher-level**: Closer to language semantics than CPU instructions
- **Portable**: Same bytecode works everywhere
- **Compact**: Smaller than source, faster to execute than interpretation

## Bytecode Architecture

### Stack-Based Model

Smog uses a stack-based architecture (like JVM, Python bytecode):

**Stack Operations:**
```
Stack: []
PUSH 5       → Stack: [5]
PUSH 3       → Stack: [5, 3]
ADD          → Stack: [8]
```

**Why stack-based?**
- Simple instruction encoding (no register allocation)
- Compact bytecode (fewer operands needed)
- Easy to implement and debug
- Natural fit for expression evaluation

### Instruction Format

Each instruction consists of:
1. **Opcode** (1 byte): What operation to perform
2. **Operand** (4 bytes, optional): Additional data

**Example:**
```
PUSH 10        → [OpPush, 10, 0, 0, 0]
STORE_LOCAL 0  → [OpStoreLocal, 0, 0, 0, 0]
```

## Instruction Set

### Stack Operations

**OpPush** - Load constant onto stack
```
Operand: index into constant pool
Example: PUSH 0  ; loads constants[0]
```

**OpPop** - Discard top stack value
```
Operand: none
Example: POP  ; removes top of stack
```

**OpDup** - Duplicate top stack value
```
Operand: none
Example: DUP  ; [5] → [5, 5]
```

### Variable Operations

**OpLoadLocal** - Load local variable onto stack
```
Operand: local variable slot index
Example: LOAD_LOCAL 0  ; loads local var at slot 0
```

**OpStoreLocal** - Store stack top to local variable
```
Operand: local variable slot index
Example: STORE_LOCAL 0  ; stores to local var at slot 0
```

**OpLoadGlobal** - Load global variable onto stack
```
Operand: index to global variable name in constant pool
Example: LOAD_GLOBAL 5  ; loads global named constants[5]
```

**OpStoreGlobal** - Store stack top to global variable
```
Operand: index to global variable name in constant pool
Example: STORE_GLOBAL 5  ; stores to global named constants[5]
```

### Message Sending

**OpSend** - Send message to receiver
```
Operand: packed value (selector index << 8 | arg count)
Example: SEND 0x0201  ; selector at index 2, 1 argument

Stack before: [receiver, arg1, arg2, ..., argN]
Stack after:  [result]
```

**OpSuperSend** - Send message to superclass
```
Same format as OpSend but dispatches to superclass
```

### Control Flow

**OpJump** - Unconditional jump
```
Operand: offset to jump to (signed)
Example: JUMP 10  ; jump forward 10 instructions
```

**OpJumpIfFalse** - Jump if top of stack is false
```
Operand: offset to jump to
Example: JUMP_IF_FALSE 5  ; jump if false
```

**OpJumpIfTrue** - Jump if top of stack is true
```
Operand: offset to jump to
Example: JUMP_IF_TRUE 5  ; jump if true
```

### Block/Closure Operations

**OpPushBlock** - Create block/closure
```
Operand: index to block bytecode
Example: PUSH_BLOCK 2  ; creates closure from block 2
```

**OpBlockReturn** - Return from block
```
Operand: none
Example: BLOCK_RETURN  ; exits block, leaves value on stack
```

### Method Operations

**OpReturn** - Return from method
```
Operand: none
Example: RETURN  ; exits method, stack top is return value
```

**OpLoadSelf** - Load 'self' onto stack
```
Operand: none
Example: LOAD_SELF  ; pushes current receiver
```

## Code Generation Examples

### Example 1: Simple Arithmetic

**Source:**
```smog
3 + 4 * 5
```

**AST:**
```
BinaryMessage(+)
├── receiver: Integer(3)
└── argument: BinaryMessage(*)
    ├── receiver: Integer(4)
    └── argument: Integer(5)
```

**Bytecode Generation:**
```
Constants: [3, 4, 5, "+", "*"]

Visit BinaryMessage(+):
  Visit Integer(3):
    PUSH 0              ; Load 3
  
  Visit BinaryMessage(*):
    Visit Integer(4):
      PUSH 1            ; Load 4
    
    Visit Integer(5):
      PUSH 2            ; Load 5
    
    SEND 4, 1           ; Send * message
  
  SEND 3, 1             ; Send + message
```

**Final Bytecode:**
```
0:  PUSH 0         ; Load 3
2:  PUSH 1         ; Load 4
4:  PUSH 2         ; Load 5
6:  SEND 4, 1      ; 4 * 5
9:  SEND 3, 1      ; 3 + result
```

**Execution Trace:**
```
Stack: []
PUSH 0         → [3]
PUSH 1         → [3, 4]
PUSH 2         → [3, 4, 5]
SEND *, 1      → [3, 20]     ; 4 * 5 = 20
SEND +, 1      → [23]         ; 3 + 20 = 23
```

### Example 2: Variable Assignment

**Source:**
```smog
| x y |
x := 10.
y := x + 5.
```

**Symbol Table:**
```
x → slot 0
y → slot 1
```

**Bytecode Generation:**
```
Constants: [10, 5, "+"]

Statement 1: x := 10
  Visit Integer(10):
    PUSH 0              ; Load 10
  STORE_LOCAL 0         ; x := 10

Statement 2: y := x + 5
  Visit BinaryMessage(+):
    Visit Identifier(x):
      LOAD_LOCAL 0      ; Load x
    Visit Integer(5):
      PUSH 1            ; Load 5
    SEND 2, 1           ; Send +
  STORE_LOCAL 1         ; y := result
```

**Final Bytecode:**
```
0:  PUSH 0         ; Load 10
2:  STORE_LOCAL 0  ; x := 10
4:  LOAD_LOCAL 0   ; Load x
6:  PUSH 1         ; Load 5
8:  SEND 2, 1      ; x + 5
11: STORE_LOCAL 1  ; y := result
```

### Example 3: Conditional (ifTrue:)

**Source:**
```smog
x > 0 ifTrue: [ 'positive' println ].
```

**Bytecode Generation:**
```
Constants: [0, ">", "positive", "println", "ifTrue:"]

Main bytecode:
  Visit BinaryMessage(>):
    LOAD_LOCAL 0        ; Load x
    PUSH 0              ; Load 0
    SEND 1, 1           ; x > 0
  
  PUSH_BLOCK 0          ; Create block
  SEND 4, 1             ; Send ifTrue:

Block 0 bytecode:
  PUSH 2                ; Load "positive"
  SEND 3, 0             ; Send println
  BLOCK_RETURN
```

**Final Bytecode:**
```
Main:
  0:  LOAD_LOCAL 0
  2:  PUSH 0
  4:  SEND 1, 1       ; >
  7:  PUSH_BLOCK 0
  9:  SEND 4, 1       ; ifTrue:

Block 0:
  0:  PUSH 2
  2:  SEND 3, 0
  5:  BLOCK_RETURN
```

### Example 4: Loop (timesRepeat:)

**Source:**
```smog
5 timesRepeat: [ count println ].
```

**Bytecode Generation:**
```
Constants: [5, "count", "println", "timesRepeat:"]

Main:
  PUSH 0                ; Load 5
  PUSH_BLOCK 0          ; Create block
  SEND 3, 1             ; Send timesRepeat:

Block 0:
  LOAD_GLOBAL 1         ; Load 'count'
  SEND 2, 0             ; Send println
  BLOCK_RETURN
```

### Example 5: Method Definition

**Source:**
```smog
Object subclass: #Counter [
    | count |
    
    increment [
        count := count + 1.
    ]
]
```

**Bytecode Generation:**
```
Class: Counter
Instance Variables: [count]

Method: increment
Constants: [1, "+"]
Bytecode:
  0:  LOAD_IVAR 0       ; Load count (instance var 0)
  2:  PUSH 0            ; Load 1
  4:  SEND 1, 1         ; Send +
  7:  STORE_IVAR 0      ; count := result
  9:  LOAD_IVAR 0       ; Return count (implicit)
  11: RETURN
```

## Constant Pool Management

The constant pool stores literal values referenced by bytecode:

**Types of Constants:**
1. Numbers (integers, floats)
2. Strings
3. Symbols
4. Message selectors
5. Block bytecode

**Pool Building:**
```go
type ConstantPool struct {
    values []interface{}
    index  map[interface{}]int  // Deduplicate constants
}

func (cp *ConstantPool) Add(value interface{}) int {
    // Check if already in pool
    if idx, exists := cp.index[value]; exists {
        return idx
    }
    
    // Add new constant
    idx := len(cp.values)
    cp.values = append(cp.values, value)
    cp.index[value] = idx
    return idx
}
```

**Example:**
```smog
x := 5 + 5.
y := 5.
```

**Constant Pool:**
```
[5, "+"]  ; Only one '5', shared by all uses
```

## Optimization Techniques

### 1. Constant Folding

Evaluate constant expressions at compile time:

**Before:**
```smog
x := 2 + 3 * 4.
```

**Unoptimized Bytecode:**
```
PUSH 2
PUSH 3
PUSH 4
SEND *, 1
SEND +, 1
STORE_LOCAL 0
```

**Optimized Bytecode:**
```
PUSH 14         ; Pre-computed
STORE_LOCAL 0
```

### 2. Peephole Optimization

Replace instruction sequences with more efficient equivalents:

**Pattern: Load then store same location**
```
LOAD_LOCAL 0
STORE_LOCAL 0
```

**Optimized:**
```
; Removed (no-op)
```

**Pattern: Duplicate load**
```
LOAD_LOCAL 0
LOAD_LOCAL 0
```

**Optimized:**
```
LOAD_LOCAL 0
DUP
```

### 3. Dead Code Elimination

Remove unreachable code:

**Source:**
```smog
method [
    ^5.
    'unreachable' println.
]
```

**Unoptimized:**
```
PUSH 5
RETURN
PUSH "unreachable"
SEND println, 0
RETURN
```

**Optimized:**
```
PUSH 5
RETURN
```

## Bytecode Verification

Before execution, bytecode should be verified:

**Checks:**
1. Valid opcodes
2. Constant pool indices in bounds
3. Local variable indices in bounds
4. Stack doesn't underflow
5. All code paths return

**Example Verifier:**
```go
func Verify(bytecode *Bytecode) error {
    stack := 0  // Simulated stack depth
    
    for _, instr := range bytecode.Instructions {
        switch instr.Opcode {
        case OpPush:
            if instr.Operand >= len(bytecode.Constants) {
                return fmt.Errorf("constant index out of bounds")
            }
            stack++
            
        case OpPop:
            if stack == 0 {
                return fmt.Errorf("stack underflow")
            }
            stack--
            
        case OpSend:
            argCount := instr.Operand & 0xFF
            if stack < argCount + 1 {  // args + receiver
                return fmt.Errorf("insufficient stack for message send")
            }
            stack -= argCount  // Consumed args and receiver, pushed result
        }
    }
    
    return nil
}
```

## Debugging Bytecode

### Disassembler

Convert bytecode to readable format:

```go
func Disassemble(bytecode *Bytecode) {
    fmt.Println("Constants:")
    for i, c := range bytecode.Constants {
        fmt.Printf("  %d: %v\n", i, c)
    }
    
    fmt.Println("\nInstructions:")
    offset := 0
    for _, instr := range bytecode.Instructions {
        fmt.Printf("%04d: %-15s", offset, instr.Opcode)
        
        if instr.Operand != 0 {
            fmt.Printf(" %d", instr.Operand)
        }
        
        fmt.Println()
        offset += 5  // opcode (1) + operand (4)
    }
}
```

**Output:**
```
Constants:
  0: 10
  1: 5
  2: "+"

Instructions:
0000: PUSH            0
0005: PUSH            1
0010: SEND            0x0201
```

### Execution Tracer

Track execution step-by-step:

```go
func TraceExecution(vm *VM, bytecode *Bytecode) {
    for !vm.Done() {
        ip := vm.IP()
        instr := bytecode.Instructions[ip]
        
        fmt.Printf("IP=%04d %-15s Stack=%v\n",
            ip, instr.Opcode, vm.StackSnapshot())
        
        vm.Step()
    }
}
```

**Output:**
```
IP=0000 PUSH            Stack=[]
IP=0005 PUSH            Stack=[10]
IP=0010 SEND            Stack=[10, 5]
IP=0015 STORE_LOCAL     Stack=[15]
```

## Best Practices

1. **Minimize constants**: Reuse values in constant pool
2. **Keep bytecode simple**: Resist premature optimization
3. **Verify bytecode**: Catch errors before execution
4. **Document opcodes**: Clear comments on each instruction
5. **Test edge cases**: Empty blocks, deep nesting, etc.
6. **Use disassembler**: Debug bytecode generation issues

## Common Patterns

### Pattern: Expression Result

Leave result on stack:
```
<expr bytecode>
; Stack top has result
```

### Pattern: Statement

Evaluate and discard result:
```
<expr bytecode>
POP             ; Discard result
```

### Pattern: Method Call

```
<receiver bytecode>
<arg1 bytecode>
<arg2 bytecode>
...
SEND selector, N
```

## Performance Characteristics

**Bytecode Size:**
- Typical program: 10-30% of source size
- Heavily commented code: More compression

**Execution Speed:**
- 10-50x slower than native code
- 2-10x faster than AST interpretation

**Memory Usage:**
- Constant pool: Shared across all instances
- Instructions: Read-only, shared
- Stack: Per-thread, temporary

## Related Documentation

- [Compiler Documentation](COMPILER.md) - How bytecode is generated
- [VM Documentation](VM_DEEP_DIVE.md) - How bytecode is executed
- [Bytecode Opcodes](../pkg/bytecode/bytecode.go) - Complete opcode reference
- [Language Specification](spec/LANGUAGE_SPEC.md) - Language semantics

## Summary

Bytecode generation transforms abstract syntax trees into compact, executable instructions for the Smog virtual machine. Understanding bytecode is essential for debugging, optimization, and extending the language. The stack-based architecture provides simplicity and portability, making Smog programs run efficiently across different platforms.

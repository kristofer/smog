# Lexical Scoping Implementation Design

## What is Lexical Scoping?

**Lexical scoping** (also called **static scoping**) means that variable bindings are determined by the structure of the source code, not by the runtime call stack.

### Example: Lexical vs Dynamic Scoping

```smog
| x |
x := 'outer'.

| makeGreeter |
makeGreeter := [
    | message |
    message := 'Hello, '.
    [ :name | message + name + ' from ' + x ]
].

| greet |
greet := makeGreeter value.

| x |  " Different x!
x := 'inner'.

greet value: 'Alice'.
" With lexical scoping: 'Hello, Alice from outer'
" With dynamic scoping: 'Hello, Alice from inner'
```

**Lexical scoping** uses the `x` from where the block was **defined** (outer).
**Dynamic scoping** uses the `x` from where the block was **called** (inner).

Smalltalk and most modern languages use **lexical scoping** because it's more predictable and allows for proper closures.

## Design Overview

### Key Concepts

1. **Environment Chain**: Each scope maintains a link to its parent scope
2. **Captured Variables**: Variables from outer scopes that are used in inner scopes
3. **Local Variables**: Variables declared in the current scope
4. **Closure**: A block + its captured environment

### Architecture

```
┌─────────────────────────────────────────────┐
│ Method/Top-Level Scope                      │
│ locals: [x, y]                              │
│ captured: []                                │
└─────────────────┬───────────────────────────┘
                  │ parent link
        ┌─────────▼──────────────────────────┐
        │ Block Scope                        │
        │ locals: [param1, temp1]            │
        │ captured: [x]  (from parent)       │
        └─────────────┬──────────────────────┘
                      │ parent link
            ┌─────────▼────────────────────┐
            │ Nested Block Scope           │
            │ locals: [param2]             │
            │ captured: [x, temp1]         │
            └──────────────────────────────┘
```

## New Opcodes

We introduce new opcodes to distinguish local vs captured variable access:

```go
// Existing opcodes (for local variables)
OpLoadLocal    // Load from current scope's local variables
OpStoreLocal   // Store to current scope's local variables

// New opcodes (for captured variables)
OpLoadCaptured   // Load from captured environment
OpStoreCaptured  // Store to captured environment

// New opcode for closure creation
OpMakeClosureWithEnv  // Create closure with captured environment
```

### Opcode Encoding

**OpLoadCaptured** / **OpStoreCaptured**:
- **Operand bits 0-15**: Index in the captured variables array
- **Operand bits 16-31**: Depth (0 = parent, 1 = grandparent, etc.)

**OpMakeClosureWithEnv**:
- **Operand bits 0-7**: Number of parameters
- **Operand bits 8-15**: Number of captured variables
- **Operand bits 16-31**: Index of block bytecode in constants

## Compiler Changes

### Environment Tracking

```go
type Compiler struct {
    // Existing fields
    instructions []bytecode.Instruction
    constants    []interface{}

    // New scoping fields
    localVars     []string           // Variables declared in THIS scope
    localCount    int                // Number of local variables

    capturedVars  []CapturedVar      // Variables captured from parent scopes
    parent        *Compiler          // Link to parent scope (for block compilation)

    // Existing fields
    fields        map[string]int
    classVars     map[string]int
    inBlock       bool
}

type CapturedVar struct {
    Name  string
    Index int   // Index in parent's locals or captured vars
    Depth int   // 0 = direct parent, 1 = grandparent, etc.
}
```

### Variable Resolution Algorithm

When compiling a variable reference:

```go
func (c *Compiler) resolveVariable(name string) VariableLocation {
    // 1. Check if it's a local variable in current scope
    if idx, ok := c.localVars[name]; ok {
        return LocalVariable{Index: idx}
    }

    // 2. Check if it's already in captured variables
    for i, captured := range c.capturedVars {
        if captured.Name == name {
            return CapturedVariable{Index: i, Depth: captured.Depth}
        }
    }

    // 3. Search in parent scope
    if c.parent != nil {
        parentLoc := c.parent.resolveVariable(name)
        if parentLoc.IsLocal() {
            // Add to captured variables
            capturedIdx := len(c.capturedVars)
            c.capturedVars = append(c.capturedVars, CapturedVar{
                Name:  name,
                Index: parentLoc.Index,
                Depth: 0,
            })
            return CapturedVariable{Index: capturedIdx, Depth: 0}
        } else if parentLoc.IsCaptured() {
            // Propagate captured variable with increased depth
            capturedIdx := len(c.capturedVars)
            c.capturedVars = append(c.capturedVars, CapturedVar{
                Name:  name,
                Index: parentLoc.Index,
                Depth: parentLoc.Depth + 1,
            })
            return CapturedVariable{Index: capturedIdx, Depth: parentLoc.Depth + 1}
        }
    }

    // 4. Check fields, class vars, globals (existing logic)
    // ...
}
```

### Block Compilation

When compiling a block:

```go
func (c *Compiler) compileBlockLiteral(block *ast.BlockLiteral) error {
    // Create new compiler for block with parent link
    blockCompiler := &Compiler{
        parent:    c,                    // Link to parent
        localVars: make([]string, 0),    // Fresh local scope
        fields:    c.fields,             // Share field map
        classVars: c.classVars,          // Share class var map
        inBlock:   true,
    }

    // Add parameters as local variables
    for _, param := range block.Parameters {
        blockCompiler.localVars = append(blockCompiler.localVars, param)
        blockCompiler.localCount++
    }

    // Compile body (this will populate capturedVars)
    for i, stmt := range block.Body {
        isLast := i == len(block.Body)-1
        if err := blockCompiler.compileStatementWithContext(stmt, isLast); err != nil {
            return err
        }
    }

    blockCompiler.emit(bytecode.OpReturn, 0)

    // Create bytecode with metadata
    blockBytecode := &bytecode.Bytecode{
        Instructions:  blockCompiler.instructions,
        Constants:     blockCompiler.constants,
        CapturedVars:  blockCompiler.capturedVars,  // NEW
        LocalCount:    blockCompiler.localCount,    // NEW
    }

    // Emit closure creation with environment
    blockIdx := c.addConstant(blockBytecode)
    operand := (blockIdx << 16) | (len(blockCompiler.capturedVars) << 8) | len(block.Parameters)
    c.emit(bytecode.OpMakeClosureWithEnv, operand)

    return nil
}
```

## VM Changes

### Closure Object

```go
type Closure struct {
    Code         *bytecode.Bytecode  // Block's compiled code
    Environment  []interface{}       // Captured variable values
    ParentEnv    *Closure            // Link to parent closure (for multi-level capture)
}
```

### OpMakeClosureWithEnv Implementation

```go
case bytecode.OpMakeClosureWithEnv:
    blockIdx := operand >> 16
    capturedCount := (operand >> 8) & 0xFF
    paramCount := operand & 0xFF

    blockBytecode := vm.constants[blockIdx].(*bytecode.Bytecode)

    // Capture variables from current frame
    environment := make([]interface{}, capturedCount)
    for i := 0; i < capturedCount; i++ {
        captured := blockBytecode.CapturedVars[i]
        if captured.Depth == 0 {
            // Capture from current frame's locals
            environment[i] = vm.locals[captured.Index]
        } else {
            // Capture from parent closure's environment
            environment[i] = vm.currentClosure.Environment[captured.Index]
        }
    }

    closure := &Closure{
        Code:        blockBytecode,
        Environment: environment,
        ParentEnv:   vm.currentClosure,
    }

    vm.push(closure)
```

### OpLoadCaptured / OpStoreCaptured Implementation

```go
case bytecode.OpLoadCaptured:
    index := operand & 0xFFFF
    depth := operand >> 16

    env := vm.currentClosure
    for d := 0; d < depth; d++ {
        env = env.ParentEnv
    }

    value := env.Environment[index]
    vm.push(value)

case bytecode.OpStoreCaptured:
    index := operand & 0xFFFF
    depth := operand >> 16

    value := vm.pop()

    env := vm.currentClosure
    for d := 0; d < depth; d++ {
        env = env.ParentEnv
    }

    env.Environment[index] = value
    vm.push(value)  // Assignment returns the value
```

## Example Compilation

### Source Code

```smog
| x |
x := 10.

| makeAdder |
makeAdder := [ :y |
    | temp |
    temp := y * 2.
    [ :z | x + temp + z ]
].

| add5 |
add5 := makeAdder value: 5.

add5 value: 3.  " Returns 10 + 10 + 3 = 23
```

### Compilation Result

**Top-level scope:**
- locals: [x, makeAdder, add5]
- captured: []

**First block `[ :y | ... ]`:**
- locals: [y, temp]
- captured: [x] (from top-level)
- Generates closure with environment [x]

**Nested block `[ :z | ... ]`:**
- locals: [z]
- captured: [x, temp] (x from grandparent, temp from parent)
- Generates closure with environment [x, temp]

### Bytecode

```
Top-level:
  PUSH 10
  STORE_LOCAL 0           ; x := 10

  MAKE_CLOSURE_WITH_ENV   ; makeAdder block
    ; captured: [x]
  STORE_LOCAL 1           ; makeAdder := [...]

  LOAD_LOCAL 1            ; makeAdder
  PUSH 5
  SEND value:, 1
  STORE_LOCAL 2           ; add5 := makeAdder value: 5

  LOAD_LOCAL 2            ; add5
  PUSH 3
  SEND value:, 1          ; add5 value: 3
  RETURN

First block bytecode:
  ; Parameters: y (local 0)
  ; Locals: temp (local 1)
  ; Captured: x (captured 0, depth 0)

  LOAD_LOCAL 0            ; y
  PUSH 2
  SEND *, 1
  STORE_LOCAL 1           ; temp := y * 2

  MAKE_CLOSURE_WITH_ENV   ; nested block
    ; captured: [x, temp]
  RETURN

Nested block bytecode:
  ; Parameters: z (local 0)
  ; Captured: x (captured 0, depth 1), temp (captured 1, depth 0)

  LOAD_CAPTURED 0, 1      ; x (from grandparent)
  LOAD_CAPTURED 1, 0      ; temp (from parent)
  SEND +, 1
  LOAD_LOCAL 0            ; z
  SEND +, 1
  RETURN
```

## Benefits

1. ✅ **Multiple variable declarations anywhere**: Each scope is independent
2. ✅ **Blocks can have local variables**: Block locals don't conflict with parent
3. ✅ **Proper closures**: Captured variables persist correctly
4. ✅ **Lexical scoping semantics**: Variables resolved where defined, not where called
5. ✅ **Nested blocks work correctly**: Multi-level capture is supported

## Implementation Plan

1. Add new opcodes to `pkg/bytecode/bytecode.go`
2. Update `Bytecode` structure to include `CapturedVars` and `LocalCount`
3. Update `Compiler` structure and variable resolution logic
4. Update `compileBlockLiteral` to use environment chains
5. Update VM to implement new opcodes and `Closure` type
6. Update tests to cover lexical scoping scenarios
7. Update documentation and examples

## Testing Strategy

Create tests for:
- Simple variable capture
- Multi-level capture (nested blocks)
- Block-local temporaries
- Multiple variable declarations at different points
- Closures that outlive their creating scope
- Shadowing (local variable with same name as captured)

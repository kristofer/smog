# Smog Architecture Overview

## System Components

### 1. Lexer/Scanner (Future)
**Location**: `pkg/lexer/`

**Responsibilities**:
- Tokenize source code
- Handle whitespace and comments
- Recognize keywords, identifiers, literals, and operators

**Key Types**:
- `Token` - Represents a lexical token
- `Lexer` - Scans source text and produces tokens

### 2. Parser
**Location**: `pkg/parser/`

**Responsibilities**:
- Parse token stream into Abstract Syntax Tree (AST)
- Enforce syntax rules
- Provide error messages for syntax errors

**Key Types**:
- `Parser` - Main parsing logic
- Uses AST nodes from `pkg/ast/`

### 3. Abstract Syntax Tree (AST)
**Location**: `pkg/ast/`

**Responsibilities**:
- Define AST node types
- Represent the structure of smog programs
- Provide visitor pattern support for traversal

**Key Types**:
- `Node` - Base interface for all AST nodes
- `Expression` - Expression nodes
- `Statement` - Statement nodes
- `Class` - Class definition
- `Method` - Method definition
- `MessageSend` - Message sending expression

### 4. Compiler
**Location**: `pkg/compiler/`

**Responsibilities**:
- Traverse AST and generate bytecode
- Manage constant pool
- Handle symbol resolution
- Optimize generated code (future)

**Key Types**:
- `Compiler` - Main compiler logic
- Symbol table management
- Constant pool management

### 5. Bytecode
**Location**: `pkg/bytecode/`

**Responsibilities**:
- Define bytecode instruction set
- Provide bytecode serialization/deserialization
- Support bytecode analysis and optimization

**Key Types**:
- `Opcode` - Instruction opcodes
- `Instruction` - Single bytecode instruction
- `Bytecode` - Complete bytecode object

### 6. Virtual Machine
**Location**: `pkg/vm/`

**Responsibilities**:
- Execute bytecode instructions
- Manage runtime stack
- Handle message dispatch
- Manage object creation and memory

**Key Types**:
- `VM` - Virtual machine state and execution loop
- Stack management
- Frame management for method calls
- Global variable storage

### 7. Runtime Objects (Future)
**Location**: `pkg/runtime/` or `internal/runtime/`

**Responsibilities**:
- Implement core object types
- Provide primitive operations
- Implement standard library classes

**Key Types**:
- `Object` - Base runtime object
- `Class` - Runtime class representation
- `Integer`, `Double`, `String`, `Array`, etc.
- `Block` - Closure implementation

### 8. Command Line Interface
**Location**: `cmd/smog/`

**Responsibilities**:
- Provide command-line interface
- Read source files
- Invoke compilation and execution pipeline
- Handle errors and display output

## Data Flow

```
Source Code (.smog file)
    ↓
[Lexer] → Tokens
    ↓
[Parser] → AST
    ↓
[Compiler] → Bytecode
    ↓
[VM] → Execution → Results
```

## Design Principles

### 1. Separation of Concerns
Each component has a clear, well-defined responsibility. This makes the codebase easier to understand, test, and maintain.

### 2. Modularity
Components are organized into separate packages with minimal dependencies. This allows components to evolve independently.

### 3. Testability
Each component should be testable in isolation with unit tests. Integration tests verify component interaction.

### 4. Smalltalk-Inspired Design
Following SOM's approach:
- Everything is an object
- Late binding through message passing
- Minimal built-in syntax
- Core functionality in standard library

### 5. Performance Considerations
- Stack-based bytecode VM for efficiency
- Future: JIT compilation opportunities
- Future: Inline caching for message dispatch
- Future: Generational garbage collection

## Package Dependencies

```
cmd/smog
  ↓
pkg/parser ← pkg/ast
  ↓
pkg/compiler ← pkg/ast, pkg/bytecode
  ↓
pkg/vm ← pkg/bytecode
```

Internal packages:
- `internal/util` - Shared utilities
- `internal/types` - Internal type definitions

## Error Handling Strategy

1. **Lexical/Syntax Errors**: Reported during parsing with line/column information
2. **Semantic Errors**: Detected during compilation (e.g., undefined variables)
3. **Runtime Errors**: Handled by VM (e.g., message not understood, stack overflow)

Error messages should be clear and actionable, pointing to the source location when possible.

## Testing Strategy

### Unit Tests
- Test individual components in isolation
- Mock dependencies where needed
- High code coverage for core logic

### Integration Tests
- Test full compilation pipeline
- Verify bytecode generation and execution
- Test standard library implementations

### End-to-End Tests
- Complete smog programs in `examples/` directory
- Verify expected output
- Performance benchmarks

## Build and Development

### Build Commands
```bash
# Build the smog interpreter
go build -o bin/smog ./cmd/smog

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run a smog file
./bin/smog examples/hello.smog
```

### Code Organization
- Public APIs in `pkg/`
- Internal implementation details in `internal/`
- Executable in `cmd/`
- Documentation in `docs/`
- Examples in `examples/`

## Future Architecture Enhancements

### 1. REPL (Read-Eval-Print Loop)
Interactive shell for executing smog code.

### 2. Debugger
Step-through debugging with breakpoints and variable inspection.

### 3. Module System
Support for organizing code into modules and managing dependencies.

### 4. JIT Compiler
Just-in-time compilation of hot code paths for improved performance.

### 5. Profiler
Performance profiling tools to identify bottlenecks.

### 6. Package Manager
Tool for discovering and installing third-party smog libraries.

## References

- [Smalltalk-80: The Language and its Implementation](https://rmod-files.lille.inria.fr/FreeBooks/BlueBook/Bluebook.pdf)
- [SOM (Simple Object Machine)](http://som-st.github.io/)
- [Crafting Interpreters](https://craftinginterpreters.com/)

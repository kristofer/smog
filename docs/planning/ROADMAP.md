# Smog Development Roadmap

## Version 0.1.0 - Foundation ✅

### Goals
Establish the basic project structure and documentation foundation.

### Completed
- [x] Initialize Go module
- [x] Create project directory structure
- [x] Basic CLI scaffolding
- [x] Initial package stubs (parser, compiler, bytecode, vm, ast)
- [x] Language specification document
- [x] Architecture documentation
- [x] Development roadmap
- [x] Implement lexer for tokenization
- [x] Implement parser for basic expressions
- [x] Create AST node types for core language features
- [x] Add example programs

## Version 0.2.0 - Core Language Features ✅

### Goals
Implement the minimum viable language interpreter with basic features.

### Parser
- [x] Literal parsing (numbers, strings, booleans, nil)
- [x] Variable declarations and assignments
- [x] Message send expressions (unary, binary, keyword)
- [ ] Block/closure syntax (completed in v0.3.0)
- [ ] Class definitions
- [ ] Method definitions

### Compiler
- [x] Compile literals to bytecode
- [x] Compile variable access/assignment
- [x] Compile message sends
- [ ] Compile blocks (completed in v0.3.0)
- [ ] Compile class definitions
- [x] Symbol table management
- [x] Constant pool generation

### Virtual Machine
- [x] Stack operations (push, pop, dup)
- [x] Local variable access
- [x] Message dispatch mechanism
- [ ] Block evaluation (completed in v0.3.0)
- [ ] Method invocation
- [x] Return handling

### Runtime
- [x] Object representation
- [ ] Class representation
- [x] Basic type system (Integer, String, Boolean, Nil)
- [x] Message lookup and dispatch (primitives)

## Version 0.3.0 - Standard Library ✅

### Goals
Implement core classes and methods for practical programming.

### Core Classes
- [x] Object class with basic methods (partial)
- [x] Integer class with arithmetic operations
- [ ] Double class for floating-point numbers
- [x] String class with manipulation methods (partial)
- [x] Array class for collections
- [x] Boolean, True, False classes
- [x] Block class for closures

### Control Flow
- [x] Conditional messages (ifTrue:, ifFalse:)
- [ ] ifTrue:ifFalse: (planned)
- [ ] Loop messages (whileTrue:, whileFalse:) (planned)
- [x] timesRepeat: for integers
- [x] Collection iteration (do:)
- [ ] collect:, select:, reject: (planned)

### Additional Features
- [x] Block/closure syntax and compilation
- [x] Array literals (#(...))
- [x] Return statements (^)
- [x] Block parameter support
- [x] Extensive teaching-quality documentation
- [x] Comprehensive test suite (48+ tests)

## Version 0.4.0 - Enhanced Features ✅

### Goals
Add features for more complex programs and better developer experience.

### Language Features
- [ ] Instance variable initialization (deferred - requires class parsing infrastructure)
- [ ] Class variables and methods (deferred - requires class parsing infrastructure)
- [x] Super message sends
- [x] Cascading messages
- [x] Array and dictionary literals

### Development Tools
- [ ] Better error messages with source locations
- [ ] Stack traces for runtime errors
- [ ] Basic debugger support
- [x] REPL (Read-Eval-Print Loop)

### Testing
- [x] Comprehensive test suite
- [x] Integration tests for example programs
- [x] Performance benchmarks

## Version 0.5.0 - Optimization

### Goals
Improve performance and add optimization capabilities.

### Compiler Optimizations
- [ ] Constant folding
- [ ] Dead code elimination
- [ ] Tail call optimization
- [ ] Inline caching hints

### VM Optimizations
- [ ] Inline caching for message dispatch
- [ ] Optimized primitive operations
- [ ] Memory pooling for objects
- [ ] Garbage collection improvements

### Performance
- [ ] Benchmark suite
- [ ] Performance profiling tools
- [ ] Comparison with other Smalltalk implementations

## Version 0.6.0 - Advanced Features

### Goals
Add advanced language features and ecosystem tools.

### Language Features
- [ ] Module system
- [ ] Exception handling (try-catch-finally)
- [ ] First-class continuations
- [ ] Meta-programming capabilities

### Ecosystem
- [ ] Package manager
- [ ] Build system
- [ ] Documentation generator
- [ ] Code formatter

## Version 1.0.0 - Production Ready

### Goals
Stabilize the language and provide production-quality implementation.

### Stability
- [ ] Complete language specification
- [ ] Comprehensive test coverage (>90%)
- [ ] Memory safety verification
- [ ] Performance benchmarks
- [ ] Security audit

### Documentation
- [ ] Complete API documentation
- [ ] Tutorial series
- [ ] Best practices guide
- [ ] Migration guides
- [ ] Example applications

### Tooling
- [ ] IDE support (LSP server)
- [ ] Syntax highlighting for popular editors
- [ ] Debugging tools
- [ ] Profiling tools

## Future Considerations (Post 1.0)

### Experimental Features
- [ ] JIT compilation
- [ ] Concurrent/parallel execution
- [ ] Foreign Function Interface (FFI)
- [ ] Native compilation
- [ ] WebAssembly target

### Ecosystem Growth
- [ ] Standard library expansion
- [ ] Community package repository
- [ ] Web framework
- [ ] Database connectors
- [ ] Network libraries

## Development Principles

### Incremental Development
Each version should be a usable increment. We avoid big-bang releases and prefer small, tested improvements.

### Test-Driven Development
Write tests before implementation when possible. All new features should have corresponding tests.

### Documentation-Driven Development
Document the design before implementation. Keep documentation up-to-date as code evolves.

### Community-Driven Priorities
Listen to early adopters and adjust priorities based on real-world usage.

## Release Schedule

- **Patch releases** (0.x.y): Bug fixes, minor improvements - As needed
- **Minor releases** (0.x.0): New features - Monthly to quarterly
- **Major releases** (x.0.0): Significant changes - When ready

## Contributing

We welcome contributions at all levels:
- Bug reports and feature requests
- Documentation improvements
- Example programs
- Core implementation
- Tool development

See CONTRIBUTING.md (to be created) for guidelines.

## Milestones

### Milestone 1: Hello World
A working interpreter that can execute: `'Hello, World!' println.`

### Milestone 2: Arithmetic
Support for basic arithmetic: `3 + 4 * 5.`

### Milestone 3: Objects
Define and use simple classes with methods.

### Milestone 4: Closures
Support for blocks and higher-order functions.

### Milestone 5: Self-Hosting
Smog compiler written in Smog (aspirational).

## Current Status

**Version**: 0.4.0  
**Status**: Enhanced features complete  
**Next Release**: 0.5.0 (optimization)

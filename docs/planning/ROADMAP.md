# Smog Development Roadmap

## Version 0.1.0 - Foundation (Current)

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

### Next Steps
- [ ] Implement lexer for tokenization
- [ ] Implement parser for basic expressions
- [ ] Create AST node types for core language features
- [ ] Add example programs

## Version 0.2.0 - Core Language Features

### Goals
Implement the minimum viable language interpreter with basic features.

### Parser
- [ ] Literal parsing (numbers, strings, booleans, nil)
- [ ] Variable declarations and assignments
- [ ] Message send expressions (unary, binary, keyword)
- [ ] Block/closure syntax
- [ ] Class definitions
- [ ] Method definitions

### Compiler
- [ ] Compile literals to bytecode
- [ ] Compile variable access/assignment
- [ ] Compile message sends
- [ ] Compile blocks
- [ ] Compile class definitions
- [ ] Symbol table management
- [ ] Constant pool generation

### Virtual Machine
- [ ] Stack operations (push, pop, dup)
- [ ] Local variable access
- [ ] Message dispatch mechanism
- [ ] Block evaluation
- [ ] Method invocation
- [ ] Return handling

### Runtime
- [ ] Object representation
- [ ] Class representation
- [ ] Basic type system (Integer, String, Boolean, Nil)
- [ ] Message lookup and dispatch

## Version 0.3.0 - Standard Library

### Goals
Implement core classes and methods for practical programming.

### Core Classes
- [ ] Object class with basic methods
- [ ] Integer class with arithmetic operations
- [ ] Double class for floating-point numbers
- [ ] String class with manipulation methods
- [ ] Array class for collections
- [ ] Boolean, True, False classes
- [ ] Block class for closures

### Control Flow
- [ ] Conditional messages (ifTrue:, ifFalse:, ifTrue:ifFalse:)
- [ ] Loop messages (whileTrue:, whileFalse:, timesRepeat:)
- [ ] Collection iteration (do:, collect:, select:, reject:)

## Version 0.4.0 - Enhanced Features

### Goals
Add features for more complex programs and better developer experience.

### Language Features
- [ ] Instance variable initialization
- [ ] Class variables and methods
- [ ] Super message sends
- [ ] Cascading messages
- [ ] Array and dictionary literals

### Development Tools
- [ ] Better error messages with source locations
- [ ] Stack traces for runtime errors
- [ ] Basic debugger support
- [ ] REPL (Read-Eval-Print Loop)

### Testing
- [ ] Comprehensive test suite
- [ ] Integration tests for example programs
- [ ] Performance benchmarks

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

**Version**: 0.1.0-dev  
**Status**: Foundation phase  
**Next Release**: 0.1.0 (documentation and structure complete)

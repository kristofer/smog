# Documentation Index

This document provides a comprehensive overview of all Smog documentation.

## Quick Start Paths

### 🌱 Path 1: Complete Beginner
If you're new to programming or Smog:
1. [Learning Guide](LEARNING_GUIDE.md) - Mental models and concepts
2. [User's Guide](USERS_GUIDE.md) - Practical examples
3. [Example Programs](../examples/) - Working code to study
4. [Language Specification](spec/LANGUAGE_SPEC.md) - Reference guide

### 💻 Path 2: Experienced Programmer
If you know other languages and want to learn Smog:
1. [User's Guide](USERS_GUIDE.md) - Jump right into examples
2. [Language Specification](spec/LANGUAGE_SPEC.md) - Complete syntax
3. [Example Programs](../examples/) - Real-world patterns
4. [Learning Guide](LEARNING_GUIDE.md) - Smog's unique philosophy

### 🔧 Path 3: Language Implementer
If you want to understand or contribute to Smog's internals:
1. [Architecture](design/ARCHITECTURE.md) - System overview
2. Compilation Pipeline:
   - [Lexer](LEXER.md) → [Parser](PARSER.md) → [Compiler](COMPILER.md) → [VM](VM_DEEP_DIVE.md)
3. [Bytecode Generation](BYTECODE_GENERATION.md) - Bytecode details
4. [Design Decisions](design/DECISIONS.md) - Why things work this way

## Documentation by Category

### 📚 Learning Resources

| Document | Purpose | Audience |
|----------|---------|----------|
| [Learning Guide](LEARNING_GUIDE.md) | Mental models and conceptual understanding | Beginners, visual learners |
| [User's Guide](USERS_GUIDE.md) | Practical programming guide with examples | All users |
| [Language Specification](spec/LANGUAGE_SPEC.md) | Complete language reference | All levels |
| [Example Programs](../examples/) | Working code samples | All levels |

### 🔬 Technical Documentation

| Document | Purpose | Audience |
|----------|---------|----------|
| [Lexer Documentation](LEXER.md) | Tokenization process | Implementers, curious users |
| [Parser Documentation](PARSER.md) | AST construction | Implementers, curious users |
| [Compiler Documentation](COMPILER.md) | Bytecode generation | Implementers, curious users |
| [Bytecode Generation Guide](BYTECODE_GENERATION.md) | Bytecode format and optimization | Implementers |
| [Bytecode Format Guide](BYTECODE_FORMAT.md) | .sg file format and CLI usage | Users, Implementers |
| [VM Deep Dive](VM_DEEP_DIVE.md) | Virtual machine internals | Implementers |
| [VM Specification](../pkg/vm/SPECIFICATION.md) | Formal VM spec | Implementers |

### 🏗️ Design & Planning

| Document | Purpose | Audience |
|----------|---------|----------|
| [Architecture](design/ARCHITECTURE.md) | System design overview | Implementers, contributors |
| [Design Decisions](design/DECISIONS.md) | Rationale for key choices | Implementers, contributors |
| [Roadmap](planning/ROADMAP.md) | Development plans | Contributors, curious users |
| [Microcontroller Porting Plan](planning/MICROCONTROLLER_PORTING_PLAN.md) | C VM and TinyGo analysis for embedded systems | Implementers, embedded developers |

### 🚀 Getting Started

| Document | Purpose | Audience |
|----------|---------|----------|
| [README](../README.md) | Project overview and quick start | Everyone |
| [Getting Started](GETTING_STARTED.md) | Installation and first steps | New users |

## Document Summaries

### For Users

**[Learning Guide](LEARNING_GUIDE.md)** (15k+ lines)
- Mental models: Everything is an object, message passing, blocks, classes
- Visual diagrams of compilation pipeline
- Common beginner mistakes and solutions
- Progressive learning path from basics to advanced
- Debugging strategies
- Conceptual analogies for programming concepts

**[User's Guide](USERS_GUIDE.md)** (16k+ lines)
- Comprehensive practical examples
- Data structures: arrays, stacks, queues
- Algorithms: sorting (bubble, quick), searching (binary, linear)
- Object-oriented patterns: encapsulation, inheritance, polymorphism, composition
- Common algorithms: factorial, fibonacci, GCD, primes
- Best practices and coding style

**[Language Specification](spec/LANGUAGE_SPEC.md)**
- Complete syntax reference
- Literals: numbers, strings, booleans, arrays
- Message types: unary, binary, keyword
- Control structures
- Classes and methods
- Blocks and closures
- Standard library

### For Implementers

**[Lexer Documentation](LEXER.md)** (13k+ lines)
- Token types and classification
- Character scanning algorithms
- Number, string, identifier, comment handling
- Error handling and recovery
- Performance considerations
- Testing strategies

**[Parser Documentation](PARSER.md)** (11k+ lines)
- Abstract Syntax Tree (AST) structure
- Grammar rules and precedence
- Recursive descent parsing
- Message precedence handling
- Error messages and recovery
- Parser API and testing

**[Compiler Documentation](COMPILER.md)** (9k+ lines)
- Constant pool management
- Variable storage (local vs global)
- Stack-based code generation
- Message sending compilation
- Blocks and closures
- Optimization opportunities

**[Bytecode Generation Guide](BYTECODE_GENERATION.md)** (12k+ lines)
- Bytecode instruction set
- Stack-based architecture
- Complete code generation examples
- Optimization techniques
- Bytecode verification
- Debugging tools (disassembler, tracer)

**[Bytecode Format Guide](BYTECODE_FORMAT.md)** (8k+ lines)
- Binary .sg file format specification
- Compilation and execution workflows
- CLI commands (compile, disassemble)
- Performance benchmarks
- Multi-file program patterns
- Version compatibility

**[Virtual Machine Deep Dive](VM_DEEP_DIVE.md)** (14k+ lines)
- VM architecture and components
- Fetch-decode-execute cycle
- Stack operations in detail
- Message dispatch mechanism
- Block/closure execution
- Control flow implementation
- Error handling
- Performance characteristics
- Optimization strategies

### For Contributors

**[Architecture](design/ARCHITECTURE.md)**
- System component overview
- Data flow through pipeline
- Design principles
- Package dependencies
- Testing strategy
- Future enhancements

**[Design Decisions](design/DECISIONS.md)**
- Stack-based vs register-based VM
- Message passing semantics
- Bytecode format choices
- Memory management approach

**[Roadmap](planning/ROADMAP.md)**
- Completed features
- Planned features
- Development timeline
- Version milestones

## Coverage Matrix

What each document covers:

| Topic | Learning | User's | Lang Spec | Lexer | Parser | Compiler | Bytecode | VM |
|-------|----------|--------|-----------|-------|--------|----------|----------|-----|
| Basic Syntax | ✅ | ✅ | ✅ | | | | | |
| Message Sending | ✅ | ✅ | ✅ | | ✅ | ✅ | ✅ | ✅ |
| Blocks/Closures | ✅ | ✅ | ✅ | | ✅ | ✅ | ✅ | ✅ |
| Classes | ✅ | ✅ | ✅ | | ✅ | ✅ | ✅ | ✅ |
| Control Flow | ✅ | ✅ | ✅ | | | | | ✅ |
| Data Structures | | ✅ | | | | | | |
| Algorithms | | ✅ | | | | | | |
| OOP Patterns | ✅ | ✅ | | | | | | |
| Tokenization | | | | ✅ | | | | |
| Parsing | | | | | ✅ | | | |
| AST | | | | | ✅ | ✅ | | |
| Code Generation | | | | | | ✅ | ✅ | |
| Bytecode Format | | | | | | | ✅ | ✅ |
| .sg Files | | ✅ | ✅ | | | ✅ | ✅ | ✅ |
| VM Execution | | | | | | | | ✅ |
| Optimization | | | | | | ✅ | ✅ | ✅ |
| Debugging | ✅ | | | ✅ | ✅ | | ✅ | ✅ |

## Document Statistics

| Document | Lines | Audience | Reading Time |
|----------|-------|----------|--------------|
| Learning Guide | ~600 | Beginner | 45-60 min |
| User's Guide | ~700 | User | 60-90 min |
| Language Spec | ~325 | User/Impl | 30-40 min |
| Lexer | ~550 | Implementer | 45-60 min |
| Parser | ~500 | Implementer | 45-60 min |
| Compiler | ~430 | Implementer | 35-45 min |
| Bytecode Guide | ~530 | Implementer | 45-60 min |
| Bytecode Format | ~330 | User/Impl | 30-40 min |
| VM Deep Dive | ~600 | Implementer | 60-75 min |

**Total Documentation**: ~4,500+ lines covering all aspects of the language

## How to Use This Documentation

### Scenario 1: "I want to write a simple program"
→ Start with [User's Guide](USERS_GUIDE.md), look at [Example Programs](../examples/)

### Scenario 2: "I don't understand how X works"
→ Check [Learning Guide](LEARNING_GUIDE.md) for mental models

### Scenario 3: "What's the exact syntax for Y?"
→ Refer to [Language Specification](spec/LANGUAGE_SPEC.md)

### Scenario 4: "How do I compile programs for distribution?"
→ Read [Bytecode Format Guide](BYTECODE_FORMAT.md) for .sg files

### Scenario 5: "How does the compiler generate bytecode?"
→ Read [Compiler Documentation](COMPILER.md) and [Bytecode Generation Guide](BYTECODE_GENERATION.md)

### Scenario 6: "Why is my code slow?"
→ Check [VM Deep Dive](VM_DEEP_DIVE.md) performance section, consider compiling to .sg

### Scenario 7: "I want to contribute a feature"
→ Read [Architecture](design/ARCHITECTURE.md), [Design Decisions](design/DECISIONS.md), and relevant technical docs

### Scenario 8: "How does message sending really work?"
→ Path: [Learning Guide](LEARNING_GUIDE.md) → [Parser](PARSER.md) → [Compiler](COMPILER.md) → [VM Deep Dive](VM_DEEP_DIVE.md)

## Feedback and Improvements

This documentation is a living resource. If you find:
- Unclear explanations
- Missing topics
- Errors or typos
- Areas needing more examples

Please open an issue or submit a pull request!

## Last Updated

This documentation index reflects the comprehensive documentation update completed in January 2026.

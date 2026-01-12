# Smog Examples

This directory contains example programs demonstrating various features of the smog language.

## Basic Examples

### hello.smog
The classic "Hello, World!" program.

```smog
'Hello, World!' println.
```

**Demonstrates**: Basic string literal and message sending.

### factorial.smog
Calculates factorial using recursion.

```smog
Object subclass: #Math [
    factorial: n [
        n <= 1 ifTrue: [ ^1 ].
        ^n * (self factorial: (n - 1))
    ]
]
```

**Demonstrates**: 
- Class definition
- Method definition
- Recursion
- Conditional execution (ifTrue:)
- Return statements (^)

## Object-Oriented Examples

### counter.smog
A simple counter class with increment/decrement operations.

**Demonstrates**:
- Instance variables
- Multiple methods
- Object state management
- Object instantiation

### point.smog
Classic 2D point class with arithmetic operations.

**Demonstrates**:
- Constructor methods
- Accessor methods
- Operator overloading (binary messages)
- Object composition
- Custom string representation

## Collection Examples

### arrays.smog
Working with arrays and collection operations.

**Demonstrates**:
- Array literals (#(...))
- Iteration (do:)
- Transformation (collect:)
- Filtering (select:)
- Accumulation patterns

## Functional Programming Examples

### blocks.smog
Block (closure) examples showing functional programming concepts.

**Demonstrates**:
- Block literals
- Block parameters
- Block evaluation
- Higher-order functions
- Closures capturing variables
- Lexical scoping

## Running Examples

Once the smog interpreter is built, you can run any example:

```bash
# Build the interpreter
go build -o bin/smog ./cmd/smog

# Run an example
./bin/smog examples/hello.smog
./bin/smog examples/factorial.smog
./bin/smog examples/counter.smog
./bin/smog examples/point.smog
./bin/smog examples/arrays.smog
./bin/smog examples/blocks.smog
```

## Learning Path

If you're new to smog, we recommend exploring the examples in this order:

1. **hello.smog** - Start with the basics
2. **factorial.smog** - Learn about classes and methods
3. **counter.smog** - Understand instance variables and state
4. **point.smog** - See more complex class interactions
5. **blocks.smog** - Explore functional programming features
6. **arrays.smog** - Work with collections

## Contributing Examples

We welcome new examples! Good examples should:
- Demonstrate one or more language features clearly
- Include comments explaining what's happening
- Be concise but complete
- Follow smog syntax conventions
- Include expected output in comments

## Future Examples

Planned examples to add:
- Control flow (loops, conditionals)
- String manipulation
- File I/O
- Error handling
- Meta-programming
- Design patterns (Singleton, Observer, etc.)
- Real-world applications

## Note on Implementation Status

**Important**: As of version 0.1.0, the smog interpreter is not yet fully implemented. These examples are written in the planned smog syntax and will run once the parser, compiler, and VM are complete. See the [roadmap](../docs/planning/ROADMAP.md) for implementation status.

These examples serve as:
1. Specification by example
2. Test cases for implementation
3. Documentation for language users
4. Validation of language design

## References

For more information about smog syntax and semantics:
- [Language Specification](../docs/spec/LANGUAGE_SPEC.md)
- [Architecture Overview](../docs/design/ARCHITECTURE.md)
- [Main README](../README.md)

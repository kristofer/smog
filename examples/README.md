# Smog Examples

This directory contains example programs demonstrating various features of the smog language.

## Runnable Examples

These examples can be executed with the current version of Smog:

### hello.smog
The classic "Hello, World!" program.

```smog
'Hello, World!' println.
```

**Demonstrates**: Basic string literal and message sending.

### arrays.smog
Working with arrays and basic collection operations.

**Demonstrates**:
- Array literals (#(...))
- Iteration (do:)
- Array access (at:)
- Array size
- Variable declarations

### blocks.smog
Block (closure) examples showing functional programming concepts.

**Demonstrates**:
- Block literals
- Block parameters
- Block evaluation
- Higher-order functions

## Syntax-Only Examples

Examples that demonstrate valid Smog syntax but require features not yet implemented (classes, object instantiation) are in the [syntax-only/](syntax-only/) directory.
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

Once the smog interpreter is built, you can run the executable examples:

```bash
# Build the interpreter
go build -o bin/smog ./cmd/smog

# Run examples directly from source
./bin/smog examples/hello.smog
./bin/smog examples/arrays.smog
./bin/smog examples/blocks.smog

# Compile to bytecode for faster execution
./bin/smog compile examples/hello.smog examples/hello.sg

# Run compiled bytecode
./bin/smog examples/hello.sg

# Inspect bytecode
./bin/smog disassemble examples/hello.sg

# Run all examples with the test script
./run_examples.sh
```

### Bytecode Files (.sg)

Some examples include pre-compiled .sg bytecode files. These provide:
- **Faster startup** - No parsing or compilation needed
- **Learning resource** - Inspect bytecode with `smog disassemble`
- **Distribution format** - Share programs without source code

To compile any example to bytecode:
```bash
./bin/smog compile examples/counter.smog examples/counter.sg
```

**Note**: Examples in the `syntax-only/` directory demonstrate valid syntax but cannot execute because classes are not yet implemented.

## Learning Path

If you're new to smog, we recommend exploring the examples in this order:

1. **hello.smog** - Start with the basics
2. **arrays.smog** - Work with collections  
3. **blocks.smog** - Explore functional programming features
4. **v0.2.0/** - Versioned examples for v0.2.0 features
5. **v0.3.0/** - Versioned examples for v0.3.0 features
6. **v0.4.0/** - Versioned examples for v0.4.0 features
7. **syntax-only/** - Advanced syntax examples (class-based, not yet executable)

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

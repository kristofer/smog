# Smog

A very simple object-oriented language inspired by Smalltalk and SOM (Simple Object Machine).

## Overview

Smog is a minimalist object-oriented programming language that follows the philosophy that "everything is an object." It features:

- **Pure Object-Oriented Design**: Everything, including classes, numbers, and booleans, is an object
- **Message Passing**: All computation happens through sending messages to objects
- **Smalltalk-Inspired Syntax**: Clean, minimal syntax based on Smalltalk
- **Bytecode Compilation**: Source code compiles to bytecode for efficient execution
- **Stack-Based VM**: Simple virtual machine for bytecode execution
- **Bytecode Object Files**: Compile to .sg files for faster loading and distribution

## Quick Start

### Building

```bash
# Build the smog interpreter
go build -o bin/smog ./cmd/smog

# Or use go run
go run ./cmd/smog [file.smog]
```

### Interactive REPL

```bash
# Start the REPL (Read-Eval-Print Loop)
./bin/smog

# Or explicitly
./bin/smog repl
```

### Running Examples

```bash
# Run hello world
./bin/smog examples/hello.smog

# Compile to bytecode for faster loading
./bin/smog compile examples/hello.smog examples/hello.sg

# Run the compiled bytecode
./bin/smog examples/hello.sg

# Disassemble bytecode to inspect it
./bin/smog disassemble examples/hello.sg

# Run other examples
./bin/smog examples/counter.smog
```

### Hello World

```smog
'Hello, World!' println.
```

## Language Features

### Everything is an Object

```smog
" Numbers are objects "
3 + 4.
5 * 2.

" Booleans are objects "
true ifTrue: [ 'yes' println ].

" Even classes are objects "
Object class.
```

### Message Passing

```smog
" Unary messages "
array size.

" Binary messages "
3 + 4.
x < y.

" Keyword messages "
array at: 1 put: 'value'.
point x: 10 y: 20.
```

### Classes and Objects

```smog
Object subclass: #Counter [
    | count |
    
    initialize [
        count := 0.
    ]
    
    increment [
        count := count + 1.
    ]
    
    value [
        ^count
    ]
]

| counter |
counter := Counter new.
counter initialize.
counter increment.
counter value println.
```

### Blocks (Closures)

```smog
" Simple block "
[ 'Hello' println ] value.

" Block with parameters "
[ :x | x * 2 ] value: 5.

" Blocks as control structures "
x > 0 ifTrue: [ 'positive' println ].

5 timesRepeat: [ 'hello' println ].

#(1 2 3) do: [ :each | each println ].
```

## Project Structure

```
smog/
├── cmd/smog/           # Main executable
├── pkg/
│   ├── ast/           # Abstract Syntax Tree definitions
│   ├── parser/        # Source code parser
│   ├── compiler/      # AST to bytecode compiler
│   ├── bytecode/      # Bytecode format and opcodes
│   └── vm/            # Virtual machine
├── stdlib/            # Standard library
│   ├── collections/   # Data structures (Set, OrderedCollection, Bag)
│   ├── core/          # Core utilities (Math, Stream)
│   ├── io/            # I/O operations (HTTP)
│   ├── crypto/        # Cryptography (AES, Hash)
│   └── compression/   # Compression (ZIP, GZIP)
├── internal/          # Internal implementation details
├── docs/
│   ├── spec/          # Language specification
│   ├── design/        # Design documents
│   └── planning/      # Development planning
└── examples/          # Example programs
    └── stdlib/        # Standard library examples
```

## Documentation

### For Users and Learners

Start here if you're new to Smog or want to learn how to use the language:

- **[Learning Guide](docs/LEARNING_GUIDE.md)** - ⭐ **START HERE** - Beginner's mental model and learning path
- **[User's Guide](docs/USERS_GUIDE.md)** - Practical guide with examples for common programming tasks
- **[Standard Library](stdlib/README.md)** - Common data structures and utilities (Set, OrderedCollection, Math, etc.)
- **[REPL Guide](docs/REPL.md)** - Interactive Read-Eval-Print Loop for experimentation
- **[Debugger Guide](docs/DEBUGGER.md)** - Interactive debugger for step-by-step execution
- **[Bytecode Format Guide](docs/BYTECODE_FORMAT.md)** - Working with .sg compiled bytecode files
- **[Language Specification](docs/spec/LANGUAGE_SPEC.md)** - Complete language reference and syntax guide
- **[Example Programs](examples/)** - Working code examples you can run and study

### For Language Implementers

Deep-dive technical documentation for understanding how Smog works internally:

- **[Lexer Documentation](docs/LEXER.md)** - How source code becomes tokens
- **[Parser Documentation](docs/PARSER.md)** - How tokens become an Abstract Syntax Tree (AST)
- **[Compiler Documentation](docs/COMPILER.md)** - How AST transforms into bytecode
- **[Bytecode Generation Guide](docs/BYTECODE_GENERATION.md)** - Detailed bytecode format and generation
- **[Virtual Machine Deep Dive](docs/VM_DEEP_DIVE.md)** - How the VM executes bytecode
- **[VM Specification](pkg/vm/SPECIFICATION.md)** - Formal VM specification

### Design and Planning

- **[Architecture](docs/design/ARCHITECTURE.md)** - System architecture overview
- **[Design Decisions](docs/design/DECISIONS.md)** - Key design decisions and rationale
- **[Roadmap](docs/planning/ROADMAP.md)** - Development roadmap and milestones
- **[Microcontroller Porting Plan](docs/planning/MICROCONTROLLER_PORTING_PLAN.md)** - Analysis of C VM and TinyGo approaches for embedded systems

## Development Status

**Current Version**: 0.5.0 (Advanced Classes)

Smog has completed advanced class features implementation. Current features:
- ✅ Complete lexer and parser
- ✅ AST-based intermediate representation
- ✅ Bytecode compiler
- ✅ Stack-based virtual machine
- ✅ **Full inheritance with method lookup**
- ✅ **Class methods (factory methods)**
- ✅ **Class variables (shared state)**
- ✅ **Complete super message send support**
- ✅ Blocks and closures
- ✅ Arrays and dictionary literals
- ✅ Cascading messages
- ✅ Self keyword
- ✅ Control flow primitives (ifTrue:, ifFalse:, timesRepeat:, do:)
- ✅ Interactive REPL
- ✅ Interactive debugger with breakpoints and stepping
- ✅ Stack traces for runtime errors
- ✅ Comprehensive documentation (teaching-quality comments)
- ✅ Extensive test suite (85+ tests, benchmarks)

### Version History
- **v0.1.0**: Foundation - project structure and documentation
- **v0.2.0**: Core language features - variables, message sends, primitives
- **v0.3.0**: Blocks, arrays, control flow, extensive documentation
- **v0.4.0**: Enhanced features - super, cascading, dictionaries, REPL
- **v0.5.0**: Advanced classes - inheritance, class methods, class variables, debugger, stack traces

See the [roadmap](docs/planning/ROADMAP.md) for planned features and timeline.

## Building and Testing

```bash
# Build the interpreter
go build -o bin/smog ./cmd/smog

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run a specific example
./bin/smog examples/hello.smog
```

## Contributing

Contributions are welcome! Areas where you can help:

- Implementing core language features
- Writing tests
- Adding example programs
- Improving documentation
- Reporting bugs and suggesting features

## Inspiration

Smog is inspired by:
- [Smalltalk-80](https://rmod-files.lille.inria.fr/FreeBooks/BlueBook/Bluebook.pdf) - The original pure OO language
- [SOM (Simple Object Machine)](http://som-st.github.io/) - A minimal Smalltalk for teaching VMs
- [Pharo](https://pharo.org/) - Modern Smalltalk environment

## License

See [LICENSE](LICENSE) file for details.

## Learn More

### New to Smog?

1. **[Learning Guide](docs/LEARNING_GUIDE.md)** - Start with the mental model for beginners
2. **[User's Guide](docs/USERS_GUIDE.md)** - Learn through practical examples
3. **[Example Programs](examples/)** - Study working code
4. **[Language Specification](docs/spec/LANGUAGE_SPEC.md)** - Deep dive into syntax

### Want to Understand the Internals?

The compilation pipeline: **Source Code** → **Lexer** → **Parser** → **Compiler** → **Bytecode** → **VM**

- [Lexer](docs/LEXER.md) - Tokenization
- [Parser](docs/PARSER.md) - AST construction
- [Compiler](docs/COMPILER.md) - Bytecode generation
- [VM Deep Dive](docs/VM_DEEP_DIVE.md) - Execution engine

### Contributing?

- Check out the [Development Roadmap](docs/planning/ROADMAP.md)
- Read the [Architecture](docs/design/ARCHITECTURE.md) guide
- Review [Design Decisions](docs/design/DECISIONS.md)

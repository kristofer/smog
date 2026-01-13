# Smog

A very simple object-oriented language inspired by Smalltalk and SOM (Simple Object Machine).

## Overview

Smog is a minimalist object-oriented programming language that follows the philosophy that "everything is an object." It features:

- **Pure Object-Oriented Design**: Everything, including classes, numbers, and booleans, is an object
- **Message Passing**: All computation happens through sending messages to objects
- **Smalltalk-Inspired Syntax**: Clean, minimal syntax based on Smalltalk
- **Bytecode Compilation**: Source code compiles to bytecode for efficient execution
- **Stack-Based VM**: Simple virtual machine for bytecode execution

## Quick Start

### Building

```bash
# Build the smog interpreter
go build -o bin/smog ./cmd/smog

# Or use go run
go run ./cmd/smog [file.smog]
```

### Running Examples

```bash
# Run hello world
./bin/smog examples/hello.smog

# Run other examples
./bin/smog examples/factorial.smog
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
├── internal/          # Internal implementation details
├── docs/
│   ├── spec/          # Language specification
│   ├── design/        # Design documents
│   └── planning/      # Development planning
└── examples/          # Example programs
```

## Documentation

- **[Language Specification](docs/spec/LANGUAGE_SPEC.md)** - Complete language reference
- **[Architecture](docs/design/ARCHITECTURE.md)** - System architecture overview
- **[Design Decisions](docs/design/DECISIONS.md)** - Key design decisions and rationale
- **[Roadmap](docs/planning/ROADMAP.md)** - Development roadmap and milestones

## Development Status

**Current Version**: 0.3.0

Smog has completed its foundational implementation. Current features:
- ✅ Complete lexer and parser
- ✅ AST-based intermediate representation
- ✅ Bytecode compiler
- ✅ Stack-based virtual machine
- ✅ Blocks and closures
- ✅ Arrays and literals
- ✅ Control flow primitives (ifTrue:, ifFalse:, timesRepeat:, do:)
- ✅ Comprehensive documentation (teaching-quality comments)
- ✅ Extensive test suite (48+ tests)

### Version History
- **v0.1.0**: Foundation - project structure and documentation
- **v0.2.0**: Core language features - variables, message sends, primitives
- **v0.3.0**: Blocks, arrays, control flow, extensive documentation

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

- Read the [Language Specification](docs/spec/LANGUAGE_SPEC.md)
- Explore [Example Programs](examples/)
- Check out the [Development Roadmap](docs/planning/ROADMAP.md)

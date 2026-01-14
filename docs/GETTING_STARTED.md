# Getting Started with Smog Development

Welcome to the smog language project! This guide will help you get started with developing smog.

## Quick Overview

Smog is a simple object-oriented language inspired by Smalltalk and SOM (Simple Object Machine). The project is currently in early development (v0.1.0-dev), with the basic project structure and documentation in place.

## Prerequisites

- **Go**: Version 1.20 or later
- **Git**: For version control
- Basic understanding of:
  - Object-oriented programming
  - Compiler/interpreter concepts
  - Go programming language

## Getting Started

### 1. Clone and Build

```bash
# Clone the repository
git clone https://github.com/kristofer/smog.git
cd smog

# Build the interpreter
go build -o bin/smog ./cmd/smog

# Run tests (when available)
go test ./...
```

### 2. Try Running Examples

```bash
# Run source files directly
./bin/smog examples/hello.smog
./bin/smog examples/counter.smog

# Compile to bytecode for faster execution
./bin/smog compile examples/hello.smog examples/hello.sg

# Run compiled bytecode
./bin/smog examples/hello.sg

# Inspect bytecode
./bin/smog disassemble examples/hello.sg
```

### 3. Explore the Codebase

```
smog/
├── cmd/smog/           # Main executable - CLI interface
├── pkg/
│   ├── parser/        # TODO: Tokenize and parse source code
│   ├── ast/           # AST node definitions
│   ├── compiler/      # TODO: Compile AST to bytecode
│   ├── bytecode/      # Bytecode opcodes and format
│   └── vm/            # Virtual machine (stub with safety checks)
├── docs/              # All documentation
│   ├── spec/         # Language specification
│   ├── design/       # Architecture and design decisions
│   └── planning/     # Roadmap and planning
└── examples/          # Example smog programs
```

## Understanding the Project

### Essential Reading

1. **[Language Specification](spec/LANGUAGE_SPEC.md)** - Learn smog syntax and semantics
2. **[Architecture Overview](design/ARCHITECTURE.md)** - Understand system components
3. **[Development Roadmap](planning/ROADMAP.md)** - See what's planned
4. **[Design Decisions](design/DECISIONS.md)** - Understand why we made certain choices

### How Smog Works

```
Source Code (.smog)
    ↓
[Lexer] → Tokens
    ↓
[Parser] → AST
    ↓
[Compiler] → Bytecode
    ↓
[Save to .sg file (optional)]
    ↓
[VM] → Execution
```

**Two Execution Paths:**

1. **Direct Execution**: `.smog` → Lexer → Parser → Compiler → VM
2. **Compiled Execution**: `.smog` → Compiler → `.sg` file → VM (faster)

## Development Workflow

### Making Changes

1. **Choose a Task**: Look at the roadmap or open issues
2. **Create a Branch**: `git checkout -b feature/your-feature`
3. **Make Changes**: Implement your feature
4. **Test**: Write and run tests
5. **Commit**: `git commit -m "Description"`
6. **Push**: `git push origin feature/your-feature`
7. **Pull Request**: Create PR for review

### Code Style

- Follow Go conventions and idioms
- Run `go fmt ./...` before committing
- Add comments for exported types and functions
- Keep functions small and focused

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/parser
```

## Current Status and Next Steps

### Completed ✅
- Project structure
- Documentation foundation
- Language specification
- Example programs
- Full lexer implementation
- Complete parser
- Bytecode compiler
- Stack-based VM
- **Bytecode file format (.sg)**
- Classes and inheritance
- Blocks and closures
- Standard primitives

### Next Priority (v0.6.0) 🚧
1. **Module System**: Import/export for .sg files
2. **Standard Library**: Comprehensive collection classes
3. **Exception Handling**: Try/catch mechanisms
4. **Debugging Support**: Source line mapping in .sg files
5. **Performance**: JIT compilation optimizations

See the [roadmap](planning/ROADMAP.md) for detailed plans.

## How to Contribute

### Areas Needing Help

1. **Core Implementation**
   - Lexer/tokenizer
   - Parser
   - Compiler
   - Runtime objects

2. **Testing**
   - Unit tests for components
   - Integration tests
   - Test utilities

3. **Documentation**
   - Tutorial writing
   - API documentation
   - Example programs

4. **Tooling**
   - Syntax highlighting
   - Editor support
   - REPL

### Good First Issues

- Add more example programs
- Improve error messages
- Write documentation
- Add unit tests

## Learning Resources

### About Smalltalk
- [Smalltalk-80: The Language](https://rmod-files.lille.inria.fr/FreeBooks/BlueBook/Bluebook.pdf)
- [Pharo by Example](https://books.pharo.org/)

### About Interpreters/Compilers
- [Crafting Interpreters](https://craftinginterpreters.com/)
- [Writing An Interpreter In Go](https://interpreterbook.com/)

### About SOM
- [SOM Homepage](http://som-st.github.io/)
- [SOM Implementations](https://github.com/SOM-st)

### Common Tasks

### Compiling and Running Programs

**Compile source to bytecode:**
```bash
# Compile with auto-generated output name
smog compile program.smog

# Compile with custom output name
smog compile program.smog output.sg
```

**Run bytecode:**
```bash
# VM automatically detects .sg files
smog program.sg
```

**Inspect bytecode:**
```bash
# View disassembled bytecode
smog disassemble program.sg
```

### Adding a New Bytecode Opcode

1. Add opcode to `pkg/bytecode/bytecode.go`
2. Implement in `pkg/vm/vm.go`
3. Add compiler support in `pkg/compiler/compiler.go`
4. Write tests
5. Update documentation

### Adding a New AST Node

1. Define node in `pkg/ast/ast.go`
2. Implement interfaces (Node, Expression/Statement)
3. Add parser support in `pkg/parser/parser.go`
4. Add compiler support
5. Write tests

### Adding an Example Program

1. Create `.smog` file in `examples/`
2. Add description to `examples/README.md`
3. Test when interpreter is ready
4. Document expected output

## Getting Help

- **Documentation**: Start with docs in the `docs/` directory
- **Examples**: Look at existing examples
- **Code**: Read package documentation and tests
- **Issues**: Open an issue on GitHub

## Project Philosophy

- **Simplicity**: Keep the language and implementation simple
- **Clarity**: Code should be easy to understand
- **Testing**: Test everything
- **Documentation**: Document decisions and design
- **Incrementality**: Small, working steps

## Tips for Success

1. **Start Small**: Begin with simple features
2. **Read First**: Understand before changing
3. **Test Often**: Write tests as you go
4. **Ask Questions**: Don't hesitate to ask for clarification
5. **Have Fun**: Enjoy the journey of building a language!

## Roadmap Milestones

### Milestone 1: Hello World ⏳
Goal: Execute `'Hello, World!' println.`

### Milestone 2: Arithmetic ⏳
Goal: Support `3 + 4 * 5.`

### Milestone 3: Objects ⏳
Goal: Define and use classes

### Milestone 4: Closures ⏳
Goal: Support blocks and higher-order functions

### Milestone 5: Self-Hosting ⏳
Goal: Smog compiler in Smog (aspirational)

## Contact and Community

- **Repository**: https://github.com/kristofer/smog
- **Issues**: Use GitHub issues for bugs and features
- **Discussions**: Use GitHub discussions for questions

---

**Welcome to the smog community! We're excited to have you here.** 🎉
